package theme

import "github.com/charmbracelet/lipgloss"

// Gruvbox dark color palette
const (
	// Background
	gruvboxBg0    = lipgloss.Color("#282828")
	gruvboxBg1    = lipgloss.Color("#3c3836")
	gruvboxBg2    = lipgloss.Color("#504945")
	gruvboxBg3    = lipgloss.Color("#665c54")

	// Foreground
	gruvboxFg0    = lipgloss.Color("#fbf1c7")
	gruvboxFg1    = lipgloss.Color("#ebdbb2")
	gruvboxFg2    = lipgloss.Color("#d5c4a1")
	gruvboxFg3    = lipgloss.Color("#bdae93")
	gruvboxFg4    = lipgloss.Color("#a89984")

	// Colors
	gruvboxRed    = lipgloss.Color("#fb4934")
	gruvboxGreen  = lipgloss.Color("#b8bb26")
	gruvboxYellow = lipgloss.Color("#fabd2f")
	gruvboxBlue   = lipgloss.Color("#83a598")
	gruvboxPurple = lipgloss.Color("#d3869b")
	gruvboxAqua   = lipgloss.Color("#8ec07c")
	gruvboxOrange = lipgloss.Color("#fe8019")
	gruvboxGray   = lipgloss.Color("#928374")
)

type gruvboxTheme struct {
	styles Styles
}

// NewGruvboxTheme creates a new Gruvbox dark color theme.
func NewGruvboxTheme() Theme {
	t := &gruvboxTheme{}
	t.styles = t.buildStyles()
	return t
}

func (t *gruvboxTheme) Name() string                   { return "gruvbox" }
func (t *gruvboxTheme) Background() lipgloss.Color     { return gruvboxBg0 }
func (t *gruvboxTheme) Foreground() lipgloss.Color     { return gruvboxFg1 }
func (t *gruvboxTheme) Primary() lipgloss.Color        { return gruvboxYellow }
func (t *gruvboxTheme) Secondary() lipgloss.Color      { return gruvboxBlue }
func (t *gruvboxTheme) Success() lipgloss.Color        { return gruvboxGreen }
func (t *gruvboxTheme) Warning() lipgloss.Color        { return gruvboxOrange }
func (t *gruvboxTheme) Error() lipgloss.Color          { return gruvboxRed }
func (t *gruvboxTheme) Info() lipgloss.Color           { return gruvboxBlue }
func (t *gruvboxTheme) Border() lipgloss.Color         { return gruvboxBg3 }
func (t *gruvboxTheme) BorderFocused() lipgloss.Color  { return gruvboxYellow }
func (t *gruvboxTheme) Muted() lipgloss.Color          { return gruvboxGray }
func (t *gruvboxTheme) UserColor() lipgloss.Color      { return gruvboxOrange }
func (t *gruvboxTheme) AssistantColor() lipgloss.Color { return gruvboxGreen }
func (t *gruvboxTheme) ToolColor() lipgloss.Color      { return gruvboxAqua }
func (t *gruvboxTheme) SystemColor() lipgloss.Color    { return gruvboxGray }
func (t *gruvboxTheme) Styles() Styles                 { return t.styles }

func (t *gruvboxTheme) buildStyles() Styles {
	return Styles{
		// Text
		Title: lipgloss.NewStyle().
			Foreground(gruvboxYellow).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(gruvboxBlue),
		Body: lipgloss.NewStyle().
			Foreground(gruvboxFg1),
		Muted: lipgloss.NewStyle().
			Foreground(gruvboxGray),
		Emphasized: lipgloss.NewStyle().
			Foreground(gruvboxFg1).
			Bold(true),

		// Status
		Success: lipgloss.NewStyle().
			Foreground(gruvboxGreen),
		Error: lipgloss.NewStyle().
			Foreground(gruvboxRed),
		Warning: lipgloss.NewStyle().
			Foreground(gruvboxOrange),
		Info: lipgloss.NewStyle().
			Foreground(gruvboxBlue),

		// Interactive
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxBg3),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxYellow),
		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxBg3).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxYellow).
			Padding(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(gruvboxFg1).
			Background(gruvboxBg2).
			Padding(0, 2),
		ButtonActive: lipgloss.NewStyle().
			Foreground(gruvboxBg0).
			Background(gruvboxYellow).
			Padding(0, 2),

		// Layout
		StatusBar: lipgloss.NewStyle().
			Foreground(gruvboxFg1).
			Background(gruvboxBg1),
		TabBar: lipgloss.NewStyle().
			Foreground(gruvboxFg1),
		TabActive: lipgloss.NewStyle().
			Foreground(gruvboxYellow).
			Bold(true).
			Underline(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(gruvboxGray),

		// Modal
		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxYellow).
			Padding(1, 2),
		ModalTitle: lipgloss.NewStyle().
			Foreground(gruvboxYellow).
			Bold(true),
		ModalFooter: lipgloss.NewStyle().
			Foreground(gruvboxGray),

		// Tool states
		ToolApproval: lipgloss.NewStyle().
			Foreground(gruvboxOrange),
		ToolExecuting: lipgloss.NewStyle().
			Foreground(gruvboxYellow),
		ToolSuccess: lipgloss.NewStyle().
			Foreground(gruvboxGreen),
		ToolError: lipgloss.NewStyle().
			Foreground(gruvboxRed),
		ToolPending: lipgloss.NewStyle().
			Foreground(gruvboxGray),

		// List
		ListItem: lipgloss.NewStyle().
			Foreground(gruvboxFg1),
		ListItemSelected: lipgloss.NewStyle().
			Foreground(gruvboxGreen).
			Bold(true),

		// Help
		HelpKey: lipgloss.NewStyle().
			Foreground(gruvboxAqua).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(gruvboxGray),
	}
}

func init() {
	Register("gruvbox", NewGruvboxTheme)
}
