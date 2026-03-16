package client

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/openziti/ziti-mcp-server/internal/tools"
)

const sessionTTL = 30 * time.Minute

type sessionEntry struct {
	ztSession      string
	httpClient     *http.Client
	authMode       tools.AuthMode
	controllerHost string
	credentialHint string
	acquiredAt     time.Time
}

var (
	sessionMu    sync.Mutex
	sessionCache = make(map[string]*sessionEntry)
)

// credentialFingerprint builds a short fingerprint that changes when credentials change.
func credentialFingerprint(authMode tools.AuthMode, hint string) string {
	return fmt.Sprintf("%s:%s", authMode, hint)
}

// getCachedSession returns a cached session if it matches the config and is still valid.
func getCachedSession(profile string, cfg tools.HandlerConfig, credHint string) (*sessionEntry, bool) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	entry := sessionCache[profile]
	if entry == nil {
		return nil, false
	}
	if entry.authMode != cfg.AuthMode {
		return nil, false
	}
	if entry.controllerHost != cfg.ZitiControllerHost {
		return nil, false
	}
	if entry.credentialHint != credHint {
		return nil, false
	}
	if time.Since(entry.acquiredAt) > sessionTTL {
		slog.Debug("session cache expired (TTL), will re-authenticate", "profile", profile)
		delete(sessionCache, profile)
		return nil, false
	}

	return entry, true
}

// storeSession caches a session entry keyed by profile name.
func storeSession(profile string, cfg tools.HandlerConfig, credHint string, httpClient *http.Client, ztSession string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	sessionCache[profile] = &sessionEntry{
		ztSession:      ztSession,
		httpClient:     httpClient,
		authMode:       cfg.AuthMode,
		controllerHost: cfg.ZitiControllerHost,
		credentialHint: credHint,
		acquiredAt:     time.Now(),
	}
}

// InvalidateSession clears the cached session for a specific profile.
func InvalidateSession(profile string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	delete(sessionCache, profile)
}

// InvalidateAllSessions clears all cached sessions.
func InvalidateAllSessions() {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	sessionCache = make(map[string]*sessionEntry)
}

// ztSessionTransport is an http.RoundTripper that injects the zt-session header.
type ztSessionTransport struct {
	base      http.RoundTripper
	ztSession string
}

func (t *ztSessionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	if t.ztSession != "" {
		req.Header.Set("zt-session", t.ztSession)
	}
	return t.base.RoundTrip(req)
}
