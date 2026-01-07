// shell/help_binding_test.go
package shell

import "testing"

func TestBindingMatchesMode(t *testing.T) {
	// Binding with no modes matches everything
	b := Binding{Key: "ctrl+c", Description: "Quit"}
	if !b.MatchesMode("chat") {
		t.Error("empty modes should match any mode")
	}
	if !b.MatchesMode("") {
		t.Error("empty modes should match empty mode")
	}

	// Binding with modes matches only those modes
	b2 := Binding{Key: "enter", Description: "Send", Modes: []string{"chat", "compose"}}
	if !b2.MatchesMode("chat") {
		t.Error("should match chat mode")
	}
	if b2.MatchesMode("history") {
		t.Error("should not match history mode")
	}
	if !b2.MatchesMode("") {
		t.Error("empty mode should show all bindings")
	}
}

func TestCategoryFilterByMode(t *testing.T) {
	cat := Category{
		Title: "Actions",
		Bindings: []Binding{
			{Key: "ctrl+c", Description: "Quit"},
			{Key: "enter", Description: "Send", Modes: []string{"chat"}},
			{Key: "enter", Description: "Select", Modes: []string{"list"}},
		},
	}

	// Filter for chat mode
	filtered := cat.FilterByMode("chat")
	if len(filtered) != 2 {
		t.Errorf("expected 2 bindings for chat, got %d", len(filtered))
	}

	// Filter for list mode
	filtered = cat.FilterByMode("list")
	if len(filtered) != 2 {
		t.Errorf("expected 2 bindings for list, got %d", len(filtered))
	}

	// Empty mode shows all
	filtered = cat.FilterByMode("")
	if len(filtered) != 3 {
		t.Errorf("expected 3 bindings for empty mode, got %d", len(filtered))
	}
}
