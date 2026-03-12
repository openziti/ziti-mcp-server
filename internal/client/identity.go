package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// authenticateIdentity performs the mTLS (certificate) auth path.
func authenticateIdentity(cfg tools.HandlerConfig, s *store.Store) (*http.Client, string, error) {
	cert := s.IdentityCert()
	key := s.IdentityKey()
	ca := s.IdentityCA()

	if cert == "" || key == "" || ca == "" {
		return nil, "", fmt.Errorf("identity certificate material is incomplete")
	}

	hint := credentialFingerprint(tools.AuthModeIdentity, lastN(cert, 16))
	if cached, ok := getCachedSession(cfg.Profile, cfg, hint); ok {
		slog.Debug("using cached zt-session (identity mode)")
		return cached.httpClient, cached.ztSession, nil
	}

	// Build mTLS transport
	transport, err := NewMTLSTransport(cert, key, ca)
	if err != nil {
		return nil, "", fmt.Errorf("creating mTLS transport: %w", err)
	}

	// Authenticate via cert method
	authURL := fmt.Sprintf("https://%s/edge/management/v1/authenticate?method=cert", cfg.ZitiControllerHost)
	slog.Debug("authenticating with certificate", "url", authURL)

	req, err := http.NewRequest(http.MethodPost, authURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("certificate authentication: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	slog.Debug("cert auth response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("certificate authentication failed (%d): %s", resp.StatusCode, string(body))
	}

	if errMsg := checkZitiError(body); errMsg != "" {
		return nil, "", fmt.Errorf("certificate authentication rejected: %s", errMsg)
	}

	ztSession := resp.Header.Get("zt-session")

	client := &http.Client{
		Transport: &ztSessionTransport{base: transport, ztSession: ztSession},
	}

	if ztSession != "" {
		storeSession(cfg.Profile, cfg, hint, client, ztSession)
	}

	return client, ztSession, nil
}
