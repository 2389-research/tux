// form/integration_test.go
package form

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestFullFormWorkflow(t *testing.T) {
	var result Values

	// Create fields separately so we can set values
	usernameField := NewInput().WithID("username").WithLabel("Username").WithValidators(Required(), MinLength(3))
	roleField := NewSelect[string]().WithID("role").WithLabel("Role").WithOptions(
		Option("Admin", "admin"),
		Option("User", "user"),
	).WithDefault("user")
	activeField := NewConfirm().WithID("active").WithLabel("Active?")

	f := New(usernameField, roleField, activeField).
		WithTheme(theme.Get("neo-terminal")).
		OnSubmit(func(v Values) {
			result = v
		})

	// Set value before Init (like TestFormValues does)
	usernameField.SetValue("alice")

	f.Init()

	// Navigate and submit
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})   // to role
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})   // to confirm
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter}) // submit

	if f.State() != StateSubmitted {
		t.Errorf("expected StateSubmitted, got %v", f.State())
	}

	if result.String("username") != "alice" {
		t.Errorf("expected username 'alice', got %s", result.String("username"))
	}
	if result.String("role") != "user" {
		t.Errorf("expected role 'user', got %s", result.String("role"))
	}
}

func TestMultiPageForm(t *testing.T) {
	f := New(
		Group("Step 1",
			NewInput().WithID("name").WithLabel("Name"),
		),
		Group("Step 2",
			NewInput().WithID("email").WithLabel("Email"),
		),
		Group("Confirm",
			NewConfirm().WithID("agree").WithLabel("I agree"),
		),
	)

	f.Init()

	if f.GroupCount() != 3 {
		t.Errorf("expected 3 groups, got %d", f.GroupCount())
	}

	// Submit first page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.CurrentGroup() != 1 {
		t.Errorf("expected group 1, got %d", f.CurrentGroup())
	}

	// Submit second page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.CurrentGroup() != 2 {
		t.Errorf("expected group 2, got %d", f.CurrentGroup())
	}

	// Submit final page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.State() != StateSubmitted {
		t.Error("form should be submitted after last page")
	}
}
