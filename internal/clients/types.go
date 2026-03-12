package clients

// ServerConfig is the MCP server configuration written to client config files.
type ServerConfig struct {
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Env          map[string]string `json:"env,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Type         string            `json:"type,omitempty"`
	DeferLoading bool              `json:"defer_loading,omitempty"`
}

// ClientConfig is the generic client configuration format.
type ClientConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
	Extra      map[string]any         `json:"-"` // preserved extra fields
}

// ClientOptions holds user-provided options for configuring a client.
type ClientOptions struct {
	Tools    []string
	ReadOnly bool
}

// ClientManager is the interface for all client configuration writers.
type ClientManager interface {
	DisplayName() string
	GetConfigPath() string
	Configure(opts ClientOptions, binaryPath string) error
}

// MCPServerName is the key used in client config files for our server.
const MCPServerName = "ziti"
