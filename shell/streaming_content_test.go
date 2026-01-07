// shell/streaming_content_test.go
package shell

import (
	"strings"
	"testing"
	"time"

	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

// mockContent implements content.Content for testing
type mockContent struct {
	text   string
	width  int
	height int
}

func (m *mockContent) Init() tea.Cmd                                 { return nil }
func (m *mockContent) Update(msg tea.Msg) (content.Content, tea.Cmd) { return m, nil }
func (m *mockContent) View() string                                  { return m.text }
func (m *mockContent) Value() any                                    { return m.text }
func (m *mockContent) SetSize(w, h int)                              { m.width = w; m.height = h }

func TestStreamingContent_Builder(t *testing.T) {
	inner := &mockContent{text: "hello"}

	sc := NewStreamingContent(inner)
	if sc == nil {
		t.Fatal("expected non-nil StreamingContent")
	}

	// Default: typewriter disabled
	if sc.typewriter {
		t.Error("expected typewriter disabled by default")
	}

	// Builder methods return self
	sc2 := sc.WithTypewriter(true).WithSpeed(50 * time.Millisecond)
	if sc2 != sc {
		t.Error("expected builder to return same instance")
	}

	if !sc.typewriter {
		t.Error("expected typewriter enabled")
	}

	if sc.typewriterSpeed != 50*time.Millisecond {
		t.Errorf("expected 50ms speed, got %v", sc.typewriterSpeed)
	}
}

func TestStreamingContent_Typewriter(t *testing.T) {
	inner := &mockContent{}
	sc := NewStreamingContent(inner).WithTypewriter(true)

	sc.SetText("Hello world")

	// Initially position is 0, should show empty or cursor only
	view := sc.View()
	if len(view) > 5 { // Just cursor character
		t.Errorf("expected minimal view initially, got %q", view)
	}

	// Advance position
	sc.position = 5
	view = sc.View()
	if !strings.HasPrefix(view, "Hello") {
		t.Errorf("expected 'Hello' prefix, got %q", view)
	}

	// Should have cursor
	if !strings.Contains(view, "│") {
		t.Errorf("expected cursor in view, got %q", view)
	}

	// Full position shows all text
	sc.position = len("Hello world")
	view = sc.View()
	if !strings.Contains(view, "Hello world") {
		t.Errorf("expected full text, got %q", view)
	}
}

func TestStreamingContent_TypewriterDisabled(t *testing.T) {
	inner := &mockContent{}
	sc := NewStreamingContent(inner) // typewriter disabled by default

	sc.SetText("Hello world")
	view := sc.View()

	// Should show full text immediately
	if view != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", view)
	}
}
