package ca

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
)

// FetchControllerCA fetches the controller's CA certificate(s) from the EST
// /.well-known/est/cacerts endpoint. The initial request uses InsecureSkipVerify
// since we don't yet trust the controller (trust-on-first-use bootstrap).
func FetchControllerCA(host string) (string, error) {
	url := fmt.Sprintf("https://%s/.well-known/est/cacerts", host)
	slog.Debug("fetching controller CA", "url", url)

	httpClient := &http.Client{
		Transport: client.NewInsecureTransport(),
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/pkcs7-mime")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching CA (controller may use publicly trusted cert): %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CA endpoint returned %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading CA response: %w", err)
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return "", fmt.Errorf("CA endpoint returned empty body")
	}

	pems, err := ParsePKCS7Certs(trimmed)
	if err != nil {
		return "", fmt.Errorf("parsing PKCS#7 response: %w", err)
	}

	if len(pems) == 0 {
		return "", fmt.Errorf("no certificates found in PKCS#7 response")
	}

	slog.Debug("extracted CA certificates", "count", len(pems))
	return strings.Join(pems, "\n"), nil
}
