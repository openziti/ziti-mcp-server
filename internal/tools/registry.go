package tools

import (
	"fmt"
	"log/slog"
)

// HandlerFunc is the signature for tool handler functions.
type HandlerFunc func(request HandlerRequest, config HandlerConfig) (HandlerResponse, error)

// ToolDef defines a single MCP tool with its metadata and handler.
type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema map[string]any  `json:"inputSchema,omitempty"`
	Meta        *ToolMeta       `json:"-"` // internal, not sent to MCP
	Annotations *ToolAnnotations `json:"annotations,omitempty"`
}

// Registry holds all registered tools and their handlers.
type Registry struct {
	tools    []ToolDef
	handlers map[string]HandlerFunc
}

// NewRegistry creates an empty tool registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]HandlerFunc),
	}
}

// Register adds a tool definition and its handler to the registry.
func (r *Registry) Register(def ToolDef, handler HandlerFunc) {
	if _, exists := r.handlers[def.Name]; exists {
		slog.Warn("duplicate tool registration", "tool", def.Name)
	}
	r.tools = append(r.tools, def)
	r.handlers[def.Name] = handler
}

// Tools returns all registered tool definitions.
func (r *Registry) Tools() []ToolDef {
	return r.tools
}

// Handler returns the handler for the given tool name.
func (r *Registry) Handler(name string) (HandlerFunc, error) {
	h, ok := r.handlers[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return h, nil
}

// ToolCount returns the number of registered tools.
func (r *Registry) ToolCount() int {
	return len(r.tools)
}
