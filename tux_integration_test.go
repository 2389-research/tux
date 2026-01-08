// tux_integration_test.go
package tux

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/2389-research/tux/theme"
)

func TestFullConversationFlow(t *testing.T) {
	// Create mock agent that simulates a conversation
	events := make(chan Event, 10)
	runStarted := make(chan struct{}, 1)
	agent := &conversationAgent{
		events:     events,
		runStarted: runStarted,
	}

	app := New(agent, WithTheme(theme.NewDraculaTheme()))

	// Submit a prompt
	app.submitInput("What is 2+2?")

	// Wait for run to start
	select {
	case <-runStarted:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run should have started")
	}

	// Simulate agent response
	events <- Event{Type: EventText, Text: "The answer is "}
	events <- Event{Type: EventText, Text: "4"}
	events <- Event{Type: EventComplete}

	// Give time for events to process
	time.Sleep(50 * time.Millisecond)

	// Verify chat has user message
	chatView := app.chat.View()
	if !strings.Contains(chatView, "What is 2+2?") {
		t.Errorf("Chat should contain user message, got: %s", chatView)
	}

	// Verify chat has assistant response
	if !strings.Contains(chatView, "The answer is") {
		t.Errorf("Chat should contain assistant response, got: %s", chatView)
	}
}

func TestConversationWithToolCall(t *testing.T) {
	events := make(chan Event, 10)
	runStarted := make(chan struct{}, 1)
	agent := &conversationAgent{
		events:     events,
		runStarted: runStarted,
	}

	app := New(agent)

	app.submitInput("Read file.txt")

	select {
	case <-runStarted:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run should have started")
	}

	// Simulate tool call flow
	events <- Event{Type: EventText, Text: "Let me read that file."}
	events <- Event{
		Type:       EventToolCall,
		ToolID:     "tool-1",
		ToolName:   "read_file",
		ToolParams: map[string]any{"path": "file.txt"},
	}
	events <- Event{
		Type:       EventToolResult,
		ToolID:     "tool-1",
		ToolOutput: "Hello, World!",
		Success:    true,
	}
	events <- Event{Type: EventText, Text: " The file contains: Hello, World!"}
	events <- Event{Type: EventComplete}

	time.Sleep(50 * time.Millisecond)

	// Verify tools tab has the call
	toolsView := app.tools.View()
	if !strings.Contains(toolsView, "read_file") {
		t.Errorf("Tools should show tool call, got: %s", toolsView)
	}
	if !strings.Contains(toolsView, "\u2713") {
		t.Errorf("Tools should show success marker, got: %s", toolsView)
	}
}

type conversationAgent struct {
	events     chan Event
	runStarted chan struct{}
}

func (a *conversationAgent) Run(ctx context.Context, prompt string) error {
	if a.runStarted != nil {
		a.runStarted <- struct{}{}
	}
	// Don't block - just return after signaling
	return nil
}

func (a *conversationAgent) Subscribe() <-chan Event {
	return a.events
}

func (a *conversationAgent) Cancel() {}

func TestConversationWithError(t *testing.T) {
	events := make(chan Event, 10)
	runStarted := make(chan struct{}, 1)
	agent := &conversationAgent{
		events:     events,
		runStarted: runStarted,
	}

	app := New(agent)

	app.submitInput("Do something")

	select {
	case <-runStarted:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run should have started")
	}

	// Simulate an error event
	events <- Event{Type: EventText, Text: "Starting..."}
	events <- Event{Type: EventError, Error: fmt.Errorf("something went wrong")}

	time.Sleep(50 * time.Millisecond)

	// For now, just verify it doesn't panic
	// Error display will be implemented in future work
	_ = app.chat.View()
}
