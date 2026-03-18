package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/openziti/ziti-mcp-server/internal/auth"
	"github.com/openziti/ziti-mcp-server/internal/ca"
	"github.com/openziti/ziti-mcp-server/internal/clients"
	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

func newInitCmd() *cobra.Command {
	var (
		authMode           string
		clientType         string
		zitiControllerHost string
		idpDomain          string
		idpClientID        string
		idpClientSecret    string
		idpAudience        string
		identityFile       string
		username           string
		password           string
		toolPatterns       string
		readOnly           bool
		profile            string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the server (authenticate and configure)",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("initializing", "auth-mode", authMode, "client", clientType, "profile", profile)

			toolPats := parseToolPatterns(toolPatterns)

			// Validate auth mode parameters
			switch authMode {
			case "identity":
				if identityFile == "" {
					return fmt.Errorf("--identity-file is required for identity auth mode")
				}
			case "updb":
				if zitiControllerHost == "" {
					return fmt.Errorf("--ziti-controller-host is required for updb auth mode")
				}
				if username == "" {
					return fmt.Errorf("--username is required for updb auth mode")
				}
				if password == "" {
					return fmt.Errorf("--password is required for updb auth mode")
				}
			case "client-credentials":
				if zitiControllerHost == "" {
					return fmt.Errorf("--ziti-controller-host is required")
				}
				if idpDomain == "" {
					return fmt.Errorf("--idp-domain is required")
				}
				if idpClientID == "" {
					return fmt.Errorf("--idp-client-id is required")
				}
				if idpClientSecret == "" {
					return fmt.Errorf("--idp-client-secret is required for client-credentials auth mode")
				}
			case "device-auth":
				if zitiControllerHost == "" {
					return fmt.Errorf("--ziti-controller-host is required")
				}
				if idpDomain == "" {
					return fmt.Errorf("--idp-domain is required")
				}
				if idpClientID == "" {
					return fmt.Errorf("--idp-client-id is required")
				}
				if idpAudience == "" {
					return fmt.Errorf("--idp-audience is required for device-auth auth mode")
				}
			default:
				return fmt.Errorf("invalid --auth-mode: %q (must be device-auth, client-credentials, identity, or updb)", authMode)
			}

			// Create/switch to the profile, then clear it
			if err := credStore.EnsureProfile(profile); err != nil {
				return fmt.Errorf("creating profile: %w", err)
			}
			credStore.ClearProfile(profile)

			// For non-identity modes, fetch the controller CA
			if authMode != "identity" {
				host := zitiControllerHost
				caPem, err := ca.FetchControllerCA(host)
				if err != nil {
					slog.Debug("controller CA fetch failed (may use public cert)", "error", err)
				} else if caPem != "" {
					if err := credStore.SetControllerCA(caPem); err != nil {
						return fmt.Errorf("storing controller CA: %w", err)
					}
					slog.Debug("controller CA certificate stored")
				}
			}

			// Run auth flow
			switch authMode {
			case "identity":
				if err := auth.RequestIdentityFileAuthorization(credStore, identityFile); err != nil {
					terminal.Error("Failed to load identity file.")
					return err
				}

			case "updb":
				if err := auth.StoreUPDBCredentials(credStore, zitiControllerHost, username, password); err != nil {
					return err
				}

			case "client-credentials":
				if err := auth.RequestClientCredentialsAuthorization(credStore, auth.ClientCredentialsConfig{
					ZitiControllerHost: zitiControllerHost,
					IDPDomain:          idpDomain,
					IDPClientID:        idpClientID,
					IDPClientSecret:    idpClientSecret,
					Audience:           idpAudience,
				}); err != nil {
					terminal.Error("Failed to authenticate with client credentials.")
					return err
				}
				if err := credStore.SetIDPClientID(idpClientID); err != nil {
					return err
				}

			case "device-auth":
				if err := auth.RequestDeviceAuthorization(credStore, auth.DeviceAuthConfig{
					IDPDomain: idpDomain,
					ClientID:  idpClientID,
					Audience:  idpAudience,
				}); err != nil {
					terminal.Error("Failed to authenticate with device auth.")
					return err
				}
				if err := credStore.SetControllerHost(zitiControllerHost); err != nil {
					return err
				}
				if err := credStore.SetIDPClientID(idpClientID); err != nil {
					return err
				}
			}

			// Configure client
			mgr, err := clients.Get(clientType)
			if err != nil {
				return err
			}

			binaryPath, _ := os.Executable()
			return mgr.Configure(clients.ClientOptions{
				Tools:    toolPats,
				ReadOnly: readOnly,
			}, binaryPath)
		},
	}

	cmd.Flags().StringVar(&authMode, "auth-mode", "", "Authentication mode: device-auth, client-credentials, identity, or updb (required)")
	cmd.Flags().StringVar(&clientType, "client", "claude", "Client to configure (claude, claude-code, cursor, opencode, windsurf, vscode, warp)")
	cmd.Flags().StringVar(&zitiControllerHost, "ziti-controller-host", "", "Ziti controller host")
	cmd.Flags().StringVar(&idpDomain, "idp-domain", "", "IdP domain")
	cmd.Flags().StringVar(&idpClientID, "idp-client-id", "", "IdP client ID")
	cmd.Flags().StringVar(&idpClientSecret, "idp-client-secret", "", "IdP client secret")
	cmd.Flags().StringVar(&idpAudience, "idp-audience", "", "IdP audience")
	cmd.Flags().StringVar(&identityFile, "identity-file", "", "Path to Ziti identity JSON file")
	cmd.Flags().StringVar(&username, "username", "", "UPDB username")
	cmd.Flags().StringVar(&password, "password", "", "UPDB password")
	cmd.Flags().StringVar(&toolPatterns, "tools", "*", "Comma-separated tool patterns")
	cmd.Flags().BoolVar(&readOnly, "read-only", false, "Only expose read-only tools")
	cmd.Flags().StringVar(&profile, "profile", "default", "Profile name for credentials")

	_ = cmd.MarkFlagRequired("auth-mode")

	return cmd
}
