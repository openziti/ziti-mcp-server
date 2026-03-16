package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	service_policy "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/service_policy"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerServicePolicies(r *tools.Registry, s *store.Store) {
	// listServicePolicies
	r.Register(tools.ToolDef{
		Name:        "listServicePolicies",
		Description: "List all Service Policies in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list service policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.ListServicePolicies(service_policy.NewListServicePoliciesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServicePolicy
	r.Register(tools.ToolDef{
		Name:        "listServicePolicy",
		Description: "Get details about a specific Ziti Service Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Service Policy Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get service policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.DetailServicePolicy(service_policy.NewDetailServicePolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createServicePolicy
	r.Register(tools.ToolDef{
		Name:        "createServicePolicy",
		Description: "Create a new Ziti Service Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":              map[string]any{"type": "string", "description": "Name of the service policy"},
				"type":              map[string]any{"type": "string", "description": "Policy type", "enum": []string{"Dial", "Bind"}},
				"semantic":          map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"identityRoles":     map[string]any{"type": "string", "description": "Comma-separated identity roles"},
				"serviceRoles":      map[string]any{"type": "string", "description": "Comma-separated service roles"},
				"postureCheckRoles": map[string]any{"type": "string", "description": "Comma-separated posture check roles"},
			},
			"required": []string{"name", "type", "semantic"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Service Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		typ, errResp, ok := RequireString(req.Parameters, "type")
		if !ok {
			return *errResp, nil
		}
		sem, errResp, ok := RequireString(req.Parameters, "semantic")
		if !ok {
			return *errResp, nil
		}

		dialBind := models.DialBind(typ)
		semantic := models.Semantic(sem)
		body := &models.ServicePolicyCreate{
			Name:     strPtr(name),
			Type:     &dialBind,
			Semantic: &semantic,
		}
		if v := OptionalString(req.Parameters, "identityRoles"); v != "" {
			body.IdentityRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "serviceRoles"); v != "" {
			body.ServiceRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "postureCheckRoles"); v != "" {
			body.PostureCheckRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "create service policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.CreateServicePolicy(service_policy.NewCreateServicePolicyParams().WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteServicePolicy
	r.Register(tools.ToolDef{
		Name:        "deleteServicePolicy",
		Description: "Delete a Ziti Service Policy.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Service Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete service policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.DeleteServicePolicy(service_policy.NewDeleteServicePolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateServicePolicy
	r.Register(tools.ToolDef{
		Name:        "updateServicePolicy",
		Description: "Update an existing Ziti Service Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                map[string]any{"type": "string", "description": "Service Policy ID"},
				"name":              map[string]any{"type": "string", "description": "New name"},
				"type":              map[string]any{"type": "string", "description": "Policy type", "enum": []string{"Dial", "Bind"}},
				"semantic":          map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"identityRoles":     map[string]any{"type": "string", "description": "Comma-separated identity roles"},
				"serviceRoles":      map[string]any{"type": "string", "description": "Comma-separated service roles"},
				"postureCheckRoles": map[string]any{"type": "string", "description": "Comma-separated posture check roles"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Service Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ServicePolicyPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = v
		}
		if v := OptionalString(req.Parameters, "type"); v != "" {
			body.Type = models.DialBind(v)
		}
		if v := OptionalString(req.Parameters, "semantic"); v != "" {
			body.Semantic = models.Semantic(v)
		}
		if v := OptionalString(req.Parameters, "identityRoles"); v != "" {
			body.IdentityRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "serviceRoles"); v != "" {
			body.ServiceRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "postureCheckRoles"); v != "" {
			body.PostureCheckRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "update service policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.PatchServicePolicy(service_policy.NewPatchServicePolicyParams().WithID(id).WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServicePolicyIdentities
	r.Register(tools.ToolDef{
		Name:        "listServicePolicyIdentities",
		Description: "List all Identities associated with a specific Service Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Policy Identities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service policy identities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.ListServicePolicyIdentities(service_policy.NewListServicePolicyIdentitiesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServicePolicyServices
	r.Register(tools.ToolDef{
		Name:        "listServicePolicyServices",
		Description: "List all Services associated with a specific Service Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Policy Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service policy services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.ListServicePolicyServices(service_policy.NewListServicePolicyServicesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServicePolicyPostureChecks
	r.Register(tools.ToolDef{
		Name:        "listServicePolicyPostureChecks",
		Description: "List all Posture Checks associated with a specific Service Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Policy Posture Checks"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service policy posture checks", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServicePolicy.ListServicePolicyPostureChecks(service_policy.NewListServicePolicyPostureChecksParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
