package handlers

import (
	"net/http"

	"github.com/openziti/ziti-mcp-server/internal/client"
	fcluster "github.com/openziti/ziti-mcp-server/internal/gen/fabric/client/cluster"
	fmodels "github.com/openziti/ziti-mcp-server/internal/gen/fabric/models"
	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/tools"
)

func registerFabricCluster(r *tools.Registry, s *store.Store) {
	r.Register(tools.ToolDef{
		Name: "listClusterMembers", Description: "List all members of the Ziti controller cluster and their current status",
		InputSchema: emptySchema(), Meta: readOnlyMeta(), Annotations: readOnlyAnnotations("List Cluster Members"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		return client.WithAuthenticatedClient(req, cfg, "list cluster members", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Cluster.ClusterListMembers(fcluster.NewClusterListMembersParams())
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "addClusterMember", Description: "Add a new member to the Ziti controller cluster",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"address": map[string]any{"type": "string", "description": "Member address"},
				"isVoter": map[string]any{"type": "boolean", "description": "Whether the member is a voter", "default": true},
			},
			"required": []string{"address"},
		},
		Meta: writeMeta(), Annotations: createAnnotations("Add Cluster Member"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		addr, errResp, ok := RequireString(req.Parameters, "address")
		if !ok {
			return *errResp, nil
		}
		isVoter := OptionalBool(req.Parameters, "isVoter", true)
		body := &fmodels.ClusterMemberAdd{
			Address: strPtr(addr),
			IsVoter: boolPtr(isVoter),
		}
		return client.WithAuthenticatedClient(req, cfg, "add cluster member", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Cluster.ClusterMemberAdd(fcluster.NewClusterMemberAddParams().WithMember(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "removeClusterMember", Description: "Remove a member from the Ziti controller cluster",
		InputSchema: idSchema(), Meta: writeMeta(), Annotations: deleteAnnotations("Remove Cluster Member"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		id, errResp, ok := RequireString(req.Parameters, "id")
		if !ok {
			return *errResp, nil
		}
		body := &fmodels.ClusterMemberRemove{
			ID: strPtr(id),
		}
		return client.WithAuthenticatedClient(req, cfg, "remove cluster member", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Cluster.ClusterMemberRemove(fcluster.NewClusterMemberRemoveParams().WithMember(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})

	r.Register(tools.ToolDef{
		Name: "transferClusterLeadership", Description: "Transfer leadership to a different member of the Ziti controller cluster",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"newLeaderId": map[string]any{"type": "string", "description": "New leader member ID (optional, auto-select if empty)"},
			},
		},
		Meta: writeMeta(), Annotations: updateAnnotations("Transfer Cluster Leadership"),
	}, func(req tools.HandlerRequest, cfg tools.HandlerConfig) (tools.HandlerResponse, error) {
		body := &fmodels.ClusterTransferLeadership{
			NewLeaderID: OptionalString(req.Parameters, "newLeaderId"),
		}
		return client.WithAuthenticatedClient(req, cfg, "transfer cluster leadership", s,
			func(httpClient *http.Client, _ string) (any, error) {
				fc := NewFabricClient(httpClient, cfg.ZitiControllerHost)
				resp, err := fc.Cluster.ClusterTransferLeadership(fcluster.NewClusterTransferLeadershipParams().WithMember(body))
				if err != nil {
					return nil, err
				}
				return ToMap(resp.Payload)
			}), nil
	})
}
