package clients

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
		Win32:  claudeWin32ConfigDir(),
		Linux:  filepath.Join(home, ".config", "Claude"),
	})
	return filepath.Join(dir, "claude_desktop_config.json")
}

// claudeWin32ConfigDir returns the Claude Desktop config directory on Windows.
// Claude installed from the Microsoft Store uses a virtualized filesystem under
// %LOCALAPPDATA%\Packages\Claude_<hash>\LocalCache\Roaming\Claude\ instead of
// the standard %APPDATA%\Claude\ path. We probe for the Store path first.
func claudeWin32ConfigDir() string {
	if runtime.GOOS != "windows" {
		return filepath.Join("{APPDATA}", "Claude")
	}
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		packagesDir := filepath.Join(localAppData, "Packages")
		entries, err := os.ReadDir(packagesDir)
		if err == nil {
			for _, e := range entries {
				if e.IsDir() && strings.HasPrefix(e.Name(), "Claude_") {
					candidate := filepath.Join(packagesDir, e.Name(), "LocalCache", "Roaming", "Claude")
					if _, err := os.Stat(candidate); err == nil {
						return candidate
					}
				}
			}
		}
	}
	return filepath.Join("{APPDATA}", "Claude")
}

func (m *ClaudeManager) Configure(opts ClientOptions, binaryPath string) error {
	return ConfigureBase(&m.BaseManager, m.GetConfigPath(), opts, binaryPath)
}
