// shell/streaming_content.go
package shell

import (
	"time"

	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

// Compile-time interface check.
var _ content.Content = (*StreamingContent)(nil)

// StreamingContent wraps a content.Content with optional typewriter effect.
type StreamingContent struct {
	inner           content.Content
	typewriter      bool
	typewriterSpeed time.Duration
	position        int
	text            string
	width           int
	height          int
}

// NewStreamingContent creates a new streaming content wrapper.
func NewStreamingContent(inner content.Content) *StreamingContent {
	return &StreamingContent{
		inner:           inner,
		typewriterSpeed: 30 * time.Millisecond,
	}
}

// WithTypewriter enables or disables typewriter effect.
func (s *StreamingContent) WithTypewriter(enabled bool) *StreamingContent {
	s.typewriter = enabled
	return s
}

// WithSpeed sets the typewriter speed.
func (s *StreamingContent) WithSpeed(d time.Duration) *StreamingContent {
	s.typewriterSpeed = d
	return s
}

// SetText updates the text to display.
func (s *StreamingContent) SetText(text string) {
	s.text = text
}

// Init implements content.Content.
func (s *StreamingContent) Init() tea.Cmd {
	if s.inner != nil {
		return s.inner.Init()
	}
	return nil
}

// typewriterTickMsg is sent to advance typewriter position.
type typewriterTickMsg struct{}

// Update implements content.Content.
func (s *StreamingContent) Update(msg tea.Msg) (content.Content, tea.Cmd) {
	switch msg.(type) {
	case typewriterTickMsg:
		if s.typewriter && s.position < len(s.text) {
			// Advance by 1-2 characters
			advance := 1
			if s.position < len(s.text)-1 {
				// Skip faster over whitespace
				if s.text[s.position] == ' ' || s.text[s.position] == '\n' {
					advance = 2
				}
			}
			s.position += advance
			if s.position > len(s.text) {
				s.position = len(s.text)
			}

			// Schedule next tick if not done
			if s.position < len(s.text) {
				return s, s.tickCmd()
			}
		}
	}

	return s, nil
}

// tickCmd returns a command that sends a typewriter tick after the configured delay.
func (s *StreamingContent) tickCmd() tea.Cmd {
	return tea.Tick(s.typewriterSpeed, func(time.Time) tea.Msg {
		return typewriterTickMsg{}
	})
}

// StartTypewriter begins the typewriter animation.
func (s *StreamingContent) StartTypewriter() tea.Cmd {
	if s.typewriter && s.position < len(s.text) {
		return s.tickCmd()
	}
	return nil
}

// View implements content.Content.
func (s *StreamingContent) View() string {
	if !s.typewriter {
		return s.text
	}

	// Typewriter mode: show text up to position + cursor
	if s.position >= len(s.text) {
		return s.text
	}

	visible := s.text[:s.position]
	cursor := "│"

	return visible + cursor
}

// Value implements content.Content.
func (s *StreamingContent) Value() any {
	return s.text
}

// SetSize implements content.Content.
func (s *StreamingContent) SetSize(width, height int) {
	s.width = width
	s.height = height
	if s.inner != nil {
		s.inner.SetSize(width, height)
	}
}
