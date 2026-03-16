package clients

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

// ClaudeCodeManager handles Claude Code configuration.
type ClaudeCodeManager struct {
	BaseManager
	scope string // "user" or "project"
}

func NewClaudeCodeManager() *ClaudeCodeManager {
	return &ClaudeCodeManager{
		BaseManager: BaseManager{
			clientType:   "claude-code",
			displayName:  "Claude Code",
			capabilities: []string{"tools"},
		},
		scope: "user",
	}
}

func (m *ClaudeCodeManager) GetConfigPath() string {
	if m.scope == "project" {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, ".mcp.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude.json")
}

func (m *ClaudeCodeManager) Configure(opts ClientOptions, binaryPath string) error {
	// Default to user scope for non-interactive Go binary
	configPath := m.GetConfigPath()
	config := ReadConfig(configPath)

	mcpServers, ok := config["mcpServers"].(map[string]any)
	if !ok {
		mcpServers = make(map[string]any)
	}

	serverConfig := m.CreateServerConfig(opts, binaryPath)
	serverConfig.Type = "stdio"
	serverConfig.DeferLoading = true

	scData, _ := json.Marshal(serverConfig)
	var scMap map[string]any
	_ = json.Unmarshal(scData, &scMap)

	mcpServers[MCPServerName] = scMap
	config["mcpServers"] = mcpServers

	if err := WriteConfig(configPath, config); err != nil {
		return err
	}

	slog.Debug("updated Claude Code config", "path", configPath)
	terminal.Success("OpenZiti MCP server configured for Claude Code at %s.",
		configPath)
	terminal.Info("%s in this project to apply changes.",
		terminal.Yellow("Start a new Claude Code session"))

	return nil
}
