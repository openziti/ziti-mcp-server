package client

import (
	"testing"
)

func TestStripNoise(t *testing.T) {
	input := map[string]any{
		"name":    "test-identity",
		"id":      "abc123",
		"_links":  map[string]any{"self": "/identities/abc123"},
		"tags":    map[string]any{},
		"envInfo": map[string]any{"os": "linux"},
		"nested": map[string]any{
			"keep":       "value",
			"authPolicy": "should-be-stripped",
		},
	}

	result := StripNoise(input)
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatal("expected map result")
	}

	if _, found := m["_links"]; found {
		t.Error("_links should be stripped")
	}
	if _, found := m["tags"]; found {
		t.Error("tags should be stripped")
	}
	if _, found := m["envInfo"]; found {
		t.Error("envInfo should be stripped")
	}
	if m["name"] != "test-identity" {
		t.Error("name should be preserved")
	}

	nested, ok := m["nested"].(map[string]any)
	if !ok {
		t.Fatal("nested should be a map")
	}
	if _, found := nested["authPolicy"]; found {
		t.Error("nested authPolicy should be stripped")
	}
	if nested["keep"] != "value" {
		t.Error("nested keep should be preserved")
	}
}

func TestStripNoiseArray(t *testing.T) {
	input := []any{
		map[string]any{"name": "a", "_links": "strip"},
		map[string]any{"name": "b", "tags": "strip"},
	}

	result := StripNoise(input)
	arr, ok := result.([]any)
	if !ok {
		t.Fatal("expected array result")
	}
	if len(arr) != 2 {
		t.Fatalf("expected 2 items, got %d", len(arr))
	}
}
