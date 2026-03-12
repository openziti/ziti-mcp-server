package client

import (
	"encoding/json"
	"log/slog"

	"github.com/openziti/ziti-mcp-server-go/internal/tools"
)

// CreateSuccessResponse strips noise from the result and wraps it as MCP content.
// If the result is an array with >1 items, each item becomes a separate text block.
func CreateSuccessResponse(result any) tools.HandlerResponse {
	before, _ := json.MarshalIndent(result, "", "  ")
	stripped := StripNoise(result)
	after, _ := json.MarshalIndent(stripped, "", "  ")

	slog.Debug("response trimmed", "before_len", len(before), "after_len", len(after))

	if arr, ok := stripped.([]any); ok && len(arr) > 1 {
		content := make([]tools.ContentItem, len(arr))
		for i, item := range arr {
			text, _ := json.MarshalIndent(item, "", "  ")
			content[i] = tools.ContentItem{Type: "text", Text: string(text)}
		}
		return tools.HandlerResponse{Content: content, IsError: false}
	}

	text, _ := json.MarshalIndent(stripped, "", "  ")
	return tools.HandlerResponse{
		Content: []tools.ContentItem{{Type: "text", Text: string(text)}},
		IsError: false,
	}
}

// CreateErrorResponse wraps an error string as an MCP error response.
func CreateErrorResponse(errMsg string) tools.HandlerResponse {
	return tools.HandlerResponse{
		Content: []tools.ContentItem{{Type: "text", Text: errMsg}},
		IsError: true,
	}
}
