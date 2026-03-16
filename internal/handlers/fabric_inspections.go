package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	finspect "github.com/openziti/ziti-mcp-server/internal/gen/fabric/client/inspect"
	fmodels "github.com/openziti/ziti-mcp-server/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerFabricInspections(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "inspect", Description: "Inspect system values from the Ziti fabric, such as stack dumps, metrics, or capability information",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"appRegex":        map[string]any{"type": "string", "description": "Regex to match application names"},
				"requestedValues": map[string]any{"type": "string", "description": "Comma-separated values to inspect"},
			},
			"required": []string{"appRegex", "requestedValues"},
		},
		Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Inspect"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		appRegex, errResp, ok := RequireString(req.Parameters, "appRegex")
		if !ok {
			return *errResp, nil
		}
		reqValues, errResp, ok := RequireString(req.Parameters, "requestedValues")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.InspectRequest{
			AppRegex:        strPtr(appRegex),
			RequestedValues: SplitCSV(reqValues),
		}
		return client.WithAuthenticatedClient(req, cfg, "inspect", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Inspect.Inspect(finspect.NewInspectParams().WithRequest(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
