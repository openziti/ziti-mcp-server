package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server-go/internal/clients"
)

func newInstallCmd() *cobra.Command {
	var (
		clientType   string
		toolPatterns string
		readOnly     bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Register the MCP server with an AI client",
		Long: `Add the Ziti MCP Server entry to an AI client's configuration file.
This does not authenticate — use the runtime login tools or 'init' for that.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("installing", "client", clientType)

			mgr, err := clients.Get(clientType)
			if err != nil {
				return err
			}

			binaryPath, _ := os.Executable()
			toolPats := parseToolPatterns(toolPatterns)

			return mgr.Configure(clients.ClientOptions{
				Tools:    toolPats,
				ReadOnly: readOnly,
			}, binaryPath)
		},
	}

	cmd.Flags().StringVar(&clientType, "client", "claude", "Client to configure (claude, claude-code, cursor, windsurf, vscode, warp)")
	cmd.Flags().StringVar(&toolPatterns, "tools", "*", "Comma-separated tool patterns")
	cmd.Flags().BoolVar(&readOnly, "read-only", false, "Only expose read-only tools")

	return cmd
}
