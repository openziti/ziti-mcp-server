package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server/internal/server"
)

func newRunCmd() *cobra.Command {
	var (
		toolPatterns string
		readOnly     bool
		profile      string
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the MCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			toolPats := parseToolPatterns(toolPatterns)

			// If a profile is specified, switch to it before starting
			if profile != "" {
				if credStore.HasProfile(profile) {
					if err := credStore.SetActiveProfile(profile); err != nil {
						return err
					}
				} else {
					if err := credStore.EnsureProfile(profile); err != nil {
						return err
					}
				}
				slog.Info("using profile", "profile", profile)
			}

			if readOnly {
				slog.Info("starting server in read-only mode")
			} else {
				slog.Info("starting server", "tools", toolPats)
			}

			return server.Start(credStore, registry, metaRegistry, server.Options{
				Tools:    toolPats,
				ReadOnly: readOnly,
			})
		},
	}

	cmd.Flags().StringVar(&toolPatterns, "tools", "*", "Comma-separated tool patterns")
	cmd.Flags().BoolVar(&readOnly, "read-only", false, "Only expose read-only tools")
	cmd.Flags().StringVar(&profile, "profile", "", "Active profile name to use")

	return cmd
}
