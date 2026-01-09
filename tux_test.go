// tux_test.go
package tux

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/2389-research/tux/theme"
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

func TestWithThemeOption(t *testing.T) {
	th := theme.NewDraculaTheme()
	opt := WithTheme(th)
	if opt == nil {
		t.Error("WithTheme should return an option")
	}
}

func TestWithTabOption(t *testing.T) {
	tab := TabDef{
		ID:    "custom",
		Label: "Custom",
	}
	opt := WithTab(tab)
	if opt == nil {
		t.Error("WithTab should return an option")
	}
}

func TestWithoutTabOption(t *testing.T) {
	opt := WithoutTab("tools")
	if opt == nil {
		t.Error("WithoutTab should return an option")
	}
}

func TestNewApp(t *testing.T) {
	events := make(chan Event, 1)
	agent := &mockAgent{events: events}

	app := New(agent)
	if app == nil {
		t.Error("New should return an App")
	}
}

func TestNewAppWithOptions(t *testing.T) {
	events := make(chan Event, 1)
	agent := &mockAgent{events: events}

	app := New(agent,
		WithTheme(theme.NewNeoTerminalTheme()),
		WithoutTab("tools"),
	)
	if app == nil {
		t.Error("New with options should return an App")
	}
}

func TestAppHasDefaultTabs(t *testing.T) {
	events := make(chan Event, 1)
	agent := &mockAgent{events: events}

	app := New(agent)

	// App should have chat content accessible
	if app.chat == nil {
		t.Error("App should have chat content")
	}
	if app.tools == nil {
		t.Error("App should have tools content")
	}
}

func TestAppWithoutToolsTab(t *testing.T) {
	events := make(chan Event, 1)
	agent := &mockAgent{events: events}

	app := New(agent, WithoutTab("tools"))

	if app.chat == nil {
		t.Error("App should still have chat content")
	}
	// tools content still exists for event routing, just not displayed as tab
}

func TestAppWithMultipleCustomTabs(t *testing.T) {
	events := make(chan Event, 1)
	agent := &mockAgent{events: events}

	app := New(agent,
		WithTab(TabDef{ID: "custom1", Label: "Custom 1"}),
		WithTab(TabDef{ID: "custom2", Label: "Custom 2"}),
	)

	if app == nil {
		t.Error("App with multiple custom tabs should not be nil")
	}
}

func TestAppRoutesTextEvents(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Process a text event
	app.processEvent(Event{Type: EventText, Text: "Hello"})

	// Chat should have the text
	view := app.chat.View()
	if !strings.Contains(view, "Hello") {
		t.Errorf("Chat should contain 'Hello', got: %s", view)
	}
}

func TestAppRoutesToolCallEvents(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Process tool call event
	app.processEvent(Event{
		Type:       EventToolCall,
		ToolID:     "tool-1",
		ToolName:   "read_file",
		ToolParams: map[string]any{"path": "/test"},
	})

	// Tools should have the call
	view := app.tools.View()
	if !strings.Contains(view, "read_file") {
		t.Errorf("Tools should contain 'read_file', got: %s", view)
	}
}

func TestAppRoutesToolResultEvents(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Add a tool call first
	app.processEvent(Event{
		Type:     EventToolCall,
		ToolID:   "tool-1",
		ToolName: "read_file",
	})

	// Then add result
	app.processEvent(Event{
		Type:       EventToolResult,
		ToolID:     "tool-1",
		ToolOutput: "file contents",
		Success:    true,
	})

	// Tools should show success
	view := app.tools.View()
	if !strings.Contains(view, "\u2713") {
		t.Errorf("Tools should show success marker, got: %s", view)
	}
}

func TestAppRoutesCompleteEvent(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Add streaming text
	app.processEvent(Event{Type: EventText, Text: "Response"})

	// Complete the message
	app.processEvent(Event{Type: EventComplete})

	// Chat should have the message finalized
	messages, ok := app.chat.Value().([]chatMessage)
	if !ok {
		t.Error("Value() should return []chatMessage")
		return
	}
	if len(messages) != 1 {
		t.Errorf("Expected 1 finalized message, got %d", len(messages))
	}
}

func TestAppRoutesErrorEvent(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Process an error event (currently a no-op, but shouldn't panic)
	app.processEvent(Event{
		Type:  EventError,
		Error: fmt.Errorf("test error"),
	})

	// For now, just verify it doesn't panic
	// Error display will be implemented in future tasks
}

func TestAppSubmitInput(t *testing.T) {
	events := make(chan Event, 10)
	runCalled := make(chan string, 1) // Channel to signal run was called with prompt
	agent := &mockAgentWithRun{
		events: events,
		onRun: func(prompt string) {
			runCalled <- prompt
		},
	}

	app := New(agent)
	app.submitInput("Hello")

	// Wait for run to be called (with timeout)
	select {
	case prompt := <-runCalled:
		if prompt != "Hello" {
			t.Errorf("Expected prompt 'Hello', got '%s'", prompt)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Agent.Run should have been called")
	}
}

func TestAppSubmitInputAddsUserMessage(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgentWithRun{events: events}

	app := New(agent)
	app.submitInput("Test prompt")

	// Chat should have user message
	view := app.chat.View()
	if !strings.Contains(view, "Test prompt") {
		t.Errorf("Chat should contain user message, got: %s", view)
	}
}

type mockAgentWithRun struct {
	events chan Event
	onRun  func(string)
}

func (m *mockAgentWithRun) Run(ctx context.Context, prompt string) error {
	if m.onRun != nil {
		m.onRun(prompt)
	}
	return nil
}

func (m *mockAgentWithRun) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgentWithRun) Cancel() {}

func TestAppCancelOnEscape(t *testing.T) {
	events := make(chan Event, 10)
	cancelCalled := make(chan struct{}, 1)
	agent := &mockAgentWithCancel{
		events: events,
		onCancel: func() {
			cancelCalled <- struct{}{}
		},
	}

	app := New(agent)

	// Start a run
	runStarted := make(chan struct{}, 1)
	agent.onRun = func() {
		runStarted <- struct{}{}
	}
	app.submitInput("Hello")

	// Wait for run to start
	select {
	case <-runStarted:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run should have started")
	}

	// Cancel
	app.cancelRun()

	// Wait for cancel to be called
	select {
	case <-cancelCalled:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Agent.Cancel should have been called")
	}
}

func TestAppIsRunning(t *testing.T) {
	events := make(chan Event, 10)
	runStarted := make(chan struct{}, 1)
	agent := &mockAgentWithCancel{
		events: events,
		onRun: func() {
			runStarted <- struct{}{}
		},
	}

	app := New(agent)

	// Initially not running
	if app.isRunning() {
		t.Error("App should not be running initially")
	}

	// Start a run
	app.submitInput("Hello")

	// Wait for run to actually start
	select {
	case <-runStarted:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run should have started")
	}

	// Should be running
	if !app.isRunning() {
		t.Error("App should be running after submitInput")
	}
}

type mockAgentWithCancel struct {
	events   chan Event
	onRun    func()
	onCancel func()
}

func (m *mockAgentWithCancel) Run(ctx context.Context, prompt string) error {
	if m.onRun != nil {
		m.onRun()
	}
	<-ctx.Done()
	return ctx.Err()
}

func (m *mockAgentWithCancel) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgentWithCancel) Cancel() {
	if m.onCancel != nil {
		m.onCancel()
	}
}

func TestAppAccumulatesErrors(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Process error events
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 2")})

	if len(app.errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(app.errors))
	}
}

func TestAppClearsErrorsOnSuccess(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Simulate a previous run with errors
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})
	app.processEvent(Event{Type: EventComplete}) // Ends run with errors (errors preserved)

	// Verify errors are still there from previous run
	if len(app.errors) != 1 {
		t.Errorf("expected 1 error after error run, got %d", len(app.errors))
	}

	// Now a successful completion (no errors in this run) should clear errors
	app.processEvent(Event{Type: EventComplete})

	if len(app.errors) != 0 {
		t.Errorf("expected 0 errors after successful run, got %d", len(app.errors))
	}
}

func TestAppKeepsErrorsOnErrorComplete(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Process error then complete (simulating an error run)
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})
	// Note: errorsInRun flag would be set automatically by processEvent
	app.processEvent(Event{Type: EventComplete})

	// Errors should NOT be cleared if there were errors in this run
	if len(app.errors) != 1 {
		t.Errorf("expected 1 error preserved, got %d", len(app.errors))
	}
}
