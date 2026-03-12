# Ziti MCP Server — Full Tool Reference

> **[← Back to README](../README.md)**

The Ziti MCP Server provides **201 Ziti API tools** plus **8 meta-tools** for managing your Ziti network through natural language. Tools are organized by resource type.

> **Tip:** Use `--read-only` or `--tools` patterns to expose only the tools you need. See [Security — Restricting Tool Access](../README.md#restricting-tool-access).

---

## Meta-Tools (Network Management)

These tools are always available regardless of `--tools` or `--read-only` filtering.

| Tool                     | Description                                                                                      |
| ------------------------ | ------------------------------------------------------------------------------------------------ |
| `loginUpdb`              | Connect using username/password authentication                                                   |
| `loginIdentity`          | Connect using a Ziti identity JSON (mTLS certificate)                                            |
| `loginClientCredentials` | Connect using OAuth2 client credentials                                                          |
| `loginDeviceAuth`        | Start OAuth2 device auth flow (returns verification URL)                                         |
| `completeLogin`          | Complete a pending device-auth login after browser approval                                      |
| `logout`                 | Disconnect from a Ziti network (clear profile credentials)                                       |
| `listNetworks`           | List all configured network profiles with connection status                                      |
| `selectNetwork`          | Switch the active network profile                                                                |

**Example prompts:**

- `Log into my Ziti network at ctrl.example.com:1280 with username admin`
- `Show me which networks I'm connected to`
- `Switch to the staging network`
- `Log out from prod`

---

## Identities

CRUD operations, relationship queries, and lifecycle management for Ziti Identities.

| Tool                                 | Description                                                                 |
| ------------------------------------ | --------------------------------------------------------------------------- |
| `listIdentities`                     | List all identities in the Ziti network                                     |
| `listIdentity`                       | Get details about a specific identity                                       |
| `createIdentity`                     | Create a new identity (name, admin, authPolicy, externalId, roleAttributes) |
| `deleteIdentity`                     | Delete an identity                                                          |
| `updateIdentity`                     | Update an existing identity                                                 |
| `listIdentityServices`               | List services accessible by an identity                                     |
| `listIdentityEdgeRouters`            | List edge routers accessible by an identity                                 |
| `listIdentityServicePolicies`        | List service policies that apply to an identity                             |
| `listIdentityEdgeRouterPolicies`     | List edge router policies that apply to an identity                         |
| `listIdentityServiceConfigs`         | List service configs associated with an identity                            |
| `getIdentityPolicyAdvice`            | Check whether an identity can dial/bind a service and get policy advice     |
| `listIdentityRoleAttributes`         | List all role attributes in use by identities                               |
| `disableIdentity`                    | Temporarily disable an identity for a specified duration                    |
| `enableIdentity`                     | Re-enable a previously disabled identity                                    |
| `getIdentityAuthenticators`          | List authenticators for an identity                                         |
| `getIdentityEnrollments`             | List enrollments for an identity                                            |
| `getIdentityFailedServiceRequests`   | List failed service requests for an identity                                |
| `getIdentityPostureData`             | Get posture data for an identity                                            |
| `removeIdentityMfa`                  | Remove MFA from an identity                                                 |
| `updateIdentityTracing`              | Update tracing configuration for an identity                                |
| `associateIdentityServiceConfigs`    | Associate service/config pairs with an identity                             |
| `disassociateIdentityServiceConfigs` | Remove service/config associations from an identity                         |

**Example prompts:**

- `Show me all Ziti identities`
- `Which identities have access to the Demo1 service?`
- `Create a new identity called 'demo-admin' and make it an admin`
- `Disable identity abc123 for 60 minutes`
- `What posture data does identity xyz have?`

---

## Services

CRUD operations and relationship queries for Ziti Services.

| Tool                                   | Description                                                                                                     |
| -------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `listServices`                         | List all services in the Ziti network                                                                           |
| `listService`                          | Get details about a specific service                                                                            |
| `createService`                        | Create a new service (name, encryptionRequired, configs, roleAttributes, terminatorStrategy, maxIdleTimeMillis) |
| `deleteService`                        | Delete a service                                                                                                |
| `updateService`                        | Update an existing service                                                                                      |
| `listServiceIdentities`                | List identities that have access to a service                                                                   |
| `listServiceEdgeRouters`               | List edge routers accessible by a service                                                                       |
| `listServiceTerminators`               | List terminators for a service                                                                                  |
| `listServiceConfig`                    | List configs associated with a service                                                                          |
| `listServiceServicePolicies`           | List service policies that apply to a service                                                                   |
| `listServiceServiceEdgeRouterPolicies` | List service edge router policies that apply to a service                                                       |
| `listServiceRoleAttributes`            | List all role attributes in use by services                                                                     |

**Example prompts:**

- `Show me all Ziti services`
- `Which identities can access the 'webapp' service?`
- `Create a new service called 'my-api' with encryption required`

---

## Edge Routers

CRUD operations, relationship queries, and re-enrollment for Ziti Edge Routers.

| Tool                                      | Description                                                                                     |
| ----------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `listEdgeRouters`                         | List all edge routers in the Ziti network                                                       |
| `listEdgeRouter`                          | Get details about a specific edge router                                                        |
| `createEdgeRouter`                        | Create a new edge router (name, isTunnelerEnabled, roleAttributes, cost, noTraversal, disabled) |
| `deleteEdgeRouter`                        | Delete an edge router                                                                           |
| `updateEdgeRouter`                        | Update an existing edge router                                                                  |
| `listEdgeRouterIdentities`                | List identities accessible by an edge router                                                    |
| `listEdgeRouterServices`                  | List services accessible by an edge router                                                      |
| `listEdgeRouterEdgeRouterPolicies`        | List edge router policies that apply to an edge router                                          |
| `listEdgeRouterServiceEdgeRouterPolicies` | List service edge router policies that apply to an edge router                                  |
| `listEdgeRouterRoleAttributes`            | List all role attributes in use by edge routers                                                 |
| `reEnrollEdgeRouter`                      | Re-enroll an edge router, generating new certificates                                           |

**Example prompts:**

- `List all edge routers and their status`
- `Which services are accessible through edge router xyz?`
- `Re-enroll edge router abc123`

---

## Edge Router Policies

CRUD operations and relationship queries for Edge Router Policies.

| Tool                              | Description                                                                      |
| --------------------------------- | -------------------------------------------------------------------------------- |
| `listEdgeRouterPolicies`          | List all edge router policies                                                    |
| `listEdgeRouterPolicy`            | Get details about a specific edge router policy                                  |
| `createEdgeRouterPolicy`          | Create a new edge router policy (name, semantic, edgeRouterRoles, identityRoles) |
| `deleteEdgeRouterPolicy`          | Delete an edge router policy                                                     |
| `updateEdgeRouterPolicy`          | Update an existing edge router policy                                            |
| `listEdgeRouterPolicyEdgeRouters` | List edge routers associated with a policy                                       |
| `listEdgeRouterPolicyIdentities`  | List identities associated with a policy                                         |

**Example prompts:**

- `List all edge router policies`
- `Which identities are covered by edge router policy xyz?`
- `Create an edge router policy that gives all identities access to all edge routers`

---

## Service Edge Router Policies

CRUD operations and relationship queries for Service Edge Router Policies.

| Tool                                     | Description                                                                             |
| ---------------------------------------- | --------------------------------------------------------------------------------------- |
| `listServiceEdgeRouterPolicies`          | List all service edge router policies                                                   |
| `listServiceEdgeRouterPolicy`            | Get details about a specific service edge router policy                                 |
| `createServiceEdgeRouterPolicy`          | Create a new service edge router policy (name, semantic, edgeRouterRoles, serviceRoles) |
| `deleteServiceEdgeRouterPolicy`          | Delete a service edge router policy                                                     |
| `updateServiceEdgeRouterPolicy`          | Update an existing service edge router policy                                           |
| `listServiceEdgeRouterPolicyEdgeRouters` | List edge routers associated with a policy                                              |
| `listServiceEdgeRouterPolicyServices`    | List services associated with a policy                                                  |

**Example prompts:**

- `Show me all service edge router policies`
- `Which edge routers are in the 'public-access' service edge router policy?`
- `Create a service edge router policy linking all services to all edge routers`

---

## Service Policies

CRUD operations and relationship queries for Service Policies (Dial/Bind).

| Tool                             | Description                                                                                                  |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `listServicePolicies`            | List all service policies                                                                                    |
| `listServicePolicy`              | Get details about a specific service policy                                                                  |
| `createServicePolicy`            | Create a new service policy (name, type Dial/Bind, semantic, identityRoles, serviceRoles, postureCheckRoles) |
| `deleteServicePolicy`            | Delete a service policy                                                                                      |
| `updateServicePolicy`            | Update an existing service policy                                                                            |
| `listServicePolicyIdentities`    | List identities associated with a policy                                                                     |
| `listServicePolicyServices`      | List services associated with a policy                                                                       |
| `listServicePolicyPostureChecks` | List posture checks associated with a policy                                                                 |

**Example prompts:**

- `Show me all Dial service policies`
- `Which identities are in the 'web-access' service policy?`
- `Create a Bind policy for the 'my-api' service`

---

## Configs

CRUD operations and relationship queries for Ziti Configs.

| Tool                 | Description                                    |
| -------------------- | ---------------------------------------------- |
| `listConfigs`        | List all configs                               |
| `listConfig`         | Get details about a specific config            |
| `createConfig`       | Create a new config (name, configTypeId, data) |
| `deleteConfig`       | Delete a config                                |
| `updateConfig`       | Update an existing config                      |
| `listConfigServices` | List services that use a specific config       |

**Example prompts:**

- `List all configs in the network`
- `Which services use config abc123?`
- `Create a new intercept.v1 config for my-service`

---

## Config Types

CRUD operations for Ziti Config Types.

| Tool                       | Description                                      |
| -------------------------- | ------------------------------------------------ |
| `listConfigTypes`          | List all config types                            |
| `listConfigType`           | Get details about a specific config type         |
| `createConfigType`         | Create a new config type (name, schema)          |
| `deleteConfigType`         | Delete a config type                             |
| `updateConfigType`         | Update an existing config type                   |
| `listConfigsForConfigType` | List all configs that use a specific config type |

**Example prompts:**

- `What config types are available?`
- `Show me all configs using the intercept.v1 config type`
- `Create a new config type with a custom JSON schema`

---

## Auth Policies

CRUD operations for Auth Policies (primary cert/extJwt/updb and secondary requirements).

| Tool               | Description                                                                               |
| ------------------ | ----------------------------------------------------------------------------------------- |
| `listAuthPolicies` | List all auth policies                                                                    |
| `listAuthPolicy`   | Get details about a specific auth policy                                                  |
| `createAuthPolicy` | Create a new auth policy (primary cert, extJwt, updb settings; secondary MFA requirement) |
| `deleteAuthPolicy` | Delete an auth policy                                                                     |
| `updateAuthPolicy` | Update an existing auth policy                                                            |

**Example prompts:**

- `List all auth policies`
- `Show me the details of the default auth policy`
- `Create an auth policy that requires MFA as a secondary factor`

---

## Authenticators

CRUD operations for Authenticators (updb/cert).

| Tool                  | Description                                                                  |
| --------------------- | ---------------------------------------------------------------------------- |
| `listAuthenticators`  | List all authenticators                                                      |
| `listAuthenticator`   | Get details about a specific authenticator                                   |
| `createAuthenticator` | Create a new authenticator (method, identityId, username, password, certPem) |
| `deleteAuthenticator` | Delete an authenticator                                                      |
| `updateAuthenticator` | Update an existing authenticator                                             |

**Example prompts:**

- `List all authenticators in the network`
- `Show me the authenticators for identity xyz`
- `Create a new updb authenticator for identity abc123`

---

## Certificate Authorities

CRUD operations, JWT retrieval, and verification for Certificate Authorities.

| Tool       | Description                                                                        |
| ---------- | ---------------------------------------------------------------------------------- |
| `listCas`  | List all certificate authorities                                                   |
| `listCa`   | Get details about a specific CA                                                    |
| `createCa` | Create a new CA (name, certPem, isAuthEnabled, enrollment settings, identityRoles) |
| `deleteCa` | Delete a CA                                                                        |
| `updateCa` | Update an existing CA                                                              |
| `getCaJwt` | Get the JWT for a CA (used for enrollment)                                         |
| `verifyCa` | Verify a CA with a signed PEM certificate                                          |

**Example prompts:**

- `List all certificate authorities`
- `Get the JWT for CA abc123`
- `Verify CA xyz with this PEM certificate`

---

## External JWT Signers

CRUD operations for External JWT Signers.

| Tool                      | Description                                                                            |
| ------------------------- | -------------------------------------------------------------------------------------- |
| `listExternalJwtSigners`  | List all external JWT signers                                                          |
| `listExternalJwtSigner`   | Get details about a specific external JWT signer                                       |
| `createExternalJwtSigner` | Create a new external JWT signer (name, issuer, audience, certPem, jwksEndpoint, etc.) |
| `deleteExternalJwtSigner` | Delete an external JWT signer                                                          |
| `updateExternalJwtSigner` | Update an existing external JWT signer                                                 |

**Example prompts:**

- `List all external JWT signers`
- `Show me the details of the Demo JWT signer`
- `Create a new external JWT signer for my IdP`

---

## Posture Checks

CRUD operations, type queries, and role attributes for Posture Checks.

| Tool                             | Description                                                                        |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `listPostureChecks`              | List all posture checks                                                            |
| `listPostureCheck`               | Get details about a specific posture check                                         |
| `createPostureCheck`             | Create a new posture check (name, typeId: OS/PROCESS/DOMAIN/MAC/MFA/PROCESS_MULTI) |
| `deletePostureCheck`             | Delete a posture check                                                             |
| `updatePostureCheck`             | Update an existing posture check                                                   |
| `listPostureCheckRoleAttributes` | List all role attributes in use by posture checks                                  |
| `listPostureCheckTypes`          | List all available posture check types                                             |
| `detailPostureCheckType`         | Get details about a specific posture check type                                    |

**Example prompts:**

- `List all posture checks`
- `What posture check types are available?`
- `Create a new MFA posture check called 'require-mfa'`

---

## Routers

CRUD operations for fabric Routers.

| Tool           | Description                                             |
| -------------- | ------------------------------------------------------- |
| `listRouters`  | List all routers                                        |
| `listRouter`   | Get details about a specific router                     |
| `createRouter` | Create a new router (name, cost, noTraversal, disabled) |
| `deleteRouter` | Delete a router                                         |
| `updateRouter` | Update an existing router                               |

**Example prompts:**

- `List all routers in the network`
- `Show me details for router xyz`
- `Create a new router with cost 100`

---

## Transit Routers

CRUD operations for Transit Routers.

| Tool                  | Description                                 |
| --------------------- | ------------------------------------------- |
| `listTransitRouters`  | List all transit routers                    |
| `listTransitRouter`   | Get details about a specific transit router |
| `createTransitRouter` | Create a new transit router                 |
| `deleteTransitRouter` | Delete a transit router                     |
| `updateTransitRouter` | Update an existing transit router           |

**Example prompts:**

- `List all transit routers`
- `Show me the details of transit router abc123`

---

## Terminators

CRUD operations for Terminators.

| Tool               | Description                                                                   |
| ------------------ | ----------------------------------------------------------------------------- |
| `listTerminators`  | List all terminators                                                          |
| `listTerminator`   | Get details about a specific terminator                                       |
| `createTerminator` | Create a new terminator (service, router, binding, address, cost, precedence) |
| `deleteTerminator` | Delete a terminator                                                           |
| `updateTerminator` | Update an existing terminator                                                 |

**Example prompts:**

- `List all terminators`
- `Show me terminators for the 'my-service' service`
- `Create a terminator binding my-service to router xyz`

---

## Enrollments

CRUD operations and refresh for Enrollments.

| Tool                | Description                                                             |
| ------------------- | ----------------------------------------------------------------------- |
| `listEnrollments`   | List all enrollments                                                    |
| `listEnrollment`    | Get details about a specific enrollment                                 |
| `createEnrollment`  | Create a new enrollment (identityId, method: ott/ottca/updb, expiresAt) |
| `deleteEnrollment`  | Delete an enrollment                                                    |
| `refreshEnrollment` | Refresh an expired enrollment with a new expiration time                |

**Example prompts:**

- `List all pending enrollments`
- `Create a new OTT enrollment for identity abc123`
- `Refresh the expired enrollment xyz with a new expiration date`

---

## Controller Settings

CRUD operations for Controller Settings (OIDC configuration).

| Tool                               | Description                                              |
| ---------------------------------- | -------------------------------------------------------- |
| `listControllerSettings`           | List all controller settings                             |
| `listControllerSetting`            | Get details about a specific controller setting          |
| `createControllerSetting`          | Create a new controller setting                          |
| `deleteControllerSetting`          | Delete a controller setting                              |
| `updateControllerSetting`          | Update an existing controller setting                    |
| `detailControllerSettingEffective` | Get the effective (merged) value of a controller setting |

**Example prompts:**

- `List all controller settings`
- `Show me the effective value of controller setting xyz`
- `Update the OIDC redirect URIs for the controller`

---

## Controllers & System Info

| Tool                         | Description                                       |
| ---------------------------- | ------------------------------------------------- |
| `listControllers`            | List all controllers in the Ziti network          |
| `listRoot`                   | Get controller version and root information       |
| `listEnumeratedCapabilities` | List all capabilities supported by the controller |
| `listSummary`                | Get a summary of entity counts across the network |

**Example prompts:**

- `What version is the Ziti controller running?`
- `Give me a summary of the network — how many identities, services, and routers exist?`
- `What capabilities does this controller support?`

---

## Identity Types

| Tool                 | Description                                |
| -------------------- | ------------------------------------------ |
| `listIdentityTypes`  | List all identity types available          |
| `detailIdentityType` | Get details about a specific identity type |

---

## Sessions

| Tool            | Description                          |
| --------------- | ------------------------------------ |
| `listSessions`  | List all sessions                    |
| `listSession`   | Get details about a specific session |
| `deleteSession` | Delete a session                     |

---

## API Sessions

| Tool               | Description                              |
| ------------------ | ---------------------------------------- |
| `listApiSessions`  | List all API sessions                    |
| `listApiSession`   | Get details about a specific API session |
| `deleteApiSession` | Delete an API session                    |

---

## Fabric Tools

The server also provides tools for managing the Ziti Fabric layer:

- **Fabric Routers** — CRUD (6 tools)
- **Fabric Services** — CRUD (6 tools)
- **Fabric Terminators** — CRUD (5 tools)
- **Fabric Circuits** — list, detail, delete (3 tools)
- **Fabric Links** — list, detail, delete, update (4 tools)
- **Fabric Cluster** — list members, add, remove, transfer leadership (4 tools)
- **Fabric Inspect** — inspect (1 tool)
- **Fabric Database** — create snapshot, check/get/fix integrity, snapshot with path (5 tools)
