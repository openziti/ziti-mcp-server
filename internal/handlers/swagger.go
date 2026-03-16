package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	edgeclient "github.com/openziti/ziti-mcp-server-go/internal/gen/edge/client"
	fabricclient "github.com/openziti/ziti-mcp-server-go/internal/gen/fabric/client"
)

// NewEdgeClient creates a go-swagger edge management client backed by the given
// authenticated *http.Client (which already injects zt-session headers).
func NewEdgeClient(httpClient *http.Client, host string) *edgeclient.ZitiEdgeManagement {
	transport := httptransport.NewWithClient(host, edgeclient.DefaultBasePath, edgeclient.DefaultSchemes, httpClient)
	return edgeclient.New(transport, strfmt.Default)
}

// NewFabricClient creates a go-swagger fabric management client backed by the given
// authenticated *http.Client.
func NewFabricClient(httpClient *http.Client, host string) *fabricclient.ZitiFabricManagement {
	transport := httptransport.NewWithClient(host, fabricclient.DefaultBasePath, fabricclient.DefaultSchemes, httpClient)
	return fabricclient.New(transport, strfmt.Default)
}

// noAuth is a no-op auth info writer. Auth is handled by the HTTP client transport.
var noAuth = runtime.ClientAuthInfoWriterFunc(func(_ runtime.ClientRequest, _ strfmt.Registry) error {
	return nil
})

// stripFieldFromDataItems removes the named field from each item in the
// top-level "data" array of a standard Ziti list response. This is useful for
// dropping expensive fields (like JSON schemas) from list endpoints while
// keeping them available on detail endpoints.
func stripFieldFromDataItems(response any, field string) {
	root, ok := response.(map[string]any)
	if !ok {
		return
	}
	data, ok := root["data"].([]any)
	if !ok {
		return
	}
	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			delete(m, field)
		}
	}
}

// ToMap converts a go-swagger model struct to map[string]any via JSON round-trip
// so that StripNoise and CreateSuccessResponse work correctly.
func ToMap(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshaling response: %w", err)
	}
	var result any
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}
	return result, nil
}
