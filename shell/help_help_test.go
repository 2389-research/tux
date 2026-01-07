package shell

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNew(t *testing.T) {
	cat := Category{
		Title: "Navigation",
		Bindings: []Binding{
			{Key: "j", Description: "Move down"},
			{Key: "k", Description: "Move up"},
		},
	}

	h := NewHelp(cat)

	if h == nil {
		t.Fatal("New() returned nil")
	}
	if len(h.categories) != 1 {
		t.Errorf("expected 1 category, got %d", len(h.categories))
	}
	if h.categories[0].Title != "Navigation" {
		t.Errorf("expected category title 'Navigation', got '%s'", h.categories[0].Title)
	}
}

func TestNewMultipleCategories(t *testing.T) {
	cat1 := Category{Title: "Navigation", Bindings: []Binding{{Key: "j", Description: "Down"}}}
	cat2 := Category{Title: "Actions", Bindings: []Binding{{Key: "enter", Description: "Select"}}}

	h := NewHelp(cat1, cat2)

	if len(h.categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(h.categories))
	}
}

func TestWithTheme(t *testing.T) {
	h := NewHelp()
	th := theme.Get("dracula")

	result := h.WithTheme(th)

	// Should return the same pointer for chaining
	if result != h {
		t.Error("WithTheme should return the same Help pointer for chaining")
	}
	// Theme should be set
	if h.theme == nil {
		t.Error("WithTheme should set the theme")
	}
}

func TestRenderEmptyMode(t *testing.T) {
	cat := Category{
		Title: "Navigation",
		Bindings: []Binding{
			{Key: "j", Description: "Move down"},
			{Key: "k", Description: "Move up"},
		},
	}
	h := NewHelp(cat).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "")

	// Should contain category title
	if !strings.Contains(result, "Navigation") {
		t.Error("Render should include category title")
	}
	// Should contain all bindings
	if !strings.Contains(result, "j") {
		t.Error("Render should include key 'j'")
	}
	if !strings.Contains(result, "Move down") {
		t.Error("Render should include description 'Move down'")
	}
	if !strings.Contains(result, "k") {
		t.Error("Render should include key 'k'")
	}
	// Should contain footer
	if !strings.Contains(result, "Press ? or esc to close") {
		t.Error("Render should include footer text")
	}
}

func TestRenderWithModeFilter(t *testing.T) {
	cat := Category{
		Title: "Actions",
		Bindings: []Binding{
			{Key: "enter", Description: "Select", Modes: []string{"normal"}},
			{Key: "i", Description: "Insert", Modes: []string{"normal"}},
			{Key: "esc", Description: "Exit insert", Modes: []string{"insert"}},
		},
	}
	h := NewHelp(cat).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "normal")

	// Should include normal mode bindings
	if !strings.Contains(result, "enter") {
		t.Error("Render with mode 'normal' should include 'enter'")
	}
	if !strings.Contains(result, "Select") {
		t.Error("Render with mode 'normal' should include 'Select'")
	}
	// Should NOT include insert mode bindings
	if strings.Contains(result, "Exit insert") {
		t.Error("Render with mode 'normal' should NOT include 'Exit insert'")
	}
}

func TestRenderHidesEmptyCategories(t *testing.T) {
	cat1 := Category{
		Title: "Normal Mode",
		Bindings: []Binding{
			{Key: "j", Description: "Down", Modes: []string{"normal"}},
		},
	}
	cat2 := Category{
		Title: "Insert Mode",
		Bindings: []Binding{
			{Key: "esc", Description: "Exit", Modes: []string{"insert"}},
		},
	}
	h := NewHelp(cat1, cat2).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "normal")

	// Should include Normal Mode category
	if !strings.Contains(result, "Normal Mode") {
		t.Error("Render should include 'Normal Mode' category when filtering by 'normal'")
	}
	// Should NOT include Insert Mode category (no visible bindings)
	if strings.Contains(result, "Insert Mode") {
		t.Error("Render should hide 'Insert Mode' category when no bindings match")
	}
}

func TestRenderHasBorder(t *testing.T) {
	cat := Category{
		Title: "Test",
		Bindings: []Binding{
			{Key: "x", Description: "Test action"},
		},
	}
	h := NewHelp(cat).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "")

	// Border should be present - check for rounded border characters
	// Lipgloss rounded border uses these characters
	hasBorder := strings.Contains(result, "\u256d") || // top-left corner
		strings.Contains(result, "\u256e") || // top-right corner
		strings.Contains(result, "\u256f") || // bottom-right corner
		strings.Contains(result, "\u2570") || // bottom-left corner
		strings.Contains(result, "\u2502") || // vertical line
		strings.Contains(result, "\u2500") // horizontal line

	if !hasBorder {
		t.Error("Render should wrap content in a border box")
	}
}

func TestRenderMultipleCategories(t *testing.T) {
	cat1 := Category{
		Title: "Navigation",
		Bindings: []Binding{
			{Key: "j", Description: "Down"},
		},
	}
	cat2 := Category{
		Title: "Actions",
		Bindings: []Binding{
			{Key: "enter", Description: "Select"},
		},
	}
	h := NewHelp(cat1, cat2).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "")

	if !strings.Contains(result, "Navigation") {
		t.Error("Render should include 'Navigation' category")
	}
	if !strings.Contains(result, "Actions") {
		t.Error("Render should include 'Actions' category")
	}
	if !strings.Contains(result, "j") {
		t.Error("Render should include key 'j'")
	}
	if !strings.Contains(result, "enter") {
		t.Error("Render should include key 'enter'")
	}
}

func TestRenderWithoutTheme(t *testing.T) {
	cat := Category{
		Title: "Test",
		Bindings: []Binding{
			{Key: "x", Description: "Test action"},
		},
	}
	h := NewHelp(cat)

	// Should not panic and should use default theme
	result := h.Render(80, "")

	if !strings.Contains(result, "Test") {
		t.Error("Render without explicit theme should still work")
	}
}

func TestRenderNoCategories(t *testing.T) {
	h := NewHelp().WithTheme(theme.Get("dracula"))

	result := h.Render(80, "")

	// Should still have footer
	if !strings.Contains(result, "Press ? or esc to close") {
		t.Error("Render with no categories should still show footer")
	}
}

func TestRenderAllBindingsFilteredOut(t *testing.T) {
	cat := Category{
		Title: "Insert Only",
		Bindings: []Binding{
			{Key: "esc", Description: "Exit", Modes: []string{"insert"}},
		},
	}
	h := NewHelp(cat).WithTheme(theme.Get("dracula"))

	result := h.Render(80, "normal")

	// Category should be hidden since no bindings match
	if strings.Contains(result, "Insert Only") {
		t.Error("Render should hide category when all bindings are filtered out")
	}
	// Footer should still be present
	if !strings.Contains(result, "Press ? or esc to close") {
		t.Error("Render should still show footer even when all categories hidden")
	}
}

func TestRenderGlobalBindings(t *testing.T) {
	cat := Category{
		Title: "Global",
		Bindings: []Binding{
			{Key: "?", Description: "Help"},           // No modes = global
			{Key: "i", Description: "Insert", Modes: []string{"normal"}},
		},
	}
	h := NewHelp(cat).WithTheme(theme.Get("dracula"))

	// Global bindings should appear in any mode
	normalResult := h.Render(80, "normal")
	insertResult := h.Render(80, "insert")

	if !strings.Contains(normalResult, "?") {
		t.Error("Global binding should appear in normal mode")
	}
	if !strings.Contains(insertResult, "?") {
		t.Error("Global binding should appear in insert mode")
	}
	// Mode-specific should only appear in its mode
	if !strings.Contains(normalResult, "Insert") {
		t.Error("Normal-mode binding should appear in normal mode")
	}
	if strings.Contains(insertResult, "Insert") {
		t.Error("Normal-mode binding should NOT appear in insert mode")
	}
}
