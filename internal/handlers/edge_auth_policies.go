package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	authpolicy "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/auth_policy"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerAuthPolicies(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listAuthPolicies", Description: "List all Auth Policies in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Auth Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list auth policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.AuthPolicy.ListAuthPolicies(
						authpolicy.NewListAuthPoliciesParams().WithLimit(&limit).WithOffset(&offset), noAuth)
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
		Name: "listAuthPolicy", Description: "Get details about a specific Ziti Auth Policy",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Auth Policy Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get auth policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.AuthPolicy.DetailAuthPolicy(authpolicy.NewDetailAuthPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createAuthPolicy", Description: "Create a new Ziti Auth Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":                           map[string]any{"type": "string", "description": "Auth policy name"},
				"primaryCertAllowed":             map[string]any{"type": "boolean", "default": false},
				"primaryCertAllowExpiredCerts":   map[string]any{"type": "boolean", "default": false},
				"primaryExtJwtAllowed":           map[string]any{"type": "boolean", "default": false},
				"primaryExtJwtAllowedSigners":    map[string]any{"type": "string", "description": "Comma-separated signer IDs"},
				"primaryUpdbAllowed":             map[string]any{"type": "boolean", "default": false},
				"primaryUpdbMinPasswordLength":   map[string]any{"type": "number", "default": 5},
				"primaryUpdbRequireMixedCase":    map[string]any{"type": "boolean", "default": false},
				"primaryUpdbRequireNumberChar":   map[string]any{"type": "boolean", "default": false},
				"primaryUpdbRequireSpecialChar":  map[string]any{"type": "boolean", "default": false},
				"primaryUpdbMaxAttempts":         map[string]any{"type": "number", "default": 0},
				"primaryUpdbLockoutDurationMinutes": map[string]any{"type": "number", "default": 0},
				"secondaryRequireTotp":           map[string]any{"type": "boolean", "default": false},
				"secondaryRequireExtJwtSigner":   map[string]any{"type": "string", "description": "External JWT signer ID"},
			},
			"required": []string{"name"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Auth Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		minPwLen := int64(5)
		if v := OptionalInt64(req.Parameters, "primaryUpdbMinPasswordLength"); v != nil {
			minPwLen = *v
		}
		maxAttempts := int64(0)
		if v := OptionalInt64(req.Parameters, "primaryUpdbMaxAttempts"); v != nil {
			maxAttempts = *v
		}
		lockoutMins := int64(0)
		if v := OptionalInt64(req.Parameters, "primaryUpdbLockoutDurationMinutes"); v != nil {
			lockoutMins = *v
		}
		reqTotp := OptionalBool(req.Parameters, "secondaryRequireTotp", false)
		body := &models.AuthPolicyCreate{
			Name: strPtr(name),
			Primary: &models.AuthPolicyPrimary{
				Cert: &models.AuthPolicyPrimaryCert{
					Allowed:           boolPtr(OptionalBool(req.Parameters, "primaryCertAllowed", false)),
					AllowExpiredCerts: boolPtr(OptionalBool(req.Parameters, "primaryCertAllowExpiredCerts", false)),
				},
				ExtJwt: &models.AuthPolicyPrimaryExtJwt{
					Allowed:        boolPtr(OptionalBool(req.Parameters, "primaryExtJwtAllowed", false)),
					AllowedSigners: SplitCSV(OptionalString(req.Parameters, "primaryExtJwtAllowedSigners")),
				},
				Updb: &models.AuthPolicyPrimaryUpdb{
					Allowed:                boolPtr(OptionalBool(req.Parameters, "primaryUpdbAllowed", false)),
					MinPasswordLength:      &minPwLen,
					RequireMixedCase:       boolPtr(OptionalBool(req.Parameters, "primaryUpdbRequireMixedCase", false)),
					RequireNumberChar:      boolPtr(OptionalBool(req.Parameters, "primaryUpdbRequireNumberChar", false)),
					RequireSpecialChar:     boolPtr(OptionalBool(req.Parameters, "primaryUpdbRequireSpecialChar", false)),
					MaxAttempts:            &maxAttempts,
					LockoutDurationMinutes: &lockoutMins,
				},
			},
			Secondary: &models.AuthPolicySecondary{
				RequireTotp: &reqTotp,
			},
		}
		if v := OptionalString(req.Parameters, "secondaryRequireExtJwtSigner"); v != "" {
			body.Secondary.RequireExtJwtSigner = strPtr(v)
		}
		return client.WithAuthenticatedClient(req, cfg, "create auth policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.AuthPolicy.CreateAuthPolicy(authpolicy.NewCreateAuthPolicyParams().WithAuthPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteAuthPolicy", Description: "Delete a Ziti Auth Policy.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Auth Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete auth policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.AuthPolicy.DeleteAuthPolicy(authpolicy.NewDeleteAuthPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateAuthPolicy", Description: "Update an existing Ziti Auth Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                             map[string]any{"type": "string", "description": "Auth policy ID"},
				"name":                           map[string]any{"type": "string"},
				"primaryCertAllowed":             map[string]any{"type": "boolean"},
				"primaryCertAllowExpiredCerts":   map[string]any{"type": "boolean"},
				"primaryExtJwtAllowed":           map[string]any{"type": "boolean"},
				"primaryExtJwtAllowedSigners":    map[string]any{"type": "string", "description": "Comma-separated signer IDs"},
				"primaryUpdbAllowed":             map[string]any{"type": "boolean"},
				"primaryUpdbMinPasswordLength":   map[string]any{"type": "number"},
				"primaryUpdbRequireMixedCase":    map[string]any{"type": "boolean"},
				"primaryUpdbRequireNumberChar":   map[string]any{"type": "boolean"},
				"primaryUpdbRequireSpecialChar":  map[string]any{"type": "boolean"},
				"primaryUpdbMaxAttempts":         map[string]any{"type": "number"},
				"primaryUpdbLockoutDurationMinutes": map[string]any{"type": "number"},
				"secondaryRequireTotp":           map[string]any{"type": "boolean"},
				"secondaryRequireExtJwtSigner":   map[string]any{"type": "string"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Auth Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.AuthPolicyPatch{
			Primary: &models.AuthPolicyPrimaryPatch{
				Cert:   &models.AuthPolicyPrimaryCertPatch{},
				ExtJwt: &models.AuthPolicyPrimaryExtJwtPatch{},
				Updb:   &models.AuthPolicyPrimaryUpdbPatch{},
			},
			Secondary: &models.AuthPolicySecondaryPatch{},
		}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = strPtr(v)
		}
		if v, exists := req.Parameters["primaryCertAllowed"]; exists && v != nil {
			body.Primary.Cert.Allowed = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["primaryCertAllowExpiredCerts"]; exists && v != nil {
			body.Primary.Cert.AllowExpiredCerts = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["primaryExtJwtAllowed"]; exists && v != nil {
			body.Primary.ExtJwt.Allowed = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "primaryExtJwtAllowedSigners"); v != "" {
			body.Primary.ExtJwt.AllowedSigners = SplitCSV(v)
		}
		if v, exists := req.Parameters["primaryUpdbAllowed"]; exists && v != nil {
			body.Primary.Updb.Allowed = boolPtr(v.(bool))
		}
		if v := OptionalInt64(req.Parameters, "primaryUpdbMinPasswordLength"); v != nil {
			body.Primary.Updb.MinPasswordLength = v
		}
		if v, exists := req.Parameters["primaryUpdbRequireMixedCase"]; exists && v != nil {
			body.Primary.Updb.RequireMixedCase = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["primaryUpdbRequireNumberChar"]; exists && v != nil {
			body.Primary.Updb.RequireNumberChar = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["primaryUpdbRequireSpecialChar"]; exists && v != nil {
			body.Primary.Updb.RequireSpecialChar = boolPtr(v.(bool))
		}
		if v := OptionalInt64(req.Parameters, "primaryUpdbMaxAttempts"); v != nil {
			body.Primary.Updb.MaxAttempts = v
		}
		if v := OptionalInt64(req.Parameters, "primaryUpdbLockoutDurationMinutes"); v != nil {
			body.Primary.Updb.LockoutDurationMinutes = v
		}
		if v, exists := req.Parameters["secondaryRequireTotp"]; exists && v != nil {
			body.Secondary.RequireTotp = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "secondaryRequireExtJwtSigner"); v != "" {
			body.Secondary.RequireExtJwtSigner = strPtr(v)
		}
		return client.WithAuthenticatedClient(req, cfg, "update auth policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.AuthPolicy.PatchAuthPolicy(authpolicy.NewPatchAuthPolicyParams().WithID(id).WithAuthPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
