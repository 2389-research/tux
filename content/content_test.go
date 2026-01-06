package content

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestViewport(t *testing.T) {
	v := NewViewport()
	v.SetSize(80, 24)

	v.SetContent("Hello\nWorld\nTest")

	if v.Content() != "Hello\nWorld\nTest" {
		t.Errorf("content mismatch")
	}
	if v.LineCount() != 3 {
		t.Errorf("expected 3 lines, got %d", v.LineCount())
	}

	v.AppendContent("\nMore")
	if v.LineCount() != 4 {
		t.Errorf("expected 4 lines after append, got %d", v.LineCount())
	}

	if v.Value() != nil {
		t.Error("viewport should return nil value")
	}
}

func TestMultiSelect(t *testing.T) {
	items := []MultiSelectItem{
		{Label: "Option A", Key: "a"},
		{Label: "Option B", Key: "b"},
		{Label: "Option C", Key: "c"},
	}

	m := NewMultiSelect(items)

	// Initially nothing selected
	selected := m.Value().([]string)
	if len(selected) != 0 {
		t.Errorf("expected 0 selected, got %d", len(selected))
	}

	// Toggle first item
	m.Toggle()
	if m.SelectedCount() != 1 {
		t.Errorf("expected 1 selected, got %d", m.SelectedCount())
	}

	// Move down and toggle
	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m.Toggle()
	if m.SelectedCount() != 2 {
		t.Errorf("expected 2 selected, got %d", m.SelectedCount())
	}

	// Select all
	m.SelectAll()
	if m.SelectedCount() != 3 {
		t.Errorf("expected 3 selected, got %d", m.SelectedCount())
	}

	// Select none
	m.SelectNone()
	if m.SelectedCount() != 0 {
		t.Errorf("expected 0 selected, got %d", m.SelectedCount())
	}

	// Value should return keys
	m.items[0].Selected = true
	m.items[2].Selected = true
	selected = m.Value().([]string)
	if len(selected) != 2 {
		t.Errorf("expected 2 in value, got %d", len(selected))
	}
	if selected[0] != "a" || selected[1] != "c" {
		t.Errorf("expected [a, c], got %v", selected)
	}
}

func TestProgress(t *testing.T) {
	p := NewProgress(ProgressConfig{
		Total:   100,
		ShowBar: true,
	})

	p.SetCurrent(50)
	if p.Percent() != 0.5 {
		t.Errorf("expected 50%%, got %.0f%%", p.Percent()*100)
	}

	p.SetMessage("Processing...")
	view := p.View()
	if !containsStr(view, "Processing") {
		t.Error("view should contain message")
	}
	if !containsStr(view, "50/100") {
		t.Error("view should contain progress count")
	}

	// Items
	p = NewProgress(ProgressConfig{ShowItems: true})
	p.AddItem(ProgressItem{Label: "Task 1", Status: ProgressComplete})
	p.AddItem(ProgressItem{Label: "Task 2", Status: ProgressRunning})
	p.AddItem(ProgressItem{Label: "Task 3", Status: ProgressPending})

	view = p.View()
	if !containsStr(view, "Task 1") {
		t.Error("view should contain Task 1")
	}
	if !containsStr(view, "✓") {
		t.Error("view should contain checkmark for complete")
	}
}

func TestTimeline(t *testing.T) {
	tl := NewTimeline()
	tl.SetSize(80, 24)

	if tl.Count() != 0 {
		t.Errorf("expected 0 items, got %d", tl.Count())
	}

	tl.AddItem(TimelineItem{
		ID:     "1",
		Title:  "First item",
		Status: TimelineSuccess,
	})
	tl.AddItem(TimelineItem{
		ID:     "2",
		Title:  "Second item",
		Status: TimelineRunning,
	})

	if tl.Count() != 2 {
		t.Errorf("expected 2 items, got %d", tl.Count())
	}

	// Update item
	tl.UpdateItem("2", TimelineItem{
		Status:  TimelineSuccess,
		Content: "Done!",
	})

	item := tl.GetItem("2")
	if item.Status != TimelineSuccess {
		t.Error("item should be updated to success")
	}

	// Clear
	tl.Clear()
	if tl.Count() != 0 {
		t.Errorf("expected 0 items after clear, got %d", tl.Count())
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
