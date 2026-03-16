package handlers

import (
	"net/http"

	"github.com/go-openapi/strfmt"

	"github.com/openziti/ziti-mcp-server/internal/client"
	authn "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/authenticator"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerAuthenticatorMgmt(r *tools.Registry, s *store.Store) {
	// reEnrollAuthenticator
	r.Register(tools.ToolDef{
		Name:        "reEnrollAuthenticator",
		Description: "Re-enroll an authenticator, generating new enrollment credentials",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":        map[string]any{"type": "string", "description": "The authenticator ID"},
				"expiresAt": map[string]any{"type": "string", "description": "Expiration date-time in ISO 8601 format"},
			},
			"required": []string{"id", "expiresAt"},
		},
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Re-Enroll Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		expiresAtStr, errResp, ok := RequireString(req.Parameters, "expiresAt")
		if !ok {
			return *errResp, nil
		}
		expiresAt, err := strfmt.ParseDateTime(expiresAtStr)
		if err != nil {
			return errResponse("Invalid expiresAt format: " + err.Error()), nil
		}

		body := &models.ReEnroll{
			ExpiresAt: &expiresAt,
		}

		return client.WithAuthenticatedClient(req, cfg, "re-enroll authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.ReEnrollAuthenticator(authn.NewReEnrollAuthenticatorParams().WithID(id).WithReEnroll(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// requestExtendAuthenticator
	r.Register(tools.ToolDef{
		Name:        "requestExtendAuthenticator",
		Description: "Request to extend/renew a certificate authenticator",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "The authenticator ID"},
				"rollKeys": map[string]any{"type": "boolean", "description": "Whether to request private key rolling", "default": false},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Request Extend Authenticator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		rollKeys := OptionalBool(req.Parameters, "rollKeys", false)

		body := &models.RequestExtendAuthenticator{
			RollKeys: rollKeys,
		}

		return client.WithAuthenticatedClient(req, cfg, "request extend authenticator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.RequestExtendAuthenticator(authn.NewRequestExtendAuthenticatorParams().WithID(id).WithRequestExtendAuthenticator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// requestExtendIdentityCertAuthenticators
	r.Register(tools.ToolDef{
		Name:        "requestExtendIdentityCertAuthenticators",
		Description: "Request to extend/renew all certificate authenticators for a specific identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "The identity ID"},
				"rollKeys": map[string]any{"type": "boolean", "description": "Whether to request private key rolling", "default": false},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Request Extend Identity Cert Authenticators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		rollKeys := OptionalBool(req.Parameters, "rollKeys", false)

		body := &models.RequestExtendAuthenticator{
			RollKeys: rollKeys,
		}

		return client.WithAuthenticatedClient(req, cfg, "request extend identity cert authenticators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Authenticator.RequestExtendAllCertAuthenticators(authn.NewRequestExtendAllCertAuthenticatorsParams().WithID(id).WithRequestExtendAuthenticator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
