package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/openziti/ziti-mcp-server-go/internal/auth"
	"github.com/openziti/ziti-mcp-server-go/internal/ca"
	"github.com/openziti/ziti-mcp-server-go/internal/client"
	"github.com/openziti/ziti-mcp-server-go/internal/config"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// pendingDeviceAuth holds state for an in-progress device-auth login.
type pendingDeviceAuth struct {
	ProfileName   string
	DeviceCode    string
	TokenEndpoint string
	ClientID      string
	Interval      time.Duration
	ExpiresAt     time.Time
}

var (
	pendingAuthMu sync.Mutex
	pendingAuth   *pendingDeviceAuth
)

// RegisterMeta registers the meta-tools (login variants, completeLogin, logout, listNetworks, selectNetwork).
func RegisterMeta(r *tools.Registry, s *store.Store) {
	registerLoginUpdb(r, s)
	registerLoginIdentity(r, s)
	registerLoginClientCredentials(r, s)
	registerLoginDeviceAuth(r, s)
	registerMetaCompleteLogin(r, s)
	registerMetaLogout(r, s)
	registerMetaListNetworks(r, s)
	registerMetaSelectNetwork(r, s)
}

// --- login shared logic ---

// prepareLogin normalizes the controller host, invalidates old session, ensures profile, clears it,
// and fetches the controller CA for non-identity modes.
func prepareLogin(s *store.Store, network, controller string, fetchCA bool) (string, error) {
	// Normalize controller host (strip protocol, trailing slash — but NOT the auth0 suffix logic)
	if controller != "" {
		controller = normalizeControllerHost(controller)
	}

	// Invalidate session cache for the old active profile
	if old := s.ActiveProfile(); old != "" {
		client.InvalidateSession(old)
	}

	// Create profile if needed, set as active, clear existing data
	if err := s.EnsureProfile(network); err != nil {
		return "", fmt.Errorf("creating profile: %w", err)
	}
	s.ClearProfile(network)

	// Fetch controller CA
	if fetchCA && controller != "" {
		caPem, err := ca.FetchControllerCA(controller)
		if err != nil {
			slog.Debug("controller CA fetch failed (may use public cert)", "error", err)
		} else if caPem != "" {
			if err := s.SetControllerCA(caPem); err != nil {
				return "", fmt.Errorf("storing controller CA: %w", err)
			}
		}
	}

	return controller, nil
}

// --- loginUpdb ---

func registerLoginUpdb(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "loginUpdb",
		Description: "Connect to a Ziti network using username/password (UPDB) authentication. Creates or updates a named profile.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name for this network (e.g. \"prod\", \"staging\", \"local\")",
				},
				"controller": map[string]any{
					"type":        "string",
					"description": "Ziti controller host or IP with optional port (e.g. \"ctrl.example.com\" or \"192.168.1.1:1280\")",
				},
				"username": map[string]any{
					"type":        "string",
					"description": "UPDB username",
				},
				"password": map[string]any{
					"type":        "string",
					"description": "UPDB password",
				},
			},
			"required": []string{"network", "controller", "username", "password"},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Login with Username/Password",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, loginUpdbHandler(s))
}

func loginUpdbHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		controller := stringParam(req.Parameters, "controller")
		username := stringParam(req.Parameters, "username")
		password := stringParam(req.Parameters, "password")

		controller, err := prepareLogin(s, network, controller, true)
		if err != nil {
			return tools.HandlerResponse{}, err
		}

		if err := auth.StoreUPDBCredentials(s, controller, username, password); err != nil {
			return tools.HandlerResponse{}, err
		}

		return jsonResponse(map[string]any{
			"status":   "connected",
			"network":  network,
			"authMode": "updb",
		})
	}
}

// --- loginIdentity ---

func registerLoginIdentity(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "loginIdentity",
		Description: "Connect to a Ziti network using a Ziti identity JSON file (mTLS certificate authentication). Creates or updates a named profile. The controller host is extracted from the identity JSON.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name for this network (e.g. \"prod\", \"staging\")",
				},
				"identityJson": map[string]any{
					"type":        "string",
					"description": "Raw Ziti identity JSON content (the full JSON object with ztAPI, id.cert, id.key, id.ca)",
				},
			},
			"required": []string{"network", "identityJson"},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Login with Identity File",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, loginIdentityHandler(s))
}

func loginIdentityHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		identityJSON := stringParam(req.Parameters, "identityJson")

		// Identity mode doesn't need CA fetch — CA is embedded in the identity JSON
		if _, err := prepareLogin(s, network, "", false); err != nil {
			return tools.HandlerResponse{}, err
		}

		if err := auth.RequestIdentityJSONAuthorization(s, identityJSON); err != nil {
			return tools.HandlerResponse{}, err
		}

		return jsonResponse(map[string]any{
			"status":   "connected",
			"network":  network,
			"authMode": "identity",
		})
	}
}

// --- loginClientCredentials ---

func registerLoginClientCredentials(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "loginClientCredentials",
		Description: "Connect to a Ziti network using OAuth2 client credentials flow. Creates or updates a named profile.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name for this network (e.g. \"prod\", \"staging\")",
				},
				"controller": map[string]any{
					"type":        "string",
					"description": "Ziti controller host or IP with optional port (e.g. \"ctrl.example.com\" or \"192.168.1.1:1280\")",
				},
				"idpDomain": map[string]any{
					"type":        "string",
					"description": "Identity Provider domain (e.g. \"myorg.auth0.com\")",
				},
				"idpClientId": map[string]any{
					"type":        "string",
					"description": "OAuth2 client ID",
				},
				"idpClientSecret": map[string]any{
					"type":        "string",
					"description": "OAuth2 client secret",
				},
				"idpAudience": map[string]any{
					"type":        "string",
					"description": "OAuth2 audience (optional, defaults to idpDomain/api/v2/)",
				},
			},
			"required": []string{"network", "controller", "idpDomain", "idpClientId", "idpClientSecret"},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Login with Client Credentials",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, loginClientCredentialsHandler(s))
}

func loginClientCredentialsHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		controller := stringParam(req.Parameters, "controller")
		idpDomain := stringParam(req.Parameters, "idpDomain")
		idpClientID := stringParam(req.Parameters, "idpClientId")
		idpClientSecret := stringParam(req.Parameters, "idpClientSecret")
		idpAudience := stringParam(req.Parameters, "idpAudience")

		controller, err := prepareLogin(s, network, controller, true)
		if err != nil {
			return tools.HandlerResponse{}, err
		}

		if err := auth.RequestClientCredentialsAuthorization(s, auth.ClientCredentialsConfig{
			ZitiControllerHost: controller,
			IDPDomain:          idpDomain,
			IDPClientID:        idpClientID,
			IDPClientSecret:    idpClientSecret,
			Audience:           idpAudience,
		}); err != nil {
			return tools.HandlerResponse{}, err
		}
		if err := s.SetIDPClientID(idpClientID); err != nil {
			return tools.HandlerResponse{}, err
		}

		return jsonResponse(map[string]any{
			"status":   "connected",
			"network":  network,
			"authMode": "client-credentials",
		})
	}
}

// --- loginDeviceAuth ---

func registerLoginDeviceAuth(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "loginDeviceAuth",
		Description: "Connect to a Ziti network using OAuth2 device authorization flow. Returns a verification URL and code for the user to approve in their browser, then call completeLogin to finish.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name for this network (e.g. \"prod\", \"staging\")",
				},
				"controller": map[string]any{
					"type":        "string",
					"description": "Ziti controller host or IP with optional port (e.g. \"ctrl.example.com\" or \"192.168.1.1:1280\")",
				},
				"idpDomain": map[string]any{
					"type":        "string",
					"description": "Identity Provider domain (e.g. \"myorg.auth0.com\")",
				},
				"idpClientId": map[string]any{
					"type":        "string",
					"description": "OAuth2 client ID",
				},
				"idpAudience": map[string]any{
					"type":        "string",
					"description": "OAuth2 audience",
				},
			},
			"required": []string{"network", "controller", "idpDomain", "idpClientId", "idpAudience"},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Login with Device Auth",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, loginDeviceAuthHandler(s))
}

func loginDeviceAuthHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		controller := stringParam(req.Parameters, "controller")
		idpDomain := stringParam(req.Parameters, "idpDomain")
		idpClientID := stringParam(req.Parameters, "idpClientId")
		idpAudience := stringParam(req.Parameters, "idpAudience")

		controller, err := prepareLogin(s, network, controller, true)
		if err != nil {
			return tools.HandlerResponse{}, err
		}

		// Store controller host and client ID now
		if err := s.SetControllerHost(controller); err != nil {
			return tools.HandlerResponse{}, err
		}
		if err := s.SetIDPClientID(idpClientID); err != nil {
			return tools.HandlerResponse{}, err
		}

		// Request device code (non-blocking step 1)
		result, err := auth.RequestDeviceCode(auth.DeviceAuthConfig{
			IDPDomain: idpDomain,
			ClientID:  idpClientID,
			Audience:  idpAudience,
		})
		if err != nil {
			return tools.HandlerResponse{}, err
		}

		// Store pending state
		interval := time.Duration(result.Interval) * time.Second
		if interval == 0 {
			interval = 5 * time.Second
		}
		pendingAuthMu.Lock()
		pendingAuth = &pendingDeviceAuth{
			ProfileName:   network,
			DeviceCode:    result.DeviceCode,
			TokenEndpoint: result.TokenEndpoint,
			ClientID:      idpClientID,
			Interval:      interval,
			ExpiresAt:     time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		}
		pendingAuthMu.Unlock()

		verifyURL := result.VerificationURIComplete
		if verifyURL == "" {
			verifyURL = result.VerificationURI
		}

		return jsonResponse(map[string]any{
			"status":          "pending_verification",
			"network":         network,
			"authMode":        "device-auth",
			"verificationUri": verifyURL,
			"userCode":        result.UserCode,
			"expiresIn":       result.ExpiresIn,
		})
	}
}

// --- completeLogin ---

func registerMetaCompleteLogin(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "completeLogin",
		Description: "Complete a pending device-auth login after the user has approved in their browser.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name (defaults to active profile)",
				},
			},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Complete Device Auth Login",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, completeLoginHandler(s))
}

func completeLoginHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		if network == "" {
			network = s.ActiveProfile()
		}

		pendingAuthMu.Lock()
		pending := pendingAuth
		pendingAuthMu.Unlock()

		if pending == nil {
			return tools.HandlerResponse{}, fmt.Errorf("no pending device-auth login")
		}
		if pending.ProfileName != network {
			return tools.HandlerResponse{}, fmt.Errorf("pending device-auth is for profile %q, not %q", pending.ProfileName, network)
		}
		if time.Now().After(pending.ExpiresAt) {
			pendingAuthMu.Lock()
			pendingAuth = nil
			pendingAuthMu.Unlock()
			return tools.HandlerResponse{}, fmt.Errorf("device code has expired, please call loginDeviceAuth again")
		}

		// Poll once
		tr, err := auth.PollDeviceTokenOnce(pending.TokenEndpoint, pending.DeviceCode, pending.ClientID)
		if err != nil {
			return tools.HandlerResponse{}, err
		}

		if tr.Error == "" {
			// Success — store tokens
			if s.ActiveProfile() != network {
				if err := s.SetActiveProfile(network); err != nil {
					return tools.HandlerResponse{}, err
				}
			}
			if err := auth.StoreTokens(s, tr); err != nil {
				return tools.HandlerResponse{}, fmt.Errorf("storing tokens: %w", err)
			}
			pendingAuthMu.Lock()
			pendingAuth = nil
			pendingAuthMu.Unlock()

			return jsonResponse(map[string]any{
				"status":  "connected",
				"network": network,
			})
		}

		switch tr.Error {
		case "authorization_pending":
			return jsonResponse(map[string]any{
				"status":  "pending",
				"network": network,
				"message": "User has not yet approved. Call completeLogin again.",
			})
		case "slow_down":
			pendingAuthMu.Lock()
			if pendingAuth != nil {
				pendingAuth.Interval *= 2
			}
			pendingAuthMu.Unlock()
			return jsonResponse(map[string]any{
				"status":  "pending",
				"network": network,
				"message": "Slow down. Call completeLogin again after a longer delay.",
			})
		case "access_denied":
			pendingAuthMu.Lock()
			pendingAuth = nil
			pendingAuthMu.Unlock()
			desc := tr.ErrorDescription
			if desc == "" {
				desc = "User denied authorization"
			}
			return tools.HandlerResponse{}, fmt.Errorf("access denied: %s", desc)
		case "expired_token":
			pendingAuthMu.Lock()
			pendingAuth = nil
			pendingAuthMu.Unlock()
			return tools.HandlerResponse{}, fmt.Errorf("device code expired, please call loginDeviceAuth again")
		default:
			pendingAuthMu.Lock()
			pendingAuth = nil
			pendingAuthMu.Unlock()
			return tools.HandlerResponse{}, fmt.Errorf("unexpected error: %s - %s", tr.Error, tr.ErrorDescription)
		}
	}
}

// --- logout ---

func registerMetaLogout(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "logout",
		Description: "Disconnect from a Ziti network by clearing the profile's credentials.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name to log out (defaults to active profile)",
				},
			},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Logout from Ziti Network",
			ReadOnlyHint:    false,
			DestructiveHint: true,
			IdempotentHint:  true,
		},
	}, logoutHandler(s))
}

func logoutHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		if network == "" {
			network = s.ActiveProfile()
		}
		if network == "" {
			return tools.HandlerResponse{}, fmt.Errorf("no active profile to log out")
		}
		if !s.HasProfile(network) {
			return tools.HandlerResponse{}, fmt.Errorf("profile %q does not exist", network)
		}

		// Best-effort revoke refresh token
		clientID := s.GetForProfile(network, store.KeyIDPClientID)
		if clientID != "" {
			if err := auth.RevokeRefreshToken(s, clientID); err != nil {
				slog.Debug("refresh token revocation failed", "error", err)
			}
		}

		// Clear the profile's data
		s.ClearProfile(network)

		// Invalidate session cache
		client.InvalidateSession(network)

		// If this was the active profile, switch to another
		if s.ActiveProfile() == network {
			names := s.ProfileNames()
			switched := ""
			for _, n := range names {
				if n != network {
					switched = n
					break
				}
			}
			if switched != "" {
				if err := s.SetActiveProfile(switched); err != nil {
					slog.Debug("failed to switch active profile after logout", "error", err)
				}
			}
		}

		return jsonResponse(map[string]any{
			"status":  "logged_out",
			"network": network,
		})
	}
}

// --- listNetworks ---

func registerMetaListNetworks(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "listNetworks",
		Description: "List all configured Ziti network profiles with their connection status.",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Meta: &tools.ToolMeta{ReadOnly: true},
		Annotations: &tools.ToolAnnotations{
			Title:           "List Ziti Networks",
			ReadOnlyHint:    true,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, listNetworksHandler(s))
}

func listNetworksHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		active := s.ActiveProfile()
		names := s.ProfileNames()

		networks := make([]map[string]any, 0, len(names))
		for _, name := range names {
			data := s.ProfileData(name)
			controllerHost := data[store.KeyControllerHost]
			authMode := detectAuthMode(data)
			status := detectProfileStatus(data, s, name)

			entry := map[string]any{
				"name":           name,
				"active":         name == active,
				"controllerHost": controllerHost,
				"authMode":       authMode,
				"status":         status,
			}
			networks = append(networks, entry)
		}

		return jsonResponse(map[string]any{
			"networks": networks,
		})
	}
}

func detectAuthMode(data map[string]string) string {
	if data[store.KeyIdentityCert] != "" {
		return "identity"
	}
	if data[store.KeyUpdbUsername] != "" {
		return "updb"
	}
	if data[store.KeyToken] != "" {
		return "token"
	}
	return "none"
}

func detectProfileStatus(data map[string]string, s *store.Store, profile string) string {
	if data[store.KeyIdentityCert] != "" && data[store.KeyControllerHost] != "" {
		return "connected"
	}
	if data[store.KeyUpdbUsername] != "" && data[store.KeyUpdbPassword] != "" && data[store.KeyControllerHost] != "" {
		return "connected"
	}
	if data[store.KeyToken] != "" {
		if config.IsTokenExpiredValue(data[store.KeyTokenExpiresAt], 300) {
			return "expired"
		}
		if data[store.KeyControllerHost] != "" || data[store.KeyDomain] != "" {
			return "connected"
		}
	}
	if data[store.KeyControllerHost] != "" {
		return "incomplete"
	}
	return "incomplete"
}

// --- selectNetwork ---

func registerMetaSelectNetwork(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name:        "selectNetwork",
		Description: "Switch the active Ziti network profile. Subsequent tool calls will target this network.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"network": map[string]any{
					"type":        "string",
					"description": "Profile name to activate",
				},
			},
			"required": []string{"network"},
		},
		Meta: &tools.ToolMeta{ReadOnly: false},
		Annotations: &tools.ToolAnnotations{
			Title:           "Select Ziti Network",
			ReadOnlyHint:    false,
			DestructiveHint: false,
			IdempotentHint:  true,
		},
	}, selectNetworkHandler(s))
}

func selectNetworkHandler(s *store.Store) tools.HandlerFunc {
	return func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		network := stringParam(req.Parameters, "network")
		if network == "" {
			return tools.HandlerResponse{}, fmt.Errorf("network is required")
		}

		if err := s.SetActiveProfile(network); err != nil {
			return tools.HandlerResponse{}, err
		}

		data := s.ProfileData(network)
		controllerHost := data[store.KeyControllerHost]
		authMode := detectAuthMode(data)

		return jsonResponse(map[string]any{
			"status":         "switched",
			"network":        network,
			"controllerHost": controllerHost,
			"authMode":       authMode,
		})
	}
}

// --- helpers ---

// normalizeControllerHost strips protocol prefix and trailing slash from a controller host.
// Unlike client.FormatDomain, it does NOT append an auth0 suffix for bare hostnames.
func normalizeControllerHost(host string) string {
	if host == "" {
		return ""
	}
	for _, prefix := range []string{"https://", "http://"} {
		if len(host) > len(prefix) && host[:len(prefix)] == prefix {
			host = host[len(prefix):]
			break
		}
	}
	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	return host
}

func stringParam(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func jsonResponse(data map[string]any) (tools.HandlerResponse, error) {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return tools.HandlerResponse{}, fmt.Errorf("encoding response: %w", err)
	}
	return tools.HandlerResponse{
		Content: []tools.ContentItem{{Type: "text", Text: string(raw)}},
	}, nil
}
