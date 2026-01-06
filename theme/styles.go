package theme

import "github.com/charmbracelet/lipgloss"

// Styles contains all the composed lipgloss styles for UI components.
type Styles struct {
	// Text
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Body       lipgloss.Style
	Muted      lipgloss.Style
	Emphasized lipgloss.Style

	// Status
	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style

	// Interactive
	Border        lipgloss.Style
	BorderFocused lipgloss.Style
	Input         lipgloss.Style
	InputFocused  lipgloss.Style
	Button        lipgloss.Style
	ButtonActive  lipgloss.Style

	// Layout
	StatusBar   lipgloss.Style
	TabBar      lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style

	// Modal
	ModalBox    lipgloss.Style
	ModalTitle  lipgloss.Style
	ModalFooter lipgloss.Style

	// Tool states
	ToolApproval  lipgloss.Style
	ToolExecuting lipgloss.Style
	ToolSuccess   lipgloss.Style
	ToolError     lipgloss.Style
	ToolPending   lipgloss.Style

	// List
	ListItem         lipgloss.Style
	ListItemSelected lipgloss.Style

	// Help
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style
}
