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

import "context"

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
