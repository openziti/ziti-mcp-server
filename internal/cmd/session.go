package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server-go/internal/terminal"
)

func newSessionCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "session",
		Short: "Display current authentication session information",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("retrieving session information")

			// Determine which profile to display
			target := profile
			if target == "" {
				target = credStore.ActiveProfile()
			}
			if target == "" {
				terminal.Output("")
				terminal.Warn("No active profile. Run %s to authenticate.",
					terminal.Cyan("ziti-mcp-server init"))
				return nil
			}

			data := credStore.ProfileData(target)
			if data == nil {
				terminal.Output("")
				terminal.Warn("Profile %q does not exist.", target)
				return nil
			}

			controllerHost := data["ziti_controller_host"]
			controllerCA := data["controller_ca"]

			terminal.Output("")
			terminal.Bold("Profile: %s", target)
			if target == credStore.ActiveProfile() {
				terminal.Info("(active)")
			}

			// Check UPDB mode
			if username := data["updb_username"]; username != "" && controllerHost != "" {
				terminal.Success("Active authentication session:")
				terminal.Output("")
				terminal.Bold("Auth mode: updb (username/password)")
				terminal.Bold("Ziti Controller Host: %s", controllerHost)
				terminal.Bold("Username: %s", terminal.MaskString(username, 8))
				if controllerCA != "" {
					terminal.Bold("Controller CA: configured")
				}
				terminal.Output("")
				terminal.Info("To use different credentials, run %s",
					terminal.Cyan("ziti-mcp-server logout --profile "+target))
				return nil
			}

			// Check identity mode
			if data["identity_cert"] != "" && controllerHost != "" {
				terminal.Success("Active authentication session:")
				terminal.Output("")
				terminal.Bold("Auth mode: identity (mTLS certificate)")
				terminal.Bold("Ziti Controller Host: %s", controllerHost)
				terminal.Output("")
				terminal.Info("To use different credentials, run %s",
					terminal.Cyan("ziti-mcp-server logout --profile "+target))
				return nil
			}

			// Token mode
			token := data["token"]
			domain := data["domain"]

			if token == "" || controllerHost == "" || domain == "" {
				terminal.Warn("No active authentication session found for profile %q.", target)
				terminal.Info("Run %s to authenticate.",
					terminal.Cyan("ziti-mcp-server init --profile "+target))
				return nil
			}

			terminal.Success("Active authentication session:")
			terminal.Output("")
			terminal.Bold("Ziti Controller Host: %s", controllerHost)
			terminal.Bold("Domain: %s", domain)

			if controllerCA != "" {
				terminal.Bold("Controller CA: configured")
			}

			if expiresAtStr := data["token_expires_at"]; expiresAtStr != "" {
				var expiresAt int64
				_, _ = fmt.Sscanf(expiresAtStr, "%d", &expiresAt)
				if expiresAt > 0 {
					now := time.Now().UnixMilli()
					expiresIn := expiresAt - now
					if expiresIn > 0 {
						hours := expiresIn / (1000 * 60 * 60)
						expTime := time.UnixMilli(expiresAt).Format(time.RFC1123)
						terminal.Bold("Token expires: in %d hours (%s)", hours, expTime)
					} else {
						expTime := time.UnixMilli(expiresAt).Format(time.RFC1123)
						terminal.Output(fmt.Sprintf("Token status: %s on %s", terminal.Red("Expired"), expTime))
					}
				}
			}

			terminal.Output("")
			terminal.Info("To use different credentials, run %s",
				terminal.Cyan("ziti-mcp-server logout --profile "+target))

			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile to display (defaults to active profile)")

	return cmd
}
