<p align="center">
  <img src="assets/readme-banner.png" alt="OpenZiti MCP Server Banner">
</p>
<div align="center">

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://go.dev/)

</div>

<div align="center">

🚀 [Getting Started](#-getting-started) • 🕸️ [Architecture](#%EF%B8%8F-architecture) • 🔐 [Authentication](#-authentication) • 🛠️ [Supported Tools](#%EF%B8%8F-supported-tools) • 🩺 [Troubleshooting](#-troubleshooting) • 📋 [Debug Logs](#-debug-logs) • 👨‍💻 [Development](#-development) • 🔒 [Security](#-security)

</div>

The Ziti MCP Server is sponsored by [NetFoundry](https://netfoundry.io) as part of its portfolio of solutions
for secure workloads and agentic computing.
NetFoundry is the creator of [OpenZiti](https://netfoundry.io/docs/openziti/)
and [zrok](https://netfoundry.io/docs/zrok/getting-started).

[MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) is an open protocol introduced by Anthropic that standardizes how large language models communicate with external tools, resources or remote services.

The Ziti MCP Server integrates with LLMs and AI agents, allowing you to perform various Ziti network management operations using natural language. For instance, you could simply ask Claude Desktop to perform Ziti management operations:

- > List which identities exist
- > Tell me if there are any exposures in the network
- > Do you see potential misconfigurations?
- > Which identities have access to the Demo1 service?
- > Create a new Ziti identity named "Demo" and get its ID
- > Log into my prod Ziti network using UPDB
- > Switch to the staging network
- > etc.

<br/>

## 🚀 Getting Started

**Prerequisites:**

- [Claude Desktop](https://claude.ai/download) or any other [MCP Client](https://modelcontextprotocol.io/clients)
- [OpenZiti](https://openziti.io/) network
- [Go 1.24+](https://go.dev/dl/) (only if building from source)

<br/>

### Install

#### Download a pre-built binary

Pre-built binaries are available for macOS, Linux, and Windows (amd64 and arm64) on the [releases page](https://github.com/openziti/ziti-mcp-server-go/releases).

**macOS (Apple Silicon)**

```bash
curl -sL https://github.com/openziti/ziti-mcp-server-go/releases/latest/download/ziti-mcp-server_darwin_arm64.tar.gz | tar xz
sudo mv ziti-mcp-server /usr/local/bin/
```

**macOS (Intel)**

```bash
curl -sL https://github.com/openziti/ziti-mcp-server-go/releases/latest/download/ziti-mcp-server_darwin_amd64.tar.gz | tar xz
sudo mv ziti-mcp-server /usr/local/bin/
```

**Linux (amd64)**

```bash
curl -sL https://github.com/openziti/ziti-mcp-server-go/releases/latest/download/ziti-mcp-server_linux_amd64.tar.gz | tar xz
sudo mv ziti-mcp-server /usr/local/bin/
```

**Linux (arm64)**

```bash
curl -sL https://github.com/openziti/ziti-mcp-server-go/releases/latest/download/ziti-mcp-server_linux_arm64.tar.gz | tar xz
sudo mv ziti-mcp-server /usr/local/bin/
```

**Windows**

Download the appropriate `.zip` from the [releases page](https://github.com/openziti/ziti-mcp-server-go/releases) and add the extracted `ziti-mcp-server.exe` to your PATH.

#### Build from source

```bash
go install github.com/openziti/ziti-mcp-server-go/cmd/ziti-mcp-server@latest
```

### Quick Start (Disconnected Mode)

The server can start with **no prior configuration**. The AI agent can log into networks at runtime using the built-in login tools:

```bash
ziti-mcp-server run
```

Then in Claude Desktop, simply ask:

> Log into my Ziti network at 192.168.1.100:1280 with username admin and password admin

The server exposes `loginUpdb`, `loginIdentity`, `loginClientCredentials`, and `loginDeviceAuth` tools that the AI agent can call directly.

### Pre-configured Setup (CLI Init)

For non-interactive or automated setups, use `init` to pre-configure credentials and register with your MCP client:

**UPDB Mode (Username/Password)**

```bash
ziti-mcp-server init \
  --auth-mode updb \
  --ziti-controller-host <your-controller-host> \
  --username <username> \
  --password <password> \
  --profile prod
```

**Device Auth Mode (Interactive Login)**

```bash
ziti-mcp-server init \
  --auth-mode device-auth \
  --ziti-controller-host <your-controller-host> \
  --idp-domain <your-idp-domain> \
  --idp-client-id <your-client-id> \
  --idp-audience <your-audience> \
  --profile prod
```

**Client Credentials Mode (Service Account)**

```bash
ziti-mcp-server init \
  --auth-mode client-credentials \
  --ziti-controller-host <your-controller-host> \
  --idp-domain <your-idp-domain> \
  --idp-client-id <your-client-id> \
  --idp-client-secret <your-client-secret> \
  --profile prod
```

**Identity File Mode (mTLS Certificate)**

```bash
ziti-mcp-server init \
  --auth-mode identity \
  --identity-file <path-to-identity.json> \
  --profile prod
```

No IdP configuration is needed — authentication uses the client certificate from the Ziti identity file.

**With read-only tools**

```bash
ziti-mcp-server init \
  --auth-mode updb \
  --ziti-controller-host <host> \
  --username <user> \
  --password <pass> \
  --read-only
```

**Client selection**

Use `--client` to auto-configure a specific MCP client (default: `claude`):

```bash
ziti-mcp-server init --client windsurf ...
ziti-mcp-server init --client cursor ...
ziti-mcp-server init --client vscode ...
ziti-mcp-server init --client claude-code ...
ziti-mcp-server init --client warp ...
```

### Multi-Profile Support

The server supports multiple named network profiles, allowing you to manage several Ziti networks simultaneously:

```bash
# Pre-configure two profiles
ziti-mcp-server init --auth-mode updb --profile prod ...
ziti-mcp-server init --auth-mode updb --profile staging ...

# Start with a specific profile active
ziti-mcp-server run --profile prod
```

At runtime, the AI agent can:
- **Log into additional networks** using `loginUpdb`, `loginIdentity`, etc.
- **List all networks** using `listNetworks`
- **Switch between networks** using `selectNetwork`
- **Log out** from a network using `logout`

Credentials are stored in `~/.config/ziti-mcp-server/config.json`.

### MCP Client Configuration

**Other MCP Clients**

To use Ziti MCP Server with any MCP Client, add this configuration and restart:

```json
{
  "mcpServers": {
    "ziti": {
      "command": "/path/to/ziti-mcp-server",
      "args": ["run"],
      "capabilities": ["tools"],
      "env": {
        "OPENZITI_MCP_DEBUG": "true"
      }
    }
  }
}
```

You can add `--tools '<pattern>'` and/or `--read-only` to the args array to control which tools are available. See [Security Best Practices](#-security-best-practices-for-tool-access) for recommended patterns.

### Verify your integration

Restart your MCP Client (Claude Desktop, Windsurf, Cursor, Warp, etc.) and ask it to help you manage your Ziti network.

## 🕸️ Architecture

The Ziti MCP Server implements the Model Context Protocol, allowing clients (like Claude) to:

1. Request a list of available Ziti tools
2. Call specific tools with parameters
3. Receive structured responses from the Ziti Management API

The server handles authentication, request validation, and secure communication with the Ziti Management API.

<div align="center">
  <img src="assets/arch.jpg" alt="Ziti MCP Server" width="800">
</div>

> [!NOTE]
> The server operates as a local process that connects to Claude Desktop, enabling secure communication without exposing your Ziti credentials.

## 🔐 Authentication

The Ziti MCP Server uses the Ziti Management API and requires authentication to access your Ziti network.

### Authentication Modes

The server supports four authentication modes:

#### UPDB Mode (Username/Password)

Use this mode for direct username/password authentication against the Ziti controller:

```bash
ziti-mcp-server init \
  --auth-mode updb \
  --ziti-controller-host <your-controller-host> \
  --username <username> \
  --password <password>
```

Or at runtime via the AI agent using the `loginUpdb` tool.

#### Device Auth Mode (Interactive Login)

Use this mode for interactive browser-based login. Recommended for development and user-facing scenarios:

```bash
ziti-mcp-server init \
  --auth-mode device-auth \
  --ziti-controller-host <your-controller-host> \
  --idp-domain <your-idp-domain> \
  --idp-client-id <your-client-id> \
  --idp-audience <your-audience>
```

Or at runtime via the `loginDeviceAuth` tool (returns a verification URL for the user, then `completeLogin` to finish).

#### Client Credentials Mode (Service Account)

Use this mode for service accounts and automation. Recommended for production environments:

> [!NOTE]
> Keep the token lifetime as minimal as possible to reduce security risks. [See more](https://auth0.com/docs/secure/tokens/access-tokens/update-access-token-lifetime)

```bash
ziti-mcp-server init \
  --auth-mode client-credentials \
  --ziti-controller-host <your-controller-host> \
  --idp-domain <your-idp-domain> \
  --idp-client-id <your-client-id> \
  --idp-client-secret <your-client-secret>
```

#### Identity File Mode (mTLS Certificate)

Use this mode for certificate-based authentication with a Ziti identity JSON file. No IdP configuration is needed:

```bash
ziti-mcp-server init \
  --auth-mode identity \
  --identity-file <path-to-identity.json>
```

The identity file is a standard Ziti identity JSON file containing `ztAPI`, `id.cert`, `id.key`, and `id.ca` fields. The certificate material is extracted and stored in the config file. The identity file may be deleted after a successful `init` (for additional security, if desired).

> [!IMPORTANT]
>
> When using CLI `init`, it needs to be run whenever:
>
> - You're setting up the MCP Server for the first time
> - You've logged out from a previous session
> - You want to switch to a different Ziti network
> - Your token has expired
>
> Alternatively, use the runtime login tools (`loginUpdb`, etc.) to authenticate without restarting the server.

### Session Management

To see information about your current authentication session:

```bash
ziti-mcp-server session
ziti-mcp-server session --profile prod
```

### Logging Out

```bash
ziti-mcp-server logout
ziti-mcp-server logout --profile prod
```

Or at runtime via the AI agent using the `logout` tool.

### Authentication Flow

The Ziti MCP server supports multiple authentication flows:

- **UPDB (username/password)** for direct authentication against the Ziti controller
- **OAuth 2.0 device authorization flow** for interactive browser-based login with an IdP
- **Client credentials flow** for service accounts and automation
- **Identity file (mTLS)** for certificate-based authentication using a Ziti identity JSON file

Credentials are stored in `~/.config/ziti-mcp-server/config.json` with 0600 permissions.


## 🛠️ Supported MCP Tools

The Ziti MCP Server provides **201 Ziti API tools** plus **8 meta-tools** for managing your Ziti network through natural language. Tools are organized by resource type.

> **Tip:** Use `--read-only` or `--tools` patterns to expose only the tools you need. See [Security Best Practices](#-security-best-practices-for-tool-access).

### Meta-Tools (Network Management)

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

### Identities

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

### Services

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

### Edge Routers

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

### Edge Router Policies

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

### Service Edge Router Policies

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

### Service Policies

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

### Configs

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

### Config Types

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

### Auth Policies

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

### Authenticators

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

### Certificate Authorities

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

### External JWT Signers

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

### Posture Checks

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

### Routers

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

### Transit Routers

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

### Terminators

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

### Enrollments

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

### Controller Settings

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

### Controllers & System Info

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

### Identity Types

| Tool                 | Description                                |
| -------------------- | ------------------------------------------ |
| `listIdentityTypes`  | List all identity types available          |
| `detailIdentityType` | Get details about a specific identity type |

### Sessions

| Tool            | Description                          |
| --------------- | ------------------------------------ |
| `listSessions`  | List all sessions                    |
| `listSession`   | Get details about a specific session |
| `deleteSession` | Delete a session                     |

### API Sessions

| Tool               | Description                              |
| ------------------ | ---------------------------------------- |
| `listApiSessions`  | List all API sessions                    |
| `listApiSession`   | Get details about a specific API session |
| `deleteApiSession` | Delete an API session                    |

### Fabric Tools

The server also provides tools for managing the Ziti Fabric layer:

- **Fabric Routers** — CRUD (6 tools)
- **Fabric Services** — CRUD (6 tools)
- **Fabric Terminators** — CRUD (5 tools)
- **Fabric Circuits** — list, detail, delete (3 tools)
- **Fabric Links** — list, detail, delete, update (4 tools)
- **Fabric Cluster** — list members, add, remove, transfer leadership (4 tools)
- **Fabric Inspect** — inspect (1 tool)
- **Fabric Database** — create snapshot, check/get/fix integrity, snapshot with path (5 tools)

### 🔒 Security Best Practices for Tool Access

When configuring the Ziti MCP Server, it's important to follow security best practices by limiting tool access based on your specific needs:

```bash
# Enable only read-only operations
ziti-mcp-server run --read-only

# Alternative way to enable only read-only operations
ziti-mcp-server run --tools 'list*,get*'

# Limit to just identity-related tools
ziti-mcp-server run --tools '*Identit*'

# Limit to read-only identity-related tools
ziti-mcp-server run --tools '*Identit*' --read-only

# Run the server with all tools enabled
ziti-mcp-server run --tools '*'
```

> [!IMPORTANT]
> When both `--read-only` and `--tools` flags are used together, the `--read-only` flag takes priority for security. Meta-tools (login, logout, listNetworks, selectNetwork) are always available regardless of filtering.

This approach offers several important benefits:

1. **Enhanced Security**: Limiting available tools reduces the potential attack surface.
2. **Better Performance**: Fewer tools means less context window usage for tool reasoning.
3. **Resource-Based Access Control**: Configure different instances with different tool sets.
4. **Simplified Auditing**: Easier to track which operations were performed.

### 🧪 Security Scanning

We recommend regularly scanning this server with community tools built to surface protocol-level risks:

- **[mcpscan.ai](https://mcpscan.ai)** — Web-based scanner for MCP endpoints
- **[mcp-scan](https://github.com/invariantlabs-ai/mcp-scan)** — CLI tool for evaluating server behavior

## 🩺 Troubleshooting

Start troubleshooting by exploring all available commands and options:

```bash
ziti-mcp-server help
```

### 🚨 Common Issues

1. **Authentication Failures**
   - Ensure you have the correct permissions in your Ziti network
   - Try re-initializing with `ziti-mcp-server init --auth-mode <mode> ...`
   - Or use the runtime login tools to re-authenticate

2. **TLS Certificate Errors**
   - The server auto-fetches the controller's CA on login via the EST `/cacerts` endpoint
   - If the CA fetch fails, add the controller CA to your system trust store
   - Or re-login to trigger a fresh CA fetch

3. **Client Can't Connect to the Server**
   - Restart your MCP client after configuration changes
   - Check that the binary path in the client config is correct

4. **Invalid Configuration Error**
   - This typically happens when no profile is active or credentials are missing
   - Use `listNetworks` to check profile status
   - Use a login tool or `ziti-mcp-server init` to authenticate

> [!TIP]
> Most connection issues can be resolved by restarting both the server and your MCP client.

## 📋 Debug logs

Enable debug mode to view detailed logs:

```sh
export OPENZITI_MCP_DEBUG=true
```

Get detailed MCP Client logs from Claude Desktop:

```sh
# Follow logs in real-time
tail -n 20 -F ~/Library/Logs/Claude/mcp*.log
```

## 👨‍💻 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/openziti/ziti-mcp-server-go.git
cd ziti-mcp-server-go

# Build
go build ./cmd/ziti-mcp-server

# Run
./ziti-mcp-server run
```

### Regenerating API Clients

The Ziti API clients in `internal/gen/` are generated from OpenAPI specs using go-swagger:

```bash
make generate
```

## 🔒 Security

The Ziti MCP Server prioritizes security:

- Credentials are stored in `~/.config/ziti-mcp-server/config.json` with 0600 permissions
- The config file is never world-readable
- Authentication supports OAuth 2.0 device authorization, client credentials, mTLS certificates, and UPDB
- Tool access can be restricted with `--read-only` and `--tools` glob patterns
- Meta-tools allow runtime login/logout without restarting the server
- Easy credential removal via `logout` command or tool

> [!IMPORTANT]
> For security best practices, always log out when you're done with a session or switching between networks.

> [!CAUTION]
> Always review the permissions requested during the authentication process to ensure they align with your security requirements.

### Reporting Issues

To provide feedback or report a bug, please [raise an issue on our issue tracker](https://github.com/openziti/ziti-mcp-server-go/issues).
