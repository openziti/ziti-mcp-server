package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

// authenticateToken performs the Bearer token (ext-jwt) auth path.
// Returns an *http.Client with zt-session injected and the zt-session value.
func authenticateToken(cfg tools.HandlerConfig, token string, s *store.Store) (*http.Client, string, error) {
	if token == "" {
		return nil, "", fmt.Errorf("missing authorization token")
	}
	if cfg.Domain == "" {
		return nil, "", fmt.Errorf("IdP domain is not configured")
	}

	hint := credentialFingerprint(tools.AuthModeToken, lastN(token, 8))
	if cached, ok := getCachedSession(cfg.Profile, cfg, hint); ok {
		slog.Debug("using cached zt-session (token mode)")
		return cached.httpClient, cached.ztSession, nil
	}

	// Build base transport (with optional CA trust)
	transport, err := buildTransport(s.ControllerCA(), "", "", "")
	if err != nil {
		return nil, "", err
	}

	// Authenticate via ext-jwt
	authURL := fmt.Sprintf("https://%s/edge/management/v1/authenticate?method=ext-jwt", cfg.ZitiControllerHost)
	slog.Debug("authenticating", "url", authURL)

	req, err := http.NewRequest(http.MethodPost, authURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("authenticating: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	slog.Debug("auth response", "status", resp.StatusCode, "body", string(body))

	// Check for Ziti API-level errors
	if errMsg := checkZitiError(body); errMsg != "" {
		return nil, "", fmt.Errorf("authentication rejected by controller: %s", errMsg)
	}

	ztSession := resp.Header.Get("zt-session")

	// Build client with zt-session header injection
	client := &http.Client{
		Transport: &ztSessionTransport{base: transport, ztSession: ztSession},
	}

	if ztSession != "" {
		storeSession(cfg.Profile, cfg, hint, client, ztSession)
	}

	return client, ztSession, nil
}

// buildTransport creates the appropriate transport based on available credentials.
func buildTransport(controllerCA, certPEM, keyPEM, caPEM string) (http.RoundTripper, error) {
	// mTLS has highest priority
	if certPEM != "" && keyPEM != "" && caPEM != "" {
		return NewMTLSTransport(certPEM, keyPEM, caPEM)
	}
	// Controller CA trust
	if controllerCA != "" {
		return NewCATrustTransport(controllerCA)
	}
	// Default transport
	return http.DefaultTransport, nil
}

// checkZitiError checks for a Ziti API-level error in the response body.
func checkZitiError(body []byte) string {
	var result struct {
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}
	if result.Error != nil {
		code := result.Error.Code
		if code == "" {
			code = "UNKNOWN"
		}
		msg := result.Error.Message
		if msg == "" {
			msg = "Authentication failed"
		}
		return fmt.Sprintf("(%s): %s", code, msg)
	}
	return ""
}

// lastN returns the last n characters of a string, or the whole string if shorter.
func lastN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

