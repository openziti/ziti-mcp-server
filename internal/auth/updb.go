package auth

import (
	"log/slog"

	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/terminal"
)

// StoreUPDBCredentials saves UPDB (username/password) credentials to the store.
func StoreUPDBCredentials(s *store.Store, controllerHost, username, password string) error {
	slog.Debug("storing UPDB credentials", "host", controllerHost, "username", username)

	if err := s.SetControllerHost(controllerHost); err != nil {
		return err
	}
	if err := s.SetUpdbUsername(username); err != nil {
		return err
	}
	if err := s.SetUpdbPassword(password); err != nil {
		return err
	}

	terminal.Success("UPDB credentials stored for controller %s.", controllerHost)
	return nil
}
