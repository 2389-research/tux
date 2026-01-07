// shell/form_confirm_test.go
package shell

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestConfirmField(t *testing.T) {
	f := NewConfirm().
		WithID("delete").
		WithLabel("Delete this file?").
		WithAffirmative("Yes, delete").
		WithNegative("Cancel")

	if f.ID() != "delete" {
		t.Errorf("expected ID 'delete', got %s", f.ID())
	}
}

func TestConfirmFieldValue(t *testing.T) {
	f := NewConfirm().WithID("test")

	// Default should be false
	if f.Value().(bool) != false {
		t.Error("default should be false")
	}

	f.SetValue(true)
	if f.Value().(bool) != true {
		t.Error("should be true after SetValue(true)")
	}
}

func TestConfirmFieldRender(t *testing.T) {
	f := NewConfirm().
		WithID("test").
		WithLabel("Confirm?")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestConfirmFieldLabel(t *testing.T) {
	f := NewConfirm().WithID("test").WithLabel("My Label")
	if f.Label() != "My Label" {
		t.Errorf("expected 'My Label', got %s", f.Label())
	}
}

func TestConfirmFieldFocusBlur(t *testing.T) {
	f := NewConfirm().WithID("test")

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

func TestConfirmFieldValidate(t *testing.T) {
	f := NewConfirm().WithID("test")

	// Confirm fields always validate successfully
	if err := f.Validate(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestConfirmFieldHandleKey(t *testing.T) {
	f := NewConfirm().WithID("test").WithLabel("Confirm?")

	// HandleKey should return false when huhField is nil
	key := tea.KeyMsg{Type: tea.KeySpace}
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
