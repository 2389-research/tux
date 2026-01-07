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

// Update implements content.Content.
func (s *StreamingContent) Update(msg tea.Msg) (content.Content, tea.Cmd) {
	return s, nil
}

// View implements content.Content.
func (s *StreamingContent) View() string {
	return s.text
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
