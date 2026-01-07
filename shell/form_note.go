// form/note.go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/2389-research/tux/theme"
)

// Compile-time interface check
var _ Field = (*NoteField)(nil)

// NoteField is a display-only field for showing information.
type NoteField struct {
	id      string
	title   string
	content string
	focused bool
}

// NewNote creates a new note field.
func NewNote() *NoteField {
	return &NoteField{}
}

// Builder methods

func (f *NoteField) WithID(id string) *NoteField {
	f.id = id
	return f
}

func (f *NoteField) WithTitle(title string) *NoteField {
	f.title = title
	return f
}

func (f *NoteField) WithContent(content string) *NoteField {
	f.content = content
	return f
}

// Field interface implementation

func (f *NoteField) ID() string    { return f.id }
func (f *NoteField) Label() string { return f.title }

func (f *NoteField) Value() any {
	return nil // Notes have no value
}

func (f *NoteField) SetValue(v any) {
	// Notes don't store values
}

func (f *NoteField) Focused() bool {
	return f.focused
}

func (f *NoteField) Focus() {
	f.focused = true
}

func (f *NoteField) Blur() {
	f.focused = false
}

func (f *NoteField) Validate() error {
	return nil // Notes don't validate
}

func (f *NoteField) Init() tea.Cmd {
	return nil
}

func (f *NoteField) HandleKey(key tea.KeyMsg) bool {
	return false // Notes don't handle input
}

func (f *NoteField) Render(width int, th theme.Theme, focused bool) string {
	styles := th.Styles()

	var output string
	if f.title != "" {
		output = styles.Title.Render(f.title) + "\n"
	}
	if f.content != "" {
		output += styles.Body.Render(f.content)
	}

	return lipgloss.NewStyle().Width(width).Render(output)
}
