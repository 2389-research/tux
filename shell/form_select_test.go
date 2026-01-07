// form/select_test.go
package shell

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestSelectField(t *testing.T) {
	f := NewSelect[string]().
		WithID("theme").
		WithLabel("Choose Theme").
		WithOptions(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		)

	if f.ID() != "theme" {
		t.Errorf("expected ID 'theme', got %s", f.ID())
	}
}

func TestSelectFieldDefault(t *testing.T) {
	f := NewSelect[string]().
		WithID("theme").
		WithOptions(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		).
		WithDefault("nord")

	if f.Value().(string) != "nord" {
		t.Errorf("expected default 'nord', got %v", f.Value())
	}
}

func TestSelectFieldValue(t *testing.T) {
	f := NewSelect[string]().WithID("test").WithOptions(
		Option("A", "a"),
		Option("B", "b"),
	)

	f.SetValue("b")
	if f.Value().(string) != "b" {
		t.Errorf("expected 'b', got %v", f.Value())
	}
}

func TestSelectFieldRender(t *testing.T) {
	f := NewSelect[string]().
		WithID("test").
		WithLabel("Pick One").
		WithOptions(Option("A", "a"))

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestSelectFieldWithValidators(t *testing.T) {
	// Create a custom validator that rejects empty string
	notEmpty := func(v any) error {
		if s, ok := v.(string); ok && s == "" {
			return errors.New("value cannot be empty")
		}
		return nil
	}

	f := NewSelect[string]().
		WithID("test").
		WithOptions(Option("A", "a"), Option("B", "")).
		WithValidators(notEmpty)

	// With a valid value
	f.SetValue("a")
	if err := f.Validate(); err != nil {
		t.Errorf("expected no error with 'a', got %v", err)
	}

	// With an empty value (should fail validation)
	f.SetValue("")
	if err := f.Validate(); err == nil {
		t.Error("expected error with empty value")
	}
}

func TestSelectFieldLabel(t *testing.T) {
	f := NewSelect[string]().WithID("test").WithLabel("My Label")
	if f.Label() != "My Label" {
		t.Errorf("expected 'My Label', got %s", f.Label())
	}
}

func TestSelectFieldFocusBlur(t *testing.T) {
	f := NewSelect[string]().WithID("test")

	// Initially not focused
	if f.Focused() {
		t.Error("should not be focused initially")
	}

	// Focus sets focused flag
	f.Focus()
	if !f.Focused() {
		t.Error("should be focused after Focus()")
	}

	// Blur clears focused flag
	f.Blur()
	if f.Focused() {
		t.Error("should not be focused after Blur()")
	}
}

func TestSelectFieldValidate(t *testing.T) {
	// Without validators, should pass
	f := NewSelect[string]().WithID("test").WithOptions(Option("A", "a"))
	f.SetValue("a")

	if err := f.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSelectFieldHandleKey(t *testing.T) {
	f := NewSelect[string]().
		WithID("test").
		WithLabel("Pick").
		WithOptions(Option("A", "a"), Option("B", "b"))

	// HandleKey should return false when huhField is nil
	key := tea.KeyMsg{Type: tea.KeyDown}
	if f.HandleKey(key) {
		t.Error("should return false when huhField is nil")
	}

	// Initialize the huh field via Render
	th := theme.Get("dracula")
	f.Render(40, th, true)

	// Now HandleKey should work (may or may not return true depending on key)
	// Just verify it doesn't panic
	f.HandleKey(key)
}
