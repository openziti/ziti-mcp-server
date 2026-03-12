package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	session "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client/session"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerSessionDetails(r *tools.Registry, s *store.Store) {
	// getSessionRoutePath
	r.Register(tools.ToolDef{
		Name:        "getSessionRoutePath",
		Description: "Get the route path for a specific session, showing the path through the network",
		InputSchema: idSchema(),
		Meta:        readOnlyMeta(),
		Annotations: readOnlyAnnotations("Get Session Route Path"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		return client.WithAuthenticatedClient(req, cfg, "get session route path", s,
			func(httpClient *http.Client, _ string) (any, error) {
				ec := NewEdgeClient(httpClient, cfg.ZitiControllerHost)
				resp, err := ec.Session.DetailSessionRoutePath(session.NewDetailSessionRoutePathParams().WithID(id), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			},
		), nil
	})
}
