package handlers

import (
	"fmt"
	"strings"

	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// RequireString extracts a required string parameter from the request.
// Returns the value, nil, true on success. On failure returns "", errResp, false.
func RequireString(params map[string]any, key string) (string, *tools.HandlerResponse, bool) {
	v, ok := params[key]
	if !ok || v == nil {
		resp := errResponse(fmt.Sprintf("Missing required parameter: %s", key))
		return "", &resp, false
	}
	s, ok := v.(string)
	if !ok || s == "" {
		resp := errResponse(fmt.Sprintf("Parameter %s must be a non-empty string", key))
		return "", &resp, false
	}
	return s, nil, true
}

// OptionalString extracts an optional string parameter, returning "" if absent.
func OptionalString(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// OptionalBool extracts an optional boolean parameter with a default value.
func OptionalBool(params map[string]any, key string, defaultVal bool) bool {
	v, ok := params[key]
	if !ok || v == nil {
		return defaultVal
	}
	b, ok := v.(bool)
	if !ok {
		return defaultVal
	}
	return b
}

// OptionalFloat extracts an optional float64 parameter, returning nil if absent.
func OptionalFloat(params map[string]any, key string) *float64 {
	v, ok := params[key]
	if !ok || v == nil {
		return nil
	}
	f, ok := v.(float64)
	if !ok {
		return nil
	}
	return &f
}

// OptionalInt64 extracts an optional int64 parameter from a JSON number, returning nil if absent.
func OptionalInt64(params map[string]any, key string) *int64 {
	v, ok := params[key]
	if !ok || v == nil {
		return nil
	}
	f, ok := v.(float64)
	if !ok {
		return nil
	}
	n := int64(f)
	return &n
}

// OptionalObject extracts an optional object parameter, returning nil if absent.
func OptionalObject(params map[string]any, key string) map[string]any {
	v, ok := params[key]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	return m
}

// RequireObject extracts a required object parameter.
func RequireObject(params map[string]any, key string) (map[string]any, *tools.HandlerResponse, bool) {
	v, ok := params[key]
	if !ok || v == nil {
		resp := errResponse(fmt.Sprintf("Missing required parameter: %s", key))
		return nil, &resp, false
	}
	m, ok := v.(map[string]any)
	if !ok {
		resp := errResponse(fmt.Sprintf("Parameter %s must be an object", key))
		return nil, &resp, false
	}
	return m, nil, true
}

// SplitCSV splits a comma-separated string into trimmed non-empty parts.
func SplitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func strPtr(s string) *string    { return &s }
func boolPtr(b bool) *bool       { return &b }
func int64Ptr(n int64) *int64    { return &n }

func errResponse(msg string) tools.HandlerResponse {
	return tools.HandlerResponse{
		Content: []tools.ContentItem{{Type: "text", Text: msg}},
		IsError: true,
	}
}

// idSchema returns a standard JSON schema for a single required "id" parameter.
func idSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string", "description": "The resource ID"},
		},
		"required": []string{"id"},
	}
}

// emptySchema returns a JSON schema with no parameters.
func emptySchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

// readOnlyMeta returns tool metadata for a read-only tool.
func readOnlyMeta() *tools.ToolMeta {
	return &tools.ToolMeta{ReadOnly: true}
}

// writeMeta returns tool metadata for a write tool.
func writeMeta() *tools.ToolMeta {
	return &tools.ToolMeta{ReadOnly: false}
}

// readOnlyAnnotations returns annotations for a read-only, idempotent tool.
func readOnlyAnnotations(title string) *tools.ToolAnnotations {
	return &tools.ToolAnnotations{
		Title:          title,
		ReadOnlyHint:   true,
		IdempotentHint: true,
	}
}

// createAnnotations returns annotations for a create (non-destructive write) tool.
func createAnnotations(title string) *tools.ToolAnnotations {
	return &tools.ToolAnnotations{
		Title:          title,
		IdempotentHint: false,
	}
}

// updateAnnotations returns annotations for an update (destructive write) tool.
func updateAnnotations(title string) *tools.ToolAnnotations {
	return &tools.ToolAnnotations{
		Title:           title,
		DestructiveHint: true,
	}
}

// deleteAnnotations returns annotations for a delete (destructive write) tool.
func deleteAnnotations(title string) *tools.ToolAnnotations {
	return &tools.ToolAnnotations{
		Title:           title,
		DestructiveHint: true,
	}
}
