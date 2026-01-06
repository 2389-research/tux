// form/input.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// Compile-time check that InputField implements Field.
var _ Field = (*InputField)(nil)

// InputField is a single-line text input.
type InputField struct {
	id          string
	label       string
	placeholder string
	value       string
	validators  []Validator
	focused     bool
	huhField    *huh.Input
	huhTheme    *huh.Theme
}

// NewInput creates a new input field.
func NewInput() *InputField {
	return &InputField{}
}

// WithID sets the field ID (builder method).
func (f *InputField) WithID(id string) *InputField {
	f.id = id
	return f
}

// WithLabel sets the field label (builder method).
func (f *InputField) WithLabel(label string) *InputField {
	f.label = label
	return f
}

// WithPlaceholder sets the placeholder text (builder method).
func (f *InputField) WithPlaceholder(placeholder string) *InputField {
	f.placeholder = placeholder
	return f
}

// WithValidators adds validators to the field (builder method).
func (f *InputField) WithValidators(validators ...Validator) *InputField {
	f.validators = append(f.validators, validators...)
	return f
}

// Field interface implementation

// ID returns the field ID.
func (f *InputField) ID() string {
	return f.id
}

// Label returns the field label.
func (f *InputField) Label() string {
	return f.label
}

// Value returns the current value.
func (f *InputField) Value() any {
	return f.value
}

// SetValue sets the field value.
func (f *InputField) SetValue(v any) {
	if s, ok := v.(string); ok {
		f.value = s
	}
}

// Focused returns whether the field is focused.
func (f *InputField) Focused() bool {
	return f.focused
}

// Focus focuses the field.
func (f *InputField) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

// Blur removes focus from the field.
func (f *InputField) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

// Validate runs all validators on the current value.
func (f *InputField) Validate() error {
	return Compose(f.validators...)(f.value)
}

// Init initializes the underlying huh input.
func (f *InputField) Init() tea.Cmd {
	f.huhField = huh.NewInput().
		Title(f.label).
		Placeholder(f.placeholder).
		Value(&f.value)

	return f.huhField.Init()
}

// HandleKey handles key input.
func (f *InputField) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

// Render renders the input field.
func (f *InputField) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}
