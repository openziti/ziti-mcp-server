package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// OIDCDiscoveryDocument represents the OpenID Connect discovery response.
type OIDCDiscoveryDocument struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint,omitempty"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint,omitempty"`
	JwksURI                          string   `json:"jwks_uri,omitempty"`
	RegistrationEndpoint             string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                  []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported           []string `json:"response_types_supported,omitempty"`
	GrantTypesSupported              []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported            []string `json:"subject_types_supported,omitempty"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`
	TokenEndpointAuthMethods         []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	RevocationEndpoint               string   `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint            string   `json:"introspection_endpoint,omitempty"`
	DeviceAuthorizationEndpoint      string   `json:"device_authorization_endpoint,omitempty"`
}

type cachedDoc struct {
	document  OIDCDiscoveryDocument
	expiresAt time.Time
}

var (
	discoveryMu    sync.RWMutex
	discoveryCache = make(map[string]cachedDoc)
	cacheTTL       = 5 * time.Minute
)

// FetchOIDCDiscoveryDocument fetches and caches the OIDC discovery document.
func FetchOIDCDiscoveryDocument(idpDomain string) (*OIDCDiscoveryDocument, error) {
	discoveryMu.RLock()
	if cached, ok := discoveryCache[idpDomain]; ok && time.Now().Before(cached.expiresAt) {
		discoveryMu.RUnlock()
		slog.Debug("using cached OIDC discovery", "domain", idpDomain)
		return &cached.document, nil
	}
	discoveryMu.RUnlock()

	normalizedDomain := strings.TrimPrefix(strings.TrimPrefix(idpDomain, "https://"), "http://")
	url := fmt.Sprintf("https://%s/.well-known/openid-configuration", normalizedDomain)
	slog.Debug("fetching OIDC discovery", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching OIDC discovery from %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OIDC discovery returned %d %s", resp.StatusCode, resp.Status)
	}

	var doc OIDCDiscoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decoding OIDC discovery: %w", err)
	}

	if doc.Issuer == "" {
		return nil, fmt.Errorf("OIDC discovery missing required field: issuer")
	}
	if doc.TokenEndpoint == "" {
		return nil, fmt.Errorf("OIDC discovery missing required field: token_endpoint")
	}

	discoveryMu.Lock()
	discoveryCache[idpDomain] = cachedDoc{document: doc, expiresAt: time.Now().Add(cacheTTL)}
	discoveryMu.Unlock()

	slog.Debug("OIDC discovery fetched", "domain", idpDomain)
	return &doc, nil
}

// GetTokenEndpoint returns the token endpoint from OIDC discovery.
func GetTokenEndpoint(idpDomain string) (string, error) {
	doc, err := FetchOIDCDiscoveryDocument(idpDomain)
	if err != nil {
		return "", err
	}
	return doc.TokenEndpoint, nil
}

// GetDeviceAuthorizationEndpoint returns the device auth endpoint from discovery.
func GetDeviceAuthorizationEndpoint(idpDomain string) (string, error) {
	doc, err := FetchOIDCDiscoveryDocument(idpDomain)
	if err != nil {
		return "", err
	}
	if doc.DeviceAuthorizationEndpoint == "" {
		return "", fmt.Errorf("IdP %s does not support device authorization flow", idpDomain)
	}
	return doc.DeviceAuthorizationEndpoint, nil
}

// GetRevocationEndpoint returns the revocation endpoint, or "" if not supported.
func GetRevocationEndpoint(idpDomain string) (string, error) {
	doc, err := FetchOIDCDiscoveryDocument(idpDomain)
	if err != nil {
		return "", err
	}
	return doc.RevocationEndpoint, nil
}

// ClearDiscoveryCache empties the OIDC discovery cache.
func ClearDiscoveryCache() {
	discoveryMu.Lock()
	discoveryCache = make(map[string]cachedDoc)
	discoveryMu.Unlock()
}
