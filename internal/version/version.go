package version

// Set via ldflags at build time:
//
//	go build -ldflags "-X github.com/openziti/ziti-mcp-server/internal/version.Version=1.0.0"
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)
