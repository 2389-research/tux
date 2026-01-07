// shell/streaming_content_test.go
package shell

import (
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
