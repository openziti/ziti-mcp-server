package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	currentidentity "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/current_identity"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerMFA(r *tools.Registry, s *store.Store) {
	// getMfaStatus
	r.Register(tools.ToolDef{
		Name:        "getMfaStatus",
		Description: "Get the MFA status for the current identity",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get MFA Status"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get MFA status", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.DetailMfa(currentidentity.NewDetailMfaParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// enrollMfa
	r.Register(tools.ToolDef{
		Name:        "enrollMfa",
		Description: "Enroll the current identity in MFA",
		InputSchema: emptySchema(),
		Meta:        writeMeta(),
		Annotations: createAnnotations("Enroll MFA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "enroll MFA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.EnrollMfa(currentidentity.NewEnrollMfaParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// verifyMfa
	r.Register(tools.ToolDef{
		Name:        "verifyMfa",
		Description: "Verify MFA enrollment with a TOTP code",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{"type": "string", "description": "TOTP code to verify"},
			},
			"required": []string{"code"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Verify MFA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		code, errResp, ok := RequireString(req.Parameters, "code")
		if !ok {
			return *errResp, nil
		}
		body := &models.MfaCode{Code: strPtr(code)}
		return client.WithAuthenticatedClient(req, cfg, "verify MFA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.VerifyMfa(currentidentity.NewVerifyMfaParams().WithMfaValidation(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteMfa
	r.Register(tools.ToolDef{
		Name:        "deleteMfa",
		Description: "Remove MFA from the current identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{"type": "string", "description": "TOTP code to authorize deletion"},
			},
			"required": []string{"code"},
		},
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete MFA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		code, errResp, ok := RequireString(req.Parameters, "code")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete MFA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.DeleteMfa(currentidentity.NewDeleteMfaParams().WithMfaValidationCode(&code), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getMfaQrCode
	r.Register(tools.ToolDef{
		Name:        "getMfaQrCode",
		Description: "Get the MFA QR code for enrollment",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get MFA QR Code"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get MFA QR code", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				_, err := ec.CurrentIdentity.DetailMfaQrCode(currentidentity.NewDetailMfaQrCodeParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return map[string]any{"status": "ok", "message": "MFA QR code retrieved successfully"}, nil
			},
		), nil
	})

	// getMfaRecoveryCodes
	r.Register(tools.ToolDef{
		Name:        "getMfaRecoveryCodes",
		Description: "Get MFA recovery codes for the current identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{"type": "string", "description": "TOTP code to authorize viewing recovery codes"},
			},
			"required": []string{"code"},
		},
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get MFA Recovery Codes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		code, errResp, ok := RequireString(req.Parameters, "code")
		if !ok {
			return *errResp, nil
		}
		body := &models.MfaCode{Code: strPtr(code)}
		return client.WithAuthenticatedClient(req, cfg, "get MFA recovery codes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.DetailMfaRecoveryCodes(
					currentidentity.NewDetailMfaRecoveryCodesParams().WithMfaValidationCode(&code).WithMfaValidation(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createMfaRecoveryCodes
	r.Register(tools.ToolDef{
		Name:        "createMfaRecoveryCodes",
		Description: "Generate new MFA recovery codes (invalidates existing codes)",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{"type": "string", "description": "TOTP code to authorize generating new recovery codes"},
			},
			"required": []string{"code"},
		},
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Create MFA Recovery Codes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		code, errResp, ok := RequireString(req.Parameters, "code")
		if !ok {
			return *errResp, nil
		}
		body := &models.MfaCode{Code: strPtr(code)}
		return client.WithAuthenticatedClient(req, cfg, "create MFA recovery codes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.CurrentIdentity.CreateMfaRecoveryCodes(
					currentidentity.NewCreateMfaRecoveryCodesParams().WithMfaValidation(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
