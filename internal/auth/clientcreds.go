package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

// ClientCredentialsConfig holds parameters for client credentials flow.
type ClientCredentialsConfig struct {
	ZitiControllerHost string
	IDPDomain          string
	IDPClientID        string
	IDPClientSecret    string
	Audience           string
	Scopes             []string
}

// RequestClientCredentialsAuthorization performs the client credentials flow.
func RequestClientCredentialsAuthorization(s *store.Store, cfg ClientCredentialsConfig) error {
	slog.Debug("initiating client credentials flow", "domain", cfg.IDPDomain)

	tokenEndpoint, err := GetTokenEndpoint(cfg.IDPDomain)
	if err != nil {
		return fmt.Errorf("OIDC discovery failed: %w", err)
	}

	slog.Debug("using token endpoint", "url", tokenEndpoint)

	audience := cfg.Audience
	if audience == "" {
		audience = cfg.IDPDomain + "/api/v2/"
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {cfg.IDPClientID},
		"client_secret": {cfg.IDPClientSecret},
		"audience":      {audience},
	}

	if len(cfg.Scopes) > 0 {
		form.Set("scope", strings.Join(cfg.Scopes, " "))
	}

	resp, err := http.PostForm(tokenEndpoint, form)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token,omitempty"`
		ExpiresIn        int64  `json:"expires_in,omitempty"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decoding token response: %w", err)
	}

	if tokenResp.Error != "" {
		return fmt.Errorf("token request failed: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("token response did not contain an access_token")
	}

	// Store credentials
	if err := s.SetToken(tokenResp.AccessToken); err != nil {
		return err
	}
	if err := s.SetControllerHost(cfg.ZitiControllerHost); err != nil {
		return err
	}
	if err := s.SetDomain(cfg.IDPDomain); err != nil {
		return err
	}

	if tokenResp.RefreshToken != "" {
		if err := s.SetRefreshToken(tokenResp.RefreshToken); err != nil {
			return err
		}
	}

	if tokenResp.ExpiresIn > 0 {
		expiresAt := time.Now().UnixMilli() + tokenResp.ExpiresIn*1000
		if err := s.SetTokenExpiresAt(expiresAt); err != nil {
			return err
		}
	}

	terminal.Success("Successfully authenticated to %s using client credentials.", cfg.IDPDomain)
	return nil
}
