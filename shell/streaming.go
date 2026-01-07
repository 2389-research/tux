// shell/streaming.go
package shell

import "time"

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
