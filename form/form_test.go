// form/form_test.go
package form

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestFormCreation(t *testing.T) {
	f := New(
		NewInput().WithID("name").WithLabel("Name"),
		NewSelect[string]().WithID("role").WithLabel("Role").WithOptions(
			Option("Admin", "admin"),
			Option("User", "user"),
		),
	)

	if f == nil {
		t.Fatal("form should not be nil")
	}
}

func TestFormValues(t *testing.T) {
	nameField := NewInput().WithID("name")
	roleField := NewSelect[string]().WithID("role").WithOptions(
		Option("Admin", "admin"),
	).WithDefault("admin")

	f := New(nameField, roleField)
	nameField.SetValue("Alice")

	values := f.Values()
	if values.String("name") != "Alice" {
		t.Errorf("expected name 'Alice', got %s", values.String("name"))
	}
	if values.String("role") != "admin" {
		t.Errorf("expected role 'admin', got %s", values.String("role"))
	}
}

func TestFormNavigation(t *testing.T) {
	f := New(
		NewInput().WithID("a"),
		NewInput().WithID("b"),
		NewInput().WithID("c"),
	)
	f.Init()

	// First field should be focused
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0, got %d", f.FocusedIndex())
	}

	// Tab moves forward
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if f.FocusedIndex() != 1 {
		t.Errorf("expected focused index 1 after Tab, got %d", f.FocusedIndex())
	}

	// Shift+Tab moves back
	f.HandleKey(tea.KeyMsg{Type: tea.KeyShiftTab})
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0 after Shift+Tab, got %d", f.FocusedIndex())
	}
}

func TestFormState(t *testing.T) {
	f := New(NewInput().WithID("test"))
	f.Init()

	if f.State() != StateActive {
		t.Error("initial state should be Active")
	}

	// Escape cancels
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})
	if f.State() != StateCancelled {
		t.Error("state should be Cancelled after Escape")
	}
}

func TestFormRender(t *testing.T) {
	f := New(
		NewInput().WithID("name").WithLabel("Name"),
	).WithTheme(theme.Get("dracula"))
	f.Init()

	output := f.Render(60, 20)
	if output == "" {
		t.Error("render should produce output")
	}
}

func TestFormGroups(t *testing.T) {
	f := New(
		Group("Page 1",
			NewInput().WithID("a"),
		),
		Group("Page 2",
			NewInput().WithID("b"),
		),
	)
	f.Init()

	if f.GroupCount() != 2 {
		t.Errorf("expected 2 groups, got %d", f.GroupCount())
	}
	if f.CurrentGroup() != 0 {
		t.Error("should start on group 0")
	}
}
