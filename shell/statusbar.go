package shell

import (
	"fmt"
	"strings"

	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

// Status represents the current status to display.
type Status struct {
	Model      string
	Connected  bool
	Streaming  bool
	TokensUsed int
	TokensMax  int
	Mode       string
	Message    string
	Hints      string
}

// StatusBar renders the status bar at the bottom of the shell.
type StatusBar struct {
	status Status
	theme  theme.Theme
}

// NewStatusBar creates a new status bar.
func NewStatusBar(th theme.Theme) *StatusBar {
	return &StatusBar{
		theme: th,
		status: Status{
			Connected: true,
		},
	}
}

// SetStatus updates the status.
func (s *StatusBar) SetStatus(status Status) {
	s.status = status
}

// View renders the status bar.
func (s *StatusBar) View(width int) string {
	styles := s.theme.Styles()

	var sections []string

	// Model
	if s.status.Model != "" {
		sections = append(sections, s.status.Model)
	}

	// Connection status
	var statusText string
	if s.status.Streaming {
		statusText = styles.Warning.Render("● streaming")
	} else if s.status.Connected {
		statusText = styles.Success.Render("● connected")
	} else {
		statusText = styles.Error.Render("○ disconnected")
	}
	sections = append(sections, statusText)

	// Tokens
	if s.status.TokensMax > 0 {
		tokenText := fmt.Sprintf("%dk/%dk", s.status.TokensUsed/1000, s.status.TokensMax/1000)
		sections = append(sections, styles.Muted.Render(tokenText))
	}

	// Mode
	if s.status.Mode != "" {
		sections = append(sections, styles.Subtitle.Render(s.status.Mode))
	}

	// Custom message
	if s.status.Message != "" {
		sections = append(sections, s.status.Message)
	}

	// Join left sections
	left := strings.Join(sections, " │ ")

	// Hints on right
	right := ""
	if s.status.Hints != "" {
		right = styles.Muted.Render(s.status.Hints)
	}

	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	spacing := width - leftWidth - rightWidth
	if spacing < 1 {
		spacing = 1
	}

	bar := left + strings.Repeat(" ", spacing) + right

	return styles.StatusBar.Width(width).Render(bar)
}

// SetModel sets the model name.
func (s *StatusBar) SetModel(model string) {
	s.status.Model = model
}

// SetConnected sets the connection status.
func (s *StatusBar) SetConnected(connected bool) {
	s.status.Connected = connected
}

// SetStreaming sets the streaming status.
func (s *StatusBar) SetStreaming(streaming bool) {
	s.status.Streaming = streaming
}

// SetTokens sets the token counts.
func (s *StatusBar) SetTokens(used, max int) {
	s.status.TokensUsed = used
	s.status.TokensMax = max
}

// SetMode sets the mode string.
func (s *StatusBar) SetMode(mode string) {
	s.status.Mode = mode
}

// SetMessage sets a temporary message.
func (s *StatusBar) SetMessage(message string) {
	s.status.Message = message
}

// SetHints sets the hints text.
func (s *StatusBar) SetHints(hints string) {
	s.status.Hints = hints
}
