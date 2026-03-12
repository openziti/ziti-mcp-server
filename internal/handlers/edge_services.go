package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/service"
	"github.com/openziti/ziti-mcp-server-go/internal/gen/edge/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerServices(r *tools.Registry, s *store.Store) {
	// listServices
	r.Register(tools.ToolDef{
		Name:        "listServices",
		Description: "List all Services in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Services"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list services", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServices(service.NewListServicesParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listService
	r.Register(tools.ToolDef{
		Name:        "listService",
		Description: "Get details about a specific Ziti service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Service Details"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.DetailService(service.NewDetailServiceParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// createService
	r.Register(tools.ToolDef{
		Name:        "createService",
		Description: "Create a new Ziti Service.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":               map[string]any{"type": "string", "description": "Name of the service"},
				"encryptionRequired": map[string]any{"type": "boolean", "description": "Whether end-to-end encryption is required", "default": true},
				"configs":            map[string]any{"type": "string", "description": "Comma-separated config IDs"},
				"roleAttributes":     map[string]any{"type": "string", "description": "Comma-separated role attributes"},
				"terminatorStrategy": map[string]any{"type": "string", "description": "Terminator strategy"},
				"maxIdleTimeMillis":  map[string]any{"type": "number", "description": "Max idle time in milliseconds"},
			},
			"required": []string{"name"},
		},
		Meta:        writeMeta(),
		Annotations: createAnnotations("Create Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		name, errResp, ok := RequireString(req.Parameters, "name")
		if !ok {
			return *errResp, nil
		}
		encReq := OptionalBool(req.Parameters, "encryptionRequired", true)
		configs := SplitCSV(OptionalString(req.Parameters, "configs"))
		roleAttrs := SplitCSV(OptionalString(req.Parameters, "roleAttributes"))
		terminatorStrategy := OptionalString(req.Parameters, "terminatorStrategy")
		maxIdle := OptionalInt64(req.Parameters, "maxIdleTimeMillis")

		body := &models.ServiceCreate{
			Name:               strPtr(name),
			EncryptionRequired: boolPtr(encReq),
		}
		if len(configs) > 0 {
			body.Configs = configs
		}
		if len(roleAttrs) > 0 {
			body.RoleAttributes = roleAttrs
		}
		if terminatorStrategy != "" {
			body.TerminatorStrategy = terminatorStrategy
		}
		if maxIdle != nil {
			body.MaxIdleTimeMillis = *maxIdle
		}

		return client.WithAuthenticatedClient(req, cfg, "create service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.CreateService(service.NewCreateServiceParams().WithService(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// deleteService
	r.Register(tools.ToolDef{
		Name:        "deleteService",
		Description: "Delete a Ziti Service.",
		InputSchema: idSchema(),
		Meta:        writeMeta(),
		Annotations: deleteAnnotations("Delete Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "delete service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.DeleteService(service.NewDeleteServiceParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// updateService
	r.Register(tools.ToolDef{
		Name:        "updateService",
		Description: "Update an existing Ziti Service.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":                 map[string]any{"type": "string", "description": "Service ID"},
				"name":              map[string]any{"type": "string", "description": "New name"},
				"encryptionRequired": map[string]any{"type": "boolean", "description": "Whether end-to-end encryption is required"},
				"configs":           map[string]any{"type": "string", "description": "Comma-separated config IDs"},
				"roleAttributes":    map[string]any{"type": "string", "description": "Comma-separated role attributes"},
				"terminatorStrategy": map[string]any{"type": "string", "description": "Terminator strategy"},
				"maxIdleTimeMillis": map[string]any{"type": "number", "description": "Max idle time in milliseconds"},
			},
			"required": []string{"id"},
		},
		Meta:        writeMeta(),
		Annotations: updateAnnotations("Update Service"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &models.ServicePatch{}
		if v := OptionalString(req.Parameters, "name"); v != "" {
			body.Name = v
		}
		if v, exists := req.Parameters["encryptionRequired"]; exists && v != nil {
			body.EncryptionRequired = v.(bool)
		}
		if v := OptionalString(req.Parameters, "configs"); v != "" {
			body.Configs = SplitCSV(v)
		}
		if v := OptionalString(req.Parameters, "roleAttributes"); v != "" {
			body.RoleAttributes = SplitCSV(v)
		}
		if v := OptionalString(req.Parameters, "terminatorStrategy"); v != "" {
			body.TerminatorStrategy = v
		}
		if v := OptionalInt64(req.Parameters, "maxIdleTimeMillis"); v != nil {
			body.MaxIdleTimeMillis = *v
		}

		return client.WithAuthenticatedClient(req, cfg, "update service", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.PatchService(service.NewPatchServiceParams().WithID(id).WithService(body), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceIdentities
	r.Register(tools.ToolDef{
		Name:        "listServiceIdentities",
		Description: "List all Identities that have access to a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Identities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service identities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceIdentities(service.NewListServiceIdentitiesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceEdgeRouters
	r.Register(tools.ToolDef{
		Name:        "listServiceEdgeRouters",
		Description: "List all Edge Routers accessible by a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Edge Routers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service edge routers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceEdgeRouters(service.NewListServiceEdgeRoutersParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceTerminators
	r.Register(tools.ToolDef{
		Name:        "listServiceTerminators",
		Description: "List all Terminators for a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Terminators"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service terminators", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceTerminators(service.NewListServiceTerminatorsParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceConfig
	r.Register(tools.ToolDef{
		Name:        "listServiceConfig",
		Description: "List all Configs associated with a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Configs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service configs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceConfig(service.NewListServiceConfigParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceServicePolicies
	r.Register(tools.ToolDef{
		Name:        "listServiceServicePolicies",
		Description: "List all Service Policies that apply to a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Service Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service service policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceServicePolicies(service.NewListServiceServicePoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceServiceEdgeRouterPolicies
	r.Register(tools.ToolDef{
		Name:        "listServiceServiceEdgeRouterPolicies",
		Description: "List all Service Edge Router Policies that apply to a specific Service",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Edge Router Policies"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "list service service edge router policies", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Service.ListServiceServiceEdgeRouterPolicies(service.NewListServiceServiceEdgeRouterPoliciesParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listServiceRoleAttributes
	r.Register(tools.ToolDef{
		Name:        "listServiceRoleAttributes",
		Description: "List all role attributes in use by Services in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Service Role Attributes"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list service role attributes", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.RoleAttributes.ListServiceRoleAttributes(nil, noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
