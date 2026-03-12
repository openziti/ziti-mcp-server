package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/router"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerRouters(r *tools.Registry, s *store.Store) {
	// listRouters
	r.Register(tools.ToolDef{
		Name:        "listRouters",
		Description: "List all Routers in the Ziti network (generic routers via Edge API)",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.ListRouters(router.NewListRoutersParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getRouter
	r.Register(tools.ToolDef{
		Name:        "getRouter",
		Description: "Get details about a specific Ziti Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Router Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.DetailRouter(router.NewDetailRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createRouter
	r.Register(tools.ToolDef{
		Name:        "createRouter",
		Description: "Create a new Ziti Router",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "Name of the router"},
				"cost":        map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal": map[string]any{"type": "boolean", "description": "Whether to disable traversal"},
				"disabled":    map[string]any{"type": "boolean", "description": "Whether the router is disabled"},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}

		body := &models.RouterCreate{
			Name: strPtr(name),
		}
		if v := OptionalInt64(req.Parameters, "cost"); v != nil {
			body.Cost = v
		}
		if v, exists := req.Parameters["noTraversal"]; exists && v != nil {
			body.NoTraversal = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["disabled"]; exists && v != nil {
			body.Disabled = boolPtr(v.(bool))
		}

		return client.WithAuthenticatedClient(req, cfg, "create router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.CreateRouter(router.NewCreateRouterParams().WithRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteRouter
	r.Register(tools.ToolDef{
		Name:        "deleteRouter",
		Description: "Delete a Ziti Router",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.DeleteRouter(router.NewDeleteRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateRouter
	r.Register(tools.ToolDef{
		Name:        "updateRouter",
		Description: "Update an existing Ziti Router",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":          map[string]any{"type": "string", "description": "Router ID"},
				"name":        map[string]any{"type": "string", "description": "New name"},
				"cost":        map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal": map[string]any{"type": "boolean", "description": "Whether to disable traversal"},
				"disabled":    map[string]any{"type": "boolean", "description": "Whether the router is disabled"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.RouterPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = v
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

		return client.WithAuthenticatedClient(req, cfg, "update router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.PatchRouter(router.NewPatchRouterParams().WithID(id).WithRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
