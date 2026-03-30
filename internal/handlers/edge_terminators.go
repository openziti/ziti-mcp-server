package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	terminator "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/terminator"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerTerminators(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listTerminators", Description: "List all Terminators in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Terminators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list terminators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.Terminator.ListTerminators(
						terminator.NewListTerminatorsParams().WithLimit(&limit).WithOffset(&offset), noAuth)
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
		Name: "listTerminator", Description: "Get details about a specific Ziti Terminator",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Terminator Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Terminator.DetailTerminator(terminator.NewDetailTerminatorParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createTerminator", Description: "Create a new Ziti Terminator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"service":    map[string]any{"type": "string", "description": "Service ID"},
				"router":     map[string]any{"type": "string", "description": "Router ID"},
				"binding":    map[string]any{"type": "string", "description": "Binding type"},
				"address":    map[string]any{"type": "string", "description": "Address"},
				"cost":       map[string]any{"type": "number", "description": "Cost value"},
				"precedence": map[string]any{"type": "string", "description": "Precedence (default, required, failed)"},
				"identity":   map[string]any{"type": "string", "description": "Identity"},
			},
			"required": []string{"service", "router", "binding", "address"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Terminator"),
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
		body := &models.TerminatorCreate{
			Service: strPtr(svc),
			Router:  strPtr(rtr),
			Binding: strPtr(binding),
			Address: strPtr(addr),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			tc := models.TerminatorCost(*c)
			body.Cost = &tc
		}
		if p := OptionalString(req.Parameters, "precedence"); p != "" {
			body.Precedence = models.TerminatorPrecedence(p)
		}
		body.Identity = OptionalString(req.Parameters, "identity")

		return client.WithAuthenticatedClient(req, cfg, "create terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Terminator.CreateTerminator(terminator.NewCreateTerminatorParams().WithTerminator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteTerminator", Description: "Delete a Ziti Terminator.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Terminator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Terminator.DeleteTerminator(terminator.NewDeleteTerminatorParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateTerminator", Description: "Update an existing Ziti Terminator.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Terminator ID"},
				"service":    map[string]any{"type": "string", "description": "Service ID"},
				"router":     map[string]any{"type": "string", "description": "Router ID"},
				"binding":    map[string]any{"type": "string", "description": "Binding type"},
				"address":    map[string]any{"type": "string", "description": "Address"},
				"cost":       map[string]any{"type": "number", "description": "Cost value"},
				"precedence": map[string]any{"type": "string", "description": "Precedence (default, required, failed)"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Terminator"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.TerminatorPatch{
			Service: OptionalString(req.Parameters, "service"),
			Router:  OptionalString(req.Parameters, "router"),
			Binding: OptionalString(req.Parameters, "binding"),
			Address: OptionalString(req.Parameters, "address"),
		}
		if c := OptionalInt64(req.Parameters, "cost"); c != nil {
			tc := models.TerminatorCost(*c)
			body.Cost = &tc
		}
		if p := OptionalString(req.Parameters, "precedence"); p != "" {
			body.Precedence = models.TerminatorPrecedence(p)
		}
		return client.WithAuthenticatedClient(req, cfg, "update terminator", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Terminator.PatchTerminator(terminator.NewPatchTerminatorParams().WithID(id).WithTerminator(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
