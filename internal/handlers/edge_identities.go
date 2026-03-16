package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	identity "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/identity"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerIdentities(r *tools.Registry, s *store.Store) {
	// listIdentities
	r.Register(tools.ToolDef{
		Name:        "listIdentities",
		Description: "List all Identities in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list identities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentities(identity.NewListIdentitiesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentity
	r.Register(tools.ToolDef{
		Name:        "listIdentity",
		Description: "Get details about a specific Ziti identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.DetailIdentity(identity.NewDetailIdentityParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createIdentity
	r.Register(tools.ToolDef{
		Name:        "createIdentity",
		Description: "Create a new Ziti Identity.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":           map[string]any{"type": "string", "description": "Name of the identity"},
				"admin":          map[string]any{"type": "boolean", "description": "Whether the identity is an admin", "default": false},
				"authPolicy":     map[string]any{"type": "string", "description": "Auth policy ID", "default": "default"},
				"externalId":     map[string]any{"type": "string", "description": "External ID for the identity"},
				"roleAttributes": map[string]any{"type": "string", "description": "Comma-separated role attributes"},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Identity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		admin := OptionalBool(req.Parameters, "admin", false)
		authPolicy := OptionalString(req.Parameters, "authPolicy")
		if authPolicy == "" {
			authPolicy = "default"
		}
		externalId := OptionalString(req.Parameters, "externalId")
		roleAttrs := SplitCSV(OptionalString(req.Parameters, "roleAttributes"))

		body := &models.IdentityCreate{
			Name:         strPtr(name),
			IsAdmin:      boolPtr(admin),
			Type:         models.IdentityType("Default").Pointer(),
			AuthPolicyID: strPtr(authPolicy),
		}
		if externalId != "" {
			body.ExternalID = strPtr(externalId)
		}
		if len(roleAttrs) > 0 {
			attrs := models.Attributes(roleAttrs)
			body.RoleAttributes = &attrs
		}

		return client.WithAuthenticatedClient(req, cfg, "create identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.CreateIdentity(identity.NewCreateIdentityParams().WithIdentity(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteIdentity
	r.Register(tools.ToolDef{
		Name:        "deleteIdentity",
		Description: "Delete a Ziti Identity.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Identity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.DeleteIdentity(identity.NewDeleteIdentityParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateIdentity
	r.Register(tools.ToolDef{
		Name:        "updateIdentity",
		Description: "Update an existing Ziti Identity.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":             map[string]any{"type": "string", "description": "Identity ID"},
				"name":           map[string]any{"type": "string", "description": "New name"},
				"admin":          map[string]any{"type": "boolean", "description": "Whether the identity is an admin"},
				"authPolicy":     map[string]any{"type": "string", "description": "Auth policy ID"},
				"externalId":     map[string]any{"type": "string", "description": "External ID"},
				"roleAttributes": map[string]any{"type": "string", "description": "Comma-separated role attributes"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Identity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.IdentityPatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = strPtr(v)
		}
		if v, exists := req.Parameters["admin"]; exists && v != nil {
			b := v.(bool)
			body.IsAdmin = &b
		}
		if v := OptionalString(req.Parameters, "authPolicy"); v != "" {
			body.AuthPolicyID = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "externalId"); v != "" {
			body.ExternalID = strPtr(v)
		}
		if v := OptionalString(req.Parameters, "roleAttributes"); v != "" {
			attrs := models.Attributes(SplitCSV(v))
			body.RoleAttributes = &attrs
		}

		return client.WithAuthenticatedClient(req, cfg, "update identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.PatchIdentity(identity.NewPatchIdentityParams().WithID(id).WithIdentity(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityServices
	r.Register(tools.ToolDef{
		Name:        "listIdentityServices",
		Description: "List all Services accessible by a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list identity services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentityServices(identity.NewListIdentityServicesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityEdgeRouters
	r.Register(tools.ToolDef{
		Name:        "listIdentityEdgeRouters",
		Description: "List all Edge Routers accessible by a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Edge Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list identity edge routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentityEdgeRouters(identity.NewListIdentityEdgeRoutersParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityServicePolicies
	r.Register(tools.ToolDef{
		Name:        "listIdentityServicePolicies",
		Description: "List all Service Policies that apply to a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Service Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list identity service policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentityServicePolicies(identity.NewListIdentityServicePoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listIdentityEdgeRouterPolicies",
		Description: "List all Edge Router Policies that apply to a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list identity edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentitysEdgeRouterPolicies(identity.NewListIdentitysEdgeRouterPoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityServiceConfigs
	r.Register(tools.ToolDef{
		Name:        "listIdentityServiceConfigs",
		Description: "List all Service Configs associated with a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Service Configs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list identity service configs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.ListIdentitysServiceConfigs(identity.NewListIdentitysServiceConfigsParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getIdentityPolicyAdvice
	r.Register(tools.ToolDef{
		Name:        "getIdentityPolicyAdvice",
		Description: "Check whether an Identity can dial or bind a specific Service and get policy advice explaining why or why not",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":        map[string]any{"type": "string", "description": "Identity ID"},
				"serviceId": map[string]any{"type": "string", "description": "Service ID"},
			},
			"required": []string{"id", "serviceId"},
		},
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Policy Advice"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		serviceID, errResp, ok := RequireString(req.Parameters, "serviceId")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity policy advice", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.GetIdentityPolicyAdvice(identity.NewGetIdentityPolicyAdviceParams().WithID(id).WithServiceID(serviceID), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listIdentityRoleAttributes
	r.Register(tools.ToolDef{
		Name:        "listIdentityRoleAttributes",
		Description: "List all role attributes in use by Identities in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Identity Role Attributes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list identity role attributes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.RoleAttributes.ListIdentityRoleAttributes(nil, noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// disableIdentity
	r.Register(tools.ToolDef{
		Name:        "disableIdentity",
		Description: "Temporarily disable a Ziti Identity for a specified duration",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":              map[string]any{"type": "string", "description": "Identity ID"},
				"durationMinutes": map[string]any{"type": "number", "description": "Duration in minutes to disable"},
			},
			"required": []string{"id", "durationMinutes"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Disable Identity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		dur := OptionalInt64(req.Parameters, "durationMinutes")
		if dur == nil {
			resp := errResponse("Missing required parameter: durationMinutes")
			return resp, nil
		}

		body := &models.DisableParams{
			DurationMinutes: dur,
		}

		return client.WithAuthenticatedClient(req, cfg, "disable identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.DisableIdentity(identity.NewDisableIdentityParams().WithID(id).WithDisable(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// enableIdentity
	r.Register(tools.ToolDef{
		Name:        "enableIdentity",
		Description: "Re-enable a previously disabled Ziti Identity",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: &tools.ToolAnnotations{Title: "Enable Identity"},
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "enable identity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.EnableIdentity(identity.NewEnableIdentityParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getIdentityAuthenticators
	r.Register(tools.ToolDef{
		Name:        "getIdentityAuthenticators",
		Description: "List all Authenticators for a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Authenticators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity authenticators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.GetIdentityAuthenticators(identity.NewGetIdentityAuthenticatorsParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getIdentityEnrollments
	r.Register(tools.ToolDef{
		Name:        "getIdentityEnrollments",
		Description: "List all Enrollments for a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Enrollments"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity enrollments", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.GetIdentityEnrollments(identity.NewGetIdentityEnrollmentsParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getIdentityFailedServiceRequests
	r.Register(tools.ToolDef{
		Name:        "getIdentityFailedServiceRequests",
		Description: "List failed service requests for a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Failed Service Requests"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity failed service requests", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.GetIdentityFailedServiceRequests(identity.NewGetIdentityFailedServiceRequestsParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// getIdentityPostureData
	r.Register(tools.ToolDef{
		Name:        "getIdentityPostureData",
		Description: "Get posture data for a specific Identity",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Identity Posture Data"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get identity posture data", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.GetIdentityPostureData(identity.NewGetIdentityPostureDataParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// removeIdentityMfa
	r.Register(tools.ToolDef{
		Name:        "removeIdentityMfa",
		Description: "Remove MFA from a specific Identity",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Remove Identity MFA"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "remove identity MFA", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.RemoveIdentityMfa(identity.NewRemoveIdentityMfaParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateIdentityTracing
	r.Register(tools.ToolDef{
		Name:        "updateIdentityTracing",
		Description: "Update tracing configuration for a specific Identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":       map[string]any{"type": "string", "description": "Identity ID"},
				"enabled":  map[string]any{"type": "boolean", "description": "Enable or disable tracing"},
				"duration": map[string]any{"type": "string", "description": "Trace duration"},
				"traceId":  map[string]any{"type": "string", "description": "Trace ID"},
				"channels": map[string]any{"type": "string", "description": "Comma-separated channel names"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Identity Tracing"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.TraceSpec{
			Enabled:  OptionalBool(req.Parameters, "enabled", false),
			Duration: OptionalString(req.Parameters, "duration"),
			TraceID:  OptionalString(req.Parameters, "traceId"),
			Channels: SplitCSV(OptionalString(req.Parameters, "channels")),
		}

		return client.WithAuthenticatedClient(req, cfg, "update identity tracing", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.UpdateIdentityTracing(identity.NewUpdateIdentityTracingParams().WithID(id).WithTraceSpec(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// associateIdentityServiceConfigs
	r.Register(tools.ToolDef{
		Name:        "associateIdentityServiceConfigs",
		Description: "Associate service configs with a specific Identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":             map[string]any{"type": "string", "description": "Identity ID"},
				"serviceConfigs": map[string]any{"type": "string", "description": "JSON array of {serviceId, configId} objects"},
			},
			"required": []string{"id", "serviceConfigs"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Associate Identity Service Configs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		scJSON, errResp, ok := RequireString(req.Parameters, "serviceConfigs")
		if !ok {
			return *errResp, nil
		}

		var assigns models.ServiceConfigsAssignList
		if err := json.Unmarshal([]byte(scJSON), &assigns); err != nil {
			return errResponse("Invalid serviceConfigs JSON: " + err.Error()), nil
		}

		return client.WithAuthenticatedClient(req, cfg, "associate identity service configs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.AssociateIdentitysServiceConfigs(
					identity.NewAssociateIdentitysServiceConfigsParams().WithID(id).WithServiceConfigs(assigns), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// disassociateIdentityServiceConfigs
	r.Register(tools.ToolDef{
		Name:        "disassociateIdentityServiceConfigs",
		Description: "Remove service config associations from a specific Identity",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":             map[string]any{"type": "string", "description": "Identity ID"},
				"serviceConfigs": map[string]any{"type": "string", "description": "JSON array of {serviceId, configId} objects"},
			},
			"required": []string{"id", "serviceConfigs"},
		},
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Disassociate Identity Service Configs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		scJSON, errResp, ok := RequireString(req.Parameters, "serviceConfigs")
		if !ok {
			return *errResp, nil
		}

		var assigns models.ServiceConfigsAssignList
		if err := json.Unmarshal([]byte(scJSON), &assigns); err != nil {
			return errResponse("Invalid serviceConfigs JSON: " + err.Error()), nil
		}

		return client.WithAuthenticatedClient(req, cfg, "disassociate identity service configs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Identity.DisassociateIdentitysServiceConfigs(
					identity.NewDisassociateIdentitysServiceConfigsParams().WithID(id).WithServiceConfigIDPairs(assigns), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
