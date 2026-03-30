package handlers

import (
	"net/http"

	"github.com/go-openapi/strfmt"

	"github.com/openziti/ziti-mcp-server/internal/client"
	ejwt "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/external_j_w_t_signer"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerExternalJwtSigners(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listExternalJwtSigners", Description: "List all External JWT Signers in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List External JWT Signers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list external JWT signers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.ExternaljwtSigner.ListExternalJwtSigners(
						ejwt.NewListExternalJwtSignersParams().WithLimit(&limit).WithOffset(&offset), noAuth)
					if err != nil {
						return nil, err
					}
					m, err := ToMap(resp.Payload)
					if err != nil {
						return nil, err
					}
					return m.(map[string]any), nil
				})
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listExternalJwtSigner", Description: "Get details about a specific Ziti External JWT Signer",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get External JWT Signer Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get external JWT signer", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ExternaljwtSigner.DetailExternalJwtSigner(ejwt.NewDetailExternalJwtSignerParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createExternalJwtSigner", Description: "Create a new Ziti External JWT Signer.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":            map[string]any{"type": "string", "description": "Signer name"},
				"issuer":          map[string]any{"type": "string", "description": "Token issuer"},
				"audience":        map[string]any{"type": "string", "description": "Token audience"},
				"enabled":         map[string]any{"type": "boolean", "description": "Whether the signer is enabled"},
				"certPem":         map[string]any{"type": "string", "description": "PEM certificate for verification"},
				"jwksEndpoint":    map[string]any{"type": "string", "description": "JWKS endpoint URL"},
				"kid":             map[string]any{"type": "string", "description": "Key ID"},
				"clientId":        map[string]any{"type": "string", "description": "OAuth client ID"},
				"externalAuthUrl": map[string]any{"type": "string", "description": "External auth URL"},
				"claimsProperty":  map[string]any{"type": "string", "description": "Claims property"},
				"useExternalId":   map[string]any{"type": "boolean", "description": "Use external ID mapping"},
				"scopes":          map[string]any{"type": "string", "description": "Comma-separated scopes"},
			},
			"required": []string{"name", "issuer", "audience", "enabled"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create External JWT Signer"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		issuer, errResp, ok := RequireString(req.Parameters, "issuer")
		if !ok {
			return *errResp, nil
		}
		audience, errResp, ok := RequireString(req.Parameters, "audience")
		if !ok {
			return *errResp, nil
		}
		enabled := OptionalBool(req.Parameters, "enabled", true)

		body := &models.ExternalJwtSignerCreate{
			Name:     strPtr(name),
			Issuer:   strPtr(issuer),
			Audience: strPtr(audience),
			Enabled:  boolPtr(enabled),
		}
		if v := OptionalString(req.Parameters, "certPem"); v != "" {
			body.CertPem = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "jwksEndpoint"); v != "" {
			u := strfmt.URI(v)
			body.JwksEndpoint = &u
		}
		if v := OptionalString(req.Parameters, "kid"); v != "" {
			body.Kid = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "clientId"); v != "" {
			body.ClientID = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "externalAuthUrl"); v != "" {
			body.ExternalAuthURL = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "claimsProperty"); v != "" {
			body.ClaimsProperty = strPtr(v)
		}
		if v, exists := req.Parameters["useExternalId"]; exists && v != nil {
			body.UseExternalID = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "scopes"); v != "" {
			body.Scopes = SplitCSV(v)
		}

		return client.WithAuthenticatedClient(req, cfg, "create external JWT signer", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ExternaljwtSigner.CreateExternalJwtSigner(ejwt.NewCreateExternalJwtSignerParams().WithExternalJwtSigner(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteExternalJwtSigner", Description: "Delete a Ziti External JWT Signer.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete External JWT Signer"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete external JWT signer", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ExternaljwtSigner.DeleteExternalJwtSigner(ejwt.NewDeleteExternalJwtSignerParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateExternalJwtSigner", Description: "Update an existing Ziti External JWT Signer.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":              map[string]any{"type": "string", "description": "Signer ID"},
				"name":            map[string]any{"type": "string"},
				"issuer":          map[string]any{"type": "string"},
				"audience":        map[string]any{"type": "string"},
				"enabled":         map[string]any{"type": "boolean"},
				"certPem":         map[string]any{"type": "string"},
				"jwksEndpoint":    map[string]any{"type": "string"},
				"kid":             map[string]any{"type": "string"},
				"clientId":        map[string]any{"type": "string"},
				"externalAuthUrl": map[string]any{"type": "string"},
				"claimsProperty":  map[string]any{"type": "string"},
				"useExternalId":   map[string]any{"type": "boolean"},
				"scopes":          map[string]any{"type": "string", "description": "Comma-separated scopes"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update External JWT Signer"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ExternalJwtSignerPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "issuer"); v != "" {
			body.Issuer = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "audience"); v != "" {
			body.Audience = strPtr(v)
		}
		if v, exists := req.Parameters["enabled"]; exists && v != nil {
			body.Enabled = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "certPem"); v != "" {
			body.CertPem = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "jwksEndpoint"); v != "" {
			u := strfmt.URI(v)
			body.JwksEndpoint = &u
		}
		if v := OptionalString(req.Parameters, "kid"); v != "" {
			body.Kid = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "clientId"); v != "" {
			body.ClientID = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "externalAuthUrl"); v != "" {
			body.ExternalAuthURL = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "claimsProperty"); v != "" {
			body.ClaimsProperty = strPtr(v)
		}
		if v, exists := req.Parameters["useExternalId"]; exists && v != nil {
			body.UseExternalID = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "scopes"); v != "" {
			body.Scopes = SplitCSV(v)
		}
		return client.WithAuthenticatedClient(req, cfg, "update external JWT signer", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ExternaljwtSigner.PatchExternalJwtSigner(ejwt.NewPatchExternalJwtSignerParams().WithID(id).WithExternalJwtSigner(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
