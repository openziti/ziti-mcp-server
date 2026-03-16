package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/client/enrollment"
	"github.com/openziti/ziti-mcp-server/internal/gen/edge/client/informational"
	wellknown "github.com/openziti/ziti-mcp-server/internal/gen/edge/client/well_known"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerNetworkInfo(r *tools.Registry, s *store.Store) {
	// getVersion
	r.Register(tools.ToolDef{
		Name:        "getVersion",
		Description: "Get the version information for the Ziti controller",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Version"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get version", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Informational.ListVersion(informational.NewListVersionParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listNetworkJwts
	r.Register(tools.ToolDef{
		Name:        "listNetworkJwts",
		Description: "List network JWTs for enrollment",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Network JWTs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list network JWTs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Enrollment.ListNetworkJWTs(enrollment.NewListNetworkJWTsParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})

	// listWellKnownCas
	r.Register(tools.ToolDef{
		Name:        "listWellKnownCas",
		Description: "List well-known Certificate Authorities trusted by the network",
		InputSchema: emptySchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("List Well-Known CAs"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list well-known CAs", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.WellKnown.ListWellKnownCas(wellknown.NewListWellKnownCasParams())
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			},
		), nil
	})
}
