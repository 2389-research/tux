package shell

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/tux/theme"
)

// ErrorModalConfig configures an ErrorModal.
type ErrorModalConfig struct {
	Errors []error
	Theme  theme.Theme
}

// ErrorModal displays a list of errors.
type ErrorModal struct {
	errors []error
	theme  theme.Theme
	width  int
	height int

	boxStyle   lipgloss.Style
	titleStyle lipgloss.Style
	errorStyle lipgloss.Style
	indexStyle lipgloss.Style
}

// NewErrorModal creates a new error modal.
func NewErrorModal(cfg ErrorModalConfig) *ErrorModal {
	th := cfg.Theme
	if th == nil {
		th = theme.NewDraculaTheme()
	}

	return &ErrorModal{
		errors: cfg.Errors,
		theme:  th,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff5555")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")).
			Bold(true),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		indexStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// ID implements Modal.
func (m *ErrorModal) ID() string { return "error-modal" }

// Title implements Modal.
func (m *ErrorModal) Title() string { return "Errors" }

// Size implements Modal.
func (m *ErrorModal) Size() Size { return SizeMedium }

// OnPush implements Modal.
func (m *ErrorModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop implements Modal.
func (m *ErrorModal) OnPop() {}

// HandleKey implements Modal.
func (m *ErrorModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type {
	case tea.KeyEscape, tea.KeyCtrlE:
		return true, func() tea.Msg { return PopMsg{} }
	}
	return false, nil
}

// Render implements Modal.
func (m *ErrorModal) Render(width, height int) string {
	var parts []string

	// Title
	title := fmt.Sprintf("Errors (%d)", len(m.errors))
	parts = append(parts, m.titleStyle.Render(title))
	parts = append(parts, "")

	// Calculate available lines for errors (title + blank + footer + blank + border padding)
	reservedLines := 6 // title, blank, footer text, blank, top/bottom border padding
	availableLines := height - reservedLines
	if availableLines < 1 {
		availableLines = 1
	}

	// Error list (truncate if too many)
	errorsToShow := m.errors
	truncated := false
	if len(errorsToShow) > availableLines {
		errorsToShow = errorsToShow[:availableLines-1] // Leave room for "...and N more"
		truncated = true
	}

	for i, err := range errorsToShow {
		index := m.indexStyle.Render(fmt.Sprintf("%d.", i+1))
		errText := m.errorStyle.Render(err.Error())
		parts = append(parts, index+" "+errText)
	}

	if truncated {
		remaining := len(m.errors) - len(errorsToShow)
		parts = append(parts, m.indexStyle.Render(fmt.Sprintf("...and %d more", remaining)))
	}

	if len(m.errors) == 0 {
		parts = append(parts, m.errorStyle.Render("No errors"))
	}

	parts = append(parts, "")
	parts = append(parts, m.indexStyle.Render("Press Esc or Ctrl+E to close"))

	content := strings.Join(parts, "\n")

	// Guard width to avoid negative values
	safeWidth := width - 4
	if safeWidth < 1 {
		safeWidth = 1
	}
	return m.boxStyle.Width(safeWidth).Render(content)
}
