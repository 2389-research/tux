// tux_test.go
package tux

import (
	"context"
	"testing"

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
