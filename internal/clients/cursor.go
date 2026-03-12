package clients

import (
	"os"
	"path/filepath"
)

// CursorManager handles Cursor configuration.
type CursorManager struct {
	BaseManager
}

func NewCursorManager() *CursorManager {
	return &CursorManager{
		BaseManager: BaseManager{
			clientType:  "cursor",
			displayName: "Cursor",
		},
	}
}

func (m *CursorManager) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	dir, _ := GetPlatformPath(PlatformPaths{
		Darwin: filepath.Join(home, ".cursor"),
		Win32:  filepath.Join("{APPDATA}", ".cursor"),
		Linux:  filepath.Join(home, ".cursor"),
	})
	return filepath.Join(dir, "mcp.json")
}

func (m *CursorManager) Configure(opts ClientOptions, binaryPath string) error {
	return ConfigureBase(&m.BaseManager, m.GetConfigPath(), opts, binaryPath)
}
