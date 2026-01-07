// form/field.go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

// Field is the interface for form fields.
type Field interface {
	// Identity
	ID() string
	Label() string

	// State
	Value() any
	SetValue(any)
	Focused() bool
	Focus()
	Blur()

	// Validation
	Validate() error

	// Rendering
	Render(width int, th theme.Theme, focused bool) string

	// Input handling
	HandleKey(key tea.KeyMsg) bool

	// Lifecycle
	Init() tea.Cmd
}

// SelectOption represents an option in a select field.
type SelectOption[T any] struct {
	Label string
	Value T
}

// Option creates a SelectOption.
func Option[T any](label string, value T) SelectOption[T] {
	return SelectOption[T]{Label: label, Value: value}
}

// Values holds form field values by ID.
type Values map[string]any

// String returns a string value or empty string if not found/wrong type.
func (v Values) String(id string) string {
	if val, ok := v[id].(string); ok {
		return val
	}
	return ""
}

// Bool returns a bool value or false if not found/wrong type.
func (v Values) Bool(id string) bool {
	if val, ok := v[id].(bool); ok {
		return val
	}
	return false
}

// Strings returns a []string value or nil if not found/wrong type.
func (v Values) Strings(id string) []string {
	if val, ok := v[id].([]string); ok {
		return val
	}
	return nil
}

// Int returns an int value or 0 if not found/wrong type.
func (v Values) Int(id string) int {
	if val, ok := v[id].(int); ok {
		return val
	}
	return 0
}
