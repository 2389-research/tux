// form/form.go
package form

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/2389-research/tux/theme"
)

// State represents the form state.
type State int

const (
	StateActive State = iota
	StateSubmitted
	StateCancelled
)

// FieldGroup represents a group of fields (a page in multi-page forms).
type FieldGroup struct {
	title  string
	fields []Field
}

// Group creates a field group with a title.
func Group(title string, fields ...Field) *FieldGroup {
	return &FieldGroup{
		title:  title,
		fields: fields,
	}
}

// Form manages a collection of fields.
type Form struct {
	groups       []*FieldGroup
	currentGroup int
	focusedIndex int
	state        State
	theme        theme.Theme
	onSubmit     func(Values)
	onCancel     func()
}

// New creates a new form. Accepts fields or groups.
func New(items ...any) *Form {
	f := &Form{
		state: StateActive,
		theme: theme.Get("dracula"), // default theme
	}

	// Collect fields into a default group, or use provided groups
	var defaultFields []Field
	for _, item := range items {
		switch v := item.(type) {
		case *FieldGroup:
			f.groups = append(f.groups, v)
		case Field:
			defaultFields = append(defaultFields, v)
		}
	}

	// If we have loose fields, put them in a default group
	if len(defaultFields) > 0 {
		f.groups = append([]*FieldGroup{{fields: defaultFields}}, f.groups...)
	}

	// Ensure at least one group
	if len(f.groups) == 0 {
		f.groups = []*FieldGroup{{}}
	}

	return f
}

// WithTheme sets the form theme.
func (f *Form) WithTheme(th theme.Theme) *Form {
	f.theme = th
	return f
}

// OnSubmit sets the submit callback.
func (f *Form) OnSubmit(fn func(Values)) *Form {
	f.onSubmit = fn
	return f
}

// OnCancel sets the cancel callback.
func (f *Form) OnCancel(fn func()) *Form {
	f.onCancel = fn
	return f
}

// Init initializes the form and focuses the first field.
func (f *Form) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, g := range f.groups {
		for _, field := range g.fields {
			if cmd := field.Init(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Focus first field
	if fields := f.currentFields(); len(fields) > 0 {
		fields[0].Focus()
	}

	return tea.Batch(cmds...)
}

// State returns the current form state.
func (f *Form) State() State {
	return f.state
}

// Values returns all field values.
func (f *Form) Values() Values {
	v := make(Values)
	for _, g := range f.groups {
		for _, field := range g.fields {
			if id := field.ID(); id != "" {
				v[id] = field.Value()
			}
		}
	}
	return v
}

// FocusedIndex returns the index of the focused field in the current group.
func (f *Form) FocusedIndex() int {
	return f.focusedIndex
}

// GroupCount returns the number of groups.
func (f *Form) GroupCount() int {
	return len(f.groups)
}

// CurrentGroup returns the current group index.
func (f *Form) CurrentGroup() int {
	return f.currentGroup
}

func (f *Form) currentFields() []Field {
	if f.currentGroup < len(f.groups) {
		return f.groups[f.currentGroup].fields
	}
	return nil
}

// HandleKey processes keyboard input.
func (f *Form) HandleKey(key tea.KeyMsg) bool {
	if f.state != StateActive {
		return false
	}

	fields := f.currentFields()

	switch key.Type {
	case tea.KeyEscape:
		f.state = StateCancelled
		if f.onCancel != nil {
			f.onCancel()
		}
		return true

	case tea.KeyTab, tea.KeyDown:
		if len(fields) > 0 {
			fields[f.focusedIndex].Blur()
			f.focusedIndex = (f.focusedIndex + 1) % len(fields)
			fields[f.focusedIndex].Focus()
		}
		return true

	case tea.KeyShiftTab, tea.KeyUp:
		if len(fields) > 0 {
			fields[f.focusedIndex].Blur()
			f.focusedIndex--
			if f.focusedIndex < 0 {
				f.focusedIndex = len(fields) - 1
			}
			fields[f.focusedIndex].Focus()
		}
		return true

	case tea.KeyEnter:
		// On last field, submit or go to next group
		if f.focusedIndex == len(fields)-1 {
			if f.currentGroup == len(f.groups)-1 {
				// Last group - submit
				f.state = StateSubmitted
				if f.onSubmit != nil {
					f.onSubmit(f.Values())
				}
			} else {
				// Next group
				fields[f.focusedIndex].Blur()
				f.currentGroup++
				f.focusedIndex = 0
				if newFields := f.currentFields(); len(newFields) > 0 {
					newFields[0].Focus()
				}
			}
			return true
		}
		// Move to next field
		if len(fields) > 0 {
			fields[f.focusedIndex].Blur()
			f.focusedIndex++
			fields[f.focusedIndex].Focus()
		}
		return true
	}

	// Delegate to focused field
	if len(fields) > 0 && f.focusedIndex < len(fields) {
		return fields[f.focusedIndex].HandleKey(key)
	}

	return false
}

// Render renders the form.
func (f *Form) Render(width, height int) string {
	if f.theme == nil {
		f.theme = theme.Get("dracula")
	}

	fields := f.currentFields()
	var parts []string

	// Group title
	if f.currentGroup < len(f.groups) && f.groups[f.currentGroup].title != "" {
		title := f.theme.Styles().Title.Render(f.groups[f.currentGroup].title)
		parts = append(parts, title, "")
	}

	// Fields
	for i, field := range fields {
		focused := i == f.focusedIndex
		parts = append(parts, field.Render(width-4, f.theme, focused))
	}

	// Page indicator for multi-group forms
	if len(f.groups) > 1 {
		indicator := f.theme.Styles().Muted.Render(
			strings.Repeat("○ ", f.currentGroup) + "● " + strings.Repeat("○ ", len(f.groups)-f.currentGroup-1),
		)
		parts = append(parts, "", indicator)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
