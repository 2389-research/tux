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

func TestFormOnCancel(t *testing.T) {
	cancelCalled := false
	f := New(NewInput().WithID("test")).
		OnCancel(func() {
			cancelCalled = true
		})
	f.Init()

	// Press Escape to cancel
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})

	if !cancelCalled {
		t.Error("OnCancel callback should have been called")
	}
	if f.State() != StateCancelled {
		t.Error("state should be Cancelled after Escape")
	}
}

func TestFormOnSubmit(t *testing.T) {
	submitCalled := false
	var submittedValues Values
	f := New(NewInput().WithID("test")).
		OnSubmit(func(v Values) {
			submitCalled = true
			submittedValues = v
		})
	f.Init()

	// Set a value
	f.Values()["test"] = "hello"

	// Press Enter on last field to submit
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if !submitCalled {
		t.Error("OnSubmit callback should have been called")
	}
	if f.State() != StateSubmitted {
		t.Error("state should be Submitted after Enter on last field")
	}
	if submittedValues == nil {
		t.Error("submittedValues should not be nil")
	}
}

func TestFormHandleKeyInactive(t *testing.T) {
	f := New(NewInput().WithID("test"))
	f.Init()

	// Cancel the form first
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})

	// Now HandleKey should return false for any key
	result := f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if result {
		t.Error("HandleKey should return false when form is not active")
	}
}

func TestFormMultiGroupNavigation(t *testing.T) {
	f := New(
		Group("Page 1",
			NewInput().WithID("a"),
		),
		Group("Page 2",
			NewInput().WithID("b"),
		),
	)
	f.Init()

	// Press Enter on last field of first group to move to next group
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if f.CurrentGroup() != 1 {
		t.Errorf("expected group 1, got %d", f.CurrentGroup())
	}
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0, got %d", f.FocusedIndex())
	}
}

func TestFormRenderWithGroupTitle(t *testing.T) {
	f := New(
		Group("Page 1",
			NewInput().WithID("a"),
		),
	).WithTheme(theme.Get("dracula"))
	f.Init()

	output := f.Render(60, 20)
	if output == "" {
		t.Error("render should produce output")
	}
}

func TestFormRenderMultiPage(t *testing.T) {
	f := New(
		Group("Page 1",
			NewInput().WithID("a"),
		),
		Group("Page 2",
			NewInput().WithID("b"),
		),
	).WithTheme(theme.Get("dracula"))
	f.Init()

	output := f.Render(60, 20)
	if output == "" {
		t.Error("render should produce output")
	}
}

func TestFormEmpty(t *testing.T) {
	// Create form with no fields
	f := New()
	f.Init()

	// HandleKey should work without panicking
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
}

func TestFormDownUpKeys(t *testing.T) {
	f := New(
		NewInput().WithID("a"),
		NewInput().WithID("b"),
	)
	f.Init()

	// Down should move forward
	f.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if f.FocusedIndex() != 1 {
		t.Errorf("expected focused index 1 after Down, got %d", f.FocusedIndex())
	}

	// Up should move back
	f.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0 after Up, got %d", f.FocusedIndex())
	}
}

func TestFormEnterNavigatesFields(t *testing.T) {
	f := New(
		NewInput().WithID("a"),
		NewInput().WithID("b"),
		NewInput().WithID("c"),
	)
	f.Init()

	// Enter on first field moves to second
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.FocusedIndex() != 1 {
		t.Errorf("expected focused index 1 after Enter, got %d", f.FocusedIndex())
	}
}

func TestFormDelegatesKeyToField(t *testing.T) {
	f := New(
		NewInput().WithID("a").WithLabel("Test"),
	).WithTheme(theme.Get("dracula"))
	f.Init()

	// Render first to initialize the huhField
	f.Render(60, 20)

	// Type a character - should delegate to field
	key := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	f.HandleKey(key)
}
