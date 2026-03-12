package ca

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"go.mozilla.org/pkcs7"
)

// ParsePKCS7Certs parses a base64-encoded PKCS#7 DER blob (from EST /cacerts)
// and returns the embedded CA certificates as PEM strings.
func ParsePKCS7Certs(base64Body string) ([]string, error) {
	der, err := base64.StdEncoding.DecodeString(base64Body)
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}

	p7, err := pkcs7.Parse(der)
	if err != nil {
		return nil, fmt.Errorf("parsing PKCS#7: %w", err)
	}

	var pems []string
	for _, cert := range p7.Certificates {
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}
		pems = append(pems, string(pem.EncodeToMemory(block)))
	}

	return pems, nil
}
