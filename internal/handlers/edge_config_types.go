package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	edgeconfig "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/config"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerConfigTypes(r *tools.Registry, s *store.Store) {
	// listConfigTypes
	r.Register(tools.ToolDef{
		Name:        "listConfigTypes",
		Description: "List all Config Types in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Config Types"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list config types", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.Config.ListConfigTypes(
						edgeconfig.NewListConfigTypesParams().WithLimit(&limit).WithOffset(&offset), noAuth)
					if err != nil {
						return nil, err
					}
					m, err := ToMap(resp.Payload)
					if err != nil {
						return nil, err
					}
					// Strip full JSON schemas from list response to reduce token usage.
					// Schemas are still available via the detail (listConfigType) endpoint.
					stripFieldFromDataItems(m, "schema")
					if mm, ok := m.(map[string]any); ok {
						return mm, nil
					}
					return nil, nil
				})
			},
		), nil
	})

	// listConfigType
	r.Register(tools.ToolDef{
		Name:        "listConfigType",
		Description: "Get details about a specific Ziti Config Type",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Config Type Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get config type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.DetailConfigType(edgeconfig.NewDetailConfigTypeParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createConfigType
	r.Register(tools.ToolDef{
		Name:        "createConfigType",
		Description: "Create a new Ziti Config Type.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":   map[string]any{"type": "string", "description": "Name of the config type"},
				"schema": map[string]any{"type": "object", "description": "JSON schema to enforce configuration against"},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Config Type"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		schema := OptionalObject(req.Parameters, "schema")

		body := &models.ConfigTypeCreate{
			Name:   strPtr(name),
			Schema: schema,
		}

		return client.WithAuthenticatedClient(req, cfg, "create config type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.CreateConfigType(edgeconfig.NewCreateConfigTypeParams().WithConfigType(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteConfigType
	r.Register(tools.ToolDef{
		Name:        "deleteConfigType",
		Description: "Delete a Ziti Config Type.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Config Type"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete config type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.DeleteConfigType(edgeconfig.NewDeleteConfigTypeParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateConfigType
	r.Register(tools.ToolDef{
		Name:        "updateConfigType",
		Description: "Update an existing Ziti Config Type.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":     map[string]any{"type": "string", "description": "Config Type ID"},
				"name":   map[string]any{"type": "string", "description": "New name"},
				"schema": map[string]any{"type": "object", "description": "New JSON schema"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Config Type"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ConfigTypePatch{
			Name:   OptionalString(req.Parameters, "name"),
			Schema: OptionalObject(req.Parameters, "schema"),
		}

		return client.WithAuthenticatedClient(req, cfg, "update config type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.PatchConfigType(edgeconfig.NewPatchConfigTypeParams().WithID(id).WithConfigType(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listConfigsForConfigType
	r.Register(tools.ToolDef{
		Name:        "listConfigsForConfigType",
		Description: "List all Configs that use a specific Config Type",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Configs For Config Type"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list configs for config type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.ListConfigsForConfigType(edgeconfig.NewListConfigsForConfigTypeParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
