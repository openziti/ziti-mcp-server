package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/client/controllers"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/client/informational"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerControllers(r *tools.Registry, s *store.Store) {
	// listControllers
	r.Register(tools.ToolDef{
		Name:        "listControllers",
		Description: "List all Controllers in the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Controllers"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list controllers", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				return fetchAllPages(func(limit, offset int64) (map[string]any, error) {
					resp, err := ec.Controllers.ListControllers(
						controllers.NewListControllersParams().WithLimit(&limit).WithOffset(&offset), noAuth)
					if err != nil {
						return nil, err
					}
					m, err := ToMap(resp.Payload)
					if err != nil {
						return nil, err
					}
					return m.(map[string]any), nil
				})
			},
		), nil
	})

	// listRoot
	r.Register(tools.ToolDef{
		Name:        "listRoot",
		Description: "Get controller version and root information",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Root Info"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list root", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Informational.ListRoot(informational.NewListRootParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listEnumeratedCapabilities
	r.Register(tools.ToolDef{
		Name:        "listEnumeratedCapabilities",
		Description: "List all capabilities supported by the Ziti controller",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Enumerated Capabilities"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list enumerated capabilities", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Informational.ListEnumeratedCapabilities(informational.NewListEnumeratedCapabilitiesParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listSummary
	r.Register(tools.ToolDef{
		Name:        "listSummary",
		Description: "Get a summary of entity counts across the Ziti network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Network Summary"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list summary", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Informational.ListSummary(informational.NewListSummaryParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
