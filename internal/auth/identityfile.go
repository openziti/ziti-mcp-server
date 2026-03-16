package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/openziti/ziti-mcp-server/internal/store"
	"github.com/openziti/ziti-mcp-server/internal/terminal"
)

// zitiIdentityFile represents the structure of a Ziti identity JSON file.
type zitiIdentityFile struct {
	ZtAPI string `json:"ztAPI"`
	ID    struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
		CA   string `json:"ca"`
	} `json:"id"`
}

// stripPemPrefix removes the "pem:" prefix that Ziti identity files use.
func stripPemPrefix(value string) string {
	return strings.TrimPrefix(value, "pem:")
}

// validatePem checks that a string contains PEM-encoded data.
func validatePem(value, fieldName string) error {
	p := stripPemPrefix(value)
	if !strings.Contains(p, "-----BEGIN ") || !strings.Contains(p, "-----END ") {
		return fmt.Errorf("invalid PEM format for %s: missing BEGIN/END markers", fieldName)
	}
	return nil
}

// parseAndStoreIdentity parses raw identity JSON and stores the certificate material.
func parseAndStoreIdentity(s *store.Store, data []byte) (string, error) {
	var identity zitiIdentityFile
	if err := json.Unmarshal(data, &identity); err != nil {
		return "", fmt.Errorf("identity file is not valid JSON: %w", err)
	}

	// Validate structure
	if identity.ZtAPI == "" {
		return "", fmt.Errorf("identity file missing required field: ztAPI")
	}
	if identity.ID.Cert == "" {
		return "", fmt.Errorf("identity file missing required field: id.cert")
	}
	if identity.ID.Key == "" {
		return "", fmt.Errorf("identity file missing required field: id.key")
	}
	if identity.ID.CA == "" {
		return "", fmt.Errorf("identity file missing required field: id.ca")
	}

	// Validate PEM format
	for _, check := range []struct{ value, name string }{
		{identity.ID.Cert, "id.cert"},
		{identity.ID.Key, "id.key"},
		{identity.ID.CA, "id.ca"},
	} {
		if err := validatePem(check.value, check.name); err != nil {
			return "", err
		}
	}

	// Extract controller host
	u, err := url.Parse(identity.ZtAPI)
	if err != nil {
		return "", fmt.Errorf("invalid ztAPI URL: %w", err)
	}
	controllerHost := u.Host
	slog.Debug("extracted controller host", "host", controllerHost)

	// Store in credential store
	if err := s.SetIdentityCert(stripPemPrefix(identity.ID.Cert)); err != nil {
		return "", err
	}
	if err := s.SetIdentityKey(stripPemPrefix(identity.ID.Key)); err != nil {
		return "", err
	}
	if err := s.SetIdentityCA(stripPemPrefix(identity.ID.CA)); err != nil {
		return "", err
	}
	if err := s.SetControllerHost(controllerHost); err != nil {
		return "", err
	}

	return controllerHost, nil
}

// RequestIdentityFileAuthorization reads a Ziti identity JSON file, validates it,
// and stores the certificate material in the credential store.
func RequestIdentityFileAuthorization(s *store.Store, identityFile string) error {
	slog.Debug("reading identity file", "path", identityFile)

	data, err := os.ReadFile(identityFile)
	if err != nil {
		return fmt.Errorf("reading identity file: %w", err)
	}

	controllerHost, err := parseAndStoreIdentity(s, data)
	if err != nil {
		return err
	}

	terminal.Success("Successfully loaded identity file for controller %s.", controllerHost)
	return nil
}

// RequestIdentityJSONAuthorization takes raw identity JSON content (not a file path),
// validates it, and stores the certificate material in the credential store.
func RequestIdentityJSONAuthorization(s *store.Store, jsonContent string) error {
	slog.Debug("processing identity JSON content")

	controllerHost, err := parseAndStoreIdentity(s, []byte(jsonContent))
	if err != nil {
		return err
	}

	slog.Info("identity JSON stored for controller", "host", controllerHost)
	return nil
}
