package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	posturecheck "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/posture_checks"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerPostureChecks(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listPostureChecks", Description: "List all Posture Checks in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Posture Checks"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list posture checks", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.ListPostureChecks(posturecheck.NewListPostureChecksParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listPostureCheck", Description: "Get details about a specific Ziti Posture Check",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Posture Check Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get posture check", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.DetailPostureCheck(posturecheck.NewDetailPostureCheckParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createPostureCheck", Description: "Create a new Ziti Posture Check.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":           map[string]any{"type": "string", "description": "Posture check name"},
				"typeId":         map[string]any{"type": "string", "description": "Type (OS, PROCESS, DOMAIN, MAC, MFA, PROCESS_MULTI)"},
				"roleAttributes": map[string]any{"type": "string", "description": "Comma-separated role attributes"},
			},
			"required": []string{"name", "typeId"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Posture Check"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		typeID, errResp, ok := RequireString(req.Parameters, "typeId")
		if !ok {
			return *errResp, nil
		}
		roleAttrs := SplitCSV(OptionalString(req.Parameters, "roleAttributes"))
		var attrs *models.Attributes
		if len(roleAttrs) > 0 {
			a := models.Attributes(roleAttrs)
			attrs = &a
		}

		// PostureCheckCreate is polymorphic; use MFA type as baseline
		mfa := &models.PostureCheckMfaCreate{}
		mfa.SetName(strPtr(name))
		mfa.SetTypeID(models.PostureCheckType(typeID))
		if attrs != nil {
			mfa.SetRoleAttributes(attrs)
		}

		return client.WithAuthenticatedClient(req, cfg, "create posture check", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.CreatePostureCheck(posturecheck.NewCreatePostureCheckParams().WithPostureCheck(mfa), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "deletePostureCheck", Description: "Delete a Ziti Posture Check.",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Delete Posture Check"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete posture check", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.DeletePostureCheck(posturecheck.NewDeletePostureCheckParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "updatePostureCheck", Description: "Update an existing Ziti Posture Check.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":             map[string]any{"type": "string", "description": "Posture check ID"},
				"typeId":         map[string]any{"type": "string", "description": "Type (OS, PROCESS, DOMAIN, MAC, MFA, PROCESS_MULTI)"},
				"name":           map[string]any{"type": "string", "description": "Posture check name"},
				"roleAttributes": map[string]any{"type": "string", "description": "Comma-separated role attributes"},
			},
			"required": []string{"id", "typeId"},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Update Posture Check"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		typeID, errResp, ok := RequireString(req.Parameters, "typeId")
		if !ok {
			return *errResp, nil
		}
		roleAttrs := SplitCSV(OptionalString(req.Parameters, "roleAttributes"))
		var attrs *models.Attributes
		if len(roleAttrs) > 0 {
			a := models.Attributes(roleAttrs)
			attrs = &a
		}

		mfa := &models.PostureCheckMfaPatch{}
		mfa.SetName(OptionalString(req.Parameters, "name"))
		mfa.SetTypeID(models.PostureCheckType(typeID))
		if attrs != nil {
			mfa.SetRoleAttributes(attrs)
		}

		return client.WithAuthenticatedClient(req, cfg, "update posture check", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.PatchPostureCheck(posturecheck.NewPatchPostureCheckParams().WithID(id).WithPostureCheck(mfa), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listPostureCheckRoleAttributes", Description: "List all role attributes in use by Posture Checks in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Posture Check Role Attributes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list posture check role attributes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.RoleAttributes.ListPostureCheckRoleAttributes(nil, noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "listPostureCheckTypes", Description: "List all available Posture Check Types in the Ziti network",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Posture Check Types"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list posture check types", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.ListPostureCheckTypes(posturecheck.NewListPostureCheckTypesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "detailPostureCheckType", Description: "Get details about a specific Posture Check Type",
		InputSchema: idSchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Posture Check Type Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get posture check type", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.PostureChecks.DetailPostureCheckType(posturecheck.NewDetailPostureCheckTypeParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
