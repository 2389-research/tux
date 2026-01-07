// shell/streaming.go
package shell

import (
	"fmt"
	"strings"
	"time"

	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

// StreamingController manages streaming state for LLM responses.
type StreamingController struct {
	text          string
	tokenCount    int
	tokenRate     float64
	lastTokenTime time.Time
	startTime     time.Time

	streaming bool
	thinking  bool
	waiting   bool

	toolCalls []ToolCall

	// Spinner animation
	spinnerFrames []string
	spinnerFrame  int
	lastSpinTime  time.Time
}

// ToolCall represents a tool call in progress.
type ToolCall struct {
	ID         string
	Name       string
	InProgress bool
}

// NewStreamingController creates a new streaming controller.
func NewStreamingController() *StreamingController {
	return &StreamingController{
		spinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start begins a streaming session.
func (s *StreamingController) Start() {
	s.streaming = true
	s.waiting = true
	s.startTime = time.Now()
	s.lastTokenTime = time.Now()
}

// End finishes a streaming session.
func (s *StreamingController) End() {
	s.streaming = false
	s.waiting = false
	s.thinking = false
}

// Reset clears all streaming state.
func (s *StreamingController) Reset() {
	s.text = ""
	s.tokenCount = 0
	s.tokenRate = 0
	s.streaming = false
	s.thinking = false
	s.waiting = false
	s.toolCalls = nil
}

// IsStreaming returns true if currently streaming.
func (s *StreamingController) IsStreaming() bool {
	return s.streaming
}

// IsWaiting returns true if streaming started but no tokens received.
func (s *StreamingController) IsWaiting() bool {
	return s.streaming && s.waiting
}

// AppendToken adds a text chunk and updates token rate.
func (s *StreamingController) AppendToken(text string) {
	now := time.Now()
	elapsed := now.Sub(s.lastTokenTime).Seconds()

	if elapsed > 0 && s.tokenCount > 0 {
		instantRate := 1.0 / elapsed
		if s.tokenRate == 0 {
			s.tokenRate = instantRate
		} else {
			// EMA with alpha = 0.3
			s.tokenRate = 0.3*instantRate + 0.7*s.tokenRate
		}
	}

	s.text += text
	s.tokenCount++
	s.lastTokenTime = now
	s.waiting = false
}

// GetText returns the accumulated text.
func (s *StreamingController) GetText() string {
	return s.text
}

// TokenCount returns the number of tokens received.
func (s *StreamingController) TokenCount() int {
	return s.tokenCount
}

// TokenRate returns the current token rate (tokens/second).
func (s *StreamingController) TokenRate() float64 {
	return s.tokenRate
}

// SetThinking sets the thinking state.
func (s *StreamingController) SetThinking(active bool) {
	s.thinking = active
	if active {
		s.lastSpinTime = time.Now()
	}
}

// IsThinking returns true if currently in thinking state.
func (s *StreamingController) IsThinking() bool {
	return s.thinking
}

// StartToolCall marks a tool call as started.
func (s *StreamingController) StartToolCall(id, name string) {
	s.toolCalls = append(s.toolCalls, ToolCall{
		ID:         id,
		Name:       name,
		InProgress: true,
	})
}

// EndToolCall marks a tool call as complete.
func (s *StreamingController) EndToolCall(id string) {
	for i := range s.toolCalls {
		if s.toolCalls[i].ID == id {
			s.toolCalls[i].InProgress = false
			break
		}
	}
}

// ActiveToolCalls returns tool calls that are still in progress.
func (s *StreamingController) ActiveToolCalls() []ToolCall {
	var active []ToolCall
	for _, tc := range s.toolCalls {
		if tc.InProgress {
			active = append(active, tc)
		}
	}
	return active
}

// RenderStatus renders the streaming status for the statusbar.
// Returns empty string if not streaming.
func (s *StreamingController) RenderStatus(th theme.Theme) string {
	if !s.streaming {
		return ""
	}

	var parts []string

	// Waiting state
	if s.waiting {
		style := lipgloss.NewStyle().Italic(true)
		if th != nil {
			style = style.Foreground(th.Muted())
		}
		return style.Render("Waiting...")
	}

	// Thinking with spinner
	if s.thinking {
		frame := s.getSpinnerFrame()
		style := lipgloss.NewStyle()
		if th != nil {
			style = style.Foreground(th.Primary())
		}
		parts = append(parts, style.Render(frame+" Thinking"))
	}

	// Active tool calls
	for _, tc := range s.toolCalls {
		if tc.InProgress {
			style := lipgloss.NewStyle()
			if th != nil {
				style = style.Foreground(th.Secondary())
			}
			parts = append(parts, style.Render("▍ "+tc.Name))
		}
	}

	// Token rate
	if s.tokenRate > 0 {
		style := lipgloss.NewStyle()
		if th != nil {
			style = style.Foreground(th.Muted())
		}
		rate := fmt.Sprintf("▸ %.0f tok/s", s.tokenRate)
		parts = append(parts, style.Render(rate))
	}

	return strings.Join(parts, "  ")
}

// getSpinnerFrame returns the current spinner frame.
func (s *StreamingController) getSpinnerFrame() string {
	if len(s.spinnerFrames) == 0 {
		return "⠋"
	}

	// Advance frame every 80ms
	now := time.Now()
	if now.Sub(s.lastSpinTime) > 80*time.Millisecond {
		s.spinnerFrame = (s.spinnerFrame + 1) % len(s.spinnerFrames)
		s.lastSpinTime = now
	}

	return s.spinnerFrames[s.spinnerFrame]
}
