package tools

import (
	"regexp"
	"strings"
)

// Glob is a simple glob pattern matcher supporting * and ? wildcards.
//   - * matches any sequence of characters (including empty string)
//   - ? matches exactly one character
type Glob struct {
	pattern string
	re      *regexp.Regexp
}

// NewGlob creates a new Glob from the given pattern string.
func NewGlob(pattern string) *Glob {
	pattern = strings.TrimSpace(pattern)
	g := &Glob{pattern: pattern}

	if pattern == "" || pattern == "*" || (!strings.Contains(pattern, "*") && !strings.Contains(pattern, "?")) {
		// No regex needed for these simple cases.
		return g
	}

	// Escape regex-special characters, then convert glob wildcards.
	escaped := regexp.QuoteMeta(pattern)
	// QuoteMeta escapes * to \* and ? to \? — convert them back to regex equivalents.
	escaped = strings.ReplaceAll(escaped, `\*`, `.*`)
	escaped = strings.ReplaceAll(escaped, `\?`, `.`)

	g.re = regexp.MustCompile(`^` + escaped + `$`)
	return g
}

// Matches tests whether str matches this glob pattern.
func (g *Glob) Matches(str string) bool {
	if g.pattern == "" {
		return str == ""
	}
	if g.pattern == "*" {
		return true
	}
	if g.re == nil {
		// No wildcards — exact match.
		return g.pattern == str
	}
	return g.re.MatchString(str)
}

// HasWildcards returns true if the pattern contains * or ?.
func (g *Glob) HasWildcards() bool {
	return strings.ContainsAny(g.pattern, "*?")
}

// String returns the original pattern.
func (g *Glob) String() string {
	return g.pattern
}

// GlobMatches is a convenience function for one-off pattern matching.
func GlobMatches(str, pattern string) bool {
	return NewGlob(pattern).Matches(str)
}
