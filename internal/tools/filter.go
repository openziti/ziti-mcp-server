package tools

import (
	"fmt"
	"log/slog"
	"strings"
)

// GetAvailableTools filters tools by glob patterns and/or readOnly flag.
// The readOnly flag takes priority over pattern matching for security.
func GetAvailableTools(allTools []ToolDef, patterns []string, readOnly bool) []ToolDef {
	filtered := allTools

	if len(patterns) > 0 {
		filtered = filterToolsByPatterns(filtered, patterns)
	}

	if readOnly {
		filtered = filterToolsByReadOnly(filtered)
	}

	return filtered
}

func filterToolsByPatterns(tools []ToolDef, patterns []string) []ToolDef {
	// Global wildcard — return everything.
	if len(patterns) == 1 && patterns[0] == "*" {
		return tools
	}

	globs := make([]*Glob, len(patterns))
	for i, p := range patterns {
		globs[i] = NewGlob(p)
	}

	enabled := make(map[string]bool)
	matchCounts := make(map[string]int)

	for _, tool := range tools {
		for _, g := range globs {
			if g.Matches(tool.Name) {
				enabled[tool.Name] = true
				matchCounts[g.String()]++
				break
			}
		}
	}

	for pattern, count := range matchCounts {
		if strings.ContainsAny(pattern, "*?") {
			slog.Debug("glob pattern matched", "pattern", pattern, "count", count)
		}
	}

	var filtered []ToolDef
	for _, tool := range tools {
		if enabled[tool.Name] {
			filtered = append(filtered, tool)
		}
	}

	slog.Debug("selected tools based on patterns", "count", len(filtered))
	return filtered
}

func filterToolsByReadOnly(tools []ToolDef) []ToolDef {
	var readOnly []ToolDef
	for _, tool := range tools {
		if tool.Meta != nil && tool.Meta.ReadOnly {
			readOnly = append(readOnly, tool)
		}
	}
	slog.Debug("filtered to read-only tools", "count", len(readOnly))
	return readOnly
}

// ValidatePatterns checks that every pattern matches at least one tool.
func ValidatePatterns(patterns []string, availableTools []ToolDef) error {
	if len(patterns) == 0 {
		return nil
	}
	if len(availableTools) == 0 {
		return fmt.Errorf("no tools available for pattern validation")
	}

	names := make([]string, len(availableTools))
	for i, t := range availableTools {
		names[i] = t.Name
	}

	for _, pattern := range patterns {
		g := NewGlob(pattern)
		found := false
		for _, name := range names {
			if g.Matches(name) {
				found = true
				break
			}
		}
		if !found {
			prefix := "Invalid tool"
			if strings.ContainsAny(pattern, "*?") {
				prefix = "No tools match the pattern"
			}
			return fmt.Errorf("%s: %s. Accepted tools are: %s", prefix, pattern, strings.Join(names, ", "))
		}
	}

	return nil
}
