package content

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSelectList(t *testing.T) {
	items := []SelectItem{
		{Label: "Option A", Value: "a"},
		{Label: "Option B", Value: "b"},
		{Label: "Option C", Value: "c"},
	}

	list := NewSelectList(items)

	// Initial selection is first item
	if list.Value() != "a" {
		t.Errorf("expected initial value 'a', got %v", list.Value())
	}
	if list.Selected() != 0 {
		t.Errorf("expected selected index 0, got %d", list.Selected())
	}

	// Move down with arrow key
	list.Update(tea.KeyMsg{Type: tea.KeyDown})
	if list.Value() != "b" {
		t.Errorf("expected value 'b' after down, got %v", list.Value())
	}

	// Move down with j key
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if list.Value() != "c" {
		t.Errorf("expected value 'c' after j, got %v", list.Value())
	}

	// Can't go past end
	list.Update(tea.KeyMsg{Type: tea.KeyDown})
	if list.Value() != "c" {
		t.Errorf("expected value to stay 'c', got %v", list.Value())
	}

	// Move up with arrow key
	list.Update(tea.KeyMsg{Type: tea.KeyUp})
	if list.Value() != "b" {
		t.Errorf("expected value 'b' after up, got %v", list.Value())
	}

	// Move up with k key
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if list.Value() != "a" {
		t.Errorf("expected value 'a' after k, got %v", list.Value())
	}

	// Can't go before start
	list.Update(tea.KeyMsg{Type: tea.KeyUp})
	if list.Value() != "a" {
		t.Errorf("expected value to stay 'a', got %v", list.Value())
	}
}

func TestSelectListJumpKeys(t *testing.T) {
	items := []SelectItem{
		{Label: "Option A", Value: "a"},
		{Label: "Option B", Value: "b"},
		{Label: "Option C", Value: "c"},
	}

	list := NewSelectList(items)

	// Jump to end with G
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if list.Selected() != 2 {
		t.Errorf("expected selected index 2 after G, got %d", list.Selected())
	}

	// Jump to start with g
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if list.Selected() != 0 {
		t.Errorf("expected selected index 0 after g, got %d", list.Selected())
	}
}

func TestSelectListView(t *testing.T) {
	items := []SelectItem{
		{Label: "Option A", Description: "First option"},
		{Label: "Option B", Value: "b"},
	}

	list := NewSelectList(items)
	view := list.View()

	// Should contain labels
	if !contains(view, "Option A") {
		t.Error("view should contain 'Option A'")
	}
	if !contains(view, "Option B") {
		t.Error("view should contain 'Option B'")
	}
	// Should contain description
	if !contains(view, "First option") {
		t.Error("view should contain description 'First option'")
	}
	// Should contain cursor
	if !contains(view, "▸") {
		t.Error("view should contain cursor")
	}
}

func TestSelectListEmpty(t *testing.T) {
	list := NewSelectList(nil)

	if list.Value() != nil {
		t.Errorf("expected nil value for empty list, got %v", list.Value())
	}

	view := list.View()
	if !contains(view, "No items") {
		t.Error("empty list should show 'No items'")
	}
}

func TestSelectListSetSelected(t *testing.T) {
	items := []SelectItem{
		{Label: "A", Value: "a"},
		{Label: "B", Value: "b"},
		{Label: "C", Value: "c"},
	}

	list := NewSelectList(items)

	list.SetSelected(2)
	if list.Selected() != 2 {
		t.Errorf("expected selected 2, got %d", list.Selected())
	}

	// Out of bounds should be ignored
	list.SetSelected(10)
	if list.Selected() != 2 {
		t.Errorf("expected selected to stay 2, got %d", list.Selected())
	}

	list.SetSelected(-1)
	if list.Selected() != 2 {
		t.Errorf("expected selected to stay 2, got %d", list.Selected())
	}
}

func TestSelectListSetItems(t *testing.T) {
	list := NewSelectList([]SelectItem{
		{Label: "A", Value: "a"},
		{Label: "B", Value: "b"},
	})

	list.SetSelected(1)

	// Replace with fewer items
	list.SetItems([]SelectItem{
		{Label: "X", Value: "x"},
	})

	if list.Selected() != 0 {
		t.Errorf("expected selected to be clamped to 0, got %d", list.Selected())
	}
	if list.Value() != "x" {
		t.Errorf("expected value 'x', got %v", list.Value())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSelectListSetSize(t *testing.T) {
	list := NewSelectList([]SelectItem{
		{Label: "A", Value: "a"},
		{Label: "B", Value: "b"},
	})

	// SetSize should not panic
	list.SetSize(100, 50)
	// Verify view still works after resize
	view := list.View()
	if view == "" {
		t.Error("View should work after SetSize")
	}
}

func TestSelectListSelectedItem(t *testing.T) {
	items := []SelectItem{
		{Label: "A", Value: "a", Description: "First"},
		{Label: "B", Value: "b", Description: "Second"},
	}
	list := NewSelectList(items)

	item := list.SelectedItem()
	if item == nil {
		t.Error("SelectedItem should not return nil")
	}
	if item.Label != "A" {
		t.Errorf("expected label 'A', got %s", item.Label)
	}

	list.SetSelected(1)
	item = list.SelectedItem()
	if item.Label != "B" {
		t.Errorf("expected label 'B', got %s", item.Label)
	}
}

func TestSelectListSelectedItemEmpty(t *testing.T) {
	list := NewSelectList(nil)
	item := list.SelectedItem()
	if item != nil {
		t.Error("SelectedItem should return nil for empty list")
	}
}

func TestSelectListItems(t *testing.T) {
	items := []SelectItem{
		{Label: "A", Value: "a"},
		{Label: "B", Value: "b"},
	}
	list := NewSelectList(items)

	retrieved := list.Items()
	if len(retrieved) != 2 {
		t.Errorf("expected 2 items, got %d", len(retrieved))
	}
	if retrieved[0].Label != "A" {
		t.Errorf("expected first item label 'A', got %s", retrieved[0].Label)
	}
}

func TestSelectListInit(t *testing.T) {
	list := NewSelectList(nil)
	cmd := list.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}
