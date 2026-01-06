package theme

import "github.com/charmbracelet/lipgloss"

// Neo-Terminal color palette - Information-rich sophistication
// Inspired by Swiss typography meets cyberdeck aesthetics
const (
	// Primary Palette
	ntDeepInk     = lipgloss.Color("#1a1b26") // Background - rich black-blue
	ntSoftPaper   = lipgloss.Color("#c0caf5") // Primary text - cool white
	ntAccentCoral = lipgloss.Color("#ff9e64") // User messages - warm
	ntAccentSage  = lipgloss.Color("#9ece6a") // Assistant - natural green
	ntAccentSky   = lipgloss.Color("#7aa2f7") // Tools/system - cool blue

	// Secondary Palette
	ntDimInk       = lipgloss.Color("#565f89") // Borders, secondary text
	ntGhost        = lipgloss.Color("#414868") // Subtle elements
	ntWarningAmber = lipgloss.Color("#e0af68") // Warnings
	ntErrorRuby    = lipgloss.Color("#f7768e") // Errors
	ntSuccessJade  = lipgloss.Color("#73daca") // Success states
)

type neoTerminalTheme struct {
	styles Styles
}

// NewNeoTerminalTheme creates a new Neo-Terminal theme.
// This theme features information-rich sophistication inspired by
// Swiss typography meets cyberdeck aesthetics.
func NewNeoTerminalTheme() Theme {
	t := &neoTerminalTheme{}
	t.styles = t.buildStyles()
	return t
}

func (t *neoTerminalTheme) Name() string                   { return "neo-terminal" }
func (t *neoTerminalTheme) Background() lipgloss.Color     { return ntDeepInk }
func (t *neoTerminalTheme) Foreground() lipgloss.Color     { return ntSoftPaper }
func (t *neoTerminalTheme) Primary() lipgloss.Color        { return ntAccentSky }
func (t *neoTerminalTheme) Secondary() lipgloss.Color      { return ntAccentCoral }
func (t *neoTerminalTheme) Success() lipgloss.Color        { return ntSuccessJade }
func (t *neoTerminalTheme) Warning() lipgloss.Color        { return ntWarningAmber }
func (t *neoTerminalTheme) Error() lipgloss.Color          { return ntErrorRuby }
func (t *neoTerminalTheme) Info() lipgloss.Color           { return ntAccentSky }
func (t *neoTerminalTheme) Border() lipgloss.Color         { return ntDimInk }
func (t *neoTerminalTheme) BorderFocused() lipgloss.Color  { return ntAccentSky }
func (t *neoTerminalTheme) Muted() lipgloss.Color          { return ntDimInk }
func (t *neoTerminalTheme) UserColor() lipgloss.Color      { return ntAccentCoral }
func (t *neoTerminalTheme) AssistantColor() lipgloss.Color { return ntAccentSage }
func (t *neoTerminalTheme) ToolColor() lipgloss.Color      { return ntAccentSky }
func (t *neoTerminalTheme) SystemColor() lipgloss.Color    { return ntDimInk }
func (t *neoTerminalTheme) Styles() Styles                 { return t.styles }

func (t *neoTerminalTheme) buildStyles() Styles {
	return Styles{
		// Text
		Title: lipgloss.NewStyle().
			Foreground(ntSoftPaper).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(ntDimInk).
			Bold(true),
		Body: lipgloss.NewStyle().
			Foreground(ntSoftPaper),
		Muted: lipgloss.NewStyle().
			Foreground(ntDimInk),
		Emphasized: lipgloss.NewStyle().
			Foreground(ntSoftPaper).
			Bold(true),

		// Status
		Success: lipgloss.NewStyle().
			Foreground(ntSuccessJade).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(ntErrorRuby).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(ntWarningAmber).
			Bold(true),
		Info: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Bold(true),

		// Interactive
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ntDimInk),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ntAccentSky),
		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ntDimInk).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ntAccentSky).
			Padding(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(ntSoftPaper).
			Background(ntGhost).
			Padding(0, 2).
			Bold(true),
		ButtonActive: lipgloss.NewStyle().
			Foreground(ntDeepInk).
			Background(ntAccentSky).
			Padding(0, 2).
			Bold(true),

		// Layout
		StatusBar: lipgloss.NewStyle().
			Foreground(ntDimInk).
			Background(ntDeepInk),
		TabBar: lipgloss.NewStyle().
			Foreground(ntSoftPaper),
		TabActive: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Bold(true).
			Underline(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(ntDimInk),

		// Modal
		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ntAccentSky).
			Padding(1, 2),
		ModalTitle: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Bold(true),
		ModalFooter: lipgloss.NewStyle().
			Foreground(ntDimInk),

		// Tool states
		ToolApproval: lipgloss.NewStyle().
			Foreground(ntAccentCoral).
			Bold(true),
		ToolExecuting: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Bold(true),
		ToolSuccess: lipgloss.NewStyle().
			Foreground(ntSuccessJade).
			Bold(true),
		ToolError: lipgloss.NewStyle().
			Foreground(ntErrorRuby).
			Bold(true),
		ToolPending: lipgloss.NewStyle().
			Foreground(ntDimInk),

		// List
		ListItem: lipgloss.NewStyle().
			Foreground(ntSoftPaper),
		ListItemSelected: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Background(ntGhost).
			Bold(true),

		// Help
		HelpKey: lipgloss.NewStyle().
			Foreground(ntAccentSky).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(ntDimInk),
	}
}

func init() {
	Register("neo-terminal", NewNeoTerminalTheme)
}
