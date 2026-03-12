package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/auth"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// APICall is a function that performs an API operation with an authenticated client.
// The httpClient has the zt-session header automatically injected.
type APICall func(httpClient *http.Client, ztSession string) (any, error)

// WithAuthenticatedClient dispatches to the appropriate auth path (token, identity, UPDB),
// obtains a zt-session, and executes the API call.
func WithAuthenticatedClient(
	request tools.HandlerRequest,
	cfg tools.HandlerConfig,
	operationName string,
	s *store.Store,
	apiCall APICall,
) tools.HandlerResponse {
	if cfg.ZitiControllerHost == "" {
		return CreateErrorResponse("Error: Ziti Controller host is not configured")
	}

	httpClient, ztSession, err := authenticate(request, cfg, s)
	if err != nil {
		// For token mode, attempt refresh and retry once
		if cfg.AuthMode == tools.AuthModeToken {
			if refreshed := tryRefreshAndRetry(s); refreshed != "" {
				slog.Info("retrying authentication with refreshed token")
				request.Token = refreshed
				httpClient, ztSession, err = authenticate(request, cfg, s)
			}
		}
		if err != nil {
			slog.Error("authentication failed", "mode", cfg.AuthMode, "error", err)
			return CreateErrorResponse(fmt.Sprintf("Error: %s", err))
		}
	}

	result, err := apiCall(httpClient, ztSession)
	if err != nil {
		slog.Error("API call failed", "operation", operationName, "error", err)
		errMsg := fmt.Sprintf("Failed to %s: %s", operationName, err)
		return CreateErrorResponse(errMsg)
	}

	if result == nil {
		return CreateSuccessResponse(map[string]any{"message": "No results found.", "logs": []any{}})
	}

	return CreateSuccessResponse(result)
}

// authenticate dispatches to the appropriate auth path based on the config.
func authenticate(request tools.HandlerRequest, cfg tools.HandlerConfig, s *store.Store) (*http.Client, string, error) {
	switch cfg.AuthMode {
	case tools.AuthModeIdentity:
		return authenticateIdentity(cfg, s)
	case tools.AuthModeUPDB:
		return authenticateUPDB(cfg, s)
	default:
		return authenticateToken(cfg, request.Token, s)
	}
}

// tryRefreshAndRetry attempts to refresh the access token using the stored refresh token.
// Returns the new access token on success, or empty string on failure.
func tryRefreshAndRetry(s *store.Store) string {
	clientID := s.IDPClientID()
	if clientID == "" {
		return ""
	}
	if s.RefreshToken() == "" {
		return ""
	}

	slog.Debug("attempting token refresh after auth failure")
	InvalidateAllSessions()

	newToken, err := auth.RefreshAccessToken(s, clientID)
	if err != nil {
		slog.Debug("token refresh failed", "error", err)
		return ""
	}

	slog.Info("access token refreshed successfully")
	return newToken
}

// FormatDomain normalizes a domain string (removes protocol, trailing slash).
func FormatDomain(domain string) string {
	if domain == "" {
		return ""
	}
	// Remove protocol
	for _, prefix := range []string{"https://", "http://"} {
		if len(domain) > len(prefix) && domain[:len(prefix)] == prefix {
			domain = domain[len(prefix):]
			break
		}
	}
	// Remove trailing slash
	if domain[len(domain)-1] == '/' {
		domain = domain[:len(domain)-1]
	}
	// Add default auth0 suffix if no dots
	hasDot := false
	for _, c := range domain {
		if c == '.' {
			hasDot = true
			break
		}
	}
	if !hasDot {
		return domain + ".your-tenant.auth0.com"
	}
	return domain
}

// ParseJSONResponse parses a JSON response body into a generic structure.
func ParseJSONResponse(data []byte) (any, error) {
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return result, nil
}
