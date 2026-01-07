// form/note_test.go
package shell

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestNoteField(t *testing.T) {
	f := NewNote().
		WithID("warning").
		WithTitle("Warning").
		WithContent("This action cannot be undone.")

	if f.ID() != "warning" {
		t.Errorf("expected ID 'warning', got %s", f.ID())
	}
}

func TestNoteFieldValue(t *testing.T) {
	f := NewNote().WithID("test").WithContent("hello")

	// Note value is always nil (display only)
	if f.Value() != nil {
		t.Error("note value should be nil")
	}
}

func TestNoteFieldRender(t *testing.T) {
	f := NewNote().
		WithID("test").
		WithTitle("Important").
		WithContent("Read this carefully.")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
	if !strings.Contains(output, "Important") && !strings.Contains(output, "Read this") {
		t.Error("render should contain title or content")
	}
}

func TestNoteFieldNoValidation(t *testing.T) {
	f := NewNote().WithID("test")

	if f.Validate() != nil {
		t.Error("note should never fail validation")
	}
}

func TestNoteFieldLabel(t *testing.T) {
	f := NewNote().WithID("test").WithTitle("My Title")
	// Note: Label() returns the title for NoteField
	if f.Label() != "My Title" {
		t.Errorf("expected 'My Title', got %s", f.Label())
	}
}

func TestNoteFieldSetValue(t *testing.T) {
	f := NewNote().WithID("test")

	// SetValue should do nothing for notes (display-only)
	f.SetValue("anything")

	// Value should still be nil
	if f.Value() != nil {
		t.Error("note value should remain nil after SetValue")
	}
}

func TestNoteFieldFocusBlur(t *testing.T) {
	f := NewNote().WithID("test")

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

func TestNoteFieldInit(t *testing.T) {
	f := NewNote().WithID("test")

	// Init should return nil for notes
	cmd := f.Init()
	if cmd != nil {
		t.Error("note Init() should return nil")
	}
}

func TestNoteFieldHandleKey(t *testing.T) {
	f := NewNote().WithID("test")

	// HandleKey should always return false for notes (display-only)
	key := tea.KeyMsg{Type: tea.KeySpace}
	if f.HandleKey(key) {
		t.Error("note HandleKey should always return false")
	}

	// Try different key types
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	if f.HandleKey(enterKey) {
		t.Error("note HandleKey should always return false for enter")
	}
}
