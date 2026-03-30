package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	fcircuit "github.com/openziti/ziti-mcp-server/internal/gen/fabric/client/circuit"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerFabricCircuits(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listCircuits", Description: "List all active Circuits in the Ziti fabric network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Circuits"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list circuits", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := fc.Circuit.ListCircuits(
						fcircuit.NewListCircuitsParams().WithLimit(&limit).WithOffset(&offset))
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
		Name: "listCircuit", Description: "Get details about a specific Ziti Circuit",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Circuit Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get circuit", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Circuit.DetailCircuit(fcircuit.NewDetailCircuitParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteCircuit", Description: "Delete/tear down an active Ziti Circuit",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Circuit"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete circuit", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Circuit.DeleteCircuit(fcircuit.NewDeleteCircuitParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
