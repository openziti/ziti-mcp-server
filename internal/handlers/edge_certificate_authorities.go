package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	ca "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/certificate_authority"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerCertificateAuthorities(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listCas", Description: "List all Certificate Authorities in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List CAs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list CAs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.ListCas(ca.NewListCasParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listCa", Description: "Get details about a specific Ziti Certificate Authority",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get CA Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get CA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.DetailCa(ca.NewDetailCaParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createCa", Description: "Create a new Ziti Certificate Authority.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":                        map[string]any{"type": "string", "description": "CA name"},
				"certPem":                     map[string]any{"type": "string", "description": "PEM-encoded certificate"},
				"isAuthEnabled":               map[string]any{"type": "boolean", "description": "Enable authentication"},
				"isAutoCaEnrollmentEnabled":   map[string]any{"type": "boolean", "description": "Enable auto CA enrollment"},
				"isOttCaEnrollmentEnabled":    map[string]any{"type": "boolean", "description": "Enable OTT CA enrollment"},
				"identityRoles":               map[string]any{"type": "string", "description": "Comma-separated identity roles"},
				"identityNameFormat":          map[string]any{"type": "string", "description": "Identity name format template"},
			},
			"required": []string{"name", "certPem", "isAuthEnabled", "isAutoCaEnrollmentEnabled", "isOttCaEnrollmentEnabled", "identityRoles"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create CA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		certPem, errResp, ok := RequireString(req.Parameters, "certPem")
		if !ok {
			return *errResp, nil
		}
		rolesStr, errResp, ok := RequireString(req.Parameters, "identityRoles")
		if !ok {
			return *errResp, nil
		}
		body := &models.CaCreate{
			Name:                       strPtr(name),
			CertPem:                    strPtr(certPem),
			IsAuthEnabled:              boolPtr(OptionalBool(req.Parameters, "isAuthEnabled", false)),
			IsAutoCaEnrollmentEnabled:  boolPtr(OptionalBool(req.Parameters, "isAutoCaEnrollmentEnabled", false)),
			IsOttCaEnrollmentEnabled:   boolPtr(OptionalBool(req.Parameters, "isOttCaEnrollmentEnabled", false)),
			IdentityRoles:              models.Roles(SplitCSV(rolesStr)),
		}
		if v := OptionalString(req.Parameters, "identityNameFormat"); v != "" {
			body.IdentityNameFormat = v
		}
		return client.WithAuthenticatedClient(req, cfg, "create CA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.CreateCa(ca.NewCreateCaParams().WithCa(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteCa", Description: "Delete a Ziti Certificate Authority.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete CA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete CA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.DeleteCa(ca.NewDeleteCaParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateCa", Description: "Update an existing Ziti Certificate Authority.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                          map[string]any{"type": "string", "description": "CA ID"},
				"name":                        map[string]any{"type": "string", "description": "CA name"},
				"isAuthEnabled":               map[string]any{"type": "boolean", "description": "Enable authentication"},
				"isAutoCaEnrollmentEnabled":   map[string]any{"type": "boolean", "description": "Enable auto CA enrollment"},
				"isOttCaEnrollmentEnabled":    map[string]any{"type": "boolean", "description": "Enable OTT CA enrollment"},
				"identityRoles":               map[string]any{"type": "string", "description": "Comma-separated identity roles"},
				"identityNameFormat":          map[string]any{"type": "string", "description": "Identity name format template"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update CA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.CaPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = strPtr(v)
		}
		if v, exists := req.Parameters["isAuthEnabled"]; exists && v != nil {
			body.IsAuthEnabled = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["isAutoCaEnrollmentEnabled"]; exists && v != nil {
			body.IsAutoCaEnrollmentEnabled = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["isOttCaEnrollmentEnabled"]; exists && v != nil {
			body.IsOttCaEnrollmentEnabled = boolPtr(v.(bool))
		}
		if v := OptionalString(req.Parameters, "identityRoles"); v != "" {
			body.IdentityRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "identityNameFormat"); v != "" {
			body.IdentityNameFormat = strPtr(v)
		}
		return client.WithAuthenticatedClient(req, cfg, "update CA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.PatchCa(ca.NewPatchCaParams().WithID(id).WithCa(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "getCaJwt", Description: "Get the JWT for a Certificate Authority, used for enrollment",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get CA JWT"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get CA JWT", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.GetCaJwt(ca.NewGetCaJwtParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "verifyCa", Description: "Verify a Certificate Authority by providing a PEM certificate signed by the CA with the common name matching the CA validation token",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":          map[string]any{"type": "string", "description": "CA ID"},
				"certificate": map[string]any{"type": "string", "description": "PEM-encoded verification certificate"},
			},
			"required": []string{"id", "certificate"},
		},
		Meta: writeMeta(), Annotations: &tools.ToolAnnotations{Title: "Verify CA"},
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		cert, errResp, ok := RequireString(req.Parameters, "certificate")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "verify CA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CertificateAuthority.VerifyCa(ca.NewVerifyCaParams().WithID(id).WithCertificate(cert), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
