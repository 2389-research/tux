// form/select.go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// Compile-time check that SelectField implements Field.
var _ Field = (*SelectField[string])(nil)

// SelectField is a single-choice selection field.
type SelectField[T comparable] struct {
	id         string
	label      string
	options    []SelectOption[T]
	value      T
	validators []Validator
	focused    bool
	huhField   *huh.Select[T]
	huhTheme   *huh.Theme
}

// NewSelect creates a new select field.
func NewSelect[T comparable]() *SelectField[T] {
	return &SelectField[T]{}
}

// Builder methods with With* prefix (matching InputField pattern)

// WithID sets the field ID (builder method).
func (f *SelectField[T]) WithID(id string) *SelectField[T] {
	f.id = id
	return f
}

// WithLabel sets the field label (builder method).
func (f *SelectField[T]) WithLabel(label string) *SelectField[T] {
	f.label = label
	return f
}

// WithOptions sets the available options (builder method).
func (f *SelectField[T]) WithOptions(options ...SelectOption[T]) *SelectField[T] {
	f.options = options
	return f
}

// WithDefault sets the default value (builder method).
func (f *SelectField[T]) WithDefault(value T) *SelectField[T] {
	f.value = value
	return f
}

// WithValidators adds validators to the field (builder method).
func (f *SelectField[T]) WithValidators(validators ...Validator) *SelectField[T] {
	f.validators = append(f.validators, validators...)
	return f
}

// Field interface implementation

// ID returns the field ID.
func (f *SelectField[T]) ID() string { return f.id }

// Label returns the field label.
func (f *SelectField[T]) Label() string { return f.label }

// Value returns the current value.
func (f *SelectField[T]) Value() any {
	return f.value
}

// SetValue sets the field value.
func (f *SelectField[T]) SetValue(v any) {
	if val, ok := v.(T); ok {
		f.value = val
	}
}

// Focused returns whether the field is focused.
func (f *SelectField[T]) Focused() bool {
	return f.focused
}

// Focus focuses the field.
func (f *SelectField[T]) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

// Blur removes focus from the field.
func (f *SelectField[T]) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

// Validate runs all validators on the current value.
func (f *SelectField[T]) Validate() error {
	return Compose(f.validators...)(f.value)
}

// Init initializes the underlying huh select.
func (f *SelectField[T]) Init() tea.Cmd {
	huhOpts := make([]huh.Option[T], len(f.options))
	for i, opt := range f.options {
		huhOpts[i] = huh.NewOption(opt.Label, opt.Value)
	}

	f.huhField = huh.NewSelect[T]().
		Title(f.label).
		Options(huhOpts...).
		Value(&f.value)

	return f.huhField.Init()
}

// HandleKey handles key input.
func (f *SelectField[T]) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

// Render renders the select field.
func (f *SelectField[T]) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}
