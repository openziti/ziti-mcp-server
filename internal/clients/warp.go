package clients

import (
	"encoding/json"

	"github.com/openziti/ziti-mcp-server-go/internal/terminal"
)

// WarpManager handles Warp configuration (prints snippet for manual paste).
type WarpManager struct {
	BaseManager
}

func NewWarpManager() *WarpManager {
	return &WarpManager{
		BaseManager: BaseManager{
			clientType:  "warp",
			displayName: "Warp",
		},
	}
}

func (m *WarpManager) GetConfigPath() string {
	return "" // Not applicable
}

func (m *WarpManager) Configure(opts ClientOptions, binaryPath string) error {
	serverConfig := m.CreateServerConfig(opts, binaryPath)
	snippet := map[string]ServerConfig{MCPServerName: serverConfig}
	data, _ := json.MarshalIndent(snippet, "", "  ")

	terminal.Success("Copy the following JSON and paste it into Warp's MCP Settings (Settings > MCP Servers > +Add):")
	terminal.Output("")
	terminal.Output("%s", string(data))
	terminal.Output("")

	return nil
}
