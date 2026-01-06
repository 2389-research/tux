package shell

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNewAutocomplete(t *testing.T) {
	ac := NewAutocomplete()
	if ac == nil {
		t.Fatal("expected non-nil autocomplete")
	}
	if ac.Active() {
		t.Error("autocomplete should not be active initially")
	}
	if ac.maxCompletions != 10 {
		t.Errorf("expected default maxCompletions 10, got %d", ac.maxCompletions)
	}
}

func TestAutocompleteRegisterProvider(t *testing.T) {
	ac := NewAutocomplete()

	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help", Description: "Show help"},
	})

	ac.RegisterProvider("command", provider)

	// Verify provider works
	ac.Show("/h", "command")
	if !ac.Active() {
		t.Error("autocomplete should be active after Show with matching provider")
	}
	if len(ac.Completions()) == 0 {
		t.Error("expected completions from registered provider")
	}
}

func TestAutocompleteUnregisterProvider(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help"},
	})

	ac.RegisterProvider("command", provider)
	ac.UnregisterProvider("command")

	ac.Show("/h", "command")
	if ac.Active() {
		t.Error("autocomplete should not be active after provider unregistered")
	}
}

func TestAutocompleteSetMaxCompletions(t *testing.T) {
	ac := NewAutocomplete()

	ac.SetMaxCompletions(5)
	if ac.maxCompletions != 5 {
		t.Errorf("expected maxCompletions 5, got %d", ac.maxCompletions)
	}

	// Zero or negative should be ignored
	ac.SetMaxCompletions(0)
	if ac.maxCompletions != 5 {
		t.Error("SetMaxCompletions(0) should be ignored")
	}

	ac.SetMaxCompletions(-1)
	if ac.maxCompletions != 5 {
		t.Error("SetMaxCompletions(-1) should be ignored")
	}
}

func TestAutocompleteSetStyles(t *testing.T) {
	ac := NewAutocomplete()
	dropdown := lipgloss.NewStyle().Background(lipgloss.Color("#000"))
	item := lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
	selected := lipgloss.NewStyle().Bold(true)
	desc := lipgloss.NewStyle().Italic(true)

	ac.SetStyles(dropdown, item, selected, desc)
	// Just verify no panic - styles are applied internally
}

func TestAutocompleteShowHide(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help"},
		{Value: "/quit", Display: "/quit"},
	})
	ac.RegisterProvider("command", provider)

	// Show with no matching provider
	ac.Show("/h", "nonexistent")
	if ac.Active() {
		t.Error("Show with nonexistent provider should not activate")
	}

	// Show with matching provider
	ac.Show("/", "command")
	if !ac.Active() {
		t.Error("Show with valid provider should activate")
	}
	if len(ac.Completions()) != 2 {
		t.Errorf("expected 2 completions, got %d", len(ac.Completions()))
	}

	// Hide
	ac.Hide()
	if ac.Active() {
		t.Error("autocomplete should not be active after Hide")
	}
	if len(ac.Completions()) != 0 {
		t.Error("completions should be cleared after Hide")
	}
}

func TestAutocompleteShowWithLimit(t *testing.T) {
	ac := NewAutocomplete()
	commands := make([]Completion, 20)
	for i := 0; i < 20; i++ {
		commands[i] = Completion{Value: "/cmd" + string(rune('a'+i)), Display: "/cmd" + string(rune('a'+i))}
	}
	provider := NewCommandProvider(commands)
	ac.RegisterProvider("command", provider)
	ac.SetMaxCompletions(5)

	ac.Show("/cmd", "command")
	if len(ac.Completions()) > 5 {
		t.Errorf("expected max 5 completions, got %d", len(ac.Completions()))
	}
}

func TestAutocompleteShowAuto(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help"},
	})
	ac.RegisterProvider("command", provider)

	ac.ShowAuto("/h")
	if !ac.Active() {
		t.Error("ShowAuto should activate for command input")
	}
}

func TestAutocompleteNavigation(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/a", Display: "/a"},
		{Value: "/b", Display: "/b"},
		{Value: "/c", Display: "/c"},
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/", "command")

	// Initial selection
	if ac.SelectedIndex() != 0 {
		t.Errorf("expected initial index 0, got %d", ac.SelectedIndex())
	}

	// Next
	ac.Next()
	if ac.SelectedIndex() != 1 {
		t.Errorf("expected index 1 after Next, got %d", ac.SelectedIndex())
	}

	// Next again
	ac.Next()
	if ac.SelectedIndex() != 2 {
		t.Errorf("expected index 2, got %d", ac.SelectedIndex())
	}

	// Next wraps around
	ac.Next()
	if ac.SelectedIndex() != 0 {
		t.Errorf("expected index 0 after wrap, got %d", ac.SelectedIndex())
	}

	// Previous wraps around
	ac.Previous()
	if ac.SelectedIndex() != 2 {
		t.Errorf("expected index 2 after Previous from 0, got %d", ac.SelectedIndex())
	}

	// Previous
	ac.Previous()
	if ac.SelectedIndex() != 1 {
		t.Errorf("expected index 1 after Previous, got %d", ac.SelectedIndex())
	}
}

func TestAutocompleteNavigationWhenInactive(t *testing.T) {
	ac := NewAutocomplete()

	// Should not panic when inactive
	ac.Next()
	ac.Previous()

	if ac.SelectedIndex() != 0 {
		t.Error("index should remain 0 when inactive")
	}
}

func TestAutocompleteNavigationEmptyCompletions(t *testing.T) {
	ac := NewAutocomplete()
	ac.active = true // Force active with no completions

	// Should not panic
	ac.Next()
	ac.Previous()
}

func TestAutocompleteGetSelected(t *testing.T) {
	ac := NewAutocomplete()

	// When inactive
	if ac.GetSelected() != nil {
		t.Error("GetSelected should return nil when inactive")
	}

	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help"},
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/h", "command")

	selected := ac.GetSelected()
	if selected == nil {
		t.Fatal("expected non-nil selected")
	}
	if selected.Value != "/help" {
		t.Errorf("expected selected value '/help', got %q", selected.Value)
	}
}

func TestAutocompleteDetectProvider(t *testing.T) {
	ac := NewAutocomplete()
	ac.RegisterProvider("command", NewCommandProvider(nil))
	ac.RegisterProvider("file", &mockFileProvider{})
	ac.RegisterProvider("history", NewHistoryProvider(nil))

	tests := []struct {
		input    string
		expected string
	}{
		{"/help", "command"},
		{"./file", "file"},
		{"~/path", "file"},
		{"/path/to", "command"}, // "/" prefix matches command first
		{"hello", "history"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ac.DetectProvider(tt.input)
			if result != tt.expected {
				t.Errorf("DetectProvider(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAutocompleteDetectProviderNoProviders(t *testing.T) {
	ac := NewAutocomplete()
	result := ac.DetectProvider("anything")
	if result != "" {
		t.Errorf("expected empty string with no providers, got %q", result)
	}
}

func TestAutocompleteDetectProviderFallback(t *testing.T) {
	ac := NewAutocomplete()
	ac.RegisterProvider("custom", NewCommandProvider(nil))

	// Should return available provider as fallback
	result := ac.DetectProvider("anything")
	if result != "custom" {
		t.Errorf("expected fallback to 'custom', got %q", result)
	}
}

func TestAutocompleteUpdate(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/a", Display: "/a"},
		{Value: "/b", Display: "/b"},
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/", "command")

	// Tab moves next
	ac.Update(tea.KeyMsg{Type: tea.KeyTab})
	if ac.SelectedIndex() != 1 {
		t.Error("Tab should move to next")
	}

	// Down also moves next
	ac.SelectedIndex() // reset check
	ac.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Shift+Tab moves previous
	ac.Update(tea.KeyMsg{Type: tea.KeyShiftTab})

	// Up also moves previous
	ac.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Esc hides
	ac.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if ac.Active() {
		t.Error("Esc should hide autocomplete")
	}
}

func TestAutocompleteUpdateWhenInactive(t *testing.T) {
	ac := NewAutocomplete()

	// Should not panic
	updated, cmd := ac.Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated != ac {
		t.Error("Update should return same autocomplete")
	}
	if cmd != nil {
		t.Error("Update when inactive should return nil cmd")
	}
}

func TestAutocompleteView(t *testing.T) {
	ac := NewAutocomplete()

	// When inactive
	if ac.View() != "" {
		t.Error("View should be empty when inactive")
	}

	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help", Description: "Show help"},
		{Value: "/quit", Display: "/quit", Description: "Exit app"},
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/", "command")

	view := ac.View()
	if view == "" {
		t.Error("View should be non-empty when active")
	}
	if !strings.Contains(view, "/help") {
		t.Error("View should contain /help")
	}
	if !strings.Contains(view, "Show help") {
		t.Error("View should contain description")
	}
}

func TestAutocompleteViewWithoutDescription(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/help", Display: "/help"}, // No description
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/", "command")

	view := ac.View()
	if !strings.Contains(view, "/help") {
		t.Error("View should contain /help")
	}
}

// CommandProvider tests

func TestCommandProvider(t *testing.T) {
	commands := []Completion{
		{Value: "/help", Display: "/help", Description: "Show help", Score: 10},
		{Value: "/history", Display: "/history", Description: "Show history", Score: 5},
		{Value: "/quit", Display: "/quit", Description: "Exit", Score: 1},
	}
	provider := NewCommandProvider(commands)

	// Prefix match
	completions := provider.GetCompletions("/h")
	if len(completions) != 2 {
		t.Errorf("expected 2 completions for '/h', got %d", len(completions))
	}

	// Exact match gets higher score
	completions = provider.GetCompletions("/help")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion for '/help', got %d", len(completions))
	}
	if completions[0].Score <= 10 {
		t.Error("exact match should have boosted score")
	}

	// No match
	completions = provider.GetCompletions("/xyz")
	if len(completions) != 0 {
		t.Errorf("expected 0 completions for '/xyz', got %d", len(completions))
	}

	// Without / prefix
	completions = provider.GetCompletions("help")
	if len(completions) != 0 {
		t.Error("should return nothing without / prefix")
	}
}

// HistoryProvider tests

func TestHistoryProvider(t *testing.T) {
	history := []string{"hello world", "help me", "goodbye"}
	provider := NewHistoryProvider(history)

	// Match in content
	completions := provider.GetCompletions("hello")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion for 'hello', got %d", len(completions))
	}

	// Partial match
	completions = provider.GetCompletions("hel")
	if len(completions) != 2 {
		t.Errorf("expected 2 completions for 'hel', got %d", len(completions))
	}

	// No match
	completions = provider.GetCompletions("xyz")
	if len(completions) != 0 {
		t.Errorf("expected 0 completions for 'xyz', got %d", len(completions))
	}

	// Empty input
	completions = provider.GetCompletions("")
	if len(completions) != 0 {
		t.Error("empty input should return no completions")
	}
}

func TestHistoryProviderAddHistory(t *testing.T) {
	provider := NewHistoryProvider(nil)

	provider.AddHistory("first")
	provider.AddHistory("second")
	provider.AddHistory("second") // Duplicate at end should be ignored

	completions := provider.GetCompletions("s")
	// Should only have one "second"
	count := 0
	for _, c := range completions {
		if c.Value == "second" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 'second' entry, got %d", count)
	}
}

func TestHistoryProviderScoring(t *testing.T) {
	provider := NewHistoryProvider([]string{"old match", "new match"})

	completions := provider.GetCompletions("match")
	if len(completions) != 2 {
		t.Fatalf("expected 2 completions, got %d", len(completions))
	}

	// More recent should have higher score (searched from end)
	// The first result should be "new match" since we search from end
	if completions[0].Value != "new match" {
		t.Error("more recent entry should come first")
	}
}

// Mock file provider for testing
type mockFileProvider struct{}

func (p *mockFileProvider) GetCompletions(input string) []Completion {
	return []Completion{{Value: input + "/file.txt", Display: "file.txt"}}
}

func TestCompletionSorting(t *testing.T) {
	ac := NewAutocomplete()
	provider := NewCommandProvider([]Completion{
		{Value: "/aaa", Display: "/aaa", Score: 1},
		{Value: "/bbb", Display: "/bbb", Score: 100},
		{Value: "/ccc", Display: "/ccc", Score: 50},
	})
	ac.RegisterProvider("command", provider)
	ac.Show("/", "command")

	completions := ac.Completions()
	if len(completions) != 3 {
		t.Fatalf("expected 3 completions, got %d", len(completions))
	}

	// Should be sorted by score descending
	// Note: CommandProvider adds bonus for prefix match, so scores will be modified
	// Just verify they're sorted
	for i := 0; i < len(completions)-1; i++ {
		if completions[i].Score < completions[i+1].Score {
			t.Errorf("completions not sorted by score: %d < %d", completions[i].Score, completions[i+1].Score)
		}
	}
}
