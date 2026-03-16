package clients

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

// VSCodeManager handles VS Code configuration.
type VSCodeManager struct {
	BaseManager
}

func NewVSCodeManager() *VSCodeManager {
	return &VSCodeManager{
		BaseManager: BaseManager{
			clientType:  "vscode",
			displayName: "VS Code",
		},
	}
}

func (m *VSCodeManager) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	dir, _ := GetPlatformPath(PlatformPaths{
		Darwin: filepath.Join(home, "Library", "Application Support", "Code", "User"),
		Win32:  filepath.Join("{APPDATA}", "Code", "User"),
		Linux:  filepath.Join(home, ".config", "Code", "User"),
	})
	return filepath.Join(dir, "mcp.json")
}

func (m *VSCodeManager) Configure(opts ClientOptions, binaryPath string) error {
	configPath := m.GetConfigPath()
	config := ReadConfig(configPath)

	servers, ok := config["servers"].(map[string]any)
	if !ok {
		servers = make(map[string]any)
	}

	serverConfig := m.CreateServerConfig(opts, binaryPath)
	scData, _ := json.Marshal(serverConfig)
	var scMap map[string]any
	_ = json.Unmarshal(scData, &scMap)

	servers[MCPServerName] = scMap
	config["servers"] = servers

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := WriteConfig(configPath, config); err != nil {
		return err
	}

	slog.Debug("updated VS Code config", "path", configPath)
	terminal.Success("OpenZiti MCP server configured globally for %s.", m.displayName)
	terminal.Info("%s to apply changes.", terminal.Yellow("Restart VS Code"))

	return nil
}
