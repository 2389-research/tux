// form/confirm.go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// Compile-time interface check
var _ Field = (*ConfirmField)(nil)

// ConfirmField is a yes/no confirmation field.
type ConfirmField struct {
	id          string
	label       string
	affirmative string
	negative    string
	value       bool
	focused     bool
	huhField    *huh.Confirm
	huhTheme    *huh.Theme
}

// NewConfirm creates a new confirm field.
func NewConfirm() *ConfirmField {
	return &ConfirmField{
		affirmative: "Yes",
		negative:    "No",
	}
}

// Builder methods with With* prefix

// WithID sets the field ID (builder method).
func (f *ConfirmField) WithID(id string) *ConfirmField {
	f.id = id
	return f
}

// WithLabel sets the field label (builder method).
func (f *ConfirmField) WithLabel(label string) *ConfirmField {
	f.label = label
	return f
}

// WithAffirmative sets the affirmative button text (builder method).
func (f *ConfirmField) WithAffirmative(text string) *ConfirmField {
	f.affirmative = text
	return f
}

// WithNegative sets the negative button text (builder method).
func (f *ConfirmField) WithNegative(text string) *ConfirmField {
	f.negative = text
	return f
}

// Field interface implementation

// ID returns the field ID.
func (f *ConfirmField) ID() string { return f.id }

// Label returns the field label.
func (f *ConfirmField) Label() string { return f.label }

// Value returns the current value.
func (f *ConfirmField) Value() any {
	return f.value
}

// SetValue sets the field value.
func (f *ConfirmField) SetValue(v any) {
	if b, ok := v.(bool); ok {
		f.value = b
	}
}

// Focused returns whether the field is focused.
func (f *ConfirmField) Focused() bool {
	return f.focused
}

// Focus focuses the field.
func (f *ConfirmField) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

// Blur removes focus from the field.
func (f *ConfirmField) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

// Validate runs validation on the field.
func (f *ConfirmField) Validate() error {
	return nil // Confirm has no validation
}

// Init initializes the underlying huh confirm.
func (f *ConfirmField) Init() tea.Cmd {
	f.huhField = huh.NewConfirm().
		Title(f.label).
		Affirmative(f.affirmative).
		Negative(f.negative).
		Value(&f.value)

	return f.huhField.Init()
}

// HandleKey handles key input.
func (f *ConfirmField) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

// Render renders the confirm field.
func (f *ConfirmField) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}
