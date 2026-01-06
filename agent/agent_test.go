package agent

import (
	"errors"
	"testing"
)

func TestToolQueue(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "Read", Risk: RiskLow},
		{ID: "2", Name: "Write", Risk: RiskMedium},
		{ID: "3", Name: "Delete", Risk: RiskHigh},
	}

	q := NewQueue(tools, nil)

	if q.Count() != 3 {
		t.Errorf("expected 3 tools, got %d", q.Count())
	}
	if q.Current() != 0 {
		t.Errorf("expected current 0, got %d", q.Current())
	}
	if q.IsComplete() {
		t.Error("queue should not be complete")
	}

	// Process first tool
	item := q.Next()
	if item.Tool.Name != "Read" {
		t.Errorf("expected Read, got %s", item.Tool.Name)
	}

	q.SetOutcome(OutcomeApproved, nil)
	q.Advance()

	// Process second tool
	item = q.Next()
	if item.Tool.Name != "Write" {
		t.Errorf("expected Write, got %s", item.Tool.Name)
	}

	q.SetOutcome(OutcomeDenied, nil)
	q.Advance()

	// Process third tool
	item = q.Next()
	if item.Tool.Name != "Delete" {
		t.Errorf("expected Delete, got %s", item.Tool.Name)
	}

	q.SetOutcome(OutcomeExecuted, &ToolResult{ToolUseID: "3", Content: "deleted"})
	q.Advance()

	// Queue should be complete
	if !q.IsComplete() {
		t.Error("queue should be complete")
	}
	if q.Next() != nil {
		t.Error("next should return nil when complete")
	}

	// Check counts
	if q.ApprovedCount() != 2 {
		t.Errorf("expected 2 approved, got %d", q.ApprovedCount())
	}
	if q.DeniedCount() != 1 {
		t.Errorf("expected 1 denied, got %d", q.DeniedCount())
	}
}

func TestToolQueueProgressHint(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
		{ID: "3", Name: "C"},
	}

	q := NewQueue(tools, nil)

	// Initial state
	hint := q.ProgressHint()
	if hint != "[●??]" {
		t.Errorf("expected [●??], got %s", hint)
	}

	// After first approved
	q.SetOutcome(OutcomeApproved, nil)
	q.Advance()
	hint = q.ProgressHint()
	if hint != "[✓●?]" {
		t.Errorf("expected [✓●?], got %s", hint)
	}

	// After second denied
	q.SetOutcome(OutcomeDenied, nil)
	q.Advance()
	hint = q.ProgressHint()
	if hint != "[✓✗●]" {
		t.Errorf("expected [✓✗●], got %s", hint)
	}

	// After complete
	q.SetOutcome(OutcomeExecuted, nil)
	q.Advance()
	hint = q.ProgressHint()
	if hint != "[✓✗✓]" {
		t.Errorf("expected [✓✗✓], got %s", hint)
	}
}

func TestToolQueueClassifier(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "Read", Risk: RiskLow},
		{ID: "2", Name: "Write", Risk: RiskMedium},
	}

	q := NewQueue(tools, RiskBasedClassifier)

	// Low risk should be auto-approved
	item := q.Next()
	if item.Action != ActionAutoApprove {
		t.Errorf("expected auto-approve for low risk, got %d", item.Action)
	}

	q.Advance()

	// Medium risk should need approval
	item = q.Next()
	if item.Action != ActionNeedsApproval {
		t.Errorf("expected needs approval for medium risk, got %d", item.Action)
	}
}

func TestToolQueueReset(t *testing.T) {
	tools := []ToolInfo{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	q := NewQueue(tools, nil)

	q.SetOutcome(OutcomeApproved, nil)
	q.Advance()
	q.SetOutcome(OutcomeDenied, nil)
	q.Advance()

	if !q.IsComplete() {
		t.Error("should be complete")
	}

	q.Reset()

	if q.IsComplete() {
		t.Error("should not be complete after reset")
	}
	if q.Current() != 0 {
		t.Errorf("expected current 0, got %d", q.Current())
	}
	if q.ApprovedCount() != 0 {
		t.Errorf("expected 0 approved, got %d", q.ApprovedCount())
	}
}

func TestToolQueueResults(t *testing.T) {
	tools := []ToolInfo{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	q := NewQueue(tools, nil)

	q.SetOutcome(OutcomeExecuted, &ToolResult{ToolUseID: "1", Content: "result1"})
	q.Advance()
	q.SetOutcome(OutcomeExecuted, &ToolResult{ToolUseID: "2", Content: "result2"})
	q.Advance()

	results := q.Results()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0].Content != "result1" {
		t.Errorf("expected result1, got %s", results[0].Content)
	}
}

func TestEventCreators(t *testing.T) {
	// Text event
	e := NewTextEvent("hello")
	if e.Type != EventText || e.Text != "hello" {
		t.Error("text event incorrect")
	}

	// Tool call event
	tool := ToolUse{ID: "1", Name: "test"}
	e = NewToolCallEvent(tool)
	if e.Type != EventToolCall || e.Tool.Name != "test" {
		t.Error("tool call event incorrect")
	}

	// Tool result event
	result := ToolResult{ToolUseID: "1", Content: "done"}
	e = NewToolResultEvent(result)
	if e.Type != EventToolResult || e.Result.Content != "done" {
		t.Error("tool result event incorrect")
	}

	// Complete event
	usage := TokenUsage{InputTokens: 100, OutputTokens: 50}
	e = NewCompleteEvent(usage)
	if e.Type != EventComplete || e.Usage.InputTokens != 100 {
		t.Error("complete event incorrect")
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		ID:      "1",
		Role:    RoleUser,
		Content: "Hello",
	}

	if msg.Role != RoleUser {
		t.Error("role should be user")
	}

	msg.ContentBlocks = []ContentBlock{
		{Type: "text", Text: "Hello"},
		{Type: "tool_use", ToolUse: &ToolUse{ID: "t1", Name: "test"}},
	}

	if len(msg.ContentBlocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(msg.ContentBlocks))
	}
}

func TestToolQueuePendingCount(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
		{ID: "3", Name: "C"},
	}
	q := NewQueue(tools, nil)

	if q.PendingCount() != 3 {
		t.Errorf("expected 3 pending, got %d", q.PendingCount())
	}

	q.SetOutcome(OutcomeApproved, nil)
	q.Advance()

	if q.PendingCount() != 2 {
		t.Errorf("expected 2 pending, got %d", q.PendingCount())
	}
}

func TestToolQueueItems(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
	}
	q := NewQueue(tools, nil)

	items := q.Items()
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if items[0].Tool.Name != "A" {
		t.Errorf("expected first item name 'A', got %s", items[0].Tool.Name)
	}
}

func TestDefaultClassifier(t *testing.T) {
	tool := ToolInfo{ID: "1", Name: "Test", Risk: RiskHigh}
	action, reason := DefaultClassifier(tool)

	if action != ActionNeedsApproval {
		t.Error("DefaultClassifier should always return ActionNeedsApproval")
	}
	if reason != "" {
		t.Error("DefaultClassifier should return empty reason")
	}
}

func TestNewErrorEvent(t *testing.T) {
	err := errors.New("test error")
	e := NewErrorEvent(err)

	if e.Type != EventError {
		t.Error("type should be EventError")
	}
	if e.Error != err {
		t.Error("error should be set")
	}
}

func TestProgressHintEmpty(t *testing.T) {
	q := NewQueue(nil, nil)
	hint := q.ProgressHint()
	if hint != "" {
		t.Errorf("expected empty hint for empty queue, got %s", hint)
	}
}

func TestProgressHintAllStates(t *testing.T) {
	tools := []ToolInfo{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
		{ID: "3", Name: "C"},
		{ID: "4", Name: "D"},
	}
	q := NewQueue(tools, nil)

	// Set first to approved
	q.SetOutcome(OutcomeApproved, nil)
	q.Advance()

	// Set second to denied
	q.SetOutcome(OutcomeDenied, nil)
	q.Advance()

	// Set third to error
	q.SetOutcome(OutcomeError, nil)
	q.Advance()

	// Fourth is current (pending)
	hint := q.ProgressHint()
	if hint != "[✓✗✗●]" {
		t.Errorf("expected [✓✗✗●], got %s", hint)
	}
}
