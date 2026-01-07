// shell/form_field_test.go
package shell

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

func TestValuesInt(t *testing.T) {
	v := Values{
		"count":       42,
		"wrongType":   "not an int",
		"anotherType": true,
	}

	// Valid int value
	if v.Int("count") != 42 {
		t.Errorf("expected 42, got %d", v.Int("count"))
	}

	// Missing key should return 0
	if v.Int("missing") != 0 {
		t.Errorf("missing key should return 0, got %d", v.Int("missing"))
	}

	// Wrong type should return 0
	if v.Int("wrongType") != 0 {
		t.Errorf("wrong type should return 0, got %d", v.Int("wrongType"))
	}

	// Another wrong type should return 0
	if v.Int("anotherType") != 0 {
		t.Errorf("another wrong type should return 0, got %d", v.Int("anotherType"))
	}
}
