// shell/streaming_test.go
package shell

import (
	"testing"
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
