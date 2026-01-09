package shell

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input is the text input component.
type Input struct {
	model           textinput.Model
	theme           theme.Theme
	prefix          string
	placeholder     string
	width           int
	historyProvider func() []string
	historyIndex    int // -1 means not navigating history
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
		model:        ti,
		theme:        th,
		prefix:       prefix,
		placeholder:  placeholder,
		historyIndex: -1,
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
		switch msg.Type {
		case tea.KeyEnter:
			value := i.model.Value()
			if value != "" {
				i.model.SetValue("")
				i.historyIndex = -1 // Reset history navigation
				return i, func() tea.Msg {
					return InputSubmitMsg{Value: value}
				}
			}
			return i, nil

		case tea.KeyUp:
			if i.historyProvider != nil {
				history := i.historyProvider()
				if len(history) > 0 {
					if i.historyIndex == -1 {
						// Start navigating from end
						i.historyIndex = len(history) - 1
					} else if i.historyIndex > 0 {
						i.historyIndex--
					}
					i.model.SetValue(history[i.historyIndex])
					i.model.CursorEnd()
				}
			}
			return i, nil

		case tea.KeyDown:
			if i.historyProvider != nil && i.historyIndex != -1 {
				history := i.historyProvider()
				if i.historyIndex < len(history)-1 {
					i.historyIndex++
					i.model.SetValue(history[i.historyIndex])
					i.model.CursorEnd()
				} else {
					// Past end - clear and reset
					i.historyIndex = -1
					i.model.SetValue("")
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

// SetHistoryProvider sets the function that provides history items.
// History should be ordered oldest to newest.
func (i *Input) SetHistoryProvider(provider func() []string) {
	i.historyProvider = provider
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
