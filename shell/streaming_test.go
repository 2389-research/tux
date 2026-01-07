// shell/streaming_test.go
package shell

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestStreamingController_Lifecycle(t *testing.T) {
	s := NewStreamingController()

	// Initially not streaming
	if s.IsStreaming() {
		t.Error("expected not streaming initially")
	}

	// Start streaming
	s.Start()
	if !s.IsStreaming() {
		t.Error("expected streaming after Start")
	}
	if !s.IsWaiting() {
		t.Error("expected waiting after Start (no tokens yet)")
	}

	// End streaming
	s.End()
	if s.IsStreaming() {
		t.Error("expected not streaming after End")
	}

	// Reset clears state
	s.Start()
	s.Reset()
	if s.IsStreaming() {
		t.Error("expected not streaming after Reset")
	}
}

func TestStreamingController_Tokens(t *testing.T) {
	s := NewStreamingController()
	s.Start()

	// Append tokens
	s.AppendToken("Hello")
	s.AppendToken(" world")

	if s.GetText() != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", s.GetText())
	}

	if s.TokenCount() != 2 {
		t.Errorf("expected 2 tokens, got %d", s.TokenCount())
	}

	// No longer waiting after first token
	if s.IsWaiting() {
		t.Error("expected not waiting after tokens received")
	}

	// Token rate should be > 0 after multiple tokens
	if s.TokenRate() <= 0 {
		t.Error("expected positive token rate")
	}
}

func TestStreamingController_Status(t *testing.T) {
	s := NewStreamingController()
	s.Start()

	// Thinking
	if s.IsThinking() {
		t.Error("expected not thinking initially")
	}
	s.SetThinking(true)
	if !s.IsThinking() {
		t.Error("expected thinking after SetThinking(true)")
	}
	s.SetThinking(false)
	if s.IsThinking() {
		t.Error("expected not thinking after SetThinking(false)")
	}

	// Tool calls
	s.StartToolCall("1", "Bash")
	s.StartToolCall("2", "Read")

	active := s.ActiveToolCalls()
	if len(active) != 2 {
		t.Errorf("expected 2 active tool calls, got %d", len(active))
	}

	s.EndToolCall("1")
	active = s.ActiveToolCalls()
	if len(active) != 1 {
		t.Errorf("expected 1 active tool call, got %d", len(active))
	}
	if active[0].Name != "Read" {
		t.Errorf("expected 'Read', got %q", active[0].Name)
	}
}

func TestStreamingController_Render(t *testing.T) {
	s := NewStreamingController()

	// Not streaming - empty render
	if s.RenderStatus(nil) != "" {
		t.Error("expected empty status when not streaming")
	}

	s.Start()

	// Waiting state
	status := s.RenderStatus(nil)
	if !strings.Contains(status, "Waiting") {
		t.Errorf("expected 'Waiting' in status, got %q", status)
	}

	// After tokens, show rate
	s.AppendToken("test")
	s.tokenRate = 42.0 // Force rate for testing
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "42") {
		t.Errorf("expected rate in status, got %q", status)
	}

	// Thinking shows spinner
	s.SetThinking(true)
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "Thinking") {
		t.Errorf("expected 'Thinking' in status, got %q", status)
	}

	// Tool calls shown
	s.StartToolCall("1", "Bash")
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "Bash") {
		t.Errorf("expected 'Bash' in status, got %q", status)
	}
}

func TestStatusBar_StreamingStatus(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)
	s := NewStreamingController()

	// No streaming status when not streaming
	sb.SetStreamingController(s, true)
	view := sb.View(80)
	if strings.Contains(view, "Thinking") || strings.Contains(view, "tok/s") {
		t.Error("expected no streaming status when not streaming")
	}

	// Start streaming and set thinking
	s.Start()
	s.AppendToken("test") // Need token to exit waiting state
	s.SetThinking(true)
	view = sb.View(80)
	if !strings.Contains(view, "Thinking") {
		t.Errorf("expected 'Thinking' in view, got %q", view)
	}

	// Disabled streaming status
	sb.SetStreamingController(s, false)
	view = sb.View(80)
	if strings.Contains(view, "Thinking") {
		t.Error("expected no streaming status when disabled")
	}
}
