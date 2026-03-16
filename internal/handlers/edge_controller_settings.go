package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	settings "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/settings"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerControllerSettings(r *tools.Registry, s *store.Store) {
	// listControllerSettings
	r.Register(tools.ToolDef{
		Name:        "listControllerSettings",
		Description: "List all Controller Settings in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Controller Settings"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list controller settings", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.ListControllerSettings(settings.NewListControllerSettingsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listControllerSetting
	r.Register(tools.ToolDef{
		Name:        "listControllerSetting",
		Description: "Get details about a specific Ziti Controller Setting",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Controller Setting Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get controller setting", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.DetailControllerSetting(settings.NewDetailControllerSettingParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createControllerSetting
	r.Register(tools.ToolDef{
		Name:        "createControllerSetting",
		Description: "Create a new Ziti Controller Setting.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"controllerId":      map[string]any{"type": "string", "description": "Controller ID to associate the setting with"},
				"oidcRedirectUris":  map[string]any{"type": "string", "description": "Comma-separated OIDC redirect URIs"},
				"oidcPostLogoutUris": map[string]any{"type": "string", "description": "Comma-separated OIDC post-logout redirect URIs"},
			},
			"required": []string{"controllerId"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Controller Setting"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		controllerID, errResp, ok := RequireString(req.Parameters, "controllerId")
		if !ok {
			return *errResp, nil
		}

		body := &models.ControllerSettingCreate{
			ControllerID: strPtr(controllerID),
		}

		redirectUris := SplitCSV(OptionalString(req.Parameters, "oidcRedirectUris"))
		postLogoutUris := SplitCSV(OptionalString(req.Parameters, "oidcPostLogoutUris"))

		if len(redirectUris) > 0 || len(postLogoutUris) > 0 {
			oidc := &models.ControllerSettingsOidc{
				RedirectUris:   redirectUris,
				PostLogoutUris: postLogoutUris,
			}
			body.Oidc = oidc
		}

		return client.WithAuthenticatedClient(req, cfg, "create controller setting", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.CreateControllerSetting(settings.NewCreateControllerSettingParams().WithControllerSetting(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteControllerSetting
	r.Register(tools.ToolDef{
		Name:        "deleteControllerSetting",
		Description: "Delete a Ziti Controller Setting.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Controller Setting"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete controller setting", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.DeleteControllerSetting(settings.NewDeleteControllerSettingParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateControllerSetting
	r.Register(tools.ToolDef{
		Name:        "updateControllerSetting",
		Description: "Update an existing Ziti Controller Setting.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                map[string]any{"type": "string", "description": "Controller Setting ID"},
				"oidcRedirectUris":  map[string]any{"type": "string", "description": "Comma-separated OIDC redirect URIs"},
				"oidcPostLogoutUris": map[string]any{"type": "string", "description": "Comma-separated OIDC post-logout redirect URIs"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Controller Setting"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}

		body := &models.ControllerSettingPatch{}

		redirectUris := SplitCSV(OptionalString(req.Parameters, "oidcRedirectUris"))
		postLogoutUris := SplitCSV(OptionalString(req.Parameters, "oidcPostLogoutUris"))

		if len(redirectUris) > 0 || len(postLogoutUris) > 0 {
			oidc := &models.ControllerSettingsOidc{
				RedirectUris:   redirectUris,
				PostLogoutUris: postLogoutUris,
			}
			body.Oidc = oidc
		}

		return client.WithAuthenticatedClient(req, cfg, "update controller setting", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.PatchControllerSetting(
					settings.NewPatchControllerSettingParams().WithID(id).WithControllerSetting(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// detailControllerSettingEffective
	r.Register(tools.ToolDef{
		Name:        "detailControllerSettingEffective",
		Description: "Get the effective (merged) value of a specific Controller Setting, reflecting all overrides",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Effective Controller Setting"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get effective controller setting", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Settings.DetailControllerSettingEffective(settings.NewDetailControllerSettingEffectiveParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
