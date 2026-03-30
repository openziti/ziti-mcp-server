package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/client/edge_router"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerEdgeRouters(r *tools.Registry, s *store.Store) {
	// listEdgeRouters
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouters",
		Description: "List all Edge Routers in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list edge routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.EdgeRouter.ListEdgeRouters(
						edge_router.NewListEdgeRoutersParams().WithLimit(&limit).WithOffset(&offset), noAuth)
					if err != nil {
						return nil, err
					}
					m, err := ToMap(resp.Payload)
					if err != nil {
						return nil, err
					}
					return m.(map[string]any), nil
				})
			},
		), nil
	})

	// listEdgeRouter
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouter",
		Description: "Get details about a specific Ziti Edge Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Edge Router Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get edge router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.DetailEdgeRouter(edge_router.NewDetailEdgeRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createEdgeRouter
	r.Register(tools.ToolDef{
		Name:        "createEdgeRouter",
		Description: "Create a new Ziti Edge Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":              map[string]any{"type": "string", "description": "Name of the edge router"},
				"isTunnelerEnabled": map[string]any{"type": "boolean", "description": "Whether tunneler mode is enabled", "default": false},
				"roleAttributes":    map[string]any{"type": "string", "description": "Comma-separated role attributes"},
				"cost":              map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal":       map[string]any{"type": "boolean", "description": "Whether to disable traversal", "default": false},
				"disabled":          map[string]any{"type": "boolean", "description": "Whether the edge router is disabled", "default": false},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Edge Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		tunneler := OptionalBool(req.Parameters, "isTunnelerEnabled", false)
		roleAttrs := SplitCSV(OptionalString(req.Parameters, "roleAttributes"))
		cost := OptionalInt64(req.Parameters, "cost")
		noTraversal := OptionalBool(req.Parameters, "noTraversal", false)
		disabled := OptionalBool(req.Parameters, "disabled", false)

		body := &models.EdgeRouterCreate{
			Name:              strPtr(name),
			IsTunnelerEnabled: tunneler,
			NoTraversal:       boolPtr(noTraversal),
			Disabled:          boolPtr(disabled),
		}
		if cost != nil {
			body.Cost = cost
		}
		if len(roleAttrs) > 0 {
			attrs := models.Attributes(roleAttrs)
			body.RoleAttributes = &attrs
		}

		return client.WithAuthenticatedClient(req, cfg, "create edge router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.CreateEdgeRouter(edge_router.NewCreateEdgeRouterParams().WithEdgeRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteEdgeRouter
	r.Register(tools.ToolDef{
		Name:        "deleteEdgeRouter",
		Description: "Delete a Ziti Edge Router.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Edge Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete edge router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.DeleteEdgeRouter(edge_router.NewDeleteEdgeRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateEdgeRouter
	r.Register(tools.ToolDef{
		Name:        "updateEdgeRouter",
		Description: "Update an existing Ziti Edge Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                map[string]any{"type": "string", "description": "Edge Router ID"},
				"name":             map[string]any{"type": "string", "description": "New name"},
				"isTunnelerEnabled": map[string]any{"type": "boolean", "description": "Whether tunneler mode is enabled"},
				"roleAttributes":   map[string]any{"type": "string", "description": "Comma-separated role attributes"},
				"cost":             map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal":      map[string]any{"type": "boolean", "description": "Whether to disable traversal"},
				"disabled":         map[string]any{"type": "boolean", "description": "Whether the edge router is disabled"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Edge Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.EdgeRouterPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = strPtr(v)
		}
		if v, exists := req.Parameters["isTunnelerEnabled"]; exists && v != nil {
			body.IsTunnelerEnabled = v.(bool)
		}
		if v := OptionalString(req.Parameters, "roleAttributes"); v != "" {
			attrs := models.Attributes(SplitCSV(v))
			body.RoleAttributes = &attrs
		}
		if v := OptionalInt64(req.Parameters, "cost"); v != nil {
			body.Cost = v
		}
		if v, exists := req.Parameters["noTraversal"]; exists && v != nil {
			b := v.(bool)
			body.NoTraversal = &b
		}
		if v, exists := req.Parameters["disabled"]; exists && v != nil {
			b := v.(bool)
			body.Disabled = &b
		}

		return client.WithAuthenticatedClient(req, cfg, "update edge router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.PatchEdgeRouter(edge_router.NewPatchEdgeRouterParams().WithID(id).WithEdgeRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterIdentities
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterIdentities",
		Description: "List all Identities accessible by a specific Edge Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Identities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router identities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.ListEdgeRouterIdentities(edge_router.NewListEdgeRouterIdentitiesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterServices
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterServices",
		Description: "List all Services accessible by a specific Edge Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.ListEdgeRouterServices(edge_router.NewListEdgeRouterServicesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterEdgeRouterPolicies",
		Description: "List all Edge Router Policies that apply to a specific Edge Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.ListEdgeRouterEdgeRouterPolicies(edge_router.NewListEdgeRouterEdgeRouterPoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterServiceEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterServiceEdgeRouterPolicies",
		Description: "List all Service Edge Router Policies that apply to a specific Edge Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Service Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list edge router service edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.ListEdgeRouterServiceEdgeRouterPolicies(edge_router.NewListEdgeRouterServiceEdgeRouterPoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEdgeRouterRoleAttributes
	r.Register(tools.ToolDef{
		Name:        "listEdgeRouterRoleAttributes",
		Description: "List all role attributes in use by Edge Routers in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Edge Router Role Attributes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list edge router role attributes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.RoleAttributes.ListEdgeRouterRoleAttributes(nil, noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// reEnrollEdgeRouter
	r.Register(tools.ToolDef{
		Name:        "reEnrollEdgeRouter",
		Description: "Re-enroll an Edge Router, generating new certificates",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Re-enroll Edge Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "re-enroll edge router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.EdgeRouter.ReEnrollEdgeRouter(edge_router.NewReEnrollEdgeRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
