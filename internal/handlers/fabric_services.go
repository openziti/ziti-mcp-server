package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	fservice "github.com/openziti/ziti-mcp-server/internal/gen/fabric/client/service"
	fmodels "github.com/openziti/ziti-mcp-server/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerFabricServices(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listFabricServices", Description: "List all Fabric Services in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Fabric Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list fabric services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := fc.Service.ListServices(
						fservice.NewListServicesParams().WithLimit(&limit).WithOffset(&offset))
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
		Name: "listFabricService", Description: "Get details about a specific Ziti Fabric Service",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Fabric Service Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get fabric service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Service.DetailService(fservice.NewDetailServiceParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createFabricService", Description: "Create a new Ziti Fabric Service.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":               map[string]any{"type": "string", "description": "Service name"},
				"terminatorStrategy": map[string]any{"type": "string", "description": "Terminator strategy"},
			},
			"required": []string{"name"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Fabric Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.ServiceCreate{
			Name:               strPtr(name),
			TerminatorStrategy: OptionalString(req.Parameters, "terminatorStrategy"),
		}
		return client.WithAuthenticatedClient(req, cfg, "create fabric service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Service.CreateService(fservice.NewCreateServiceParams().WithService(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteFabricService", Description: "Delete a Ziti Fabric Service.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Fabric Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete fabric service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Service.DeleteService(fservice.NewDeleteServiceParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateFabricService", Description: "Update an existing Ziti Fabric Service.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                 map[string]any{"type": "string", "description": "Service ID"},
				"name":               map[string]any{"type": "string", "description": "New name"},
				"terminatorStrategy": map[string]any{"type": "string", "description": "Terminator strategy"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Fabric Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.ServicePatch{
			Name:               OptionalString(req.Parameters, "name"),
			TerminatorStrategy: OptionalString(req.Parameters, "terminatorStrategy"),
		}
		return client.WithAuthenticatedClient(req, cfg, "update fabric service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Service.PatchService(fservice.NewPatchServiceParams().WithID(id).WithService(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listFabricServiceTerminators", Description: "List all Terminators assigned to a specific Fabric Service",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Fabric Service Terminators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list fabric service terminators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := fc.Service.ListServiceTerminators(
						fservice.NewListServiceTerminatorsParams().WithID(id).WithLimit(&limit).WithOffset(&offset))
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
}
