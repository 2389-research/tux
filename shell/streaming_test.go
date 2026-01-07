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
