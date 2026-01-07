// form/input_test.go
package form

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestInputField(t *testing.T) {
	f := NewInput().
		WithID("username").
		WithLabel("Username").
		WithPlaceholder("Enter username...")

	if f.ID() != "username" {
		t.Errorf("expected ID 'username', got %s", f.ID())
	}
	if f.Label() != "Username" {
		t.Errorf("expected label 'Username', got %s", f.Label())
	}
}

func TestInputFieldValue(t *testing.T) {
	f := NewInput().WithID("test")

	f.SetValue("hello")
	if f.Value().(string) != "hello" {
		t.Errorf("expected 'hello', got %v", f.Value())
	}
}

func TestInputFieldValidation(t *testing.T) {
	f := NewInput().
		WithID("email").
		WithValidators(Required(), Email())

	f.SetValue("")
	if f.Validate() == nil {
		t.Error("empty should fail validation")
	}

	f.SetValue("not-email")
	if f.Validate() == nil {
		t.Error("invalid email should fail")
	}

	f.SetValue("test@example.com")
	if f.Validate() != nil {
		t.Error("valid email should pass")
	}
}

func TestInputFieldRender(t *testing.T) {
	f := NewInput().
		WithID("test").
		WithLabel("Test Input")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestInputFieldFocus(t *testing.T) {
	f := NewInput().WithID("test")

	if f.Focused() {
		t.Error("should not be focused initially")
	}

	f.Focus()
	if !f.Focused() {
		t.Error("should be focused after Focus()")
	}

	f.Blur()
	if f.Focused() {
		t.Error("should not be focused after Blur()")
	}
}

func TestInputFieldHandleKey(t *testing.T) {
	f := NewInput().WithID("test").WithLabel("Input")

	// HandleKey should return false when huhField is nil
	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
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
