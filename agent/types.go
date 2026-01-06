package agent

import "time"

// Message represents a conversation message.
type Message struct {
	ID            string
	Role          Role
	Content       string
	ContentBlocks []ContentBlock
	Timestamp     time.Time
}

// Role represents the role of a message sender.
type Role string

const (
	// RoleUser indicates a user message.
	RoleUser Role = "user"
	// RoleAssistant indicates an assistant message.
	RoleAssistant Role = "assistant"
	// RoleTool indicates a tool result message.
	RoleTool Role = "tool"
	// RoleSystem indicates a system message.
	RoleSystem Role = "system"
)

// ContentBlock represents a block of content within a message.
type ContentBlock struct {
	Type       string // "text", "tool_use", "tool_result"
	Text       string
	ToolUse    *ToolUse
	ToolResult *ToolResult
}

// ToolUse represents a tool invocation.
type ToolUse struct {
	ID    string
	Name  string
	Input map[string]any
}

// ToolResult represents the result of a tool invocation.
type ToolResult struct {
	ToolUseID string
	Content   string
	IsError   bool
}

// Event represents an event from the agent backend.
type Event struct {
	Type   EventType
	Text   string       // For EventText
	Tool   *ToolUse     // For EventToolCall
	Result *ToolResult  // For EventToolResult
	Error  error        // For EventError
	Usage  *TokenUsage  // For EventComplete
}

// EventType identifies the type of agent event.
type EventType int

const (
	// EventText indicates streaming text content.
	EventText EventType = iota
	// EventToolCall indicates a tool invocation request.
	EventToolCall
	// EventToolResult indicates a tool result.
	EventToolResult
	// EventComplete indicates the turn is complete.
	EventComplete
	// EventError indicates an error occurred.
	EventError
)

// TokenUsage contains token usage statistics.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	CacheHits    int
}

// NewTextEvent creates a text event.
func NewTextEvent(text string) Event {
	return Event{Type: EventText, Text: text}
}

// NewToolCallEvent creates a tool call event.
func NewToolCallEvent(tool ToolUse) Event {
	return Event{Type: EventToolCall, Tool: &tool}
}

// NewToolResultEvent creates a tool result event.
func NewToolResultEvent(result ToolResult) Event {
	return Event{Type: EventToolResult, Result: &result}
}

// NewCompleteEvent creates a completion event.
func NewCompleteEvent(usage TokenUsage) Event {
	return Event{Type: EventComplete, Usage: &usage}
}

// NewErrorEvent creates an error event.
func NewErrorEvent(err error) Event {
	return Event{Type: EventError, Error: err}
}
