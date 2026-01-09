// Package tux provides a shared TUI library for building
// multi-agent terminal interfaces.
//
// The library provides:
//   - Shell: top-level container with tabs, modals, input, status
//   - Tabs: switchable content panes
//   - Modals: overlays including wizards, forms, approvals
//   - Content: composable primitives (viewport, lists, forms)
//   - Agent: backend interface for streaming, tool execution
package tux

import (
	"context"
	"sync"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/shell"
	"github.com/2389-research/tux/theme"
)

const Version = "0.1.0"

// Agent is the interface that agent implementations must satisfy.
// This includes orchestrator agents and any custom agent implementations.
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
	Text       string         // For EventText
	ToolName   string         // For EventToolCall, EventToolResult
	ToolID     string         // For EventToolCall, EventToolResult
	ToolParams map[string]any // For EventToolCall
	ToolOutput string         // For EventToolResult
	Success    bool           // For EventToolResult
	Error      error          // For EventError
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

// TabDef defines a custom tab.
type TabDef struct {
	ID       string
	Label    string
	Shortcut string          // e.g., "ctrl+d"
	Hidden   bool            // If true, tab is not shown in tab bar
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

// defaultAppConfig returns the default configuration with Dracula theme
// and no tabs removed.
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

// App is the main agent TUI application.
type App struct {
	agent  Agent
	shell  *shell.Shell
	config *appConfig

	// Content
	chat  *ChatContent
	tools *ToolsContent

	// Runtime state
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new App with the given agent and options.
func New(agent Agent, opts ...Option) *App {
	if agent == nil {
		panic("tux.New: agent cannot be nil")
	}

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

// Run starts the App.
func (a *App) Run() error {
	return a.shell.Run()
}

// submitInput starts an agent run with the given prompt.
func (a *App) submitInput(prompt string) {
	// Add user message to chat
	a.chat.AddUserMessage(prompt)

	// Create cancellable context (protected by mutex)
	a.mu.Lock()
	a.ctx, a.cancel = context.WithCancel(context.Background())
	ctx := a.ctx
	a.mu.Unlock()

	// Subscribe to events
	events := a.agent.Subscribe()

	// Run agent in goroutine
	go func() {
		// Start event processing
		go a.processEvents(events)

		// Run agent
		if err := a.agent.Run(ctx, prompt); err != nil {
			a.mu.Lock()
			cancelled := a.ctx.Err() != nil
			a.mu.Unlock()
			if !cancelled {
				// Not cancelled, real error
				a.processEvent(Event{Type: EventError, Error: err})
			}
		}
	}()
}

// processEvents reads events from the channel and routes them.
func (a *App) processEvents(events <-chan Event) {
	for event := range events {
		a.processEvent(event)
	}
}

// processEvent routes an agent event to the appropriate content.
// This is called for each event received from the agent's event channel.
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

// cancelRun cancels the current agent run.
// Thread-safe: acquires mutex before accessing ctx/cancel.
func (a *App) cancelRun() {
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
	a.mu.Unlock()

	// Call agent.Cancel() outside the mutex to avoid blocking while holding lock
	a.agent.Cancel()
}

// isRunning returns true if an agent run is in progress.
// Thread-safe: acquires mutex before accessing ctx.
func (a *App) isRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.ctx != nil && a.ctx.Err() == nil
}
