package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	flink "github.com/openziti/ziti-mcp-server/internal/gen/fabric/client/link"
	fmodels "github.com/openziti/ziti-mcp-server/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerFabricLinks(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listLinks", Description: "List all Links (router-to-router connections) in the Ziti fabric network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Links"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list links", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Link.ListLinks(flink.NewListLinksParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listLink", Description: "Get details about a specific Ziti Link",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Link Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get link", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Link.DetailLink(flink.NewDetailLinkParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deleteLink", Description: "Delete a Ziti Link between routers",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Link"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete link", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Link.DeleteLink(flink.NewDeleteLinkParams().WithID(id))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updateLink", Description: "Update an existing Ziti Link",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "description": "Link ID"},
				"down":       map[string]any{"type": "boolean", "description": "Mark link as down"},
				"staticCost": map[string]any{"type": "number", "description": "Static cost value"},
			},
			"required": []string{"id"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Link"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.LinkPatch{
			Down: OptionalBool(req.Parameters, "down", false),
		}
		if c := OptionalInt64(req.Parameters, "staticCost"); c != nil {
			body.StaticCost = *c
		}
		return client.WithAuthenticatedClient(req, cfg, "update link", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Link.PatchLink(flink.NewPatchLinkParams().WithID(id).WithLink(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
