package clients

import (
	"os"
	"path/filepath"
	"strings"
)

// OpenCodeManager handles OpenCode configuration.
// OpenCode uses a different MCP config format: { "mcp": { "name": { "type": "local", "command": [...] } } }
type OpenCodeManager struct {
	BaseManager
}

func NewOpenCodeManager() *OpenCodeManager {
	return &OpenCodeManager{
		BaseManager: BaseManager{
			clientType:  "opencode",
			displayName: "OpenCode",
		},
	}
}

func (m *OpenCodeManager) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

func (m *OpenCodeManager) Configure(opts ClientOptions, binaryPath string) error {
	configPath := m.GetConfigPath()
	config := ReadConfig(configPath)

	// OpenCode uses "mcp" key instead of "mcpServers"
	mcp, ok := config["mcp"].(map[string]any)
	if !ok {
		mcp = make(map[string]any)
	}

	// Build command array: [binary, "run", "--tools", "...", ...]
	args := []string{"run", "--tools", strings.Join(opts.Tools, ",")}
	if opts.ReadOnly {
		args = append(args, "--read-only")
	}
	command := append([]string{binaryPath}, args...)

	serverEntry := map[string]any{
		"type":    "local",
		"command": command,
	}

	mcp[MCPServerName] = serverEntry
	config["mcp"] = mcp

	// Remove mcpServers key if it was defaulted in by ReadConfig
	if servers, ok := config["mcpServers"]; ok {
		if m, ok := servers.(map[string]any); ok && len(m) == 0 {
			delete(config, "mcpServers")
		}
	}

	return WriteConfig(configPath, config)
}
