package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	edge_router_policy "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/edge_router_policy"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerEdgeRouterPolicies(r *tools.Registry, s *store.Store) {
	// listEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterPolicies",
		Description: "List all Edge Router Policies in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.ListEdgeRouterPolicies(edge_router_policy.NewListEdgeRouterPoliciesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterPolicy",
		Description: "Get details about a specific Ziti Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Edge Router Policy Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.DetailEdgeRouterPolicy(edge_router_policy.NewDetailEdgeRouterPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "createEdgeRouterPolicy",
		Description: "Create a new Ziti Edge Router Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":            map[string]any{"type": "string", "description": "Name of the edge router policy"},
				"semantic":        map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"edgeRouterRoles": map[string]any{"type": "string", "description": "Comma-separated edge router roles"},
				"identityRoles":   map[string]any{"type": "string", "description": "Comma-separated identity roles"},
			},
			"required": []string{"name", "semantic"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Edge Router Policy"),
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
		body := &models.EdgeRouterPolicyCreate{
			Name:     strPtr(name),
			Semantic: &semantic,
		}
		if v := OptionalString(req.Parameters, "edgeRouterRoles"); v != "" {
			body.EdgeRouterRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "identityRoles"); v != "" {
			body.IdentityRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "create edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.CreateEdgeRouterPolicy(edge_router_policy.NewCreateEdgeRouterPolicyParams().WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "deleteEdgeRouterPolicy",
		Description: "Delete a Ziti Edge Router Policy.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Edge Router Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.DeleteEdgeRouterPolicy(edge_router_policy.NewDeleteEdgeRouterPolicyParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateEdgeRouterPolicy
	r.Register(tools.ToolDef{
		Name:        "updateEdgeRouterPolicy",
		Description: "Update an existing Ziti Edge Router Policy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":              map[string]any{"type": "string", "description": "Edge Router Policy ID"},
				"name":            map[string]any{"type": "string", "description": "New name"},
				"semantic":        map[string]any{"type": "string", "description": "Semantic", "enum": []string{"AllOf", "AnyOf"}},
				"edgeRouterRoles": map[string]any{"type": "string", "description": "Comma-separated edge router roles"},
				"identityRoles":   map[string]any{"type": "string", "description": "Comma-separated identity roles"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Edge Router Policy"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.EdgeRouterPolicyPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = v
		}
		if v := OptionalString(req.Parameters, "semantic"); v != "" {
			body.Semantic = models.Semantic(v)
		}
		if v := OptionalString(req.Parameters, "edgeRouterRoles"); v != "" {
			body.EdgeRouterRoles = models.Roles(SplitCSV(v))
		}
		if v := OptionalString(req.Parameters, "identityRoles"); v != "" {
			body.IdentityRoles = models.Roles(SplitCSV(v))
		}

		return client.WithAuthenticatedClient(req, cfg, "update edge router policy", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.PatchEdgeRouterPolicy(edge_router_policy.NewPatchEdgeRouterPolicyParams().WithID(id).WithPolicy(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterPolicyEdgeRouters
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterPolicyEdgeRouters",
		Description: "List all Edge Routers associated with a specific Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Policy Edge Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router policy edge routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.ListEdgeRouterPolicyEdgeRouters(edge_router_policy.NewListEdgeRouterPolicyEdgeRoutersParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterPolicyIdentities
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterPolicyIdentities",
		Description: "List all Identities associated with a specific Edge Router Policy",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Policy Identities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router policy identities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouterPolicy.ListEdgeRouterPolicyIdentities(edge_router_policy.NewListEdgeRouterPolicyIdentitiesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
