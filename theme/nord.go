package theme

import "github.com/charmbracelet/lipgloss"

// Nord color palette (Polar Night, Snow Storm, Frost, Aurora)
const (
	// Polar Night
	nordNight0 = lipgloss.Color("#2e3440")
	nordNight1 = lipgloss.Color("#3b4252")
	nordNight2 = lipgloss.Color("#434c5e")
	nordNight3 = lipgloss.Color("#4c566a")

	// Snow Storm
	nordSnow0 = lipgloss.Color("#d8dee9")
	nordSnow1 = lipgloss.Color("#e5e9f0")
	nordSnow2 = lipgloss.Color("#eceff4")

	// Frost
	nordFrost0 = lipgloss.Color("#8fbcbb")
	nordFrost1 = lipgloss.Color("#88c0d0")
	nordFrost2 = lipgloss.Color("#81a1c1")
	nordFrost3 = lipgloss.Color("#5e81ac")

	// Aurora
	nordRed    = lipgloss.Color("#bf616a")
	nordOrange = lipgloss.Color("#d08770")
	nordYellow = lipgloss.Color("#ebcb8b")
	nordGreen  = lipgloss.Color("#a3be8c")
	nordPurple = lipgloss.Color("#b48ead")
)

type nordTheme struct {
	styles Styles
}

// NewNordTheme creates a new Nord color theme.
func NewNordTheme() Theme {
	t := &nordTheme{}
	t.styles = t.buildStyles()
	return t
}

func (t *nordTheme) Name() string                   { return "nord" }
func (t *nordTheme) Background() lipgloss.Color     { return nordNight0 }
func (t *nordTheme) Foreground() lipgloss.Color     { return nordSnow2 }
func (t *nordTheme) Primary() lipgloss.Color        { return nordFrost1 }
func (t *nordTheme) Secondary() lipgloss.Color      { return nordFrost2 }
func (t *nordTheme) Success() lipgloss.Color        { return nordGreen }
func (t *nordTheme) Warning() lipgloss.Color        { return nordYellow }
func (t *nordTheme) Error() lipgloss.Color          { return nordRed }
func (t *nordTheme) Info() lipgloss.Color           { return nordFrost1 }
func (t *nordTheme) Border() lipgloss.Color         { return nordNight3 }
func (t *nordTheme) BorderFocused() lipgloss.Color  { return nordFrost1 }
func (t *nordTheme) Muted() lipgloss.Color          { return nordNight3 }
func (t *nordTheme) UserColor() lipgloss.Color      { return nordOrange }
func (t *nordTheme) AssistantColor() lipgloss.Color { return nordGreen }
func (t *nordTheme) ToolColor() lipgloss.Color      { return nordFrost0 }
func (t *nordTheme) SystemColor() lipgloss.Color    { return nordNight3 }
func (t *nordTheme) Styles() Styles                 { return t.styles }

func (t *nordTheme) buildStyles() Styles {
	return Styles{
		// Text
		Title: lipgloss.NewStyle().
			Foreground(nordFrost1).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(nordFrost2),
		Body: lipgloss.NewStyle().
			Foreground(nordSnow2),
		Muted: lipgloss.NewStyle().
			Foreground(nordNight3),
		Emphasized: lipgloss.NewStyle().
			Foreground(nordSnow2).
			Bold(true),

		// Status
		Success: lipgloss.NewStyle().
			Foreground(nordGreen),
		Error: lipgloss.NewStyle().
			Foreground(nordRed),
		Warning: lipgloss.NewStyle().
			Foreground(nordYellow),
		Info: lipgloss.NewStyle().
			Foreground(nordFrost1),

		// Interactive
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nordNight3),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nordFrost1),
		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nordNight3).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nordFrost1).
			Padding(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(nordSnow2).
			Background(nordNight2).
			Padding(0, 2),
		ButtonActive: lipgloss.NewStyle().
			Foreground(nordNight0).
			Background(nordFrost1).
			Padding(0, 2),

		// Layout
		StatusBar: lipgloss.NewStyle().
			Foreground(nordSnow2).
			Background(nordNight1),
		TabBar: lipgloss.NewStyle().
			Foreground(nordSnow2),
		TabActive: lipgloss.NewStyle().
			Foreground(nordFrost1).
			Bold(true).
			Underline(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(nordNight3),

		// Modal
		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nordFrost1).
			Padding(1, 2),
		ModalTitle: lipgloss.NewStyle().
			Foreground(nordFrost1).
			Bold(true),
		ModalFooter: lipgloss.NewStyle().
			Foreground(nordNight3),

		// Tool states
		ToolApproval: lipgloss.NewStyle().
			Foreground(nordOrange),
		ToolExecuting: lipgloss.NewStyle().
			Foreground(nordYellow),
		ToolSuccess: lipgloss.NewStyle().
			Foreground(nordGreen),
		ToolError: lipgloss.NewStyle().
			Foreground(nordRed),
		ToolPending: lipgloss.NewStyle().
			Foreground(nordNight3),

		// List
		ListItem: lipgloss.NewStyle().
			Foreground(nordSnow2),
		ListItemSelected: lipgloss.NewStyle().
			Foreground(nordGreen).
			Bold(true),

		// Help
		HelpKey: lipgloss.NewStyle().
			Foreground(nordFrost1).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(nordNight3),
	}
}

func init() {
	Register("nord", NewNordTheme)
}
