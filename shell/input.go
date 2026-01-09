package shell

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input is the text input component.
type Input struct {
	model       textinput.Model
	theme       theme.Theme
	prefix      string
	placeholder string
	width       int
}

// NewInput creates a new input component.
func NewInput(th theme.Theme, prefix, placeholder string) *Input {
	ti := textinput.New()
	ti.Prompt = prefix
	ti.Placeholder = placeholder
	ti.Focus()

	styles := th.Styles()
	ti.PromptStyle = styles.Muted
	ti.TextStyle = styles.Body

	return &Input{
		model:       ti,
		theme:       th,
		prefix:      prefix,
		placeholder: placeholder,
	}
}

// Init initializes the input.
func (i *Input) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input messages.
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			value := i.model.Value()
			if value != "" {
				i.model.SetValue("")
				return i, func() tea.Msg {
					return InputSubmitMsg{Value: value}
				}
			}
			return i, nil
		}
	}

	var cmd tea.Cmd
	i.model, cmd = i.model.Update(msg)
	return i, cmd
}

// View renders the input.
func (i *Input) View() string {
	styles := i.theme.Styles()
	return styles.Input.Width(i.width - 4).Render(i.model.View())
}

// Value returns the current input value.
func (i *Input) Value() string {
	return i.model.Value()
}

// SetValue sets the input value.
func (i *Input) SetValue(value string) {
	i.model.SetValue(value)
}

// SetWidth sets the input width.
func (i *Input) SetWidth(width int) {
	i.width = width
	i.model.Width = width - 4 - len(i.prefix)
}

// Focus focuses the input.
func (i *Input) Focus() tea.Cmd {
	return i.model.Focus()
}

// Blur blurs the input.
func (i *Input) Blur() {
	i.model.Blur()
}

// Focused returns whether the input is focused.
func (i *Input) Focused() bool {
	return i.model.Focused()
}

// InputSubmitMsg is sent when the user submits input.
type InputSubmitMsg struct {
	Value string
}
