package handlers

import (
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// RegisterAll registers all MCP tool handlers with the given registry.
func RegisterAll(r *tools.Registry, s *store.Store) {
	// Edge Management API tools
	registerIdentities(r, s)
	registerServices(r, s)
	registerEdgeRouters(r, s)
	registerServicePolicies(r, s)
	registerEdgeRouterPolicies(r, s)
	registerServiceEdgeRouterPolicies(r, s)
	registerConfigs(r, s)
	registerConfigTypes(r, s)
	registerCertificateAuthorities(r, s)
	registerAuthPolicies(r, s)
	registerAuthenticators(r, s)
	registerAuthenticatorMgmt(r, s)
	registerTerminators(r, s)
	registerPostureChecks(r, s)
	registerExternalJwtSigners(r, s)
	registerEnrollments(r, s)
	registerSessions(r, s)
	registerAPISessions(r, s)
	registerCurrentIdentity(r, s)
	registerMFA(r, s)
	registerControllerSettings(r, s)
	registerControllers(r, s)
	registerRouters(r, s)
	registerTransitRouters(r, s)
	registerNetworkInfo(r, s)
	registerSessionDetails(r, s)

	// Fabric Management API tools (includes database tools)
	registerFabricRouters(r, s)
	registerFabricServices(r, s)
	registerFabricTerminators(r, s)
	registerFabricCircuits(r, s)
	registerFabricLinks(r, s)
	registerFabricCluster(r, s)
	registerFabricInspections(r, s)
	registerFabricDatabase(r, s)
}
