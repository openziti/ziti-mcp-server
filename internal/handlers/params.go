package handlers

import (
	"fmt"
	"strings"

	"github.com/openziti/ziti-mcp-server/internal/tools"
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

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
func int64Ptr(n int64) *int64 { return &n }

// listPageSize is the number of items requested per page when auto-paginating.
// OpenZiti supports a maximum of 500 per page.
const listPageSize = int64(500)

// listPageFunc is a function that fetches one page from a list endpoint.
type listPageFunc func(limit, offset int64) (map[string]any, error)

// fetchAllPages calls fetchPage repeatedly with increasing offsets until all
// items have been collected, then returns a single merged result map whose
// "data" field contains the complete item slice.  The "meta" field from the
// last page is preserved so callers still see pagination metadata.
func fetchAllPages(fetchPage listPageFunc) (map[string]any, error) {
	first, err := fetchPage(listPageSize, 0)
	if err != nil {
		return nil, err
	}

	data, _ := first["data"].([]any)
	total := paginationTotalCount(first)

	// Single page — nothing else to fetch.
	if int64(len(data)) >= total {
		return first, nil
	}

	// Accumulate remaining pages.
	offset := int64(len(data))
	for offset < total {
		page, err := fetchPage(listPageSize, offset)
		if err != nil {
			return nil, err
		}
		pageData, _ := page["data"].([]any)
		if len(pageData) == 0 {
			break
		}
		data = append(data, pageData...)
		offset += int64(len(pageData))
	}

	first["data"] = data
	return first, nil
}

// paginationTotalCount extracts meta.pagination.totalCount from a list response map.
func paginationTotalCount(m map[string]any) int64 {
	meta, _ := m["meta"].(map[string]any)
	if meta == nil {
		return 0
	}
	pagination, _ := meta["pagination"].(map[string]any)
	if pagination == nil {
		return 0
	}
	switch v := pagination["totalCount"].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	}
	return 0
}

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
