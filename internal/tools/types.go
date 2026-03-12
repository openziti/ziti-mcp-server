package tools

// AuthMode represents the authentication mode for the current session.
type AuthMode string

const (
	AuthModeToken    AuthMode = "token"
	AuthModeIdentity AuthMode = "identity"
	AuthModeUPDB     AuthMode = "updb"
)

// ToolAnnotations provides hints about a tool's behavior to MCP clients.
type ToolAnnotations struct {
	DestructiveHint bool   `json:"destructiveHint,omitempty"`
	IdempotentHint  bool   `json:"idempotentHint,omitempty"`
	OpenWorldHint   bool   `json:"openWorldHint,omitempty"`
	ReadOnlyHint    bool   `json:"readOnlyHint,omitempty"`
	Title           string `json:"title,omitempty"`
}

// ToolMeta contains internal metadata not exposed to MCP clients.
type ToolMeta struct {
	RequiredScopes []string `json:"requiredScopes"`
	ReadOnly       bool     `json:"readOnly,omitempty"`
}

// HandlerRequest is passed to tool handlers with auth context and parameters.
type HandlerRequest struct {
	Token      string
	Parameters map[string]any
}

// HandlerConfig provides the connection context for tool handlers.
type HandlerConfig struct {
	ZitiControllerHost string
	Domain             string
	AuthMode           AuthMode
	Profile            string
}

// ContentItem represents a single content block in an MCP response.
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// HandlerResponse is the return value from tool handlers.
type HandlerResponse struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError"`
}
