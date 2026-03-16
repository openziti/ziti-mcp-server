package handlers

import (
	"net/http"

	"github.com/go-openapi/strfmt"

	"github.com/openziti/ziti-mcp-server/internal/client"
	enrollment "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/enrollment"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerEnrollments(r *tools.Registry, s *store.Store) {
	// listEnrollments
	r.Register(tools.ToolDef{
		Name:        "listEnrollments",
		Description: "List all Enrollments in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Enrollments"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list enrollments", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.ListEnrollments(enrollment.NewListEnrollmentsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEnrollment
	r.Register(tools.ToolDef{
		Name:        "listEnrollment",
		Description: "Get details about a specific Ziti Enrollment",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Enrollment Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get enrollment", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.DetailEnrollment(enrollment.NewDetailEnrollmentParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createEnrollment
	r.Register(tools.ToolDef{
		Name:        "createEnrollment",
		Description: "Create a new Ziti Enrollment.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"identityId": map[string]any{"type": "string", "description": "The identity ID to create the enrollment for"},
				"method":     map[string]any{"type": "string", "description": "Enrollment method", "enum": []string{"ott", "ottca", "updb"}},
				"expiresAt":  map[string]any{"type": "string", "description": "Expiration date-time in ISO 8601 format"},
				"caId":       map[string]any{"type": "string", "description": "CA ID (required for ottca method)"},
				"username":   map[string]any{"type": "string", "description": "Username (required for updb method)"},
			},
			"required": []string{"identityId", "method", "expiresAt"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Enrollment"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		identityID, errResp, ok := RequireString(req.Parameters, "identityId")
		if !ok {
			return *errResp, nil
		}
		method, errResp, ok := RequireString(req.Parameters, "method")
		if !ok {
			return *errResp, nil
		}
		expiresAtStr, errResp, ok := RequireString(req.Parameters, "expiresAt")
		if !ok {
			return *errResp, nil
		}
		expiresAt, err := strfmt.ParseDateTime(expiresAtStr)
		if err != nil {
			return errResponse("Invalid expiresAt format: " + err.Error()), nil
		}

		body := &models.EnrollmentCreate{
			IdentityID: strPtr(identityID),
			Method:     strPtr(method),
			ExpiresAt:  &expiresAt,
		}
		if v := OptionalString(req.Parameters, "caId"); v != "" {
			body.CaID = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "username"); v != "" {
			body.Username = strPtr(v)
		}

		return client.WithAuthenticatedClient(req, cfg, "create enrollment", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.CreateEnrollment(enrollment.NewCreateEnrollmentParams().WithEnrollment(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteEnrollment
	r.Register(tools.ToolDef{
		Name:        "deleteEnrollment",
		Description: "Delete a Ziti Enrollment.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Enrollment"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete enrollment", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.DeleteEnrollment(enrollment.NewDeleteEnrollmentParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// refreshEnrollment
	r.Register(tools.ToolDef{
		Name:        "refreshEnrollment",
		Description: "Refresh an expired enrollment, extending its expiration time",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":        map[string]any{"type": "string", "description": "The enrollment ID"},
				"expiresAt": map[string]any{"type": "string", "description": "New expiration date-time in ISO 8601 format"},
			},
			"required": []string{"id", "expiresAt"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Refresh Enrollment"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		expiresAtStr, errResp, ok := RequireString(req.Parameters, "expiresAt")
		if !ok {
			return *errResp, nil
		}
		expiresAt, err := strfmt.ParseDateTime(expiresAtStr)
		if err != nil {
			return errResponse("Invalid expiresAt format: " + err.Error()), nil
		}

		body := &models.EnrollmentRefresh{
			ExpiresAt: &expiresAt,
		}

		return client.WithAuthenticatedClient(req, cfg, "refresh enrollment", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.RefreshEnrollment(enrollment.NewRefreshEnrollmentParams().WithID(id).WithRefresh(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
