package theme

import "github.com/charmbracelet/lipgloss"

// Dracula color palette
const (
	draculaBackground  = lipgloss.Color("#282a36")
	draculaCurrentLine = lipgloss.Color("#44475a")
	draculaForeground  = lipgloss.Color("#f8f8f2")
	draculaComment     = lipgloss.Color("#6272a4")
	draculaCyan        = lipgloss.Color("#8be9fd")
	draculaGreen       = lipgloss.Color("#50fa7b")
	draculaOrange      = lipgloss.Color("#ffb86c")
	draculaPink        = lipgloss.Color("#ff79c6")
	draculaPurple      = lipgloss.Color("#bd93f9")
	draculaRed         = lipgloss.Color("#ff5555")
	draculaYellow      = lipgloss.Color("#f1fa8c")
)

type draculaTheme struct {
	styles Styles
}

// NewDraculaTheme creates a new Dracula color theme.
func NewDraculaTheme() Theme {
	t := &draculaTheme{}
	t.styles = t.buildStyles()
	return t
}

func (t *draculaTheme) Name() string            { return "dracula" }
func (t *draculaTheme) Background() lipgloss.Color  { return draculaBackground }
func (t *draculaTheme) Foreground() lipgloss.Color  { return draculaForeground }
func (t *draculaTheme) Primary() lipgloss.Color     { return draculaPurple }
func (t *draculaTheme) Secondary() lipgloss.Color   { return draculaCyan }
func (t *draculaTheme) Success() lipgloss.Color     { return draculaGreen }
func (t *draculaTheme) Warning() lipgloss.Color     { return draculaOrange }
func (t *draculaTheme) Error() lipgloss.Color       { return draculaRed }
func (t *draculaTheme) Info() lipgloss.Color        { return draculaCyan }
func (t *draculaTheme) Border() lipgloss.Color      { return draculaComment }
func (t *draculaTheme) BorderFocused() lipgloss.Color { return draculaPurple }
func (t *draculaTheme) Muted() lipgloss.Color       { return draculaComment }
func (t *draculaTheme) UserColor() lipgloss.Color   { return draculaOrange }
func (t *draculaTheme) AssistantColor() lipgloss.Color { return draculaGreen }
func (t *draculaTheme) ToolColor() lipgloss.Color   { return draculaCyan }
func (t *draculaTheme) SystemColor() lipgloss.Color { return draculaComment }
func (t *draculaTheme) Styles() Styles              { return t.styles }

func (t *draculaTheme) buildStyles() Styles {
	return Styles{
		// Text
		Title: lipgloss.NewStyle().
			Foreground(draculaPurple).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(draculaCyan),
		Body: lipgloss.NewStyle().
			Foreground(draculaForeground),
		Muted: lipgloss.NewStyle().
			Foreground(draculaComment),
		Emphasized: lipgloss.NewStyle().
			Foreground(draculaForeground).
			Bold(true),

		// Status
		Success: lipgloss.NewStyle().
			Foreground(draculaGreen),
		Error: lipgloss.NewStyle().
			Foreground(draculaRed),
		Warning: lipgloss.NewStyle().
			Foreground(draculaOrange),
		Info: lipgloss.NewStyle().
			Foreground(draculaCyan),

		// Interactive
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(draculaComment),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(draculaPurple),
		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(draculaComment).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(draculaPurple).
			Padding(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(draculaForeground).
			Background(draculaCurrentLine).
			Padding(0, 2),
		ButtonActive: lipgloss.NewStyle().
			Foreground(draculaBackground).
			Background(draculaPurple).
			Padding(0, 2),

		// Layout
		StatusBar: lipgloss.NewStyle().
			Foreground(draculaForeground).
			Background(draculaCurrentLine),
		TabBar: lipgloss.NewStyle().
			Foreground(draculaForeground),
		TabActive: lipgloss.NewStyle().
			Foreground(draculaPurple).
			Bold(true).
			Underline(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(draculaComment),

		// Modal
		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(draculaPurple).
			Padding(1, 2),
		ModalTitle: lipgloss.NewStyle().
			Foreground(draculaPurple).
			Bold(true),
		ModalFooter: lipgloss.NewStyle().
			Foreground(draculaComment),

		// Tool states
		ToolApproval: lipgloss.NewStyle().
			Foreground(draculaOrange),
		ToolExecuting: lipgloss.NewStyle().
			Foreground(draculaYellow),
		ToolSuccess: lipgloss.NewStyle().
			Foreground(draculaGreen),
		ToolError: lipgloss.NewStyle().
			Foreground(draculaRed),
		ToolPending: lipgloss.NewStyle().
			Foreground(draculaComment),

		// List
		ListItem: lipgloss.NewStyle().
			Foreground(draculaForeground),
		ListItemSelected: lipgloss.NewStyle().
			Foreground(draculaGreen).
			Bold(true),

		// Help
		HelpKey: lipgloss.NewStyle().
			Foreground(draculaCyan).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(draculaComment),
	}
}
