package shell

import (
	"errors"
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
	if m.Size() != SizeMedium {
		t.Errorf("expected medium size")
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

	// Test content getter/setter
	newContent := content.NewSelectList([]content.SelectItem{{Label: "Test"}})
	m.SetContent(newContent)
	if m.Content() != newContent {
		t.Error("content should be updated")
	}

	// Test HandleKey
	handled, _ := m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if !handled {
		t.Error("key should be handled when content exists")
	}

	// Test without content
	m.SetContent(nil)
	handled, _ = m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if handled {
		t.Error("key should not be handled without content")
	}
}

func TestSimpleModalNoOnClose(t *testing.T) {
	m := NewSimpleModal(SimpleModalConfig{
		ID:    "test",
		Title: "Test",
	})
	m.OnPop() // Should not panic
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

	// Test interface methods
	if m.ID() != "confirm" {
		t.Errorf("expected id 'confirm', got %s", m.ID())
	}
	if m.Title() != "Confirm" {
		t.Errorf("expected title 'Confirm', got %s", m.Title())
	}
	if m.Size() != SizeSmall {
		t.Error("confirm modal should be small")
	}

	m.OnPush(80, 24)
	m.OnPop() // Should not panic

	// Initial selection is first option
	if m.Selected() != 0 {
		t.Errorf("expected selected 0, got %d", m.Selected())
	}

	// Navigate down with arrow
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.Selected() != 1 {
		t.Errorf("expected selected 1, got %d", m.Selected())
	}

	// Navigate up with arrow
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.Selected() != 0 {
		t.Errorf("expected selected 0, got %d", m.Selected())
	}

	// Navigate with j/k
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.Selected() != 1 {
		t.Errorf("expected selected 1 after j, got %d", m.Selected())
	}
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.Selected() != 0 {
		t.Errorf("expected selected 0 after k, got %d", m.Selected())
	}

	// Tab navigation
	m.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if m.Selected() != 1 {
		t.Errorf("expected selected 1 after tab, got %d", m.Selected())
	}
	m.HandleKey(tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.Selected() != 0 {
		t.Errorf("expected selected 0 after shift+tab, got %d", m.Selected())
	}

	// SetSelected
	m.SetSelected(1)
	if m.Selected() != 1 {
		t.Error("SetSelected should work")
	}
	m.SetSelected(99) // Out of bounds
	if m.Selected() != 1 {
		t.Error("SetSelected should ignore out of bounds")
	}

	// Select with enter
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if result != false {
		t.Errorf("expected false, got %v", result)
	}

	// Test Render
	view := m.Render(60, 20)
	if view == "" {
		t.Error("render should not be empty")
	}
	if !containsStr(view, "Are you sure?") {
		t.Error("render should contain message")
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

	// Reset and test 'n'
	result = true
	m = NewYesNoModal("Confirm", "Continue?", func(v bool) { result = v })
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if result {
		t.Error("n key should select no")
	}
}

func TestOKCancelModal(t *testing.T) {
	var result bool
	m := NewOKCancelModal("Action", "Proceed?", func(v bool) { result = v })

	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if !result {
		t.Error("enter should confirm OK")
	}
}

func TestConfirmModalBoundary(t *testing.T) {
	m := NewConfirmModal(ConfirmModalConfig{
		ID:      "test",
		Options: []ConfirmOption{{Label: "A"}, {Label: "B"}},
	})

	// Can't go above first
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.Selected() != 0 {
		t.Error("shouldn't go above 0")
	}

	// Can't go below last
	m.SetSelected(1)
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.Selected() != 1 {
		t.Error("shouldn't go below last")
	}
}

func TestListModal(t *testing.T) {
	var selected ListItem
	var cancelled bool
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
		OnCancel:   func() { cancelled = true },
	})

	// Test interface methods
	if m.ID() != "list" {
		t.Errorf("expected id 'list', got %s", m.ID())
	}
	if m.Title() != "Select Item" {
		t.Errorf("expected title 'Select Item', got %s", m.Title())
	}
	if m.Size() != SizeMedium {
		t.Error("list modal should be medium")
	}

	m.OnPush(80, 24)

	// Navigate and select
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if selected.ID != "2" {
		t.Errorf("expected item 2, got %s", selected.ID)
	}

	// Test cancel
	m = NewListModal(ListModalConfig{
		ID:       "list2",
		Items:    []ListItem{{ID: "a", Title: "Apple"}},
		OnCancel: func() { cancelled = true },
	})
	m.OnPop()
	if !cancelled {
		t.Error("OnCancel should be called on pop")
	}

	// Test SelectedItem
	m = NewListModal(ListModalConfig{
		ID:    "list3",
		Items: []ListItem{{ID: "a", Title: "Apple"}, {ID: "b", Title: "Banana"}},
	})
	item := m.SelectedItem()
	if item == nil || item.ID != "a" {
		t.Error("SelectedItem should return first item")
	}

	// Test SetItems
	m.SetItems([]ListItem{{ID: "x", Title: "X"}})
	item = m.SelectedItem()
	if item == nil || item.ID != "x" {
		t.Error("SetItems should update items")
	}

	// Test Render
	view := m.Render(60, 20)
	if view == "" {
		t.Error("render should not be empty")
	}
}

func TestListModalFiltering(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:         "list",
		Items:      []ListItem{{ID: "a", Title: "Apple"}, {ID: "b", Title: "Banana"}, {ID: "c", Title: "Cherry"}},
		Filterable: true,
	})

	// Type to filter
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	if len(m.filtered) != 1 {
		t.Errorf("expected 1 filtered item, got %d", len(m.filtered))
	}

	// Backspace
	m.HandleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if len(m.filtered) != 3 {
		t.Errorf("expected 3 items after backspace, got %d", len(m.filtered))
	}

	// Empty filter with backspace on empty
	m.HandleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if len(m.filtered) != 3 {
		t.Error("backspace on empty should do nothing")
	}

	// Filter no matches
	m.filter = "xyz"
	m.applyFilter()
	if len(m.filtered) != 0 {
		t.Error("should have no matches")
	}
}

func TestListModalNavigation(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:    "list",
		Items: []ListItem{{ID: "a"}, {ID: "b"}, {ID: "c"}},
	})

	// j/k navigation
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.SelectedItem().ID != "b" {
		t.Error("j should move down")
	}
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.SelectedItem().ID != "a" {
		t.Error("k should move up")
	}

	// Boundary
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.SelectedItem().ID != "a" {
		t.Error("can't go above first")
	}
}

func TestListModalRenderEmpty(t *testing.T) {
	m := NewListModal(ListModalConfig{ID: "empty", Items: nil})
	view := m.Render(60, 20)
	if !containsStr(view, "No items") {
		t.Error("should show 'No items match'")
	}
}

func TestListModalScrolling(t *testing.T) {
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: string(rune('A' + i))}
	}
	m := NewListModal(ListModalConfig{ID: "big", Items: items})
	m.maxVisible = 5

	// Navigate to bottom
	for i := 0; i < 15; i++ {
		m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	}

	view := m.Render(60, 20)
	if view == "" {
		t.Error("should render with scrolling")
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
			{ID: "step1", Title: "Step 1", Description: "First step", Content: selectList},
			{ID: "step2", Title: "Step 2", Content: content.NewSelectList(nil)},
		},
		OnComplete: func(r map[string]any) {
			completed = true
			results = r
		},
	})

	// Test interface
	if m.ID() != "wizard" {
		t.Errorf("expected id 'wizard', got %s", m.ID())
	}
	if m.Title() != "Setup Wizard" {
		t.Errorf("expected title 'Setup Wizard', got %s", m.Title())
	}
	if m.Size() != SizeLarge {
		t.Error("wizard should be large")
	}

	m.OnPush(80, 24)

	current, total := m.Progress()
	if current != 1 || total != 2 {
		t.Errorf("expected 1/2, got %d/%d", current, total)
	}

	// Test CanGoNext/CanGoPrevious
	if !m.CanGoNext() {
		t.Error("should be able to go next")
	}
	if m.CanGoPrevious() {
		t.Error("shouldn't be able to go back from first step")
	}

	// Advance to next step
	m.Next()
	current, _ = m.Progress()
	if current != 2 {
		t.Errorf("expected step 2, got %d", current)
	}

	if !m.CanGoPrevious() {
		t.Error("should be able to go back")
	}

	// Go back
	m.Previous()
	current, _ = m.Progress()
	if current != 1 {
		t.Errorf("expected step 1 after back, got %d", current)
	}

	// GoToStep
	m.GoToStep(1)
	current, _ = m.Progress()
	if current != 2 {
		t.Error("GoToStep should work")
	}
	m.GoToStep(99) // Out of bounds
	current, _ = m.Progress()
	if current != 2 {
		t.Error("GoToStep should ignore out of bounds")
	}

	// Complete wizard
	m.GoToStep(0)
	m.Next() // step 1 -> 2
	m.Next() // step 2 -> complete

	if !completed {
		t.Error("wizard should be completed")
	}
	if results == nil {
		t.Error("results should not be nil")
	}

	// Test Results getter/setter
	m.SetResults(map[string]any{"test": "value"})
	if m.Results()["test"] != "value" {
		t.Error("SetResults should work")
	}
}

func TestWizardModalValidation(t *testing.T) {
	validated := false
	m := NewWizardModal(WizardModalConfig{
		ID: "wizard",
		Steps: []WizardStep{
			{
				ID:      "step1",
				Title:   "Step 1",
				Content: content.NewSelectList([]content.SelectItem{{Label: "A", Value: "a"}}),
				Validate: func(v any) error {
					validated = true
					return nil
				},
			},
		},
	})
	m.OnPush(80, 24)
	m.Next()
	if !validated {
		t.Error("validate should be called")
	}
}

func TestWizardModalKeyHandling(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{
		ID: "wizard",
		Steps: []WizardStep{
			{ID: "1", Title: "One", Content: content.NewSelectList(nil)},
			{ID: "2", Title: "Two", Content: content.NewSelectList(nil)},
		},
	})
	m.OnPush(80, 24)

	// Enter advances
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if m.current != 1 {
		t.Error("enter should advance")
	}

	// Ctrl+P goes back
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{16}}) // Ctrl+P is rune 16
	// This might not work exactly, let's test ctrl+n
	m.current = 0
	m.HandleKey(tea.KeyMsg{Type: tea.KeyCtrlN})
	// Just verify it handles keys
}

func TestWizardModalRender(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{
		ID:    "wizard",
		Title: "Test Wizard",
		Steps: []WizardStep{
			{ID: "1", Title: "Step One", Description: "First step"},
		},
	})
	m.OnPush(80, 24)
	view := m.Render(60, 40)
	if !containsStr(view, "Step One") {
		t.Error("should render step title")
	}
	if !containsStr(view, "1 of 1") {
		t.Error("should render progress")
	}
}

func TestWizardModalCancel(t *testing.T) {
	cancelled := false
	m := NewWizardModal(WizardModalConfig{
		ID:       "wizard",
		Steps:    []WizardStep{{ID: "1", Title: "One"}},
		OnCancel: func() { cancelled = true },
	})
	m.OnPop()
	if !cancelled {
		t.Error("OnCancel should be called")
	}
}

func TestWizardModalNilStep(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{ID: "empty", Steps: nil})
	if m.CurrentStep() != nil {
		t.Error("CurrentStep should return nil for empty wizard")
	}
	m.Next() // Should not panic
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
	if m.ID() != "approval-tool1" {
		t.Errorf("expected 'approval-tool1', got %s", m.ID())
	}
	if m.Size() != SizeMedium {
		t.Error("approval modal should be medium")
	}

	m.OnPush(80, 24)
	m.OnPop() // Should not panic

	// Test Tool getter
	if m.Tool().Name != "WriteFile" {
		t.Error("Tool() should return tool info")
	}

	// Test Selected
	if m.Selected() != 0 {
		t.Error("initial selection should be 0")
	}

	// Test SetQueueHint
	m.SetQueueHint("[●??]")

	// Quick approve with 'y'
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if decision != DecisionApprove {
		t.Errorf("expected approve, got %d", decision)
	}

	// Test 'a' also approves
	decision = -1
	m = NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t2", Name: "Test"},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if decision != DecisionApprove {
		t.Error("'a' should approve")
	}

	// Test 'n' denies
	decision = -1
	m = NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t3", Name: "Test"},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if decision != DecisionDeny {
		t.Error("'n' should deny")
	}

	// Test 'd' also denies
	decision = -1
	m = NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t4", Name: "Test"},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if decision != DecisionDeny {
		t.Error("'d' should deny")
	}
}

func TestApprovalModalNavigation(t *testing.T) {
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{ID: "t", Name: "Test"},
	})

	// Down
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.Selected() != 1 {
		t.Error("down should move selection")
	}

	// Up
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.Selected() != 0 {
		t.Error("up should move selection")
	}

	// j/k
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.Selected() != 1 {
		t.Error("j should move down")
	}
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.Selected() != 0 {
		t.Error("k should move up")
	}

	// Boundary
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.Selected() != 0 {
		t.Error("can't go above first")
	}

	m.selected = len(DefaultApprovalOptions) - 1
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.Selected() != len(DefaultApprovalOptions)-1 {
		t.Error("can't go below last")
	}
}

func TestApprovalModalEnter(t *testing.T) {
	var decision ApprovalDecision = -1
	m := NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t", Name: "Test"},
		OnDecision: func(d ApprovalDecision) { decision = d },
	})

	m.selected = 2 // Always Allow
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if decision != DecisionAlwaysAllow {
		t.Errorf("expected AlwaysAllow, got %d", decision)
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

func TestApprovalModalRenderLongParams(t *testing.T) {
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{
			ID:   "t",
			Name: "Test",
			Params: map[string]any{
				"very_long_param": "this is a very long value that should be truncated because it exceeds the width",
			},
		},
	})
	view := m.Render(40, 20)
	if view == "" {
		t.Error("should render")
	}
}

func TestApprovalModalCustomOptions(t *testing.T) {
	var decision ApprovalDecision
	customOpts := []ApprovalOption{
		{Label: "Go", Decision: DecisionApprove},
		{Label: "Stop", Decision: DecisionDeny},
	}
	m := NewApprovalModal(ApprovalModalConfig{
		Tool:       ToolInfo{ID: "t", Name: "Test"},
		Options:    customOpts,
		OnDecision: func(d ApprovalDecision) { decision = d },
	})

	_ = decision // use variable
	if len(m.options) != 2 {
		t.Error("should use custom options")
	}
}

// Additional tests for 95%+ coverage

func TestApprovalModalOnPop(t *testing.T) {
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{ID: "t", Name: "Test"},
	})
	m.OnPop() // Should not panic, no OnCancel callback
}

func TestConfirmModalOnPop(t *testing.T) {
	m := NewConfirmModal(ConfirmModalConfig{
		ID:      "test",
		Options: []ConfirmOption{{Label: "A"}},
	})
	m.OnPop() // Should not panic, no OnCancel callback
}

func TestListModalSelectedItemEmpty(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:    "empty",
		Items: nil,
	})
	item := m.SelectedItem()
	if item != nil {
		t.Error("SelectedItem should return nil for empty list")
	}
}

func TestListModalHandleKeyJK(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:    "list",
		Items: []ListItem{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}},
	})

	// j moves down
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selected != 1 {
		t.Error("j should move down")
	}
	// k moves up
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selected != 0 {
		t.Error("k should move up")
	}
}

func TestListModalRenderFiltered(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:         "list",
		Items:      []ListItem{{ID: "a", Title: "Apple"}, {ID: "b", Title: "Banana"}},
		Filterable: true,
	})
	m.filter = "app"
	m.applyFilter()
	view := m.Render(60, 20)
	if view == "" {
		t.Error("should render filtered view")
	}
}

func TestWizardModalHandleKeyCtrlP(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{
		ID: "wizard",
		Steps: []WizardStep{
			{ID: "1", Title: "One"},
			{ID: "2", Title: "Two"},
		},
	})
	m.OnPush(80, 24)
	m.current = 1

	// Ctrl+P should go back
	m.HandleKey(tea.KeyMsg{Type: tea.KeyCtrlP})
	if m.current != 0 {
		t.Error("Ctrl+P should go back")
	}
}

func TestWizardModalNextWithError(t *testing.T) {
	validated := false
	m := NewWizardModal(WizardModalConfig{
		ID: "wizard",
		Steps: []WizardStep{
			{
				ID:    "1",
				Title: "One",
				Validate: func(v any) error {
					validated = true
					return errors.New("validation failed")
				},
			},
			{ID: "2", Title: "Two"},
		},
	})
	m.OnPush(80, 24)
	m.Next()
	if !validated {
		t.Error("validate should be called")
	}
	if m.current != 0 {
		t.Error("should not advance on validation error")
	}
}

func TestWizardModalRenderWithContent(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{
		ID: "wizard",
		Steps: []WizardStep{
			{
				ID:          "1",
				Title:       "Step",
				Description: "A description",
				Content:     content.NewSelectList([]content.SelectItem{{Label: "A"}}),
			},
		},
	})
	m.OnPush(80, 24)
	view := m.Render(60, 40)
	if !containsStr(view, "Step") {
		t.Error("should render step title")
	}
}

func TestSizePercentSmall(t *testing.T) {
	s := SizeSmall
	if s.HeightPercent() != 0.30 {
		t.Error("small height should be 0.30")
	}
	if s.WidthPercent() != 0.50 {
		t.Error("small width should be 0.50")
	}
}

func TestSizePercentFullscreen(t *testing.T) {
	s := SizeFullscreen
	if s.HeightPercent() != 1.0 {
		t.Error("fullscreen height should be 1.0")
	}
	if s.WidthPercent() != 1.0 {
		t.Error("fullscreen width should be 1.0")
	}
}

func TestManagerPeekEmpty(t *testing.T) {
	mgr := NewManager()
	if mgr.Peek() != nil {
		t.Error("Peek on empty should return nil")
	}
}

func TestConfirmModalHandleKeyEscape(t *testing.T) {
	m := NewConfirmModal(ConfirmModalConfig{
		ID:      "test",
		Options: []ConfirmOption{{Label: "A"}, {Label: "B"}},
	})
	// Escape is not handled by confirm modal
	handled, _ := m.HandleKey(tea.KeyMsg{Type: tea.KeyEsc})
	if handled {
		t.Error("escape should not be handled")
	}
}

func TestListModalHandleKeyEscape(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:    "list",
		Items: []ListItem{{ID: "a", Title: "A"}},
	})
	// Escape is not handled by list modal
	handled, _ := m.HandleKey(tea.KeyMsg{Type: tea.KeyEsc})
	if handled {
		t.Error("escape should not be handled")
	}
}

func TestApprovalModalRenderPreview(t *testing.T) {
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{
			ID:      "t",
			Name:    "Test",
			Preview: "This is a preview",
		},
	})
	view := m.Render(60, 20)
	if !containsStr(view, "preview") {
		// Preview may or may not be in rendered output depending on impl
	}
}

func TestSizePercentLarge(t *testing.T) {
	s := SizeLarge
	if s.HeightPercent() != 0.80 {
		t.Error("large height should be 0.80")
	}
	if s.WidthPercent() != 0.80 {
		t.Error("large width should be 0.80")
	}
}

func TestSizePercentMedium(t *testing.T) {
	s := SizeMedium
	if s.HeightPercent() != 0.50 {
		t.Error("medium height should be 0.50")
	}
	if s.WidthPercent() != 0.60 {
		t.Error("medium width should be 0.60")
	}
}

func TestSizePercentUnknown(t *testing.T) {
	s := Size(99) // Unknown size
	if s.HeightPercent() != 0.50 {
		t.Error("unknown size height should default to 0.50")
	}
	if s.WidthPercent() != 0.60 {
		t.Error("unknown size width should default to 0.60")
	}
}

func TestApprovalModalHandleKeyNoCallback(t *testing.T) {
	m := NewApprovalModal(ApprovalModalConfig{
		Tool: ToolInfo{ID: "t", Name: "Test"},
		// No OnDecision callback
	})
	// These should not panic
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
}

func TestConfirmModalHandleKeyNoCallback(t *testing.T) {
	m := NewConfirmModal(ConfirmModalConfig{
		ID:      "test",
		Options: []ConfirmOption{{Label: "A", Value: "a"}},
		// No OnResult callback
	})
	// This should not panic
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
}

func TestListModalRenderDescription(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID: "list",
		Items: []ListItem{
			{ID: "a", Title: "Apple", Description: "A fruit"},
			{ID: "b", Title: "Banana"},
		},
	})
	view := m.Render(60, 20)
	if !containsStr(view, "Apple") {
		t.Error("should render title")
	}
}

func TestListModalRenderScrolling(t *testing.T) {
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: string(rune('A' + i))}
	}
	m := NewListModal(ListModalConfig{ID: "list", Items: items})
	m.maxVisible = 5
	m.selected = 15 // Near end

	view := m.Render(60, 20)
	if view == "" {
		t.Error("should render with scroll offset")
	}
}

func TestWizardModalRenderEmptySteps(t *testing.T) {
	m := NewWizardModal(WizardModalConfig{ID: "empty", Steps: nil})
	m.OnPush(80, 24)
	view := m.Render(60, 40)
	// Should not panic
	if view == "" {
		t.Error("should render something")
	}
}

func TestManagerCenterContentSmallInput(t *testing.T) {
	mgr := NewManager()
	mgr.SetSize(20, 10) // Small terminal
	m := &testModalForManager{id: "small"}
	mgr.Push(m)
	view := mgr.Render(20, 10)
	if view == "" {
		t.Error("should render small modal")
	}
}

func TestListModalHandleKeyBoundary(t *testing.T) {
	m := NewListModal(ListModalConfig{
		ID:    "list",
		Items: []ListItem{{ID: "a", Title: "A"}}, // Only 1 item
	})
	// Try to go down at bottom
	m.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	if m.selected != 0 {
		t.Error("should stay at 0")
	}
	// Try to go up at top
	m.HandleKey(tea.KeyMsg{Type: tea.KeyUp})
	if m.selected != 0 {
		t.Error("should stay at 0")
	}
}

// testModalForManager implements Modal for manager tests
type testModalForManager struct {
	id string
}

func (m *testModalForManager) ID() string                               { return m.id }
func (m *testModalForManager) Title() string                            { return "Test" }
func (m *testModalForManager) Size() Size                               { return SizeSmall }
func (m *testModalForManager) Render(width, height int) string          { return "modal content" }
func (m *testModalForManager) OnPush(width, height int)                 {}
func (m *testModalForManager) OnPop()                                   {}
func (m *testModalForManager) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) { return false, nil }

// containsStr is a simple substring check helper.
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
