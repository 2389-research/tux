// modal/form_test.go
package shell

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFormModal(t *testing.T) {
	f := NewForm(
		NewInputField().WithID("name").WithLabel("Name"),
	)

	m := NewFormModal(FormModalConfig{
		ID:    "test-form",
		Title: "Test Form",
		Form:  f,
	})

	if m.ID() != "test-form" {
		t.Errorf("expected ID 'test-form', got %s", m.ID())
	}
	if m.Title() != "Test Form" {
		t.Errorf("expected title 'Test Form', got %s", m.Title())
	}
}

func TestFormModalSize(t *testing.T) {
	f := NewForm(NewInputField().WithID("test"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		Size: SizeLarge,
	})

	if m.Size() != SizeLarge {
		t.Errorf("expected SizeLarge, got %v", m.Size())
	}
}

func TestFormModalOnSubmit(t *testing.T) {
	var submitted bool
	var values Values

	f := NewForm(NewInputField().WithID("name"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		OnSubmit: func(v Values) {
			submitted = true
			values = v
		},
	})

	m.OnPush(80, 24)

	// Simulate Enter to submit (single field form)
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if !submitted {
		t.Error("OnSubmit should have been called")
	}
	if values == nil {
		t.Error("values should not be nil")
	}
}

func TestFormModalOnCancel(t *testing.T) {
	var cancelled bool

	f := NewForm(NewInputField().WithID("name"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		OnCancel: func() {
			cancelled = true
		},
	})

	m.OnPush(80, 24)

	// Simulate Escape to cancel
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})

	if !cancelled {
		t.Error("OnCancel should have been called")
	}
}

func TestFormModalRender(t *testing.T) {
	f := NewForm(NewInputField().WithID("name").WithLabel("Name"))

	m := NewFormModal(FormModalConfig{
		Title: "Test",
		Form:  f,
	})

	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestFormModalDefaultID(t *testing.T) {
	f := NewForm(NewInputField().WithID("test"))
	m := NewFormModal(FormModalConfig{
		Form: f,
		// No ID set
	})
	if m.ID() != "form-modal" {
		t.Errorf("expected default 'form-modal', got %s", m.ID())
	}
}

func TestFormModalOnPop(t *testing.T) {
	f := NewForm(NewInputField().WithID("test"))
	m := NewFormModal(FormModalConfig{Form: f})
	m.OnPop() // Should not panic
}

func TestFormModalHandleKeyNilForm(t *testing.T) {
	m := NewFormModal(FormModalConfig{})
	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if handled || cmd != nil {
		t.Error("nil form should return false, nil")
	}
}

func TestFormModalHandleKeyNoStateChange(t *testing.T) {
	// Create a form with multiple fields so a single key doesn't complete it
	f := NewForm(
		NewInputField().WithID("field1").WithLabel("Field 1"),
		NewInputField().WithID("field2").WithLabel("Field 2"),
	)

	m := NewFormModal(FormModalConfig{
		Form: f,
	})

	m.OnPush(80, 24)

	// Send a regular character key that won't complete the form
	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	// The key should be handled by the form but shouldn't return PopMsg
	if !handled {
		t.Error("key should be handled by form")
	}
	if cmd != nil {
		t.Error("cmd should be nil when form state doesn't change")
	}
}
