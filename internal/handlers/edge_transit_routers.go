package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/router"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerTransitRouters(r *tools.Registry, s *store.Store) {
	// listTransitRouters
	r.Register(tools.ToolDef{
		Name:        "listTransitRouters",
		Description: "List all Transit Routers in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Transit Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list transit routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.ListTransitRouters(router.NewListTransitRoutersParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listTransitRouter
	r.Register(tools.ToolDef{
		Name:        "listTransitRouter",
		Description: "Get details about a specific Ziti Transit Router",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Transit Router Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get transit router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.DetailTransitRouter(router.NewDetailTransitRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createTransitRouter
	r.Register(tools.ToolDef{
		Name:        "createTransitRouter",
		Description: "Create a new Ziti Transit Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "Name of the transit router"},
				"cost":        map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal": map[string]any{"type": "boolean", "description": "Whether to disable traversal"},
				"disabled":    map[string]any{"type": "boolean", "description": "Whether the transit router is disabled"},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Transit Router"),
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

		return client.WithAuthenticatedClient(req, cfg, "create transit router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.CreateTransitRouter(router.NewCreateTransitRouterParams().WithRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteTransitRouter
	r.Register(tools.ToolDef{
		Name:        "deleteTransitRouter",
		Description: "Delete a Ziti Transit Router.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Transit Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete transit router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.DeleteTransitRouter(router.NewDeleteTransitRouterParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateTransitRouter
	r.Register(tools.ToolDef{
		Name:        "updateTransitRouter",
		Description: "Update an existing Ziti Transit Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":          map[string]any{"type": "string", "description": "Transit Router ID"},
				"name":        map[string]any{"type": "string", "description": "New name"},
				"cost":        map[string]any{"type": "number", "description": "Cost for routing (0-65535)"},
				"noTraversal": map[string]any{"type": "boolean", "description": "Whether to disable traversal"},
				"disabled":    map[string]any{"type": "boolean", "description": "Whether the transit router is disabled"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Transit Router"),
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

		return client.WithAuthenticatedClient(req, cfg, "update transit router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Router.PatchTransitRouter(router.NewPatchTransitRouterParams().WithID(id).WithRouter(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
