package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	currentapisession "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/current_api_session"
	currentidentity "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/current_identity"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerCurrentIdentity(r *tools.Registry, s *store.Store) {
	// getCurrentIdentity
	r.Register(tools.ToolDef{
		Name:        "getCurrentIdentity",
		Description: "Get details about the currently authenticated identity",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Current Identity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get current identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.GetCurrentIdentity(currentidentity.NewGetCurrentIdentityParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getCurrentApiSession
	r.Register(tools.ToolDef{
		Name:        "getCurrentApiSession",
		Description: "Get details about the current API session",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Current API Session"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get current API session", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.GetCurrentAPISession(currentapisession.NewGetCurrentAPISessionParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteCurrentApiSession
	r.Register(tools.ToolDef{
		Name:        "deleteCurrentApiSession",
		Description: "Delete/logout the current API session",
		InputSchema: emptySchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Current API Session"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "delete current API session", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.DeleteCurrentAPISession(currentapisession.NewDeleteCurrentAPISessionParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listCurrentIdentityAuthenticators
	r.Register(tools.ToolDef{
		Name:        "listCurrentIdentityAuthenticators",
		Description: "List all authenticators for the current identity",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Current Identity Authenticators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list current identity authenticators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.ListCurrentIdentityAuthenticators(currentapisession.NewListCurrentIdentityAuthenticatorsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getCurrentIdentityAuthenticator
	r.Register(tools.ToolDef{
		Name:        "getCurrentIdentityAuthenticator",
		Description: "Get details about a specific authenticator for the current identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Current Identity Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get current identity authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.DetailCurrentIdentityAuthenticator(
					currentapisession.NewDetailCurrentIdentityAuthenticatorParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateCurrentIdentityAuthenticator
	r.Register(tools.ToolDef{
		Name:        "updateCurrentIdentityAuthenticator",
		Description: "Update an authenticator for the current identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":              map[string]any{"type": "string", "description": "Authenticator ID"},
				"currentPassword": map[string]any{"type": "string", "description": "Current password for verification", "writeOnly": true, "format": "password"},
				"password":        map[string]any{"type": "string", "description": "New password", "writeOnly": true, "format": "password"},
				"username":        map[string]any{"type": "string", "description": "New username"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Current Identity Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.AuthenticatorPatchWithCurrent{}
		if v := OptionalString(req.Parameters, "currentPassword"); v != "" {
			pw := models.Password(v)
			body.CurrentPassword = &pw
		}
		if v := OptionalString(req.Parameters, "password"); v != "" {
			pw := models.PasswordNullable(v)
			body.Password = &pw
		}
		if v := OptionalString(req.Parameters, "username"); v != "" {
			un := models.UsernameNullable(v)
			body.Username = &un
		}

		return client.WithAuthenticatedClient(req, cfg, "update current identity authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.PatchCurrentIdentityAuthenticator(
					currentapisession.NewPatchCurrentIdentityAuthenticatorParams().WithID(id).WithAuthenticator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// extendCurrentIdentityAuthenticator
	r.Register(tools.ToolDef{
		Name:        "extendCurrentIdentityAuthenticator",
		Description: "Extend/renew a certificate authenticator for the current identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":            map[string]any{"type": "string", "description": "Authenticator ID"},
				"clientCertCsr": map[string]any{"type": "string", "description": "PEM-encoded certificate signing request"},
			},
			"required": []string{"id", "clientCertCsr"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Extend Current Identity Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		csr, errResp, ok := RequireString(req.Parameters, "clientCertCsr")
		if !ok {
			return *errResp, nil
		}
		body := &models.IdentityExtendEnrollmentRequest{
			ClientCertCsr: strPtr(csr),
		}

		return client.WithAuthenticatedClient(req, cfg, "extend current identity authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.ExtendCurrentIdentityAuthenticator(
					currentapisession.NewExtendCurrentIdentityAuthenticatorParams().WithID(id).WithExtend(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// extendVerifyCurrentIdentityAuthenticator
	r.Register(tools.ToolDef{
		Name:        "extendVerifyCurrentIdentityAuthenticator",
		Description: "Verify and complete the extension of a certificate authenticator",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Authenticator ID"},
				"clientCert": map[string]any{"type": "string", "description": "PEM-encoded client certificate"},
			},
			"required": []string{"id", "clientCert"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Extend Verify Current Identity Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		cert, errResp, ok := RequireString(req.Parameters, "clientCert")
		if !ok {
			return *errResp, nil
		}
		body := &models.IdentityExtendValidateEnrollmentRequest{
			ClientCert: strPtr(cert),
		}

		return client.WithAuthenticatedClient(req, cfg, "extend verify current identity authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentAPISession.ExtendVerifyCurrentIdentityAuthenticator(
					currentapisession.NewExtendVerifyCurrentIdentityAuthenticatorParams().WithID(id).WithExtend(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
