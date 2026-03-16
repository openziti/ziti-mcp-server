package clients

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

// BaseManager provides common functionality for client config writers.
type BaseManager struct {
	clientType   string
	displayName  string
	capabilities []string
}

func (b *BaseManager) DisplayName() string { return b.displayName }

// CreateServerConfig builds the MCP server config entry for a client config file.
func (b *BaseManager) CreateServerConfig(opts ClientOptions, binaryPath string) ServerConfig {
	args := []string{"run", "--tools", strings.Join(opts.Tools, ",")}
	if opts.ReadOnly {
		args = append(args, "--read-only")
	}

	cfg := ServerConfig{
		Command: binaryPath,
		Args:    args,
		Env: map[string]string{
			"OPENZITI_MCP_DEBUG": "true",
		},
	}

	if len(b.capabilities) > 0 {
		cfg.Capabilities = b.capabilities
	}

	return cfg
}

// ReadConfig reads a JSON config file, preserving unknown fields.
func ReadConfig(configPath string) map[string]any {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("could not read config", "path", configPath, "error", err)
		}
		return map[string]any{"mcpServers": map[string]any{}}
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		slog.Warn("could not parse config", "path", configPath, "error", err)
		return map[string]any{"mcpServers": map[string]any{}}
	}

	return config
}

// WriteConfig writes a JSON config file with 2-space indent.
func WriteConfig(configPath string, config map[string]any) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// ConfigureBase is the standard config flow: read, update mcpServers, write.
func ConfigureBase(mgr *BaseManager, configPath string, opts ClientOptions, binaryPath string) error {
	config := ReadConfig(configPath)

	mcpServers, ok := config["mcpServers"].(map[string]any)
	if !ok {
		mcpServers = make(map[string]any)
	}

	serverConfig := mgr.CreateServerConfig(opts, binaryPath)

	// Convert to map for JSON serialization
	scData, _ := json.Marshal(serverConfig)
	var scMap map[string]any
	_ = json.Unmarshal(scData, &scMap)

	mcpServers[MCPServerName] = scMap
	config["mcpServers"] = mcpServers

	if err := WriteConfig(configPath, config); err != nil {
		return err
	}

	slog.Debug("updated config", "client", mgr.displayName, "path", configPath)
	terminal.Success("OpenZiti MCP server configured. %s to apply changes.",
		terminal.Yellow("Restart "+mgr.displayName))

	return nil
}
