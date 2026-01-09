# Streaming Display Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement StreamingController for state management and StreamingContent for typewriter effects, with auto-display in statusbar.

**Architecture:** StreamingController lives in shell/, accessed via shell.Streaming(). StatusBar reads from controller on render. StreamingContent wraps content.Content with typewriter effect.

**Tech Stack:** Go, Bubble Tea, tux theme system

---

### Task 1: StreamingController Struct and Lifecycle

**Files:**
- Create: `shell/streaming.go`
- Create: `shell/streaming_test.go`

**Step 1: Write the failing test**

```go
// shell/streaming_test.go
package shell

import (
	"testing"
)

func TestStreamingController_Lifecycle(t *testing.T) {
	s := NewStreamingController()

	// Initially not streaming
	if s.IsStreaming() {
		t.Error("expected not streaming initially")
	}

	// Start streaming
	s.Start()
	if !s.IsStreaming() {
		t.Error("expected streaming after Start")
	}
	if !s.IsWaiting() {
		t.Error("expected waiting after Start (no tokens yet)")
	}

	// End streaming
	s.End()
	if s.IsStreaming() {
		t.Error("expected not streaming after End")
	}

	// Reset clears state
	s.Start()
	s.Reset()
	if s.IsStreaming() {
		t.Error("expected not streaming after Reset")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingController_Lifecycle -v`
Expected: FAIL with "undefined: NewStreamingController"

**Step 3: Write minimal implementation**

```go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingController_Lifecycle -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming.go shell/streaming_test.go
git commit -m "feat(shell): add StreamingController with lifecycle methods"
```

---

### Task 2: StreamingController Token Methods

**Files:**
- Modify: `shell/streaming.go`
- Modify: `shell/streaming_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/streaming_test.go

func TestStreamingController_Tokens(t *testing.T) {
	s := NewStreamingController()
	s.Start()

	// Append tokens
	s.AppendToken("Hello")
	s.AppendToken(" world")

	if s.GetText() != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", s.GetText())
	}

	if s.TokenCount() != 2 {
		t.Errorf("expected 2 tokens, got %d", s.TokenCount())
	}

	// No longer waiting after first token
	if s.IsWaiting() {
		t.Error("expected not waiting after tokens received")
	}

	// Token rate should be > 0 after multiple tokens
	if s.TokenRate() <= 0 {
		t.Error("expected positive token rate")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingController_Tokens -v`
Expected: FAIL with "undefined: s.AppendToken"

**Step 3: Write minimal implementation**

```go
// Add to shell/streaming.go

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingController_Tokens -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming.go shell/streaming_test.go
git commit -m "feat(shell): add token methods to StreamingController"
```

---

### Task 3: StreamingController Status Methods

**Files:**
- Modify: `shell/streaming.go`
- Modify: `shell/streaming_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/streaming_test.go

func TestStreamingController_Status(t *testing.T) {
	s := NewStreamingController()
	s.Start()

	// Thinking
	if s.IsThinking() {
		t.Error("expected not thinking initially")
	}
	s.SetThinking(true)
	if !s.IsThinking() {
		t.Error("expected thinking after SetThinking(true)")
	}
	s.SetThinking(false)
	if s.IsThinking() {
		t.Error("expected not thinking after SetThinking(false)")
	}

	// Tool calls
	s.StartToolCall("1", "Bash")
	s.StartToolCall("2", "Read")

	active := s.ActiveToolCalls()
	if len(active) != 2 {
		t.Errorf("expected 2 active tool calls, got %d", len(active))
	}

	s.EndToolCall("1")
	active = s.ActiveToolCalls()
	if len(active) != 1 {
		t.Errorf("expected 1 active tool call, got %d", len(active))
	}
	if active[0].Name != "Read" {
		t.Errorf("expected 'Read', got %q", active[0].Name)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingController_Status -v`
Expected: FAIL with "undefined: s.IsThinking"

**Step 3: Write minimal implementation**

```go
// Add to shell/streaming.go

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingController_Status -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming.go shell/streaming_test.go
git commit -m "feat(shell): add status methods to StreamingController"
```

---

### Task 4: StreamingController Render Methods

**Files:**
- Modify: `shell/streaming.go`
- Modify: `shell/streaming_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/streaming_test.go

import "strings"

func TestStreamingController_Render(t *testing.T) {
	s := NewStreamingController()

	// Not streaming - empty render
	if s.RenderStatus(nil) != "" {
		t.Error("expected empty status when not streaming")
	}

	s.Start()

	// Waiting state
	status := s.RenderStatus(nil)
	if !strings.Contains(status, "Waiting") {
		t.Errorf("expected 'Waiting' in status, got %q", status)
	}

	// After tokens, show rate
	s.AppendToken("test")
	s.tokenRate = 42.0 // Force rate for testing
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "42") {
		t.Errorf("expected rate in status, got %q", status)
	}

	// Thinking shows spinner
	s.SetThinking(true)
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "Thinking") {
		t.Errorf("expected 'Thinking' in status, got %q", status)
	}

	// Tool calls shown
	s.StartToolCall("1", "Bash")
	status = s.RenderStatus(nil)
	if !strings.Contains(status, "Bash") {
		t.Errorf("expected 'Bash' in status, got %q", status)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingController_Render -v`
Expected: FAIL with "undefined: s.RenderStatus"

**Step 3: Write minimal implementation**

```go
// Add to shell/streaming.go

import (
	"fmt"
	"strings"
	"time"

	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingController_Render -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming.go shell/streaming_test.go
git commit -m "feat(shell): add render methods to StreamingController"
```

---

### Task 5: Shell Integration

**Files:**
- Modify: `shell/shell.go`
- Modify: `shell/statusbar.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

func TestShell_Streaming(t *testing.T) {
	sh := New(nil, DefaultConfig())

	// Streaming() returns controller
	s := sh.Streaming()
	if s == nil {
		t.Fatal("expected non-nil StreamingController")
	}

	// Same instance each time
	if sh.Streaming() != s {
		t.Error("expected same StreamingController instance")
	}

	// Can disable streaming status
	sh.SetStreamingStatusVisible(false)
	// (visibility tested in statusbar render)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShell_Streaming -v`
Expected: FAIL with "undefined: sh.Streaming"

**Step 3: Write minimal implementation**

```go
// Add to shell/shell.go struct:
// streaming *StreamingController
// streamingStatusVisible bool

// In New():
// streaming: NewStreamingController(),
// streamingStatusVisible: true,

// Add methods to shell/shell.go:

// Streaming returns the streaming controller.
func (s *Shell) Streaming() *StreamingController {
	return s.streaming
}

// SetStreamingStatusVisible controls whether streaming status appears in statusbar.
func (s *Shell) SetStreamingStatusVisible(visible bool) {
	s.streamingStatusVisible = visible
}
```

Update Shell struct in shell.go:

```go
type Shell struct {
	// Components
	tabs         *TabBar
	input        *Input
	statusBar    *StatusBar
	modalManager *Manager
	streaming    *StreamingController

	// State
	width                  int
	height                 int
	focused                FocusTarget
	ready                  bool
	streamingStatusVisible bool

	// Configuration
	theme  theme.Theme
	config Config
}
```

Update New() in shell.go:

```go
func New(th theme.Theme, cfg Config) *Shell {
	if th == nil {
		th = theme.NewDraculaTheme()
	}

	s := &Shell{
		theme:                  th,
		config:                 cfg,
		tabs:                   NewTabBar(th),
		input:                  NewInput(th, cfg.InputPrefix, cfg.InputPlaceholder),
		statusBar:              NewStatusBar(th),
		modalManager:           NewManager(),
		streaming:              NewStreamingController(),
		focused:                FocusInput,
		streamingStatusVisible: true,
	}

	return s
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShell_Streaming -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): add Streaming() method to Shell"
```

---

### Task 6: StatusBar Streaming Integration

**Files:**
- Modify: `shell/statusbar.go`
- Add to: `shell/streaming_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/streaming_test.go

func TestStatusBar_StreamingStatus(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)
	s := NewStreamingController()

	// No streaming status when not streaming
	sb.SetStreamingController(s, true)
	view := sb.View(80)
	if strings.Contains(view, "Thinking") || strings.Contains(view, "tok/s") {
		t.Error("expected no streaming status when not streaming")
	}

	// Start streaming and set thinking
	s.Start()
	s.SetThinking(true)
	view = sb.View(80)
	if !strings.Contains(view, "Thinking") {
		t.Errorf("expected 'Thinking' in view, got %q", view)
	}

	// Disabled streaming status
	sb.SetStreamingController(s, false)
	view = sb.View(80)
	if strings.Contains(view, "Thinking") {
		t.Error("expected no streaming status when disabled")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStatusBar_StreamingStatus -v`
Expected: FAIL with "undefined: sb.SetStreamingController"

**Step 3: Write minimal implementation**

```go
// Add to StatusBar struct in shell/statusbar.go:
// streaming        *StreamingController
// streamingVisible bool

// Add method:

// SetStreamingController sets the streaming controller for status display.
func (s *StatusBar) SetStreamingController(sc *StreamingController, visible bool) {
	s.streaming = sc
	s.streamingVisible = visible
}

// Update View() to include streaming status after connection status:

// In View(), after the connection status section, add:
	// Streaming status
	if s.streamingVisible && s.streaming != nil {
		streamingStatus := s.streaming.RenderStatus(s.theme)
		if streamingStatus != "" {
			sections = append(sections, streamingStatus)
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStatusBar_StreamingStatus -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/statusbar.go shell/streaming_test.go
git commit -m "feat(shell): integrate streaming status in statusbar"
```

---

### Task 7: Wire StatusBar in Shell.View

**Files:**
- Modify: `shell/shell.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

func TestShell_StreamingInView(t *testing.T) {
	sh := New(nil, DefaultConfig())

	// Simulate window size
	sh.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Start streaming with thinking
	sh.Streaming().Start()
	sh.Streaming().SetThinking(true)

	view := sh.View()
	if !strings.Contains(view, "Thinking") {
		t.Errorf("expected 'Thinking' in shell view")
	}

	// Disable streaming status
	sh.SetStreamingStatusVisible(false)
	view = sh.View()
	if strings.Contains(view, "Thinking") {
		t.Error("expected no 'Thinking' when streaming status disabled")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShell_StreamingInView -v`
Expected: FAIL (streaming not wired to statusbar)

**Step 3: Write minimal implementation**

Update Shell.View() to pass streaming controller to statusbar before rendering:

```go
// In View(), before rendering statusbar:
	// Status bar
	if s.config.ShowStatusBar {
		s.statusBar.SetStreamingController(s.streaming, s.streamingStatusVisible)
		sections = append(sections, s.statusBar.View(s.width))
	}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShell_StreamingInView -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): wire streaming controller to statusbar in View"
```

---

### Task 8: StreamingContent Struct and Builder

**Files:**
- Create: `shell/streaming_content.go`
- Create: `shell/streaming_content_test.go`

**Step 1: Write the failing test**

```go
// shell/streaming_content_test.go
package shell

import (
	"testing"
	"time"

	"github.com/2389-research/tux/content"
)

// mockContent implements content.Content for testing
type mockContent struct {
	text   string
	width  int
	height int
}

func (m *mockContent) Init() tea.Cmd                           { return nil }
func (m *mockContent) Update(msg tea.Msg) (content.Content, tea.Cmd) { return m, nil }
func (m *mockContent) View() string                            { return m.text }
func (m *mockContent) Value() any                              { return m.text }
func (m *mockContent) SetSize(w, h int)                        { m.width = w; m.height = h }

func TestStreamingContent_Builder(t *testing.T) {
	inner := &mockContent{text: "hello"}

	sc := NewStreamingContent(inner)
	if sc == nil {
		t.Fatal("expected non-nil StreamingContent")
	}

	// Default: typewriter disabled
	if sc.typewriter {
		t.Error("expected typewriter disabled by default")
	}

	// Builder methods return self
	sc2 := sc.WithTypewriter(true).WithSpeed(50 * time.Millisecond)
	if sc2 != sc {
		t.Error("expected builder to return same instance")
	}

	if !sc.typewriter {
		t.Error("expected typewriter enabled")
	}

	if sc.typewriterSpeed != 50*time.Millisecond {
		t.Errorf("expected 50ms speed, got %v", sc.typewriterSpeed)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingContent_Builder -v`
Expected: FAIL with "undefined: NewStreamingContent"

**Step 3: Write minimal implementation**

```go
// shell/streaming_content.go
package shell

import (
	"time"

	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingContent_Builder -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming_content.go shell/streaming_content_test.go
git commit -m "feat(shell): add StreamingContent struct with builder"
```

---

### Task 9: StreamingContent Typewriter Effect

**Files:**
- Modify: `shell/streaming_content.go`
- Modify: `shell/streaming_content_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/streaming_content_test.go

func TestStreamingContent_Typewriter(t *testing.T) {
	inner := &mockContent{}
	sc := NewStreamingContent(inner).WithTypewriter(true)

	sc.SetText("Hello world")

	// Initially position is 0, should show empty or cursor only
	view := sc.View()
	if len(view) > 5 { // Just cursor character
		t.Errorf("expected minimal view initially, got %q", view)
	}

	// Advance position
	sc.position = 5
	view = sc.View()
	if !strings.HasPrefix(view, "Hello") {
		t.Errorf("expected 'Hello' prefix, got %q", view)
	}

	// Should have cursor
	if !strings.Contains(view, "│") {
		t.Errorf("expected cursor in view, got %q", view)
	}

	// Full position shows all text
	sc.position = len("Hello world")
	view = sc.View()
	if !strings.Contains(view, "Hello world") {
		t.Errorf("expected full text, got %q", view)
	}
}

func TestStreamingContent_TypewriterDisabled(t *testing.T) {
	inner := &mockContent{}
	sc := NewStreamingContent(inner) // typewriter disabled by default

	sc.SetText("Hello world")
	view := sc.View()

	// Should show full text immediately
	if view != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", view)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStreamingContent_Typewriter -v`
Expected: FAIL (View doesn't implement typewriter logic)

**Step 3: Write minimal implementation**

```go
// Update View() in shell/streaming_content.go:

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

// Add typewriter tick message and update handling:

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStreamingContent_Typewriter -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/streaming_content.go shell/streaming_content_test.go
git commit -m "feat(shell): add typewriter effect to StreamingContent"
```

---

### Task 10: Run All Tests and Final Verification

**Files:**
- None (verification only)

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: All tests pass

**Step 2: Verify no lint issues**

Run: `go vet ./...`
Expected: No issues

**Step 3: Final commit if any cleanup needed**

```bash
git status
# If any uncommitted changes:
git add -A
git commit -m "chore: final cleanup for streaming display"
```

---

## Summary

| Task | Component | Description |
|------|-----------|-------------|
| 1 | StreamingController | Struct + lifecycle (Start/End/Reset) |
| 2 | StreamingController | Token methods (AppendToken, GetText, TokenRate) |
| 3 | StreamingController | Status methods (SetThinking, tool calls) |
| 4 | StreamingController | RenderStatus for statusbar |
| 5 | Shell | Add Streaming() method |
| 6 | StatusBar | SetStreamingController integration |
| 7 | Shell | Wire streaming to statusbar in View() |
| 8 | StreamingContent | Struct + builder methods |
| 9 | StreamingContent | Typewriter effect implementation |
| 10 | All | Final verification |
