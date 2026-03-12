package clients

import (
	"os"
	"path/filepath"
)

// ClaudeManager handles Claude Desktop configuration.
type ClaudeManager struct {
	BaseManager
}

func NewClaudeManager() *ClaudeManager {
	return &ClaudeManager{
		BaseManager: BaseManager{
			clientType:   "claude",
			displayName:  "Claude Desktop",
			capabilities: []string{"tools"},
		},
	}
}

func (m *ClaudeManager) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	dir, _ := GetPlatformPath(PlatformPaths{
		Darwin: filepath.Join(home, "Library", "Application Support", "Claude"),
		Win32:  filepath.Join("{APPDATA}", "Claude"),
		Linux:  filepath.Join(home, ".config", "Claude"),
	})
	return filepath.Join(dir, "claude_desktop_config.json")
}

func (m *ClaudeManager) Configure(opts ClientOptions, binaryPath string) error {
	return ConfigureBase(&m.BaseManager, m.GetConfigPath(), opts, binaryPath)
}
