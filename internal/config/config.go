package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// ZitiConfig represents the essential connection info needed for API operations.
type ZitiConfig struct {
	ZitiControllerHost string
	Token              string
	Domain             string
	TenantName         string
	AuthMode           tools.AuthMode
}

// LoadConfig reads stored credentials and auto-detects the auth mode.
// Priority: identity cert > UPDB username > token.
func LoadConfig(s *store.Store) *ZitiConfig {
	host := s.ControllerHost()

	// Identity mode
	if s.IdentityCert() != "" {
		slog.Debug("detected identity mode (certificate material found)")
		return &ZitiConfig{
			ZitiControllerHost: host,
			AuthMode:           tools.AuthModeIdentity,
			TenantName:         "identity",
		}
	}

	// UPDB mode
	if s.UpdbUsername() != "" {
		slog.Debug("detected UPDB mode (username/password found)")
		return &ZitiConfig{
			ZitiControllerHost: host,
			AuthMode:           tools.AuthModeUPDB,
			TenantName:         "updb",
		}
	}

	// Token mode
	token := s.Token()
	if IsTokenExpired(s, 300) {
		slog.Debug("token is expired or expiring soon")
		token = ""
	}
	domain := s.Domain()
	tenantName := domain
	if tenantName == "" {
		tenantName = "default"
	}

	return &ZitiConfig{
		ZitiControllerHost: host,
		Token:              token,
		Domain:             domain,
		TenantName:         tenantName,
		AuthMode:           tools.AuthModeToken,
	}
}

// ValidateConfig checks that the config has everything needed for its auth mode.
func ValidateConfig(cfg *ZitiConfig, s *store.Store) bool {
	if cfg == nil {
		slog.Debug("configuration is nil")
		return false
	}

	if cfg.ZitiControllerHost == "" {
		slog.Debug("ziti controller host is missing")
		return false
	}

	switch cfg.AuthMode {
	case tools.AuthModeIdentity:
		if s.IdentityCert() == "" {
			slog.Debug("identity certificate is missing")
			return false
		}
		return true

	case tools.AuthModeUPDB:
		if s.UpdbUsername() == "" {
			slog.Debug("UPDB username is missing")
			return false
		}
		if s.UpdbPassword() == "" {
			slog.Debug("UPDB password is missing")
			return false
		}
		return true

	default: // token mode
		if cfg.Token == "" {
			slog.Debug("token is missing")
			return false
		}
		if cfg.Domain == "" {
			slog.Debug("domain is missing")
			return false
		}
		if IsTokenExpired(s, 300) {
			slog.Debug("token is expired")
			return false
		}
		return true
	}
}

// LoadConfigOrNil loads and validates config, returning nil if invalid.
// Used by the server to detect disconnected state.
func LoadConfigOrNil(s *store.Store) *ZitiConfig {
	cfg := LoadConfig(s)
	if !ValidateConfig(cfg, s) {
		return nil
	}
	return cfg
}

// IsTokenExpired checks whether the stored token has expired (with a buffer in seconds).
func IsTokenExpired(s *store.Store, bufferSecs int64) bool {
	expiresAt := s.TokenExpiresAt()
	if expiresAt == 0 {
		return true
	}
	now := time.Now().UnixMilli()
	return now+bufferSecs*1000 >= expiresAt
}

// IsTokenExpiredValue checks whether a token expiry value (as string millis) has expired.
func IsTokenExpiredValue(expiresAtStr string, bufferSecs int64) bool {
	if expiresAtStr == "" {
		return true
	}
	var expiresAt int64
	_, _ = fmt.Sscanf(expiresAtStr, "%d", &expiresAt)
	if expiresAt == 0 {
		return true
	}
	now := time.Now().UnixMilli()
	return now+bufferSecs*1000 >= expiresAt
}

// LoadConfigFromData creates a ZitiConfig from raw profile data (for display/status).
func LoadConfigFromData(data map[string]string) *ZitiConfig {
	host := data[store.KeyControllerHost]
	if data[store.KeyIdentityCert] != "" {
		return &ZitiConfig{
			ZitiControllerHost: host,
			AuthMode:           tools.AuthModeIdentity,
		}
	}
	if data[store.KeyUpdbUsername] != "" {
		return &ZitiConfig{
			ZitiControllerHost: host,
			AuthMode:           tools.AuthModeUPDB,
		}
	}
	return &ZitiConfig{
		ZitiControllerHost: host,
		Token:              data[store.KeyToken],
		Domain:             data[store.KeyDomain],
		AuthMode:           tools.AuthModeToken,
	}
}
