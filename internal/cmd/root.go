package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server-go/internal/handlers"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
	"github.com/openziti/ziti-mcp-server-go/internal/version"
)

var (
	credStore    *store.Store
	registry     *tools.Registry
	metaRegistry *tools.Registry
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ziti-mcp-server",
		Short: "OpenZiti MCP Server",
		Long: `OpenZiti MCP Server — a Model Context Protocol server that gives AI assistants
controlled access to the OpenZiti Controller Management API.`,
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initLogging()
		},
	}

	root.AddCommand(
		newInitCmd(),
		newRunCmd(),
		newLogoutCmd(),
		newSessionCmd(),
	)

	return root
}

// Execute runs the root command.
func Execute() error {
	// Initialize credential store
	credStore = store.New("")
	if err := credStore.Load(); err != nil {
		slog.Warn("could not load credential store", "error", err)
	}

	// Initialize tool registry and register all MCP tools
	registry = tools.NewRegistry()
	handlers.RegisterAll(registry, credStore)

	// Initialize meta-tool registry
	metaRegistry = tools.NewRegistry()
	handlers.RegisterMeta(metaRegistry, credStore)

	return newRootCmd().Execute()
}

func initLogging() {
	debug := os.Getenv("OPENZITI_MCP_DEBUG") == "true" ||
		strings.Contains(os.Getenv("DEBUG"), "ziti-mcp-server")

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

// parseToolPatterns parses comma-separated tool patterns.
func parseToolPatterns(value string) []string {
	if value == "" {
		return []string{"*"}
	}
	parts := strings.Split(value, ",")
	var patterns []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			patterns = append(patterns, p)
		}
	}
	if len(patterns) == 0 {
		return []string{"*"}
	}
	return patterns
}
