package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
)

// NewMTLSTransport creates an http.Transport configured for mutual TLS
// using client cert/key and a CA certificate pool.
func NewMTLSTransport(certPEM, keyPEM, caPEM string) (*http.Transport, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return nil, fmt.Errorf("loading client certificate: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM([]byte(caPEM)) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caPool,
		},
	}, nil
}

// NewCATrustTransport creates an http.Transport that trusts the given CA
// certificate(s) but does not present a client certificate.
func NewCATrustTransport(caPEM string) (*http.Transport, error) {
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM([]byte(caPEM)) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caPool,
		},
	}, nil
}

// NewInsecureTransport creates an http.Transport that skips TLS verification.
// Used only for bootstrap operations like fetching the controller CA.
func NewInsecureTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // intentional for CA bootstrap
		},
	}
}
