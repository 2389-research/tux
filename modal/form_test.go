// modal/form_test.go
package modal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/form"
)

func TestFormModal(t *testing.T) {
	f := form.New(
		form.NewInput().WithID("name").WithLabel("Name"),
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
	f := form.New(form.NewInput().WithID("test"))

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
	var values form.Values

	f := form.New(form.NewInput().WithID("name"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		OnSubmit: func(v form.Values) {
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

	f := form.New(form.NewInput().WithID("name"))

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
	f := form.New(form.NewInput().WithID("name").WithLabel("Name"))

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
