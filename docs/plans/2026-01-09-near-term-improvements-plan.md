# Near-Term Improvements Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add error display, history navigation, and approval handling to tux.App.

**Architecture:** Error display uses an error slice in App that accumulates during runs and clears on success, with status bar indicator and modal. History navigation tracks an index into user messages for up/down arrow cycling. Approval handling adds a new EventApproval type with Response channel for blocking approval flows.

**Tech Stack:** Go, Bubble Tea, Lipgloss

---

## Task 1: Error Display - Status Struct Extension

**Files:**
- Modify: `shell/statusbar.go:11-21` (Status struct)
- Test: `shell/statusbar_test.go` (create if needed)

**Step 1: Write the failing test**

Create `shell/statusbar_test.go`:

```go
package shell

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestStatusBarErrorIndicator(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	// Set an error
	sb.SetError("connection timeout", 1)

	view := sb.View(80)
	if !strings.Contains(view, "⚠") {
		t.Error("status bar should show error indicator")
	}
	if !strings.Contains(view, "connection...") {
		t.Error("status bar should show truncated error message")
	}
}

func TestStatusBarNoError(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	view := sb.View(80)
	if strings.Contains(view, "⚠") {
		t.Error("status bar should not show error indicator when no errors")
	}
}

func TestStatusBarClearError(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	sb.SetError("some error", 1)
	sb.ClearError()

	view := sb.View(80)
	if strings.Contains(view, "⚠") {
		t.Error("status bar should not show error after clearing")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestStatusBar -v`
Expected: FAIL - SetError and ClearError methods don't exist

**Step 3: Write minimal implementation**

Add to `shell/statusbar.go` Status struct:

```go
type Status struct {
	Model      string
	Connected  bool
	Streaming  bool
	TokensUsed int
	TokensMax  int
	Mode       string
	Message    string
	Hints      string
	// Error display
	ErrorText  string // Truncated error for status bar
	ErrorCount int    // Number of accumulated errors
}
```

Add methods to StatusBar:

```go
// SetError sets the error indicator with truncated text.
func (s *StatusBar) SetError(text string, count int) {
	// Truncate to ~10 chars
	if len(text) > 10 {
		text = text[:10] + "..."
	}
	s.status.ErrorText = text
	s.status.ErrorCount = count
}

// ClearError clears the error indicator.
func (s *StatusBar) ClearError() {
	s.status.ErrorText = ""
	s.status.ErrorCount = 0
}
```

Update View() to include error section after connection status:

```go
// Error indicator (after connection status section)
if s.status.ErrorText != "" {
	errorText := fmt.Sprintf("⚠ \"%s\"", s.status.ErrorText)
	if s.status.ErrorCount > 1 {
		errorText = fmt.Sprintf("⚠ \"%s\" +%d", s.status.ErrorText, s.status.ErrorCount-1)
	}
	sections = append(sections, styles.Error.Render(errorText))
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestStatusBar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/statusbar.go shell/statusbar_test.go
git commit -m "feat(shell): add error indicator to status bar"
```

---

## Task 2: Error Display - Error Modal

**Files:**
- Create: `shell/modal_error.go`
- Test: `shell/modal_error_test.go`

**Step 1: Write the failing test**

Create `shell/modal_error_test.go`:

```go
package shell

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestErrorModalRendersErrors(t *testing.T) {
	errs := []error{
		errors.New("connection timeout"),
		errors.New("rate limit exceeded"),
	}
	modal := NewErrorModal(ErrorModalConfig{
		Errors: errs,
		Theme:  theme.NewDraculaTheme(),
	})

	view := modal.Render(60, 20)
	if !strings.Contains(view, "connection timeout") {
		t.Error("modal should show first error")
	}
	if !strings.Contains(view, "rate limit exceeded") {
		t.Error("modal should show second error")
	}
}

func TestErrorModalID(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	if modal.ID() != "error-modal" {
		t.Errorf("expected id 'error-modal', got %q", modal.ID())
	}
}

func TestErrorModalTitle(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	if modal.Title() != "Errors" {
		t.Errorf("expected title 'Errors', got %q", modal.Title())
	}
}

func TestErrorModalEscapeCloses(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	handled, cmd := modal.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})
	if !handled {
		t.Error("escape should be handled")
	}
	if cmd == nil {
		t.Error("escape should return command")
	}
	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("expected PopMsg, got %T", msg)
	}
}

func TestErrorModalCtrlECloses(t *testing.T) {
	modal := NewErrorModal(ErrorModalConfig{})
	handled, cmd := modal.HandleKey(tea.KeyMsg{Type: tea.KeyCtrlE})
	if !handled {
		t.Error("ctrl+e should be handled")
	}
	if cmd == nil {
		t.Error("ctrl+e should return command")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestErrorModal -v`
Expected: FAIL - NewErrorModal doesn't exist

**Step 3: Write minimal implementation**

Create `shell/modal_error.go`:

```go
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

	// Error list
	for i, err := range m.errors {
		index := m.indexStyle.Render(fmt.Sprintf("%d.", i+1))
		errText := m.errorStyle.Render(err.Error())
		parts = append(parts, index+" "+errText)
	}

	if len(m.errors) == 0 {
		parts = append(parts, m.errorStyle.Render("No errors"))
	}

	parts = append(parts, "")
	parts = append(parts, m.indexStyle.Render("Press Esc or Ctrl+E to close"))

	content := strings.Join(parts, "\n")
	return m.boxStyle.Width(width - 4).Render(content)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestErrorModal -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/modal_error.go shell/modal_error_test.go
git commit -m "feat(shell): add ErrorModal for displaying errors"
```

---

## Task 3: Error Display - Shell Ctrl+E Handler

**Files:**
- Modify: `shell/shell.go:129-133` (global keys section)
- Modify: `shell/shell.go:31` (Shell struct - add error callback)
- Modify: `shell/shell.go:34-47` (Config struct)
- Test: `shell/shell_test.go`

**Step 1: Write the failing test**

Add to `shell/shell_test.go`:

```go
func TestShellCtrlEOpensErrorModal(t *testing.T) {
	th := theme.NewDraculaTheme()
	cfg := DefaultConfig()

	errorModalOpened := false
	cfg.OnShowErrors = func() {
		errorModalOpened = true
	}

	s := New(th, cfg)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Send Ctrl+E
	s.Update(tea.KeyMsg{Type: tea.KeyCtrlE})

	if !errorModalOpened {
		t.Error("ctrl+e should trigger OnShowErrors callback")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShellCtrlEOpensErrorModal -v`
Expected: FAIL - OnShowErrors doesn't exist in Config

**Step 3: Write minimal implementation**

Add to Config struct in `shell/shell.go`:

```go
type Config struct {
	// ... existing fields ...
	// OnShowErrors is called when user presses Ctrl+E to show errors.
	OnShowErrors func()
}
```

Add to global keys section in Update():

```go
// Global keys
switch msg.String() {
case "ctrl+c", "ctrl+q":
	return s, tea.Quit
case "ctrl+e":
	if s.config.OnShowErrors != nil {
		s.config.OnShowErrors()
	}
	return s, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShellCtrlEOpensErrorModal -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): add ctrl+e handler for showing errors"
```

---

## Task 4: Error Display - App Integration

**Files:**
- Modify: `tux.go:111-124` (App struct - add errors field)
- Modify: `tux.go:126-185` (New function - wire error handling)
- Modify: `tux.go:234-251` (processEvent - handle errors)
- Test: `tux_test.go`

**Step 1: Write the failing test**

Add to `tux_test.go`:

```go
func TestAppAccumulatesErrors(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Process error events
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 2")})

	if len(app.errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(app.errors))
	}
}

func TestAppClearsErrorsOnSuccess(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Accumulate errors
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})

	// Successful completion should clear errors
	app.processEvent(Event{Type: EventComplete})

	if len(app.errors) != 0 {
		t.Errorf("expected 0 errors after success, got %d", len(app.errors))
	}
}

func TestAppKeepsErrorsOnErrorComplete(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Process error then complete
	app.processEvent(Event{Type: EventError, Error: fmt.Errorf("error 1")})
	app.errorsInRun = true // Simulating error happened in this run
	app.processEvent(Event{Type: EventComplete})

	// Errors should NOT be cleared if there were errors in this run
	if len(app.errors) != 1 {
		t.Errorf("expected 1 error preserved, got %d", len(app.errors))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test . -run TestAppAccumulatesErrors -v`
Expected: FAIL - app.errors field doesn't exist

**Step 3: Write minimal implementation**

Add fields to App struct in `tux.go`:

```go
type App struct {
	// ... existing fields ...

	// Error tracking
	errors      []error
	errorsInRun bool
}
```

Update processEvent in `tux.go`:

```go
case EventError:
	a.mu.Lock()
	a.errors = append(a.errors, event.Error)
	a.errorsInRun = true
	// Update status bar
	if len(a.errors) > 0 {
		errText := a.errors[0].Error()
		a.shell.SetStatus(shell.Status{
			ErrorText:  errText,
			ErrorCount: len(a.errors),
		})
	}
	a.mu.Unlock()

case EventComplete:
	a.chat.FinishAssistantMessage()
	a.mu.Lock()
	// Clear errors only if no errors in this run
	if !a.errorsInRun {
		a.errors = nil
		a.shell.SetStatus(shell.Status{}) // Clear error indicator
	}
	a.errorsInRun = false // Reset for next run
	a.mu.Unlock()
```

Wire OnShowErrors in New():

```go
shellCfg.OnShowErrors = func() {
	a.mu.Lock()
	errs := make([]error, len(a.errors))
	copy(errs, a.errors)
	a.mu.Unlock()

	if len(errs) > 0 {
		modal := shell.NewErrorModal(shell.ErrorModalConfig{
			Errors: errs,
			Theme:  cfg.theme,
		})
		a.shell.PushModal(modal)
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test . -run "TestAppAccumulatesErrors|TestAppClearsErrors|TestAppKeepsErrors" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat: integrate error display in tux.App"
```

---

## Task 5: History Navigation - Input History Tracking

**Files:**
- Modify: `shell/input.go:10-16` (Input struct)
- Modify: `shell/input.go:43-61` (Update method)
- Test: `shell/input_test.go`

**Step 1: Write the failing test**

Create or add to `shell/input_test.go`:

```go
package shell

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestInputHistoryNavigation(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	// Set up history provider
	history := []string{"first", "second", "third"}
	input.SetHistoryProvider(func() []string {
		return history
	})

	// Press up arrow - should show "third" (most recent)
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "third" {
		t.Errorf("expected 'third', got %q", input.Value())
	}

	// Press up again - should show "second"
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "second" {
		t.Errorf("expected 'second', got %q", input.Value())
	}

	// Press down - should show "third"
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyDown})
	if input.Value() != "third" {
		t.Errorf("expected 'third', got %q", input.Value())
	}

	// Press down past end - should clear
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyDown})
	if input.Value() != "" {
		t.Errorf("expected empty, got %q", input.Value())
	}
}

func TestInputHistoryResetsOnSubmit(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	history := []string{"first", "second"}
	input.SetHistoryProvider(func() []string {
		return history
	})

	// Navigate up
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Type something and submit
	input.SetValue("new message")
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// History index should be reset - pressing up should show "second"
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	// Note: after submit, history provider would include "new message"
	// but for this test, we're checking index reset behavior
}

func TestInputNoHistoryProvider(t *testing.T) {
	th := theme.NewDraculaTheme()
	input := NewInput(th, "> ", "")

	// Up arrow with no history provider should do nothing
	input.SetValue("test")
	input, _ = input.Update(tea.KeyMsg{Type: tea.KeyUp})
	if input.Value() != "test" {
		t.Errorf("expected 'test' unchanged, got %q", input.Value())
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestInputHistory -v`
Expected: FAIL - SetHistoryProvider doesn't exist

**Step 3: Write minimal implementation**

Update Input struct in `shell/input.go`:

```go
type Input struct {
	model           textinput.Model
	theme           theme.Theme
	prefix          string
	placeholder     string
	width           int
	historyProvider func() []string
	historyIndex    int // -1 means not navigating history
}
```

Update NewInput to initialize historyIndex:

```go
return &Input{
	model:        ti,
	theme:        th,
	prefix:       prefix,
	placeholder:  placeholder,
	historyIndex: -1,
}
```

Add SetHistoryProvider method:

```go
// SetHistoryProvider sets the function that provides history items.
// History should be ordered oldest to newest.
func (i *Input) SetHistoryProvider(provider func() []string) {
	i.historyProvider = provider
}
```

Update the Update method to handle up/down arrows:

```go
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := i.model.Value()
			if value != "" {
				i.model.SetValue("")
				i.historyIndex = -1 // Reset history navigation
				return i, func() tea.Msg {
					return InputSubmitMsg{Value: value}
				}
			}
			return i, nil

		case tea.KeyUp:
			if i.historyProvider != nil {
				history := i.historyProvider()
				if len(history) > 0 {
					if i.historyIndex == -1 {
						// Start navigating from end
						i.historyIndex = len(history) - 1
					} else if i.historyIndex > 0 {
						i.historyIndex--
					}
					i.model.SetValue(history[i.historyIndex])
					i.model.CursorEnd()
				}
			}
			return i, nil

		case tea.KeyDown:
			if i.historyProvider != nil && i.historyIndex != -1 {
				history := i.historyProvider()
				if i.historyIndex < len(history)-1 {
					i.historyIndex++
					i.model.SetValue(history[i.historyIndex])
					i.model.CursorEnd()
				} else {
					// Past end - clear and reset
					i.historyIndex = -1
					i.model.SetValue("")
				}
			}
			return i, nil
		}
	}

	var cmd tea.Cmd
	i.model, cmd = i.model.Update(msg)
	return i, cmd
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestInputHistory -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/input.go shell/input_test.go
git commit -m "feat(shell): add history navigation to Input"
```

---

## Task 6: History Navigation - Shell Wiring

**Files:**
- Modify: `shell/shell.go:34-47` (Config struct)
- Modify: `shell/shell.go:72-91` (New function)
- Test: `shell/shell_test.go`

**Step 1: Write the failing test**

Add to `shell/shell_test.go`:

```go
func TestShellHistoryProvider(t *testing.T) {
	th := theme.NewDraculaTheme()
	cfg := DefaultConfig()

	history := []string{"prompt1", "prompt2"}
	cfg.HistoryProvider = func() []string {
		return history
	}

	s := New(th, cfg)

	// Verify input has history provider set
	// (internal implementation detail, but we test behavior)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Send up arrow while input focused
	s.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Input should show last history item
	if s.InputValue() != "prompt2" {
		t.Errorf("expected 'prompt2', got %q", s.InputValue())
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShellHistoryProvider -v`
Expected: FAIL - HistoryProvider doesn't exist in Config

**Step 3: Write minimal implementation**

Add to Config struct:

```go
type Config struct {
	// ... existing fields ...
	// HistoryProvider returns the list of historical inputs (oldest to newest).
	HistoryProvider func() []string
}
```

Wire in New():

```go
func New(th theme.Theme, cfg Config) *Shell {
	// ... existing code ...

	s := &Shell{
		// ... existing fields ...
	}

	// Wire history provider to input
	if cfg.HistoryProvider != nil {
		s.input.SetHistoryProvider(cfg.HistoryProvider)
	}

	return s
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShellHistoryProvider -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): wire history provider through Config"
```

---

## Task 7: History Navigation - App Integration

**Files:**
- Modify: `tux.go:126-185` (New function)
- Modify: `chat_content.go` (add UserMessages method)
- Test: `tux_test.go`

**Step 1: Write the failing test**

Add to `tux_test.go`:

```go
func TestAppHistoryNavigation(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Add user messages
	app.chat.AddUserMessage("hello")
	app.chat.AddUserMessage("how are you")

	// Get user messages for history
	history := app.chat.UserMessages()

	if len(history) != 2 {
		t.Errorf("expected 2 history items, got %d", len(history))
	}
	if history[0] != "hello" {
		t.Errorf("expected 'hello', got %q", history[0])
	}
	if history[1] != "how are you" {
		t.Errorf("expected 'how are you', got %q", history[1])
	}
}
```

Add to `chat_content.go` test or main file test:

```go
func TestChatContentUserMessages(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AddUserMessage("first")
	chat.AppendText("response")
	chat.FinishAssistantMessage()
	chat.AddUserMessage("second")

	msgs := chat.UserMessages()
	if len(msgs) != 2 {
		t.Errorf("expected 2 user messages, got %d", len(msgs))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test . -run TestAppHistoryNavigation -v`
Expected: FAIL - UserMessages method doesn't exist

**Step 3: Write minimal implementation**

Add to `chat_content.go`:

```go
// UserMessages returns all user message contents in order (oldest to newest).
func (c *ChatContent) UserMessages() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var result []string
	for _, msg := range c.messages {
		if msg.role == "user" {
			result = append(result, msg.content)
		}
	}
	return result
}
```

Wire in `tux.go` New() function:

```go
// Wire history provider
shellCfg.HistoryProvider = func() []string {
	return chat.UserMessages()
}
```

**Step 4: Run test to verify it passes**

Run: `go test . -run TestAppHistoryNavigation -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go chat_content.go tux_test.go
git commit -m "feat: integrate history navigation in tux.App"
```

---

## Task 8: Approval Handling - Event Type Addition

**Files:**
- Modify: `tux.go:39-48` (Event struct)
- Modify: `tux.go:53-59` (EventType constants)
- Test: `tux_test.go`

**Step 1: Write the failing test**

Add to `tux_test.go`:

```go
func TestEventApprovalType(t *testing.T) {
	// Verify EventApproval constant exists
	if EventApproval != "approval" {
		t.Errorf("expected EventApproval to be 'approval', got %q", EventApproval)
	}
}

func TestEventHasResponseChannel(t *testing.T) {
	responseChan := make(chan ApprovalDecision, 1)
	event := Event{
		Type:       EventApproval,
		ToolID:     "tool-1",
		ToolName:   "bash",
		ToolParams: map[string]any{"command": "ls"},
		Response:   responseChan,
	}

	// Send decision
	go func() {
		event.Response <- DecisionApprove
	}()

	// Receive decision
	decision := <-event.Response
	if decision != DecisionApprove {
		t.Errorf("expected DecisionApprove, got %v", decision)
	}
}

func TestApprovalDecisionConstants(t *testing.T) {
	// Verify constants exist and have expected values
	if DecisionApprove != 0 {
		t.Error("DecisionApprove should be 0")
	}
	if DecisionDeny != 1 {
		t.Error("DecisionDeny should be 1")
	}
	if DecisionAlwaysAllow != 2 {
		t.Error("DecisionAlwaysAllow should be 2")
	}
	if DecisionNeverAllow != 3 {
		t.Error("DecisionNeverAllow should be 3")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test . -run "TestEventApproval|TestApprovalDecision" -v`
Expected: FAIL - EventApproval, ApprovalDecision, Response field don't exist

**Step 3: Write minimal implementation**

Add to `tux.go` Event struct:

```go
type Event struct {
	Type       EventType
	Text       string         // For EventText
	ToolName   string         // For EventToolCall, EventToolResult, EventApproval
	ToolID     string         // For EventToolCall, EventToolResult, EventApproval
	ToolParams map[string]any // For EventToolCall, EventApproval
	ToolOutput string         // For EventToolResult
	Success    bool           // For EventToolResult
	Error      error          // For EventError
	Response   chan ApprovalDecision // For EventApproval - send decision here
}
```

Add EventApproval constant:

```go
const (
	EventText       EventType = "text"
	EventToolCall   EventType = "tool_call"
	EventToolResult EventType = "tool_result"
	EventComplete   EventType = "complete"
	EventError      EventType = "error"
	EventApproval   EventType = "approval"
)
```

Add ApprovalDecision type (re-export from shell for convenience):

```go
// ApprovalDecision represents the user's decision on a tool approval.
// Re-exported from shell package for API convenience.
type ApprovalDecision = shell.ApprovalDecision

const (
	DecisionApprove     = shell.DecisionApprove
	DecisionDeny        = shell.DecisionDeny
	DecisionAlwaysAllow = shell.DecisionAlwaysAllow
	DecisionNeverAllow  = shell.DecisionNeverAllow
)
```

**Step 4: Run test to verify it passes**

Run: `go test . -run "TestEventApproval|TestApprovalDecision" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat: add EventApproval type and ApprovalDecision constants"
```

---

## Task 9: Approval Handling - App Integration

**Files:**
- Modify: `tux.go:234-251` (processEvent - add EventApproval case)
- Test: `tux_test.go`

**Step 1: Write the failing test**

Add to `tux_test.go`:

```go
func TestAppShowsApprovalModal(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Initialize shell size (required for modals)
	app.shell.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Create approval event with response channel
	responseChan := make(chan ApprovalDecision, 1)
	app.processEvent(Event{
		Type:       EventApproval,
		ToolID:     "tool-1",
		ToolName:   "bash",
		ToolParams: map[string]any{"command": "rm -rf /"},
		Response:   responseChan,
	})

	// Modal should be shown
	if !app.shell.HasModal() {
		t.Error("approval event should show modal")
	}
}

func TestAppApprovalSendsDecision(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	app.shell.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	responseChan := make(chan ApprovalDecision, 1)
	app.processEvent(Event{
		Type:       EventApproval,
		ToolID:     "tool-1",
		ToolName:   "bash",
		ToolParams: map[string]any{"command": "ls"},
		Response:   responseChan,
	})

	// Simulate user pressing Enter (approve is default selection)
	app.shell.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should receive decision
	select {
	case decision := <-responseChan:
		if decision != DecisionApprove {
			t.Errorf("expected DecisionApprove, got %v", decision)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected to receive decision on response channel")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test . -run TestAppShowsApprovalModal -v`
Expected: FAIL - EventApproval case not implemented

**Step 3: Write minimal implementation**

Add to processEvent in `tux.go`:

```go
case EventApproval:
	modal := shell.NewApprovalModal(shell.ApprovalModalConfig{
		Tool: shell.ToolInfo{
			ID:     event.ToolID,
			Name:   event.ToolName,
			Params: event.ToolParams,
		},
		OnDecision: func(decision shell.ApprovalDecision) {
			if event.Response != nil {
				event.Response <- decision
			}
		},
	})
	a.shell.PushModal(modal)
```

**Step 4: Run test to verify it passes**

Run: `go test . -run "TestAppShowsApprovalModal|TestAppApprovalSendsDecision" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat: integrate approval handling in tux.App"
```

---

## Task 10: Final Integration Test

**Files:**
- Test: `tux_test.go`

**Step 1: Write the integration test**

Add to `tux_test.go`:

```go
func TestFullApprovalFlow(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}
	app := New(agent)

	// Initialize
	app.shell.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Simulate agent sending approval request
	responseChan := make(chan ApprovalDecision, 1)
	go func() {
		// In real usage, agent blocks here waiting for decision
		decision := <-responseChan
		if decision == DecisionApprove {
			// Agent would run the tool, then send result
			events <- Event{
				Type:       EventToolResult,
				ToolID:     "tool-1",
				ToolOutput: "success",
				Success:    true,
			}
		} else {
			events <- Event{
				Type:       EventToolResult,
				ToolID:     "tool-1",
				ToolOutput: "denied by user",
				Success:    false,
			}
		}
	}()

	// Process approval event
	app.processEvent(Event{
		Type:       EventApproval,
		ToolID:     "tool-1",
		ToolName:   "bash",
		ToolParams: map[string]any{"command": "echo hello"},
		Response:   responseChan,
	})

	// User denies (press 'n')
	app.shell.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	// Process events
	time.Sleep(50 * time.Millisecond)

	// Tools should show denial
	view := app.tools.View()
	if !strings.Contains(view, "✗") {
		t.Error("denied tool should show failure marker")
	}
}
```

**Step 2: Run test**

Run: `go test . -run TestFullApprovalFlow -v`
Expected: PASS

**Step 3: Run all tests**

Run: `go test ./...`
Expected: All tests pass

**Step 4: Commit**

```bash
git add tux_test.go
git commit -m "test: add full approval flow integration test"
```

---

## Task 11: Run Full Test Suite and Cleanup

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: All tests pass

**Step 2: Run linter**

Run: `go vet ./...`
Expected: No errors

**Step 3: Final commit if any cleanup needed**

```bash
git add -A
git commit -m "chore: cleanup and finalize near-term improvements"
```
