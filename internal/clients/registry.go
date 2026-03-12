package clients

import "fmt"

// Registry maps client type names to their managers.
var Registry = map[string]ClientManager{
	"claude":      NewClaudeManager(),
	"claude-code": NewClaudeCodeManager(),
	"cursor":      NewCursorManager(),
	"windsurf":    NewWindsurfManager(),
	"vscode":      NewVSCodeManager(),
	"warp":        NewWarpManager(),
}

// ValidClientTypes returns the list of supported client type names.
func ValidClientTypes() []string {
	return []string{"claude", "claude-code", "cursor", "windsurf", "vscode", "warp"}
}

// Get returns the client manager for the given type, or an error.
func Get(clientType string) (ClientManager, error) {
	mgr, ok := Registry[clientType]
	if !ok {
		return nil, fmt.Errorf("invalid client type: %s (available: claude, claude-code, cursor, windsurf, vscode, warp)", clientType)
	}
	return mgr, nil
}
