package shell

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerType defines the visual style of the spinner.
type SpinnerType int

const (
	SpinnerDefault   SpinnerType = iota // Dot spinner
	SpinnerExecution                    // Points, shows elapsed time
	SpinnerStreaming                    // MiniDot, shows token rate
	SpinnerLoading                      // Line spinner
)

// Spinner provides a loading indicator with different visual styles.
type Spinner struct {
	spinnerType SpinnerType
	spinner     spinner.Model
	message     string
	startTime   time.Time
	tokenRate   float64
	active      bool
	style       lipgloss.Style
}

// NewSpinner creates a new spinner with the specified type.
func NewSpinner(spinnerType SpinnerType) *Spinner {
	s := &Spinner{
		spinnerType: spinnerType,
		spinner:     spinner.New(),
		style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")),
	}

	// Set spinner style based on type
	switch spinnerType {
	case SpinnerDefault:
		s.spinner.Spinner = spinner.Dot
	case SpinnerExecution:
		s.spinner.Spinner = spinner.Points
	case SpinnerStreaming:
		s.spinner.Spinner = spinner.MiniDot
	case SpinnerLoading:
		s.spinner.Spinner = spinner.Line
	}

	return s
}

// SetMessage sets the message displayed alongside the spinner.
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// SetTokenRate sets the token rate for streaming spinners.
func (s *Spinner) SetTokenRate(rate float64) {
	s.tokenRate = rate
}

// SetStyle sets the style for the spinner text.
func (s *Spinner) SetStyle(style lipgloss.Style) {
	s.style = style
}

// Start activates the spinner and returns a command to begin animation.
func (s *Spinner) Start() tea.Cmd {
	s.active = true
	s.startTime = time.Now()
	return s.spinner.Tick
}

// Stop deactivates the spinner.
func (s *Spinner) Stop() {
	s.active = false
}

// Active returns whether the spinner is currently running.
func (s *Spinner) Active() bool {
	return s.active
}

// Update handles spinner tick messages.
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	if !s.active {
		return s, nil
	}

	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View renders the spinner.
func (s *Spinner) View() string {
	if !s.active {
		return ""
	}

	spinnerView := s.spinner.View()

	switch s.spinnerType {
	case SpinnerExecution:
		elapsed := time.Since(s.startTime).Round(time.Second)
		if s.message != "" {
			return s.style.Render(fmt.Sprintf("%s %s (%s)", spinnerView, s.message, elapsed))
		}
		return s.style.Render(fmt.Sprintf("%s %s", spinnerView, elapsed))

	case SpinnerStreaming:
		if s.tokenRate > 0 {
			if s.message != "" {
				return s.style.Render(fmt.Sprintf("%s %s (%.1f tok/s)", spinnerView, s.message, s.tokenRate))
			}
			return s.style.Render(fmt.Sprintf("%s %.1f tok/s", spinnerView, s.tokenRate))
		}
		if s.message != "" {
			return s.style.Render(fmt.Sprintf("%s %s", spinnerView, s.message))
		}
		return s.style.Render(spinnerView)

	default:
		if s.message != "" {
			return s.style.Render(fmt.Sprintf("%s %s", spinnerView, s.message))
		}
		return s.style.Render(spinnerView)
	}
}

// Elapsed returns the duration since the spinner was started.
func (s *Spinner) Elapsed() time.Duration {
	if !s.active {
		return 0
	}
	return time.Since(s.startTime)
}

// Type returns the spinner type.
func (s *Spinner) Type() SpinnerType {
	return s.spinnerType
}

// Message returns the current message.
func (s *Spinner) Message() string {
	return s.message
}

// TokenRate returns the current token rate.
func (s *Spinner) TokenRate() float64 {
	return s.tokenRate
}
