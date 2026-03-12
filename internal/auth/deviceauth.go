package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/browser"

	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/terminal"
)

// DeviceAuthConfig holds parameters for the device authorization flow.
type DeviceAuthConfig struct {
	IDPDomain string
	ClientID  string
	Audience  string
	Scopes    string
}

// DeviceCodeResult holds the result of requesting a device code (non-blocking step 1).
type DeviceCodeResult struct {
	DeviceCode                  string
	UserCode                    string
	VerificationURI             string
	VerificationURIComplete     string
	TokenEndpoint               string
	ClientID                    string
	Interval                    int
	ExpiresIn                   int
}

// RequestDeviceAuthorization runs the full device authorization flow:
// discover endpoints, request device code, open browser, poll for token.
func RequestDeviceAuthorization(s *store.Store, cfg DeviceAuthConfig) error {
	slog.Debug("discovering OIDC endpoints", "domain", cfg.IDPDomain)
	doc, err := FetchOIDCDiscoveryDocument(cfg.IDPDomain)
	if err != nil {
		return fmt.Errorf("OIDC discovery failed: %w", err)
	}

	if doc.DeviceAuthorizationEndpoint == "" {
		return fmt.Errorf("IdP %s does not support device authorization flow", cfg.IDPDomain)
	}

	// Request device code
	form := url.Values{"client_id": {cfg.ClientID}}
	if cfg.Audience != "" {
		form.Set("audience", cfg.Audience)
	}
	if cfg.Scopes != "" {
		form.Set("scope", cfg.Scopes)
	}

	slog.Debug("requesting device code", "endpoint", doc.DeviceAuthorizationEndpoint)
	resp, err := http.PostForm(doc.DeviceAuthorizationEndpoint, form)
	if err != nil {
		return fmt.Errorf("device code request failed: %w", err)
	}
	defer resp.Body.Close()

	var deviceResp struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		Interval                int    `json:"interval"`
		ExpiresIn               int    `json:"expires_in"`
		Error                   string `json:"error"`
		ErrorDescription        string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return fmt.Errorf("decoding device code response: %w", err)
	}

	if deviceResp.Error != "" {
		return fmt.Errorf("device code error: %s - %s", deviceResp.Error, deviceResp.ErrorDescription)
	}

	verifyURL := deviceResp.VerificationURIComplete
	if verifyURL == "" {
		verifyURL = deviceResp.VerificationURI
	}

	terminal.Info("")
	terminal.Success("Verify this code on screen: %s", deviceResp.UserCode)
	terminal.Info("Verify at this URL: %s", verifyURL)
	terminal.Info("")

	if err := browser.OpenURL(verifyURL); err != nil {
		slog.Debug("failed to open browser", "error", err)
	}

	// Poll for token
	interval := time.Duration(deviceResp.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second
	}

	terminal.Info("Waiting for authorization...")

	for {
		time.Sleep(interval)

		tokenResp, err := pollForToken(doc.TokenEndpoint, deviceResp.DeviceCode, cfg.ClientID)
		if err != nil {
			return err
		}

		if tokenResp.Error == "" {
			// Success
			if err := storeTokens(s, tokenResp); err != nil {
				return fmt.Errorf("storing tokens: %w", err)
			}
			terminal.Success("Successfully authenticated!")
			return nil
		}

		switch tokenResp.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval *= 2
			continue
		case "access_denied":
			desc := tokenResp.ErrorDescription
			if desc == "" {
				desc = "User denied authorization or IdP configuration issue"
			}
			return fmt.Errorf("access denied: %s", desc)
		case "expired_token":
			return fmt.Errorf("device code expired, please try again")
		default:
			return fmt.Errorf("unexpected error: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
		}
	}
}

// RequestDeviceCode performs step 1 of the device auth flow: requests a device code
// without blocking. Returns the result for the caller to present to the user.
func RequestDeviceCode(cfg DeviceAuthConfig) (*DeviceCodeResult, error) {
	slog.Debug("discovering OIDC endpoints for device code", "domain", cfg.IDPDomain)
	doc, err := FetchOIDCDiscoveryDocument(cfg.IDPDomain)
	if err != nil {
		return nil, fmt.Errorf("OIDC discovery failed: %w", err)
	}

	if doc.DeviceAuthorizationEndpoint == "" {
		return nil, fmt.Errorf("IdP %s does not support device authorization flow", cfg.IDPDomain)
	}

	form := url.Values{"client_id": {cfg.ClientID}}
	if cfg.Audience != "" {
		form.Set("audience", cfg.Audience)
	}
	if cfg.Scopes != "" {
		form.Set("scope", cfg.Scopes)
	}

	resp, err := http.PostForm(doc.DeviceAuthorizationEndpoint, form)
	if err != nil {
		return nil, fmt.Errorf("device code request failed: %w", err)
	}
	defer resp.Body.Close()

	var deviceResp struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		Interval                int    `json:"interval"`
		ExpiresIn               int    `json:"expires_in"`
		Error                   string `json:"error"`
		ErrorDescription        string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("decoding device code response: %w", err)
	}

	if deviceResp.Error != "" {
		return nil, fmt.Errorf("device code error: %s - %s", deviceResp.Error, deviceResp.ErrorDescription)
	}

	return &DeviceCodeResult{
		DeviceCode:              deviceResp.DeviceCode,
		UserCode:                deviceResp.UserCode,
		VerificationURI:         deviceResp.VerificationURI,
		VerificationURIComplete: deviceResp.VerificationURIComplete,
		TokenEndpoint:           doc.TokenEndpoint,
		ClientID:                cfg.ClientID,
		Interval:                deviceResp.Interval,
		ExpiresIn:               deviceResp.ExpiresIn,
	}, nil
}

// PollDeviceTokenOnce polls the token endpoint once (non-blocking).
// Returns the token response. Check .Error for "authorization_pending".
func PollDeviceTokenOnce(tokenEndpoint, deviceCode, clientID string) (*TokenResponse, error) {
	tr, err := pollForToken(tokenEndpoint, deviceCode, clientID)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:      tr.AccessToken,
		RefreshToken:     tr.RefreshToken,
		ExpiresIn:        tr.ExpiresIn,
		Error:            tr.Error,
		ErrorDescription: tr.ErrorDescription,
	}, nil
}

// TokenResponse is the exported version of tokenResponse for use by meta-tools.
type TokenResponse struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int64
	Error            string
	ErrorDescription string
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func pollForToken(tokenEndpoint, deviceCode, clientID string) (*tokenResponse, error) {
	form := url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {deviceCode},
		"client_id":   {clientID},
	}

	resp, err := http.PostForm(tokenEndpoint, form)
	if err != nil {
		return nil, fmt.Errorf("polling token endpoint: %w", err)
	}
	defer resp.Body.Close()

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	return &tr, nil
}

// StoreTokens saves token data from a token response into the credential store.
func StoreTokens(s *store.Store, tr *TokenResponse) error {
	tenant := terminal.GetTenantFromToken(tr.AccessToken)

	if err := s.SetToken(tr.AccessToken); err != nil {
		return err
	}
	if err := s.SetDomain(tenant); err != nil {
		return err
	}

	if tr.RefreshToken != "" {
		if err := s.SetRefreshToken(tr.RefreshToken); err != nil {
			return err
		}
		slog.Debug("refresh token stored")
	}

	if tr.ExpiresIn > 0 {
		expiresAt := time.Now().UnixMilli() + tr.ExpiresIn*1000
		if err := s.SetTokenExpiresAt(expiresAt); err != nil {
			return err
		}
		slog.Debug("token expires at", "time", time.UnixMilli(expiresAt).Format(time.RFC3339))
	}

	return nil
}

// storeTokens is the internal version using the unexported tokenResponse.
func storeTokens(s *store.Store, tr *tokenResponse) error {
	return StoreTokens(s, &TokenResponse{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresIn:    tr.ExpiresIn,
	})
}

// RevokeRefreshToken revokes the stored refresh token via the IdP's revocation endpoint.
func RevokeRefreshToken(s *store.Store, idpClientID string) error {
	refreshToken := s.RefreshToken()
	if refreshToken == "" {
		slog.Debug("no refresh token to revoke")
		return nil
	}

	domain := s.Domain()
	if domain == "" {
		slog.Debug("no domain for token revocation")
		return nil
	}

	if idpClientID == "" {
		slog.Debug("client ID required for revocation")
		return nil
	}

	revocationEndpoint, err := GetRevocationEndpoint(domain)
	if err != nil || revocationEndpoint == "" {
		slog.Debug("IdP does not support token revocation")
		return nil
	}

	form := url.Values{
		"client_id": {idpClientID},
		"token":     {refreshToken},
	}

	resp, err := http.PostForm(revocationEndpoint, form)
	if err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		slog.Debug("refresh token revoked")
	} else {
		slog.Debug("revocation endpoint returned", "status", resp.StatusCode)
	}

	return nil
}

// RefreshAccessToken uses a stored refresh token to obtain a new access token.
func RefreshAccessToken(s *store.Store, idpClientID string) (string, error) {
	refreshToken := s.RefreshToken()
	if refreshToken == "" {
		return "", fmt.Errorf("no refresh token found")
	}

	domain := s.Domain()
	if domain == "" {
		return "", fmt.Errorf("no domain for token refresh")
	}

	if idpClientID == "" {
		return "", fmt.Errorf("client ID required for token refresh")
	}

	tokenEndpoint, err := GetTokenEndpoint(domain)
	if err != nil {
		return "", err
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {idpClientID},
		"refresh_token": {refreshToken},
	}

	resp, err := http.PostForm(tokenEndpoint, form)
	if err != nil {
		return "", fmt.Errorf("refreshing token: %w", err)
	}
	defer resp.Body.Close()

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("decoding refresh response: %w", err)
	}

	if tr.Error != "" {
		return "", fmt.Errorf("refresh error: %s", tr.Error)
	}

	if err := storeTokens(s, &tr); err != nil {
		return "", err
	}

	return tr.AccessToken, nil
}

// formatDomain normalizes a domain string.
func formatDomain(domain string) string {
	if domain == "" {
		return ""
	}
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")
	if !strings.Contains(domain, ".") {
		return domain + ".your-tenant.auth0.com"
	}
	return domain
}
