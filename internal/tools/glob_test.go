package tools

import "testing"

func TestGlob(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"*", "anything", true},
		{"*", "", true},
		{"", "", true},
		{"", "x", false},
		{"list*", "listIdentities", true},
		{"list*", "listServices", true},
		{"list*", "deleteIdentities", false},
		{"*Identities", "listIdentities", true},
		{"*Identities", "deleteIdentities", true},
		{"file.?s", "file.js", true},
		{"file.?s", "file.ts", true},
		{"file.?s", "file.jsx", false},
		{"exact", "exact", true},
		{"exact", "other", false},
		{"test*", "testing", true},
		{"test*", "contest", false},
	}

	for _, tt := range tests {
		g := NewGlob(tt.pattern)
		got := g.Matches(tt.input)
		if got != tt.want {
			t.Errorf("Glob(%q).Matches(%q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
		}
	}
}

func TestGlobHasWildcards(t *testing.T) {
	if !NewGlob("test*").HasWildcards() {
		t.Error("expected HasWildcards for test*")
	}
	if !NewGlob("file.?s").HasWildcards() {
		t.Error("expected HasWildcards for file.?s")
	}
	if NewGlob("exact").HasWildcards() {
		t.Error("did not expect HasWildcards for exact")
	}
}
