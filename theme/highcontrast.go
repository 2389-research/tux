package theme

import "github.com/charmbracelet/lipgloss"

// High contrast color palette for accessibility
const (
	hcBlack   = lipgloss.Color("#000000")
	hcWhite   = lipgloss.Color("#ffffff")
	hcYellow  = lipgloss.Color("#ffff00")
	hcCyan    = lipgloss.Color("#00ffff")
	hcGreen   = lipgloss.Color("#00ff00")
	hcRed     = lipgloss.Color("#ff0000")
	hcMagenta = lipgloss.Color("#ff00ff")
	hcBlue    = lipgloss.Color("#0080ff")
	hcOrange  = lipgloss.Color("#ff8000")
	hcGray    = lipgloss.Color("#808080")
)

type highContrastTheme struct {
	styles Styles
}

// NewHighContrastTheme creates a new high contrast theme for accessibility.
func NewHighContrastTheme() Theme {
	t := &highContrastTheme{}
	t.styles = t.buildStyles()
	return t
}

func (t *highContrastTheme) Name() string                   { return "high-contrast" }
func (t *highContrastTheme) Background() lipgloss.Color     { return hcBlack }
func (t *highContrastTheme) Foreground() lipgloss.Color     { return hcWhite }
func (t *highContrastTheme) Primary() lipgloss.Color        { return hcYellow }
func (t *highContrastTheme) Secondary() lipgloss.Color      { return hcCyan }
func (t *highContrastTheme) Success() lipgloss.Color        { return hcGreen }
func (t *highContrastTheme) Warning() lipgloss.Color        { return hcOrange }
func (t *highContrastTheme) Error() lipgloss.Color          { return hcRed }
func (t *highContrastTheme) Info() lipgloss.Color           { return hcCyan }
func (t *highContrastTheme) Border() lipgloss.Color         { return hcWhite }
func (t *highContrastTheme) BorderFocused() lipgloss.Color  { return hcYellow }
func (t *highContrastTheme) Muted() lipgloss.Color          { return hcGray }
func (t *highContrastTheme) UserColor() lipgloss.Color      { return hcOrange }
func (t *highContrastTheme) AssistantColor() lipgloss.Color { return hcGreen }
func (t *highContrastTheme) ToolColor() lipgloss.Color      { return hcCyan }
func (t *highContrastTheme) SystemColor() lipgloss.Color    { return hcGray }
func (t *highContrastTheme) Styles() Styles                 { return t.styles }

func (t *highContrastTheme) buildStyles() Styles {
	return Styles{
		// Text
		Title: lipgloss.NewStyle().
			Foreground(hcYellow).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(hcCyan),
		Body: lipgloss.NewStyle().
			Foreground(hcWhite),
		Muted: lipgloss.NewStyle().
			Foreground(hcGray),
		Emphasized: lipgloss.NewStyle().
			Foreground(hcWhite).
			Bold(true),

		// Status
		Success: lipgloss.NewStyle().
			Foreground(hcGreen),
		Error: lipgloss.NewStyle().
			Foreground(hcRed),
		Warning: lipgloss.NewStyle().
			Foreground(hcOrange),
		Info: lipgloss.NewStyle().
			Foreground(hcCyan),

		// Interactive
		Border: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(hcWhite),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(hcYellow),
		Input: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(hcWhite).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(hcYellow).
			Padding(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(hcBlack).
			Background(hcWhite).
			Padding(0, 2),
		ButtonActive: lipgloss.NewStyle().
			Foreground(hcBlack).
			Background(hcYellow).
			Padding(0, 2).
			Bold(true),

		// Layout
		StatusBar: lipgloss.NewStyle().
			Foreground(hcBlack).
			Background(hcWhite),
		TabBar: lipgloss.NewStyle().
			Foreground(hcWhite),
		TabActive: lipgloss.NewStyle().
			Foreground(hcYellow).
			Bold(true).
			Underline(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(hcGray),

		// Modal
		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(hcYellow).
			Padding(1, 2),
		ModalTitle: lipgloss.NewStyle().
			Foreground(hcYellow).
			Bold(true),
		ModalFooter: lipgloss.NewStyle().
			Foreground(hcGray),

		// Tool states
		ToolApproval: lipgloss.NewStyle().
			Foreground(hcOrange).
			Bold(true),
		ToolExecuting: lipgloss.NewStyle().
			Foreground(hcYellow).
			Bold(true),
		ToolSuccess: lipgloss.NewStyle().
			Foreground(hcGreen).
			Bold(true),
		ToolError: lipgloss.NewStyle().
			Foreground(hcRed).
			Bold(true),
		ToolPending: lipgloss.NewStyle().
			Foreground(hcGray),

		// List
		ListItem: lipgloss.NewStyle().
			Foreground(hcWhite),
		ListItemSelected: lipgloss.NewStyle().
			Foreground(hcGreen).
			Bold(true),

		// Help
		HelpKey: lipgloss.NewStyle().
			Foreground(hcCyan).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(hcWhite),
	}
}

func init() {
	Register("high-contrast", NewHighContrastTheme)
}
