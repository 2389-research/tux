package shell

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestInputHistoryNavigation(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	// Set up history provider
	history := []string{"first", "second", "third"}
	input.SetHistoryProvider(func() []string {
		return history
	})

	// Press up arrow - should show "third" (most recent)
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "third" {
		t.Errorf("expected 'third', got %q", input.Value())
	}

	// Press up again - should show "second"
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "second" {
		t.Errorf("expected 'second', got %q", input.Value())
	}

	// Press down - should show "third"
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyDown})
	if input.Value() != "third" {
		t.Errorf("expected 'third', got %q", input.Value())
	}

	// Press down past end - should clear
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyDown})
	if input.Value() != "" {
		t.Errorf("expected empty, got %q", input.Value())
	}
}

func TestInputHistoryResetsOnSubmit(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	history := []string{"first", "second"}
	input.SetHistoryProvider(func() []string {
		return history
	})

	// Navigate up
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Type something and submit
	input.SetValue("new message")
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// History index should be reset - pressing up should show "second" again
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "second" {
		t.Errorf("expected 'second' after reset, got %q", input.Value())
	}
}

func TestInputNoHistoryProvider(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	// Up arrow with no history provider should do nothing
	input.SetValue("test")
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "test" {
		t.Errorf("expected 'test' unchanged, got %q", input.Value())
	}
}
