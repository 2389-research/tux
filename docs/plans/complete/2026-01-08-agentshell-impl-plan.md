# AgentShell Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement `tux.New(agent)` that creates a fully-wired agent TUI with Chat tab, Tools tab, approval flow, and event routing.

**Architecture:** `tux.New()` creates an `App` struct that wraps `shell.Shell` and manages agent integration. The App handles event routing from agent to tabs, approval flow bridging, and input submission.

**Tech Stack:** Go, Bubble Tea, existing tux shell components, mux orchestrator interface

---

## Task Overview

1. Define Agent interface in tux package
2. Define functional options (WithTheme, WithTab, etc.)
3. Create App struct and tux.New() constructor
4. Create ChatContent for the Chat tab
5. Create ToolsContent for the Tools tab
6. Wire event routing (agent events → tabs)
7. Wire approval flow (blocking → event-driven bridge)
8. Handle input submission and cancellation
9. Integration test with mock agent

---

### Task 1: Define Agent Interface

**Files:**
- Modify: `tux.go`
- Create: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go
package tux

import (
	"context"
	"testing"
)

// mockAgent implements Agent for testing
type mockAgent struct {
	events chan Event
}

func (m *mockAgent) Run(ctx context.Context, prompt string) error {
	return nil
}

func (m *mockAgent) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgent) Cancel() {}

func TestAgentInterfaceExists(t *testing.T) {
	events := make(chan Event)
	var agent Agent = &mockAgent{events: events}
	if agent == nil {
		t.Error("Agent interface should be implementable")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run TestAgentInterfaceExists -v`
Expected: FAIL with "undefined: Agent" or "undefined: Event"

**Step 3: Write minimal implementation**

```go
// tux.go - add to existing file
package tux

import "context"

const Version = "0.1.0"

// Agent is the interface that agent implementations must satisfy.
// mux's *agent.Agent satisfies this interface.
type Agent interface {
	// Run starts the agent with the given prompt.
	// It runs until completion or context cancellation.
	Run(ctx context.Context, prompt string) error

	// Subscribe returns a channel of events from the agent.
	// The channel is closed when the agent completes.
	Subscribe() <-chan Event

	// Cancel cancels the current agent run.
	Cancel()
}

// Event represents an event from the agent.
type Event struct {
	Type       EventType
	Text       string            // For EventText
	ToolName   string            // For EventToolCall, EventToolResult
	ToolID     string            // For EventToolCall, EventToolResult
	ToolParams map[string]any    // For EventToolCall
	ToolOutput string            // For EventToolResult
	Success    bool              // For EventToolResult
	Error      error             // For EventError
}

// EventType identifies the type of agent event.
type EventType string

const (
	EventText       EventType = "text"
	EventToolCall   EventType = "tool_call"
	EventToolResult EventType = "tool_result"
	EventComplete   EventType = "complete"
	EventError      EventType = "error"
)
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run TestAgentInterfaceExists -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): define Agent interface and Event types"
```

---

### Task 2: Define Functional Options

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add to existing
func TestWithThemeOption(t *testing.T) {
	th := theme.NewDraculaTheme()
	opt := WithTheme(th)
	if opt == nil {
		t.Error("WithTheme should return an option")
	}
}

func TestWithTabOption(t *testing.T) {
	tab := TabDef{
		ID:    "custom",
		Label: "Custom",
	}
	opt := WithTab(tab)
	if opt == nil {
		t.Error("WithTab should return an option")
	}
}

func TestWithoutTabOption(t *testing.T) {
	opt := WithoutTab("tools")
	if opt == nil {
		t.Error("WithoutTab should return an option")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestWith" -v`
Expected: FAIL with "undefined: WithTheme"

**Step 3: Write minimal implementation**

```go
// tux.go - add after Event types

import (
	"context"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
)

// TabDef defines a custom tab.
type TabDef struct {
	ID       string
	Label    string
	Shortcut string          // e.g., "ctrl+d"
	Hidden   bool
	Content  content.Content
}

// Option configures the App.
type Option func(*appConfig)

// appConfig holds configuration for the App.
type appConfig struct {
	theme       theme.Theme
	customTabs  []TabDef
	removedTabs map[string]bool
}

func defaultAppConfig() *appConfig {
	return &appConfig{
		theme:       theme.NewDraculaTheme(),
		removedTabs: make(map[string]bool),
	}
}

// WithTheme sets the theme for the App.
func WithTheme(th theme.Theme) Option {
	return func(c *appConfig) {
		c.theme = th
	}
}

// WithTab adds a custom tab to the App.
func WithTab(tab TabDef) Option {
	return func(c *appConfig) {
		c.customTabs = append(c.customTabs, tab)
	}
}

// WithoutTab removes a default tab from the App.
func WithoutTab(id string) Option {
	return func(c *appConfig) {
		c.removedTabs[id] = true
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestWith" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): add functional options (WithTheme, WithTab, WithoutTab)"
```

---

### Task 3: Create App Struct and Constructor

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add to existing
func TestNewApp(t *testing.T) {
	events := make(chan Event)
	agent := &mockAgent{events: events}

	app := New(agent)
	if app == nil {
		t.Error("New should return an App")
	}
}

func TestNewAppWithOptions(t *testing.T) {
	events := make(chan Event)
	agent := &mockAgent{events: events}

	app := New(agent,
		WithTheme(theme.NewNeoTerminalTheme()),
		WithoutTab("tools"),
	)
	if app == nil {
		t.Error("New with options should return an App")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestNewApp" -v`
Expected: FAIL with "undefined: New"

**Step 3: Write minimal implementation**

```go
// tux.go - add after options

import (
	"github.com/2389-research/tux/shell"
)

// App is the main agent TUI application.
type App struct {
	agent  Agent
	shell  *shell.Shell
	config *appConfig

	// Runtime state
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new App with the given agent and options.
func New(agent Agent, opts ...Option) *App {
	cfg := defaultAppConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	shellCfg := shell.DefaultConfig()
	sh := shell.New(cfg.theme, shellCfg)

	app := &App{
		agent:  agent,
		shell:  sh,
		config: cfg,
	}

	// Add default tabs (unless removed)
	if !cfg.removedTabs["chat"] {
		// Chat tab will be added in Task 4
	}
	if !cfg.removedTabs["tools"] {
		// Tools tab will be added in Task 5
	}

	// Add custom tabs
	for _, tab := range cfg.customTabs {
		sh.AddTab(shell.Tab{
			ID:       tab.ID,
			Label:    tab.Label,
			Shortcut: tab.Shortcut,
			Hidden:   tab.Hidden,
			Content:  tab.Content,
		})
	}

	return app
}

// Run starts the App.
func (a *App) Run() error {
	return a.shell.Run()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestNewApp" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): add App struct and New() constructor"
```

---

### Task 4: Create ChatContent

**Files:**
- Create: `chat_content.go`
- Create: `chat_content_test.go`

**Step 1: Write the failing test**

```go
// chat_content_test.go
package tux

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNewChatContent(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)
	if chat == nil {
		t.Error("NewChatContent should return a ChatContent")
	}
}

func TestChatContentAppendText(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AppendText("Hello ")
	chat.AppendText("World")

	view := chat.View()
	if view == "" {
		t.Error("View should contain text")
	}
}

func TestChatContentAddUserMessage(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AddUserMessage("What is 2+2?")

	view := chat.View()
	if view == "" {
		t.Error("View should contain user message")
	}
}

func TestChatContentFinishAssistant(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AppendText("The answer is 4")
	chat.FinishAssistantMessage()

	// Should be able to append new message
	chat.AppendText("New response")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestChatContent" -v`
Expected: FAIL with "undefined: NewChatContent"

**Step 3: Write minimal implementation**

```go
// chat_content.go
package tux

import (
	"strings"

	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatContent displays the conversation in the Chat tab.
type ChatContent struct {
	theme    theme.Theme
	messages []chatMessage
	current  strings.Builder // Current streaming message
	width    int
	height   int
}

type chatMessage struct {
	role    string // "user" or "assistant"
	content string
}

// NewChatContent creates a new ChatContent.
func NewChatContent(th theme.Theme) *ChatContent {
	return &ChatContent{
		theme:    th,
		messages: make([]chatMessage, 0),
	}
}

// Init implements content.Content.
func (c *ChatContent) Init() tea.Cmd {
	return nil
}

// Update implements content.Content.
func (c *ChatContent) Update(msg tea.Msg) (interface{}, tea.Cmd) {
	return c, nil
}

// View implements content.Content.
func (c *ChatContent) View() string {
	var parts []string

	styles := c.theme.Styles()

	for _, msg := range c.messages {
		var style lipgloss.Style
		if msg.role == "user" {
			style = styles.UserMessage
		} else {
			style = styles.AssistantMessage
		}
		parts = append(parts, style.Render(msg.content))
	}

	// Add current streaming message if any
	if c.current.Len() > 0 {
		parts = append(parts, styles.AssistantMessage.Render(c.current.String()))
	}

	return strings.Join(parts, "\n\n")
}

// Value implements content.Content.
func (c *ChatContent) Value() any {
	return c.messages
}

// SetSize implements content.Content.
func (c *ChatContent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// AppendText appends streaming text to the current assistant message.
func (c *ChatContent) AppendText(text string) {
	c.current.WriteString(text)
}

// AddUserMessage adds a user message to the conversation.
func (c *ChatContent) AddUserMessage(content string) {
	c.messages = append(c.messages, chatMessage{
		role:    "user",
		content: content,
	})
}

// FinishAssistantMessage completes the current streaming message.
func (c *ChatContent) FinishAssistantMessage() {
	if c.current.Len() > 0 {
		c.messages = append(c.messages, chatMessage{
			role:    "assistant",
			content: c.current.String(),
		})
		c.current.Reset()
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestChatContent" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add chat_content.go chat_content_test.go
git commit -m "feat(tux): add ChatContent for streaming conversation display"
```

---

### Task 5: Create ToolsContent

**Files:**
- Create: `tools_content.go`
- Create: `tools_content_test.go`

**Step 1: Write the failing test**

```go
// tools_content_test.go
package tux

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNewToolsContent(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)
	if tools == nil {
		t.Error("NewToolsContent should return a ToolsContent")
	}
}

func TestToolsContentAddCall(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})

	view := tools.View()
	if view == "" {
		t.Error("View should contain tool call")
	}
}

func TestToolsContentAddResult(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})
	tools.AddToolResult("tool-1", "file contents here", true)

	view := tools.View()
	if view == "" {
		t.Error("View should contain tool result")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestToolsContent" -v`
Expected: FAIL with "undefined: NewToolsContent"

**Step 3: Write minimal implementation**

```go
// tools_content.go
package tux

import (
	"fmt"
	"strings"
	"time"

	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// ToolsContent displays the tool call timeline in the Tools tab.
type ToolsContent struct {
	theme  theme.Theme
	items  []toolItem
	width  int
	height int
}

type toolItem struct {
	id        string
	name      string
	params    map[string]any
	output    string
	success   bool
	completed bool
	timestamp time.Time
}

// NewToolsContent creates a new ToolsContent.
func NewToolsContent(th theme.Theme) *ToolsContent {
	return &ToolsContent{
		theme: th,
		items: make([]toolItem, 0),
	}
}

// Init implements content.Content.
func (c *ToolsContent) Init() tea.Cmd {
	return nil
}

// Update implements content.Content.
func (c *ToolsContent) Update(msg tea.Msg) (interface{}, tea.Cmd) {
	return c, nil
}

// View implements content.Content.
func (c *ToolsContent) View() string {
	if len(c.items) == 0 {
		return "No tool calls yet"
	}

	var parts []string
	styles := c.theme.Styles()

	for _, item := range c.items {
		var status string
		if item.completed {
			if item.success {
				status = "✓"
			} else {
				status = "✗"
			}
		} else {
			status = "⋯"
		}

		line := fmt.Sprintf("%s %s", status, item.name)
		if item.completed && item.output != "" {
			// Truncate output for display
			output := item.output
			if len(output) > 50 {
				output = output[:47] + "..."
			}
			line += fmt.Sprintf(" → %s", output)
		}

		parts = append(parts, styles.Text.Render(line))
	}

	return strings.Join(parts, "\n")
}

// Value implements content.Content.
func (c *ToolsContent) Value() any {
	return c.items
}

// SetSize implements content.Content.
func (c *ToolsContent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// AddToolCall adds a tool call to the timeline.
func (c *ToolsContent) AddToolCall(id, name string, params map[string]any) {
	c.items = append(c.items, toolItem{
		id:        id,
		name:      name,
		params:    params,
		timestamp: time.Now(),
	})
}

// AddToolResult adds a result to an existing tool call.
func (c *ToolsContent) AddToolResult(id, output string, success bool) {
	for i := range c.items {
		if c.items[i].id == id {
			c.items[i].output = output
			c.items[i].success = success
			c.items[i].completed = true
			return
		}
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestToolsContent" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tools_content.go tools_content_test.go
git commit -m "feat(tux): add ToolsContent for tool call timeline"
```

---

### Task 6: Wire Default Tabs in App

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add to existing
func TestAppHasDefaultTabs(t *testing.T) {
	events := make(chan Event)
	agent := &mockAgent{events: events}

	app := New(agent)

	// App should have chat content accessible
	if app.chat == nil {
		t.Error("App should have chat content")
	}
	if app.tools == nil {
		t.Error("App should have tools content")
	}
}

func TestAppWithoutToolsTab(t *testing.T) {
	events := make(chan Event)
	agent := &mockAgent{events: events}

	app := New(agent, WithoutTab("tools"))

	if app.chat == nil {
		t.Error("App should still have chat content")
	}
	// tools tab should be removed but tools content exists for event routing
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestAppHas" -v`
Expected: FAIL with "app.chat undefined"

**Step 3: Write minimal implementation**

Update the App struct and New() to include chat and tools content:

```go
// tux.go - update App struct
type App struct {
	agent  Agent
	shell  *shell.Shell
	config *appConfig

	// Content
	chat  *ChatContent
	tools *ToolsContent

	// Runtime state
	ctx    context.Context
	cancel context.CancelFunc
}

// Update New() function
func New(agent Agent, opts ...Option) *App {
	cfg := defaultAppConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	shellCfg := shell.DefaultConfig()
	sh := shell.New(cfg.theme, shellCfg)

	// Create content
	chat := NewChatContent(cfg.theme)
	tools := NewToolsContent(cfg.theme)

	app := &App{
		agent:  agent,
		shell:  sh,
		config: cfg,
		chat:   chat,
		tools:  tools,
	}

	// Add default tabs (unless removed)
	if !cfg.removedTabs["chat"] {
		sh.AddTab(shell.Tab{
			ID:      "chat",
			Label:   "Chat",
			Content: chat,
		})
	}
	if !cfg.removedTabs["tools"] {
		sh.AddTab(shell.Tab{
			ID:       "tools",
			Label:    "Tools",
			Shortcut: "ctrl+o",
			Content:  tools,
		})
	}

	// Add custom tabs
	for _, tab := range cfg.customTabs {
		sh.AddTab(shell.Tab{
			ID:       tab.ID,
			Label:    tab.Label,
			Shortcut: tab.Shortcut,
			Hidden:   tab.Hidden,
			Content:  tab.Content,
		})
	}

	return app
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestAppHas" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): wire default Chat and Tools tabs in App"
```

---

### Task 7: Wire Event Routing

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add to existing
func TestAppRoutesTextEvents(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Send a text event
	events <- Event{Type: EventText, Text: "Hello"}

	// Process the event
	app.processEvent(Event{Type: EventText, Text: "Hello"})

	// Chat should have the text
	view := app.chat.View()
	if !strings.Contains(view, "Hello") {
		t.Errorf("Chat should contain 'Hello', got: %s", view)
	}
}

func TestAppRoutesToolCallEvents(t *testing.T) {
	events := make(chan Event, 10)
	agent := &mockAgent{events: events}

	app := New(agent)

	// Process tool call event
	app.processEvent(Event{
		Type:       EventToolCall,
		ToolID:     "tool-1",
		ToolName:   "read_file",
		ToolParams: map[string]any{"path": "/test"},
	})

	// Tools should have the call
	view := app.tools.View()
	if !strings.Contains(view, "read_file") {
		t.Errorf("Tools should contain 'read_file', got: %s", view)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestAppRoutes" -v`
Expected: FAIL with "app.processEvent undefined"

**Step 3: Write minimal implementation**

```go
// tux.go - add processEvent method
func (a *App) processEvent(event Event) {
	switch event.Type {
	case EventText:
		a.chat.AppendText(event.Text)

	case EventToolCall:
		a.tools.AddToolCall(event.ToolID, event.ToolName, event.ToolParams)

	case EventToolResult:
		a.tools.AddToolResult(event.ToolID, event.ToolOutput, event.Success)

	case EventComplete:
		a.chat.FinishAssistantMessage()

	case EventError:
		// TODO: Show error in status bar or modal
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestAppRoutes" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): add event routing from agent to tabs"
```

---

### Task 8: Handle Input Submission

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add
func TestAppSubmitInput(t *testing.T) {
	events := make(chan Event, 10)
	runCalled := false
	agent := &mockAgentWithRun{
		events: events,
		onRun: func(prompt string) {
			runCalled = true
			if prompt != "Hello" {
				t.Errorf("Expected prompt 'Hello', got '%s'", prompt)
			}
		},
	}

	app := New(agent)
	app.submitInput("Hello")

	// Give goroutine time to start
	time.Sleep(10 * time.Millisecond)

	if !runCalled {
		t.Error("Agent.Run should have been called")
	}
}

type mockAgentWithRun struct {
	events chan Event
	onRun  func(string)
}

func (m *mockAgentWithRun) Run(ctx context.Context, prompt string) error {
	if m.onRun != nil {
		m.onRun(prompt)
	}
	return nil
}

func (m *mockAgentWithRun) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgentWithRun) Cancel() {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestAppSubmitInput" -v`
Expected: FAIL with "app.submitInput undefined"

**Step 3: Write minimal implementation**

```go
// tux.go - add submitInput method and update Run
func (a *App) submitInput(prompt string) {
	// Add user message to chat
	a.chat.AddUserMessage(prompt)

	// Create cancellable context
	a.ctx, a.cancel = context.WithCancel(context.Background())

	// Subscribe to events
	events := a.agent.Subscribe()

	// Run agent in goroutine
	go func() {
		// Start event processing
		go a.processEvents(events)

		// Run agent
		if err := a.agent.Run(a.ctx, prompt); err != nil {
			if a.ctx.Err() == nil {
				// Not cancelled, real error
				a.processEvent(Event{Type: EventError, Error: err})
			}
		}
	}()
}

func (a *App) processEvents(events <-chan Event) {
	for event := range events {
		a.processEvent(event)
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestAppSubmitInput" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): add input submission with agent.Run"
```

---

### Task 9: Handle Cancellation

**Files:**
- Modify: `tux.go`
- Modify: `tux_test.go`

**Step 1: Write the failing test**

```go
// tux_test.go - add
func TestAppCancelOnEscape(t *testing.T) {
	events := make(chan Event, 10)
	cancelCalled := false
	agent := &mockAgentWithCancel{
		events: events,
		onCancel: func() {
			cancelCalled = true
		},
	}

	app := New(agent)

	// Start a run
	app.submitInput("Hello")
	time.Sleep(10 * time.Millisecond)

	// Cancel
	app.cancelRun()

	if !cancelCalled {
		t.Error("Agent.Cancel should have been called")
	}
}

type mockAgentWithCancel struct {
	events   chan Event
	onCancel func()
}

func (m *mockAgentWithCancel) Run(ctx context.Context, prompt string) error {
	<-ctx.Done()
	return ctx.Err()
}

func (m *mockAgentWithCancel) Subscribe() <-chan Event {
	return m.events
}

func (m *mockAgentWithCancel) Cancel() {
	if m.onCancel != nil {
		m.onCancel()
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestAppCancel" -v`
Expected: FAIL with "app.cancelRun undefined"

**Step 3: Write minimal implementation**

```go
// tux.go - add cancelRun method
func (a *App) cancelRun() {
	if a.cancel != nil {
		a.cancel()
	}
	a.agent.Cancel()
}

// isRunning returns true if an agent run is in progress
func (a *App) isRunning() bool {
	return a.ctx != nil && a.ctx.Err() == nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run "TestAppCancel" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tux.go tux_test.go
git commit -m "feat(tux): add cancellation support"
```

---

### Task 10: Integration Test

**Files:**
- Create: `tux_integration_test.go`

**Step 1: Write the integration test**

```go
// tux_integration_test.go
package tux

import (
	"context"
	"testing"
	"time"

	"github.com/2389-research/tux/theme"
)

func TestFullConversationFlow(t *testing.T) {
	// Create mock agent that simulates a conversation
	events := make(chan Event, 10)
	agent := &conversationAgent{events: events}

	app := New(agent, WithTheme(theme.NewDraculaTheme()))

	// Submit a prompt
	app.submitInput("What is 2+2?")

	// Simulate agent response
	go func() {
		time.Sleep(10 * time.Millisecond)
		events <- Event{Type: EventText, Text: "The answer is "}
		events <- Event{Type: EventText, Text: "4"}
		events <- Event{Type: EventComplete}
		close(events)
	}()

	// Wait for processing
	time.Sleep(50 * time.Millisecond)

	// Verify chat has content
	chatView := app.chat.View()
	if chatView == "" {
		t.Error("Chat should have content")
	}
}

type conversationAgent struct {
	events chan Event
}

func (a *conversationAgent) Run(ctx context.Context, prompt string) error {
	// Events are sent by test
	return nil
}

func (a *conversationAgent) Subscribe() <-chan Event {
	return a.events
}

func (a *conversationAgent) Cancel() {}
```

**Step 2: Run test**

Run: `go test ./... -run "TestFullConversation" -v`
Expected: PASS

**Step 3: Commit**

```bash
git add tux_integration_test.go
git commit -m "test(tux): add integration test for conversation flow"
```

---

### Task 11: Final Verification

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: All PASS

**Step 2: Run build**

Run: `go build ./...`
Expected: Success, no errors

**Step 3: Final commit if needed**

If any cleanup needed, commit it.

---

## Summary

After completing all tasks, `tux.New(agent)` will:

1. Create an App with Chat and Tools tabs wired
2. Route agent events to appropriate tabs
3. Handle input submission (calls agent.Run)
4. Handle cancellation (Escape key)
5. Support customization via functional options

**Not yet implemented (future work):**
- Approval flow (blocking → event-driven bridge)
- Help modal with Ctrl+H
- Onboarding wizard
- Mouse support enhancements
