package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	frouter "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/client/router"
	fmodels "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerFabricRouters(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listFabricRouters", Description: "List all Fabric Routers in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Fabric Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list fabric routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.ListRouters(frouter.NewListRoutersParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listFabricRouter", Description: "Get details about a specific Ziti Fabric Router",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Fabric Router Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get fabric router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.DetailRouter(frouter.NewDetailRouterParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createFabricRouter", Description: "Create a new Ziti Fabric Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":          map[string]any{"type": "string", "description": "Router ID"},
				"name":        map[string]any{"type": "string", "description": "Router name"},
				"cost":        map[string]any{"type": "number", "description": "Cost value"},
				"noTraversal": map[string]any{"type": "boolean", "description": "Disable traversal", "default": false},
				"disabled":    map[string]any{"type": "boolean", "description": "Disable router", "default": false},
			},
			"required": []string{"id", "name"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Fabric Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		routerID, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		noTraversal := OptionalBool(req.Parameters, "noTraversal", false)
		body := &fmodels.RouterCreate{
			ID:          strPtr(routerID),
			Name:        strPtr(name),
			NoTraversal: boolPtr(noTraversal),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			body.Cost = c
		} else {
			body.Cost = int64Ptr(0)
		}
		if v, exists := req.Parameters["disabled"]; exists && v != nil {
			body.Disabled = boolPtr(v.(bool))
		}
		return client.WithAuthenticatedClient(req, cfg, "create fabric router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.CreateRouter(frouter.NewCreateRouterParams().WithRouter(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteFabricRouter", Description: "Delete a Ziti Fabric Router.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Fabric Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete fabric router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.DeleteRouter(frouter.NewDeleteRouterParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateFabricRouter", Description: "Update an existing Ziti Fabric Router.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":          map[string]any{"type": "string", "description": "Router ID"},
				"name":        map[string]any{"type": "string", "description": "New name"},
				"cost":        map[string]any{"type": "number", "description": "Cost value"},
				"noTraversal": map[string]any{"type": "boolean"},
				"disabled":    map[string]any{"type": "boolean"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Fabric Router"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.RouterPatch{
			Name: OptionalString(req.Parameters, "name"),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			body.Cost = c
		}
		if v, exists := req.Parameters["noTraversal"]; exists && v != nil {
			body.NoTraversal = boolPtr(v.(bool))
		}
		if v, exists := req.Parameters["disabled"]; exists && v != nil {
			body.Disabled = boolPtr(v.(bool))
		}
		return client.WithAuthenticatedClient(req, cfg, "update fabric router", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.PatchRouter(frouter.NewPatchRouterParams().WithID(id).WithRouter(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listFabricRouterTerminators", Description: "List all Terminators assigned to a specific Fabric Router",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Fabric Router Terminators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list fabric router terminators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Router.ListRouterTerminators(frouter.NewListRouterTerminatorsParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
