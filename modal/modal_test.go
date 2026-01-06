package modal

import (
	"testing"

	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

func TestSimpleModal(t *testing.T) {
	closed := false
	m := NewSimpleModal(SimpleModalConfig{
		ID:      "test",
		Title:   "Test Modal",
		Content: content.NewSelectList(nil),
		Footer:  "Press Enter",
		OnClose: func() { closed = true },
	})

	if m.ID() != "test" {
		t.Errorf("expected id 'test', got %s", m.ID())
	}
	if m.Title() != "Test Modal" {
		t.Errorf("expected title 'Test Modal', got %s", m.Title())
	}

	m.OnPush(80, 24)
	m.OnPop()
	if !closed {
		t.Error("OnClose should have been called")
	}

	view := m.Render(60, 20)
	if view == "" {
		t.Error("render should not be empty")
	}
}

func TestConfirmModal(t *testing.T) {
	var result any
	m := NewConfirmModal(ConfirmModalConfig{
		ID:      "confirm",
		Title:   "Confirm",
		Message: "Are you sure?",
		Options: []ConfirmOption{
			{Label: "Yes", Value: true},
			{Label: "No", Value: false},
		},
		OnResult: func(v any) { result = v },
	})

	// Initial selection is first option
	if m.Selected() != 0 {
		t.Errorf("expected selected 0, got %d", m.Selected())
	}

	// Navigate down
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.Selected() != 1 {
		t.Errorf("expected selected 1, got %d", m.Selected())
	}

	// Select with enter
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestYesNoModal(t *testing.T) {
	var result bool
	m := NewYesNoModal("Confirm", "Continue?", func(v bool) { result = v })

	// Quick select with 'y'
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if !result {
		t.Error("y key should select yes")
	}
}

func TestListModal(t *testing.T) {
	var selected ListItem
	m := NewListModal(ListModalConfig{
		ID:    "list",
		Title: "Select Item",
		Items: []ListItem{
			{ID: "1", Title: "First", Description: "The first item"},
			{ID: "2", Title: "Second"},
			{ID: "3", Title: "Third"},
		},
		Filterable: true,
		OnSelect:   func(item ListItem) { selected = item },
	})

	// Navigate and select
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if selected.ID != "2" {
		t.Errorf("expected item 2, got %s", selected.ID)
	}

	// Test filtering
	m = NewListModal(ListModalConfig{
		ID:         "list2",
		Items:      []ListItem{{ID: "a", Title: "Apple"}, {ID: "b", Title: "Banana"}},
		Filterable: true,
	})

	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	if len(m.filtered) != 1 {
		t.Errorf("expected 1 filtered item, got %d", len(m.filtered))
	}
	if m.filtered[0].ID != "b" {
		t.Errorf("expected banana, got %s", m.filtered[0].Title)
	}
}

func TestWizardModal(t *testing.T) {
	var completed bool
	var results map[string]any

	selectList := content.NewSelectList([]content.SelectItem{
		{Label: "Option A", Value: "a"},
		{Label: "Option B", Value: "b"},
	})

	m := NewWizardModal(WizardModalConfig{
		ID:    "wizard",
		Title: "Setup Wizard",
		Steps: []WizardStep{
			{ID: "step1", Title: "Step 1", Content: selectList},
			{ID: "step2", Title: "Step 2", Content: content.NewSelectList(nil)},
		},
		OnComplete: func(r map[string]any) {
			completed = true
			results = r
		},
	})

	current, total := m.Progress()
	if current != 1 || total != 2 {
		t.Errorf("expected 1/2, got %d/%d", current, total)
	}

	// Advance to next step
	m.Next()
	current, _ = m.Progress()
	if current != 2 {
		t.Errorf("expected step 2, got %d", current)
	}

	// Go back
	m.Previous()
	current, _ = m.Progress()
	if current != 1 {
		t.Errorf("expected step 1 after back, got %d", current)
	}

	// Complete wizard
	m.Next() // step 1 -> 2
	m.Next() // step 2 -> complete

	if !completed {
		t.Error("wizard should be completed")
	}
	if results == nil {
		t.Error("results should not be nil")
	}
}

func TestApprovalModal(t *testing.T) {
	var decision ApprovalDecision
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{
			ID:      "tool1",
			Name:    "WriteFile",
			Params:  map[string]any{"path": "/tmp/test.txt"},
			Preview: "Will create a test file",
			Risk:    RiskMedium,
		},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})

	if m.Title() != "Tool Approval" {
		t.Errorf("expected 'Tool Approval', got %s", m.Title())
	}

	// Quick approve with 'y'
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if decision != DecisionApprove {
		t.Errorf("expected approve, got %d", decision)
	}

	// Reset and test deny
	decision = -1
	m = NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t2", Name: "Delete"},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})

	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if decision != DecisionDeny {
		t.Errorf("expected deny, got %d", decision)
	}
}

func TestApprovalModalRisk(t *testing.T) {
	tests := []struct {
		risk RiskLevel
	}{
		{RiskLow},
		{RiskMedium},
		{RiskHigh},
	}

	for _, tt := range tests {
		m := NewApprovalModal(ApprovalModalConfig{
			Tool: ToolInfo{ID: "t", Name: "Test", Risk: tt.risk},
		})
		m.OnPush(80, 24)
		view := m.Render(60, 20)
		if view == "" {
			t.Errorf("render should not be empty for risk %d", tt.risk)
		}
	}
}
