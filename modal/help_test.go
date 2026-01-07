// modal/help_test.go
package modal

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/2389-research/tux/help"
)

func TestHelpModal(t *testing.T) {
	h := help.New(
		help.Category{
			Title:    "Navigation",
			Bindings: []help.Binding{{Key: "j", Description: "Move down"}},
		},
	)

	m := NewHelpModal(HelpModalConfig{
		ID:    "test-help",
		Title: "Test Help",
		Help:  h,
	})

	if m.ID() != "test-help" {
		t.Errorf("expected ID 'test-help', got %s", m.ID())
	}
	if m.Title() != "Test Help" {
		t.Errorf("expected title 'Test Help', got %s", m.Title())
	}
}

func TestHelpModalDefaults(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{
		Help: h,
	})

	// Default ID
	if m.ID() != "help-modal" {
		t.Errorf("expected default ID 'help-modal', got %s", m.ID())
	}

	// Default title
	if m.Title() != "Help" {
		t.Errorf("expected default title 'Help', got %s", m.Title())
	}

	// Default size
	if m.Size() != SizeMedium {
		t.Errorf("expected default size SizeMedium, got %v", m.Size())
	}
}

func TestHelpModalSize(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{
		Help: h,
		Size: SizeLarge,
	})

	if m.Size() != SizeLarge {
		t.Errorf("expected SizeLarge, got %v", m.Size())
	}
}

func TestHelpModalOnPush(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})

	m.OnPush(80, 24)

	// OnPush should store dimensions (verify via render using them)
	// This is a basic test - the dimensions are stored internally
}

func TestHelpModalOnPop(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPop() // Should not panic
}

func TestHelpModalHandleKeyEscape(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)

	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})

	if !handled {
		t.Error("Escape should be handled")
	}
	if cmd == nil {
		t.Error("Escape should return a command")
	}

	// Execute the command and verify it returns PopMsg
	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("expected PopMsg, got %T", msg)
	}
}

func TestHelpModalHandleKeyQuestionMark(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)

	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

	if !handled {
		t.Error("? should be handled (toggle close)")
	}
	if cmd == nil {
		t.Error("? should return a command")
	}

	// Execute the command and verify it returns PopMsg
	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("expected PopMsg, got %T", msg)
	}
}

func TestHelpModalHandleKeyOther(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)

	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	if handled {
		t.Error("other keys should not be handled")
	}
	if cmd != nil {
		t.Error("other keys should return nil cmd")
	}
}

func TestHelpModalRender(t *testing.T) {
	h := help.New(
		help.Category{
			Title: "Navigation",
			Bindings: []help.Binding{
				{Key: "j", Description: "Move down"},
				{Key: "k", Description: "Move up"},
			},
		},
	)

	m := NewHelpModal(HelpModalConfig{
		Help: h,
	})

	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render should produce output")
	}

	// Verify help content is present
	if !strings.Contains(output, "j") {
		t.Error("render should contain binding key 'j'")
	}
	if !strings.Contains(output, "Move down") {
		t.Error("render should contain binding description")
	}
}

func TestHelpModalRenderWithMode(t *testing.T) {
	h := help.New(
		help.Category{
			Title:    "Navigation",
			Bindings: []help.Binding{{Key: "j", Description: "Move down", Modes: []string{"normal"}}},
		},
	)

	m := NewHelpModal(HelpModalConfig{
		Help: h,
		Mode: "normal",
	})

	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestHelpModalNilHelp(t *testing.T) {
	m := NewHelpModal(HelpModalConfig{})

	// Should not panic with nil help
	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render with nil help should still produce output")
	}
}
