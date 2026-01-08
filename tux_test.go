// tux_test.go
package tux

import (
	"context"
	"testing"
)

// mockAgent implements Agent for testing
type mockAgent struct {
	events chan Event
}

func (m *mockAgent) Run(ctx context.Context, prompt string) error {
	return nil
}

func (m *mockAgent) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgent) Cancel() {}

func TestAgentInterface(t *testing.T) {
	events := make(chan Event, 1) // buffered to allow non-blocking send
	agent := &mockAgent{events: events}

	// Verify Agent interface is satisfied
	var _ Agent = agent

	// Verify Run() can be called without panicking
	err := agent.Run(context.Background(), "test prompt")
	if err != nil {
		t.Errorf("Run() returned unexpected error: %v", err)
	}

	// Verify Subscribe() returns a usable channel
	ch := agent.Subscribe()
	if ch == nil {
		t.Fatal("Subscribe() returned nil channel")
	}

	// Verify events can be sent through the channel
	testEvent := Event{Type: EventText, Text: "hello"}
	events <- testEvent

	// Verify we can receive the event
	select {
	case received := <-ch:
		if received.Type != EventText {
			t.Errorf("expected EventText, got %v", received.Type)
		}
		if received.Text != "hello" {
			t.Errorf("expected 'hello', got %q", received.Text)
		}
	default:
		t.Error("expected to receive event from channel")
	}

	// Verify Cancel() can be called without panicking
	agent.Cancel()
}
