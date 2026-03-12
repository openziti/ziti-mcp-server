package terminal

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"strings"
)

// GetTenantFromToken extracts the issuer (tenant/domain) from a JWT access token
// by decoding the payload without verification.
func GetTenantFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		slog.Debug("token is not a valid JWT (expected 3 parts)")
		return ""
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		slog.Debug("failed to decode JWT payload", "error", err)
		return ""
	}

	var claims struct {
		Iss string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		slog.Debug("failed to parse JWT claims", "error", err)
		return ""
	}

	// Strip protocol prefix
	tenant := strings.TrimPrefix(claims.Iss, "https://")
	tenant = strings.TrimPrefix(tenant, "http://")
	tenant = strings.TrimSuffix(tenant, "/")

	return tenant
}
