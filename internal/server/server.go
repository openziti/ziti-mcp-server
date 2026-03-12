package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/openziti/ziti-mcp-server-go/internal/client"
	"github.com/openziti/ziti-mcp-server-go/internal/config"
	"github.com/openziti/ziti-mcp-server-go/internal/store"
	"github.com/openziti/ziti-mcp-server-go/internal/tools"
	"github.com/openziti/ziti-mcp-server-go/internal/version"
)

// Options configures the MCP server.
type Options struct {
	Tools    []string
	ReadOnly bool
}

// Start initializes and runs the MCP server over stdio.
// The metaRegistry parameter may be nil if no meta-tools are needed.
func Start(s *store.Store, registry *tools.Registry, metaRegistry *tools.Registry, opts Options) error {
	slog.Debug("initializing OpenZiti MCP server")

	// Filter Ziti API tools (meta-tools are exempt from filtering)
	allTools := registry.Tools()
	availableTools := tools.GetAvailableTools(allTools, opts.Tools, opts.ReadOnly)

	// Create MCP server
	srv := mcpserver.NewMCPServer(
		"ziti",
		version.Version,
		mcpserver.WithLogging(),
	)

	// Register each available Ziti API tool
	for _, toolDef := range availableTools {
		td := toolDef // capture loop variable
		mcpTool := buildMCPTool(td)
		handler := createToolHandler(td.Name, registry, s)
		srv.AddTool(mcpTool, handler)
	}

	// Register meta-tools (always available, not filtered)
	metaToolCount := 0
	if metaRegistry != nil {
		for _, toolDef := range metaRegistry.Tools() {
			td := toolDef
			mcpTool := buildMCPTool(td)
			handler := createMetaToolHandler(td.Name, metaRegistry, s)
			srv.AddTool(mcpTool, handler)
			metaToolCount++
		}
	}

	enabledCount := len(availableTools)
	totalCount := len(allTools)
	slog.Info("OpenZiti MCP Server running on stdio",
		"version", version.Version,
		"tools", fmt.Sprintf("%d/%d", enabledCount, totalCount),
		"meta-tools", metaToolCount)

	// Start stdio transport
	stdioSrv := mcpserver.NewStdioServer(srv)
	return stdioSrv.Listen(context.Background(), os.Stdin, os.Stdout)
}

// buildMCPTool creates an mcp.Tool from a ToolDef.
func buildMCPTool(td tools.ToolDef) mcp.Tool {
	mcpTool := mcp.Tool{
		Name:        td.Name,
		Description: td.Description,
	}

	// Set input schema using RawInputSchema to avoid mcp-go's custom
	// marshaler which outputs "required":[] for tools with no required
	// fields, which violates JSON Schema draft-04 and causes issues
	// with Claude Desktop's schema validation.
	if td.InputSchema != nil {
		schema := make(map[string]any, len(td.InputSchema)+1)
		for k, v := range td.InputSchema {
			schema[k] = v
		}
		if _, hasType := schema["type"]; !hasType {
			schema["type"] = "object"
		}
		rawSchema, _ := json.Marshal(schema)
		mcpTool.RawInputSchema = rawSchema
	} else {
		mcpTool.RawInputSchema = json.RawMessage(`{"type":"object","properties":{}}`)
	}

	// Set annotations
	if td.Annotations != nil {
		mcpTool.Annotations = mcp.ToolAnnotation{
			Title:           td.Annotations.Title,
			ReadOnlyHint:    boolPtr(td.Annotations.ReadOnlyHint),
			DestructiveHint: boolPtr(td.Annotations.DestructiveHint),
			IdempotentHint:  boolPtr(td.Annotations.IdempotentHint),
			OpenWorldHint:   boolPtr(td.Annotations.OpenWorldHint),
		}
	}

	return mcpTool
}

// createToolHandler wraps a registered Ziti API tool handler for mcp-go.
// Validates config on each call — returns an error if disconnected.
func createToolHandler(
	toolName string,
	registry *tools.Registry,
	s *store.Store,
) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		slog.Debug("tool call received", "tool", toolName)

		handler, err := registry.Handler(toolName)
		if err != nil {
			return errorResult(fmt.Sprintf("Unknown tool: %s", toolName)), nil
		}

		// Reload config on each call
		currentCfg := config.LoadConfig(s)
		if !config.ValidateConfig(currentCfg, s) {
			return errorResult("OpenZiti configuration is invalid or missing. Use the 'login' tool to connect to a Ziti network."), nil
		}

		// Build handler request
		params := request.GetArguments()
		if params == nil {
			params = make(map[string]any)
		}

		handlerReq := tools.HandlerRequest{
			Token:      currentCfg.Token,
			Parameters: params,
		}

		handlerCfg := tools.HandlerConfig{
			ZitiControllerHost: currentCfg.ZitiControllerHost,
			AuthMode:           currentCfg.AuthMode,
			Profile:            s.ActiveProfile(),
		}

		if currentCfg.AuthMode != tools.AuthModeIdentity && currentCfg.AuthMode != tools.AuthModeUPDB {
			handlerCfg.Domain = client.FormatDomain(currentCfg.Domain)
		}

		// Execute handler
		result, err := handler(handlerReq, handlerCfg)
		if err != nil {
			slog.Error("handler error", "tool", toolName, "error", err)
			return errorResult(fmt.Sprintf("Error: %s", err)), nil
		}

		// Convert to MCP result
		mcpResult := &mcp.CallToolResult{
			IsError: result.IsError,
		}
		for _, item := range result.Content {
			mcpResult.Content = append(mcpResult.Content, mcp.TextContent{
				Type: "text",
				Text: item.Text,
			})
		}

		slog.Debug("tool call completed", "tool", toolName)
		return mcpResult, nil
	}
}

// createMetaToolHandler wraps a meta-tool handler for mcp-go.
// Does NOT validate Ziti config — meta-tools operate on the store directly.
func createMetaToolHandler(
	toolName string,
	registry *tools.Registry,
	s *store.Store,
) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		slog.Debug("meta-tool call received", "tool", toolName)

		handler, err := registry.Handler(toolName)
		if err != nil {
			return errorResult(fmt.Sprintf("Unknown meta-tool: %s", toolName)), nil
		}

		params := request.GetArguments()
		if params == nil {
			params = make(map[string]any)
		}

		handlerReq := tools.HandlerRequest{
			Parameters: params,
		}
		handlerCfg := tools.HandlerConfig{
			Profile: s.ActiveProfile(),
		}

		result, err := handler(handlerReq, handlerCfg)
		if err != nil {
			slog.Error("meta-tool error", "tool", toolName, "error", err)
			return errorResult(fmt.Sprintf("Error: %s", err)), nil
		}

		mcpResult := &mcp.CallToolResult{
			IsError: result.IsError,
		}
		for _, item := range result.Content {
			mcpResult.Content = append(mcpResult.Content, mcp.TextContent{
				Type: "text",
				Text: item.Text,
			})
		}

		slog.Debug("meta-tool call completed", "tool", toolName)
		return mcpResult, nil
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: msg},
		},
	}
}

func boolPtr(b bool) *bool {
	return &b
}
