// Package agent provides the backend interface for agent integration.
package agent

import (
	"context"
)

// Backend is the interface that agent implementations must satisfy.
// It abstracts the underlying LLM/agent system (like mux) from the TUI.
type Backend interface {
	// Stream starts a conversation turn and returns a channel of events.
	// The channel is closed when the turn completes.
	Stream(ctx context.Context, messages []Message) (<-chan Event, error)

	// ExecuteTool executes a tool after approval.
	ExecuteTool(ctx context.Context, tool ToolUse) (ToolResult, error)

	// ApprovalRequests returns a channel of approval requests from the backend.
	ApprovalRequests() <-chan ApprovalRequest

	// RespondToApproval sends an approval decision back to the backend.
	RespondToApproval(requestID string, decision ApprovalDecision) error

	// DescribeTool returns metadata about a tool for display.
	DescribeTool(name string) ToolDescription

	// Cancel cancels the current operation.
	Cancel()
}

// ApprovalRequest represents a request for user approval.
type ApprovalRequest struct {
	ID       string
	Tool     ToolInfo
	Response chan<- ApprovalDecision
}

// ApprovalDecision represents the user's decision on a tool approval.
type ApprovalDecision int

const (
	// Approve allows the tool to run this time.
	Approve ApprovalDecision = iota
	// Deny skips the tool this time.
	Deny
	// AlwaysAllow remembers to always allow this tool.
	AlwaysAllow
	// NeverAllow remembers to never allow this tool.
	NeverAllow
)

// ToolInfo contains information about a tool.
type ToolInfo struct {
	ID      string
	Name    string
	Params  map[string]any
	Preview string
	Risk    RiskLevel
}

// RiskLevel indicates the risk level of a tool.
type RiskLevel int

const (
	// RiskLow indicates minimal risk.
	RiskLow RiskLevel = iota
	// RiskMedium indicates moderate risk.
	RiskMedium
	// RiskHigh indicates significant risk.
	RiskHigh
)

// ToolDescription contains metadata about a tool.
type ToolDescription struct {
	Name        string
	Description string
	Category    string
	Risk        RiskLevel
}
