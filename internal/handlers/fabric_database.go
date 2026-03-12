package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	fdatabase "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/client/database"
	fmodels "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

func registerFabricDatabase(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "checkDataIntegrity", Description: "Start a data integrity scan on the Ziti datastore. Only one scan can run at a time.",
		InputSchema: emptySchema(), Meta: writeMeta(), Annotations: &tools.ToolAnnotations{Title: "Check Data Integrity"},
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "check data integrity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Database.CheckDataIntegrity(fdatabase.NewCheckDataIntegrityParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "getDataIntegrityResults", Description: "Get the results from an in-progress or completed data integrity check",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("Get Data Integrity Results"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "get data integrity results", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Database.DataIntegrityResults(fdatabase.NewDataIntegrityResultsParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "fixDataIntegrity", Description: "Run a data integrity scan and attempt to fix any issues found. Only one scan can run at a time.",
		InputSchema: emptySchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Fix Data Integrity"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "fix data integrity", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Database.FixDataIntegrity(fdatabase.NewFixDataIntegrityParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createDatabaseSnapshot", Description: "Create a new database snapshot at the default location",
		InputSchema: emptySchema(), Meta: writeMeta(), Annotations: createAnnotations("Create Database Snapshot"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "create database snapshot", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Database.CreateDatabaseSnapshot(fdatabase.NewCreateDatabaseSnapshotParams(), noAuth)
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "createDatabaseSnapshotWithPath", Description: "Create a new database snapshot at a specified path",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string", "description": "File path for the snapshot"},
			},
			"required": []string{"path"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Create Database Snapshot With Path"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		path, errResp, ok := RequireString(req.Parameters, "path")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.DatabaseSnapshotCreate{
			Path: path,
		}
		return client.WithAuthenticatedClient(req, cfg, "create database snapshot with path", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Database.CreateDatabaseSnapshotWithPath(fdatabase.NewCreateDatabaseSnapshotWithPathParams().WithSnapshot(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
