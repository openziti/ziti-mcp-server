package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	edgeconfig "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/config"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerConfigs(r *tools.Registry, s *store.Store) {
	// listConfigs
	r.Register(tools.ToolDef{
		Name:        "listConfigs",
		Description: "List all Configs in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Configs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list configs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.ListConfigs(edgeconfig.NewListConfigsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listConfig
	r.Register(tools.ToolDef{
		Name:        "listConfig",
		Description: "Get details about a specific Ziti Config",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Config Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get config", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.DetailConfig(edgeconfig.NewDetailConfigParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createConfig
	r.Register(tools.ToolDef{
		Name:        "createConfig",
		Description: "Create a new Ziti Config.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":         map[string]any{"type": "string", "description": "Name of the config"},
				"configTypeId": map[string]any{"type": "string", "description": "Config type ID"},
				"data":         map[string]any{"type": "object", "description": "Config data payload"},
			},
			"required": []string{"name", "configTypeId", "data"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Config"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		configTypeId, errResp, ok := RequireString(req.Parameters, "configTypeId")
		if !ok {
			return *errResp, nil
		}
		data, errResp, ok := RequireObject(req.Parameters, "data")
		if !ok {
			return *errResp, nil
		}

		body := &models.ConfigCreate{
			Name:         strPtr(name),
			ConfigTypeID: strPtr(configTypeId),
			Data:         data,
		}

		return client.WithAuthenticatedClient(req, cfg, "create config", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.CreateConfig(edgeconfig.NewCreateConfigParams().WithConfig(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteConfig
	r.Register(tools.ToolDef{
		Name:        "deleteConfig",
		Description: "Delete a Ziti Config.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Config"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete config", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.DeleteConfig(edgeconfig.NewDeleteConfigParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateConfig
	r.Register(tools.ToolDef{
		Name:        "updateConfig",
		Description: "Update an existing Ziti Config.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":   map[string]any{"type": "string", "description": "Config ID"},
				"name": map[string]any{"type": "string", "description": "New name"},
				"data": map[string]any{"type": "object", "description": "New config data payload"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Config"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ConfigPatch{
			Name: OptionalString(req.Parameters, "name"),
			Data: OptionalObject(req.Parameters, "data"),
		}

		return client.WithAuthenticatedClient(req, cfg, "update config", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.PatchConfig(edgeconfig.NewPatchConfigParams().WithID(id).WithConfig(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listConfigServices
	r.Register(tools.ToolDef{
		Name:        "listConfigServices",
		Description: "List all Services that use a specific Config",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Config Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list config services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Config.ListConfigServices(edgeconfig.NewListConfigServicesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
