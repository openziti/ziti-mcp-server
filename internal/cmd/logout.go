package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server-go/internal/auth"
	"github.com/openziti/ziti-mcp-server-go/internal/terminal"
)

func newLogoutCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials for a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine which profile to clear
			target := profile
			if target == "" {
				target = credStore.ActiveProfile()
			}
			if target == "" {
				terminal.Warn("No active profile to log out.")
				return nil
			}

			slog.Debug("removing credentials", "profile", target)
			terminal.Info("Clearing authentication data for profile %q...", target)
			terminal.Output("")

			// Revoke refresh token before clearing credentials
			if clientID := credStore.GetForProfile(target, "idp_client_id"); clientID != "" {
				if err := auth.RevokeRefreshToken(credStore, clientID); err != nil {
					slog.Debug("refresh token revocation failed", "error", err)
				}
			}

			results := credStore.ClearProfile(target)

			successCount := 0
			for _, r := range results {
				if r.Success {
					successCount++
				}
			}

			if successCount > 0 {
				terminal.Success("Successfully removed %d credential(s) from profile %q.", successCount, target)
			} else {
				terminal.Warn("No OpenZiti MCP authentication data was found in profile %q.", target)
			}

			var failCount int
			for _, r := range results {
				if r.Error != nil {
					failCount++
				}
			}
			if failCount > 0 {
				terminal.Warn("Some credentials could not be removed (%d failures).", failCount)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile to log out (defaults to active profile)")

	return cmd
}
