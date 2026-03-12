package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// authenticateUPDB performs the username/password auth path.
func authenticateUPDB(cfg tools.HandlerConfig, s *store.Store) (*http.Client, string, error) {
	username := s.UpdbUsername()
	password := s.UpdbPassword()

	if username == "" || password == "" {
		return nil, "", fmt.Errorf("UPDB credentials are incomplete")
	}

	hint := credentialFingerprint(tools.AuthModeUPDB, username)
	if cached, ok := getCachedSession(cfg.Profile, cfg, hint); ok {
		slog.Debug("using cached zt-session (UPDB mode)")
		return cached.httpClient, cached.ztSession, nil
	}

	// Build transport (with optional CA trust)
	transport, err := buildTransport(s.ControllerCA(), "", "", "")
	if err != nil {
		return nil, "", err
	}

	// Authenticate via password method
	authURL := fmt.Sprintf("https://%s/edge/v1/authenticate?method=password", cfg.ZitiControllerHost)
	slog.Debug("authenticating with UPDB", "url", authURL)

	body, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("UPDB authentication: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	slog.Debug("UPDB auth response", "status", resp.StatusCode, "body", string(respBody))

	if resp.StatusCode >= 400 {
		if errMsg := checkZitiError(respBody); errMsg != "" {
			return nil, "", fmt.Errorf("UPDB authentication rejected: %s", errMsg)
		}
		return nil, "", fmt.Errorf("UPDB authentication failed (%d): %s", resp.StatusCode, string(respBody))
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
