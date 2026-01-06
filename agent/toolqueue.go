package agent

import "strings"

// ToolAction indicates how a tool should be handled.
type ToolAction int

const (
	// ActionNeedsApproval requires user approval.
	ActionNeedsApproval ToolAction = iota
	// ActionAutoApprove automatically approves.
	ActionAutoApprove
	// ActionAutoDeny automatically denies.
	ActionAutoDeny
)

// ToolOutcome indicates the result of tool processing.
type ToolOutcome int

const (
	// OutcomePending indicates the tool hasn't been processed.
	OutcomePending ToolOutcome = iota
	// OutcomeApproved indicates the tool was approved.
	OutcomeApproved
	// OutcomeDenied indicates the tool was denied.
	OutcomeDenied
	// OutcomeExecuted indicates the tool was executed.
	OutcomeExecuted
	// OutcomeError indicates an error occurred.
	OutcomeError
)

// ToolDisposition represents a tool and its processing state.
type ToolDisposition struct {
	Tool    ToolInfo
	Action  ToolAction
	Outcome ToolOutcome
	Result  *ToolResult
	Reason  string
}

// Classifier determines how a tool should be handled.
type Classifier func(tool ToolInfo) (ToolAction, string)

// Queue manages sequential tool processing with pre-classification.
type Queue struct {
	items      []ToolDisposition
	current    int
	classifier Classifier
}

// NewQueue creates a new tool queue.
func NewQueue(tools []ToolInfo, classifier Classifier) *Queue {
	q := &Queue{
		items:      make([]ToolDisposition, len(tools)),
		current:    0,
		classifier: classifier,
	}

	// Pre-classify all tools
	for i, tool := range tools {
		action, reason := ActionNeedsApproval, ""
		if classifier != nil {
			action, reason = classifier(tool)
		}
		q.items[i] = ToolDisposition{
			Tool:    tool,
			Action:  action,
			Outcome: OutcomePending,
			Reason:  reason,
		}
	}

	return q
}

// Next returns the next tool disposition to process.
// Returns nil if all tools have been processed.
func (q *Queue) Next() *ToolDisposition {
	if q.current >= len(q.items) {
		return nil
	}
	return &q.items[q.current]
}

// Advance moves to the next tool.
func (q *Queue) Advance() {
	if q.current < len(q.items) {
		q.current++
	}
}

// SetOutcome sets the outcome for the current tool.
func (q *Queue) SetOutcome(outcome ToolOutcome, result *ToolResult) {
	if q.current < len(q.items) {
		q.items[q.current].Outcome = outcome
		q.items[q.current].Result = result
	}
}

// IsComplete returns true if all tools have been processed.
func (q *Queue) IsComplete() bool {
	return q.current >= len(q.items)
}

// ProgressHint returns a string showing queue progress.
// Format: "[●✓✗?]" where:
// ● = current, ✓ = approved/executed, ✗ = denied/error, ? = pending
func (q *Queue) ProgressHint() string {
	if len(q.items) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("[")

	for i, item := range q.items {
		if i == q.current {
			b.WriteString("●")
		} else {
			switch item.Outcome {
			case OutcomePending:
				b.WriteString("?")
			case OutcomeApproved, OutcomeExecuted:
				b.WriteString("✓")
			case OutcomeDenied, OutcomeError:
				b.WriteString("✗")
			}
		}
	}

	b.WriteString("]")
	return b.String()
}

// Results returns all tool results.
func (q *Queue) Results() []ToolResult {
	var results []ToolResult
	for _, item := range q.items {
		if item.Result != nil {
			results = append(results, *item.Result)
		}
	}
	return results
}

// Count returns the total number of tools.
func (q *Queue) Count() int {
	return len(q.items)
}

// Current returns the current index.
func (q *Queue) Current() int {
	return q.current
}

// ApprovedCount returns the number of approved/executed tools.
func (q *Queue) ApprovedCount() int {
	count := 0
	for _, item := range q.items {
		if item.Outcome == OutcomeApproved || item.Outcome == OutcomeExecuted {
			count++
		}
	}
	return count
}

// DeniedCount returns the number of denied tools.
func (q *Queue) DeniedCount() int {
	count := 0
	for _, item := range q.items {
		if item.Outcome == OutcomeDenied {
			count++
		}
	}
	return count
}

// PendingCount returns the number of pending tools.
func (q *Queue) PendingCount() int {
	count := 0
	for _, item := range q.items {
		if item.Outcome == OutcomePending {
			count++
		}
	}
	return count
}

// Items returns all items in the queue.
func (q *Queue) Items() []ToolDisposition {
	return q.items
}

// Reset resets the queue to the beginning.
func (q *Queue) Reset() {
	q.current = 0
	for i := range q.items {
		q.items[i].Outcome = OutcomePending
		q.items[i].Result = nil
	}
}

// DefaultClassifier is a simple classifier that requires approval for all tools.
func DefaultClassifier(tool ToolInfo) (ToolAction, string) {
	return ActionNeedsApproval, ""
}

// RiskBasedClassifier auto-approves low-risk tools.
func RiskBasedClassifier(tool ToolInfo) (ToolAction, string) {
	switch tool.Risk {
	case RiskLow:
		return ActionAutoApprove, "low risk"
	default:
		return ActionNeedsApproval, ""
	}
}
