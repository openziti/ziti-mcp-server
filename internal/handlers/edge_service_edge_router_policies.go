package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	service_edge_router_policy "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/service_edge_router_policy"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerServiceEdgeRouterPolicies(r *tools.Registry, s *store.Store) {
	// listServiceEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listServiceEdgeRouterPolicies",
		Description: "List all Service Edge Router Policies in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list service edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.ListServiceEdgeRouterPolicies(service_edge_router_policy.NewListServiceEdgeRouterPoliciesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "listServiceEdgeRouterPolicy",
		Description: "Get details about a specific Ziti Service Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Service Edge Router Policy Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get service edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.DetailServiceEdgeRouterPolicy(service_edge_router_policy.NewDetailServiceEdgeRouterPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createServiceEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "createServiceEdgeRouterPolicy",
		Description: "Create a new Ziti Service Edge Router Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":            map[string]any{"type": "string", "description": "Name of the service edge router policy"},
				"semantic":        map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"edgeRouterRoles": map[string]any{"type": "string", "description": "Comma-separated edge router roles"},
				"serviceRoles":    map[string]any{"type": "string", "description": "Comma-separated service roles"},
			},
			"required": []string{"name", "semantic"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Service Edge Router Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		sem, errResp, ok := RequireString(req.Parameters, "semantic")
		if !ok {
			return *errResp, nil
		}

		semantic := models.Semantic(sem)
		body := &models.ServiceEdgeRouterPolicyCreate{
			Name:     strPtr(name),
			Semantic: &semantic,
		}
		if v := OptionalString(req.Parameters, "edgeRouterRoles"); v != "" {
			body.EdgeRouterRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "serviceRoles"); v != "" {
			body.ServiceRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "create service edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.CreateServiceEdgeRouterPolicy(service_edge_router_policy.NewCreateServiceEdgeRouterPolicyParams().WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteServiceEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "deleteServiceEdgeRouterPolicy",
		Description: "Delete a Ziti Service Edge Router Policy.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Service Edge Router Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete service edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.DeleteServiceEdgeRouterPolicy(service_edge_router_policy.NewDeleteServiceEdgeRouterPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateServiceEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "updateServiceEdgeRouterPolicy",
		Description: "Update an existing Ziti Service Edge Router Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":              map[string]any{"type": "string", "description": "Service Edge Router Policy ID"},
				"name":            map[string]any{"type": "string", "description": "New name"},
				"semantic":        map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"edgeRouterRoles": map[string]any{"type": "string", "description": "Comma-separated edge router roles"},
				"serviceRoles":    map[string]any{"type": "string", "description": "Comma-separated service roles"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Service Edge Router Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ServiceEdgeRouterPolicyPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = v
		}
		if v := OptionalString(req.Parameters, "semantic"); v != "" {
			body.Semantic = models.Semantic(v)
		}
		if v := OptionalString(req.Parameters, "edgeRouterRoles"); v != "" {
			body.EdgeRouterRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "serviceRoles"); v != "" {
			body.ServiceRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "update service edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.PatchServiceEdgeRouterPolicy(service_edge_router_policy.NewPatchServiceEdgeRouterPolicyParams().WithID(id).WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceEdgeRouterPolicyEdgeRouters
	r.Register(tools.ToolDef{
		Name:        "listServiceEdgeRouterPolicyEdgeRouters",
		Description: "List all Edge Routers associated with a specific Service Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Edge Router Policy Edge Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service edge router policy edge routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.ListServiceEdgeRouterPolicyEdgeRouters(service_edge_router_policy.NewListServiceEdgeRouterPolicyEdgeRoutersParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceEdgeRouterPolicyServices
	r.Register(tools.ToolDef{
		Name:        "listServiceEdgeRouterPolicyServices",
		Description: "List all Services associated with a specific Service Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Edge Router Policy Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service edge router policy services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.ServiceEdgeRouterPolicy.ListServiceEdgeRouterPolicyServices(service_edge_router_policy.NewListServiceEdgeRouterPolicyServicesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
