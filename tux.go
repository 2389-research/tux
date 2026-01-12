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
	"fmt"
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
	Text       string                  // For EventText
	ToolName   string                  // For EventToolCall, EventToolResult, EventApproval
	ToolID     string                  // For EventToolCall, EventToolResult, EventApproval
	ToolParams map[string]any          // For EventToolCall, EventApproval
	ToolOutput string                  // For EventToolResult
	Success    bool                    // For EventToolResult
	Error      error                   // For EventError
	Response   chan ApprovalDecision   // For EventApproval - send decision here
}

// EventType identifies the type of agent event.
type EventType string

const (
	EventText       EventType = "text"
	EventToolCall   EventType = "tool_call"
	EventToolResult EventType = "tool_result"
	EventComplete   EventType = "complete"
	EventError      EventType = "error"
	EventApproval   EventType = "approval"
)

// ApprovalDecision represents the user's decision on a tool approval.
// Re-exported from shell package for API convenience.
type ApprovalDecision = shell.ApprovalDecision

const (
	DecisionApprove     = shell.DecisionApprove
	DecisionDeny        = shell.DecisionDeny
	DecisionAlwaysAllow = shell.DecisionAlwaysAllow
	DecisionNeverAllow  = shell.DecisionNeverAllow
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
	theme          theme.Theme
	customTabs     []TabDef
	removedTabs    map[string]bool
	helpCategories []shell.Category
	autocomplete   *shell.Autocomplete
	suggestions    *shell.Suggestions
	onQuickActions func()
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

// HelpCategory is a re-export of shell.Category for API convenience.
type HelpCategory = shell.Category

// HelpBinding is a re-export of shell.Binding for API convenience.
type HelpBinding = shell.Binding

// WithHelpCategories sets the help overlay categories.
// When set, pressing '?' shows the help overlay with these keybindings.
func WithHelpCategories(categories ...HelpCategory) Option {
	return func(c *appConfig) {
		c.helpCategories = categories
	}
}

// Autocomplete is a re-export of shell.Autocomplete for API convenience.
type Autocomplete = shell.Autocomplete

// Completion is a re-export of shell.Completion for API convenience.
type Completion = shell.Completion

// CompletionProvider is a re-export of shell.CompletionProvider for API convenience.
type CompletionProvider = shell.CompletionProvider

// ListItem is a re-export of shell.ListItem for API convenience.
type ListItem = shell.ListItem

// ListModal is a re-export of shell.ListModal for API convenience.
type ListModal = shell.ListModal

// ListModalConfig is a re-export of shell.ListModalConfig for API convenience.
type ListModalConfig = shell.ListModalConfig

// NewListModal creates a new list modal for command palette functionality.
func NewListModal(cfg ListModalConfig) *ListModal {
	return shell.NewListModal(cfg)
}

// Suggestion is a re-export of shell.Suggestion for API convenience.
type Suggestion = shell.Suggestion

// SuggestionProvider is a re-export of shell.SuggestionProvider for API convenience.
type SuggestionProvider = shell.SuggestionProvider

// Suggestions is a re-export of shell.Suggestions for API convenience.
type Suggestions = shell.Suggestions

// NewSuggestions creates a new suggestions component.
func NewSuggestions() *Suggestions {
	return shell.NewSuggestions()
}

// NewAutocomplete creates a new autocomplete component.
func NewAutocomplete() *Autocomplete {
	return shell.NewAutocomplete()
}

// NewCommandProvider creates a command completion provider.
func NewCommandProvider(commands []Completion) *shell.CommandProvider {
	return shell.NewCommandProvider(commands)
}

// NewHistoryProvider creates a history completion provider.
func NewHistoryProvider(history []string) *shell.HistoryProvider {
	return shell.NewHistoryProvider(history)
}

// WithAutocomplete sets the autocomplete component for the input.
// When set, Tab triggers completion suggestions.
func WithAutocomplete(ac *Autocomplete) Option {
	return func(c *appConfig) {
		c.autocomplete = ac
	}
}

// WithQuickActions sets the callback for the ':' key quick actions.
// When set, pressing ':' with empty input opens quick actions.
func WithQuickActions(fn func()) Option {
	return func(c *appConfig) {
		c.onQuickActions = fn
	}
}

// WithSuggestions sets the suggestions component for the input.
// When set, suggestions are analyzed on each input change.
func WithSuggestions(s *Suggestions) Option {
	return func(c *appConfig) {
		c.suggestions = s
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

	// Error tracking
	errors      []error
	errorsInRun bool
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

	// Create content
	chat := NewChatContent(cfg.theme)
	tools := NewToolsContent(cfg.theme)

	app := &App{
		agent:  agent,
		config: cfg,
		chat:   chat,
		tools:  tools,
	}

	// Wire input submission to agent
	shellCfg.OnInputSubmit = app.submitInput

	// Wire history provider
	shellCfg.HistoryProvider = func() []string {
		return chat.UserMessages()
	}

	// Wire error display
	shellCfg.OnShowErrors = func() {
		app.mu.Lock()
		errs := make([]error, len(app.errors))
		copy(errs, app.errors)
		app.mu.Unlock()

		if len(errs) > 0 {
			modal := shell.NewErrorModal(shell.ErrorModalConfig{
				Errors: errs,
				Theme:  cfg.theme,
			})
			app.shell.PushModal(modal)
		}
	}

	// Wire help categories
	shellCfg.HelpCategories = cfg.helpCategories

	// Wire autocomplete
	shellCfg.Autocomplete = cfg.autocomplete

	// Wire suggestions
	shellCfg.Suggestions = cfg.suggestions

	// Wire quick actions
	shellCfg.OnQuickActions = cfg.onQuickActions

	sh := shell.New(cfg.theme, shellCfg)
	app.shell = sh

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

	// Cancel any existing run before starting a new one
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel()
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	ctx := a.ctx
	a.mu.Unlock()

	// Subscribe to events
	events := a.agent.Subscribe()

	// Start streaming display
	a.shell.Streaming().Start()

	// Run agent in goroutine
	go func() {
		// Start event processing
		go a.processEvents(events)

		// Run agent
		if err := a.agent.Run(ctx, prompt); err != nil {
			// Use local ctx (not a.ctx) to check if THIS run was cancelled
			if ctx.Err() == nil {
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
	streaming := a.shell.Streaming()

	switch event.Type {
	case EventText:
		a.chat.AppendText(event.Text)
		streaming.AppendToken(event.Text)

	case EventToolCall:
		a.tools.AddToolCall(event.ToolID, event.ToolName, event.ToolParams)
		streaming.StartToolCall(event.ToolID, event.ToolName)

	case EventToolResult:
		a.tools.AddToolResult(event.ToolID, event.ToolOutput, event.Success)
		streaming.EndToolCall(event.ToolID)

	case EventComplete:
		a.chat.FinishAssistantMessage()
		streaming.End()
		a.mu.Lock()
		// Clear errors only if no errors in this run
		if !a.errorsInRun {
			a.errors = nil
			a.shell.SetStatus(shell.Status{}) // Clear error indicator
		}
		a.errorsInRun = false // Reset for next run
		a.mu.Unlock()

	case EventError:
		a.mu.Lock()
		// Handle nil error defensively
		err := event.Error
		if err == nil {
			err = fmt.Errorf("unknown error")
		}
		a.errors = append(a.errors, err)
		a.errorsInRun = true
		errCount := len(a.errors)
		errText := a.errors[0].Error()
		a.mu.Unlock()
		// Update status bar outside mutex
		a.shell.SetStatus(shell.Status{
			ErrorText:  errText,
			ErrorCount: errCount,
		})

	case EventApproval:
		modal := shell.NewApprovalModal(shell.ApprovalModalConfig{
			Tool: shell.ToolInfo{
				ID:     event.ToolID,
				Name:   event.ToolName,
				Params: event.ToolParams,
			},
			OnDecision: func(decision shell.ApprovalDecision) {
				if event.Response != nil {
					// Send asynchronously to avoid blocking UI thread
					go func(d shell.ApprovalDecision) {
						event.Response <- d
					}(decision)
				}
			},
		})
		a.shell.PushModal(modal)
	}

	// Trigger UI refresh after external state change
	a.shell.Send(shell.RefreshMsg{})
}

// cancelRun cancels the current agent run.
// Thread-safe: acquires mutex before accessing ctx/cancel.
func (a *App) cancelRun() {
	a.mu.Lock()
	shouldCancel := a.cancel != nil
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
	a.mu.Unlock()

	// Call agent.Cancel() outside the mutex to avoid blocking while holding lock
	if shouldCancel {
		a.agent.Cancel()
	}
}

// isRunning returns true if an agent run is in progress.
// Thread-safe: acquires mutex before accessing ctx.
func (a *App) isRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.ctx != nil && a.ctx.Err() == nil
}

// ClearChat clears all messages from the chat display.
// Use this when starting a new session or before loading a different session.
func (a *App) ClearChat() {
	a.chat.Clear()
}

// AddChatUserMessage adds a user message to the chat display.
// Use this when restoring a previous session.
func (a *App) AddChatUserMessage(content string) {
	a.chat.AddUserMessage(content)
}

// AddChatAssistantMessage adds an assistant message to the chat display.
// Use this when restoring a previous session.
func (a *App) AddChatAssistantMessage(content string) {
	a.chat.AddAssistantMessage(content)
}

// PushModal pushes a modal onto the modal stack.
// Use this to show command palettes, dialogs, or other modal overlays.
func (a *App) PushModal(m shell.Modal) {
	a.shell.PushModal(m)
}

// PopModal pops the top modal from the stack.
func (a *App) PopModal() {
	a.shell.PopModal()
}

// SetInputValue sets the input text.
// Use this when applying quick action values to the input.
func (a *App) SetInputValue(value string) {
	a.shell.SetInputValue(value)
}
