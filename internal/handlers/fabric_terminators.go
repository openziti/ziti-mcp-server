package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	fterminator "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/client/terminator"
	fmodels "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerFabricTerminators(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listFabricTerminators", Description: "List all Fabric Terminators in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Fabric Terminators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list fabric terminators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Terminator.ListTerminators(fterminator.NewListTerminatorsParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listFabricTerminator", Description: "Get details about a specific Ziti Fabric Terminator",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Fabric Terminator Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get fabric terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Terminator.DetailTerminator(fterminator.NewDetailTerminatorParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createFabricTerminator", Description: "Create a new Ziti Fabric Terminator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"service":    map[string]any{"type": "string", "description": "Service ID"},
				"router":     map[string]any{"type": "string", "description": "Router ID"},
				"binding":    map[string]any{"type": "string", "description": "Binding type"},
				"address":    map[string]any{"type": "string", "description": "Address"},
				"cost":       map[string]any{"type": "number", "description": "Cost value"},
				"precedence": map[string]any{"type": "string", "description": "Precedence (default, required, failed)"},
			},
			"required": []string{"service", "router", "binding", "address"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Fabric Terminator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		svc, errResp, ok := RequireString(req.Parameters, "service")
		if !ok {
			return *errResp, nil
		}
		rtr, errResp, ok := RequireString(req.Parameters, "router")
		if !ok {
			return *errResp, nil
		}
		binding, errResp, ok := RequireString(req.Parameters, "binding")
		if !ok {
			return *errResp, nil
		}
		addr, errResp, ok := RequireString(req.Parameters, "address")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.TerminatorCreate{
			Service: strPtr(svc),
			Router:  strPtr(rtr),
			Binding: strPtr(binding),
			Address: strPtr(addr),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			tc := fmodels.TerminatorCost(*c)
			body.Cost = &tc
		}
		if p := OptionalString(req.Parameters, "precedence"); p != "" {
			body.Precedence = fmodels.TerminatorPrecedence(p)
		}
		return client.WithAuthenticatedClient(req, cfg, "create fabric terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Terminator.CreateTerminator(fterminator.NewCreateTerminatorParams().WithTerminator(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteFabricTerminator", Description: "Delete a Ziti Fabric Terminator.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Fabric Terminator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete fabric terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Terminator.DeleteTerminator(fterminator.NewDeleteTerminatorParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateFabricTerminator", Description: "Update an existing Ziti Fabric Terminator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Terminator ID"},
				"service":    map[string]any{"type": "string"},
				"router":     map[string]any{"type": "string"},
				"binding":    map[string]any{"type": "string"},
				"address":    map[string]any{"type": "string"},
				"cost":       map[string]any{"type": "number"},
				"precedence": map[string]any{"type": "string"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Fabric Terminator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.TerminatorPatch{
			Service: OptionalString(req.Parameters, "service"),
			Router:  OptionalString(req.Parameters, "router"),
			Binding: OptionalString(req.Parameters, "binding"),
			Address: OptionalString(req.Parameters, "address"),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			tc := fmodels.TerminatorCost(*c)
			body.Cost = &tc
		}
		if p := OptionalString(req.Parameters, "precedence"); p != "" {
			body.Precedence = fmodels.TerminatorPrecedence(p)
		}
		return client.WithAuthenticatedClient(req, cfg, "update fabric terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Terminator.PatchTerminator(fterminator.NewPatchTerminatorParams().WithID(id).WithTerminator(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
