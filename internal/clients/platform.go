package clients

import (
	"fmt"
	"os"
	"runtime"
)

// PlatformPaths holds per-OS path templates.
type PlatformPaths struct {
	Darwin string
	Win32  string
	Linux  string
}

// GetPlatformPath resolves the path for the current OS.
// On Windows, {APPDATA} is replaced with the APPDATA env var.
func GetPlatformPath(paths PlatformPaths) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return paths.Darwin, nil
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		// Simple string replacement for {APPDATA} placeholder
		result := paths.Win32
		if len(result) > 9 && result[:9] == "{APPDATA}" {
			result = appdata + result[9:]
		}
		return result, nil
	case "linux":
		return paths.Linux, nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
