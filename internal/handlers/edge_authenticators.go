package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	authn "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/authenticator"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerAuthenticators(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listAuthenticators", Description: "List all Authenticators in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Authenticators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list authenticators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.ListAuthenticators(authn.NewListAuthenticatorsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listAuthenticator", Description: "Get details about a specific Ziti Authenticator",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Authenticator Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.DetailAuthenticator(authn.NewDetailAuthenticatorParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createAuthenticator", Description: "Create a new Ziti Authenticator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"method":     map[string]any{"type": "string", "description": "Authentication method (updb or cert)"},
				"identityId": map[string]any{"type": "string", "description": "Identity ID"},
				"username":   map[string]any{"type": "string", "description": "Username (for updb)"},
				"password":   map[string]any{"type": "string", "description": "Password (for updb)", "writeOnly": true, "format": "password"},
				"certPem":    map[string]any{"type": "string", "description": "PEM certificate (for cert)"},
			},
			"required": []string{"method", "identityId"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		method, errResp, ok := RequireString(req.Parameters, "method")
		if !ok {
			return *errResp, nil
		}
		identityID, errResp, ok := RequireString(req.Parameters, "identityId")
		if !ok {
			return *errResp, nil
		}
		body := &models.AuthenticatorCreate{
			Method:     strPtr(method),
			IdentityID: strPtr(identityID),
			Username:   OptionalString(req.Parameters, "username"),
			Password:   OptionalString(req.Parameters, "password"),
			CertPem:    OptionalString(req.Parameters, "certPem"),
		}
		return client.WithAuthenticatedClient(req, cfg, "create authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.CreateAuthenticator(authn.NewCreateAuthenticatorParams().WithAuthenticator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteAuthenticator", Description: "Delete a Ziti Authenticator.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.DeleteAuthenticator(authn.NewDeleteAuthenticatorParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateAuthenticator", Description: "Update an existing Ziti Authenticator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "Authenticator ID"},
				"username": map[string]any{"type": "string", "description": "New username"},
				"password": map[string]any{"type": "string", "description": "New password", "writeOnly": true, "format": "password"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.AuthenticatorPatch{}
		if v := OptionalString(req.Parameters, "username"); v != "" {
			u := models.UsernameNullable(v)
			body.Username = &u
		}
		if v := OptionalString(req.Parameters, "password"); v != "" {
			p := models.PasswordNullable(v)
			body.Password = &p
		}
		return client.WithAuthenticatedClient(req, cfg, "update authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.PatchAuthenticator(authn.NewPatchAuthenticatorParams().WithID(id).WithAuthenticator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
