package shell

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestErrorModalRendersErrors(t *testing.T) {
	errs := []error{
		errors.New("connection timeout"),
		errors.New("rate limit exceeded"),
	}
	modal := NewErrorModal(ErrorModalConfig{
		Errors: errs,
		Theme:  theme.NewDraculaTheme(),
	})

	view := modal.Render(60, 20)
	if !strings.Contains(view, "connection timeout") {
		t.Error("modal should show first error")
	}
	if !strings.Contains(view, "rate limit exceeded") {
		t.Error("modal should show second error")
	}
}

func TestErrorModalID(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	if modal.ID() != "error-modal" {
		t.Errorf("expected id 'error-modal', got %q", modal.ID())
	}
}

func TestErrorModalTitle(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	if modal.Title() != "Errors" {
		t.Errorf("expected title 'Errors', got %q", modal.Title())
	}
}

func TestErrorModalEscapeCloses(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	handled, cmd := modal.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})
	if !handled {
		t.Error("escape should be handled")
	}
	if cmd == nil {
		t.Error("escape should return command")
	}
	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("expected PopMsg, got %T", msg)
	}
}

func TestErrorModalCtrlECloses(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	handled, cmd := modal.HandleKey(tea.KeyMsg{Type: tea.KeyCtrlE})
	if !handled {
		t.Error("ctrl+e should be handled")
	}
	if cmd == nil {
		t.Error("ctrl+e should return command")
	}
	// Verify the command produces PopMsg (same as Escape)
	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("expected PopMsg, got %T", msg)
	}
}
