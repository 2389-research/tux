// form/note_test.go
package form

import (
	"strings"
	"testing"

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
