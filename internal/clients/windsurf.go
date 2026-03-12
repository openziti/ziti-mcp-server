package clients

import (
	"os"
	"path/filepath"
)

// WindsurfManager handles Windsurf configuration.
type WindsurfManager struct {
	BaseManager
}

func NewWindsurfManager() *WindsurfManager {
	return &WindsurfManager{
		BaseManager: BaseManager{
			clientType:  "windsurf",
			displayName: "Windsurf",
		},
	}
}

func (m *WindsurfManager) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	dir, _ := GetPlatformPath(PlatformPaths{
		Darwin: filepath.Join(home, ".codeium", "windsurf"),
		Win32:  filepath.Join("{APPDATA}", ".codeium", "windsurf"),
		Linux:  filepath.Join(home, ".codeium", "windsurf"),
	})
	return filepath.Join(dir, "mcp_config.json")
}

func (m *WindsurfManager) Configure(opts ClientOptions, binaryPath string) error {
	return ConfigureBase(&m.BaseManager, m.GetConfigPath(), opts, binaryPath)
}
