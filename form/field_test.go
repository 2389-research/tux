// form/field_test.go
package form

import "testing"

func TestFieldInterface(t *testing.T) {
	// Verify Option helper works
	opt := Option("Display", "value")
	if opt.Label != "Display" {
		t.Errorf("expected label 'Display', got %s", opt.Label)
	}
	if opt.Value != "value" {
		t.Errorf("expected value 'value', got %v", opt.Value)
	}
}

func TestValuesAccessors(t *testing.T) {
	v := Values{
		"name":     "Alice",
		"active":   true,
		"features": []string{"a", "b"},
	}

	if v.String("name") != "Alice" {
		t.Errorf("expected 'Alice', got %s", v.String("name"))
	}
	if v.String("missing") != "" {
		t.Error("missing key should return empty string")
	}
	if !v.Bool("active") {
		t.Error("expected true")
	}
	if v.Bool("missing") {
		t.Error("missing bool should return false")
	}
	strs := v.Strings("features")
	if len(strs) != 2 || strs[0] != "a" {
		t.Errorf("expected [a, b], got %v", strs)
	}
}
