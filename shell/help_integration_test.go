// shell/help_integration_test.go
package shell_test

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/shell"
	"github.com/2389-research/tux/theme"
)

func TestHelpFullWorkflow(t *testing.T) {
	h := shell.NewHelp(
		shell.Category{
			Title: "Navigation",
			Bindings: []shell.Binding{
				{Key: "↑↓", Description: "Navigate", Modes: []string{"list", "history"}},
				{Key: "enter", Description: "Select", Modes: []string{"list"}},
				{Key: "enter", Description: "Send message", Modes: []string{"chat"}},
			},
		},
		shell.Category{
			Title: "General",
			Bindings: []shell.Binding{
				{Key: "ctrl+c", Description: "Quit"},
				{Key: "?", Description: "Toggle help"},
				{Key: "esc", Description: "Close/cancel"},
			},
		},
	).WithTheme(theme.Get("dracula"))

	t.Run("chat mode", func(t *testing.T) {
		output := h.Render(80, "chat")
		if !strings.Contains(output, "Send message") {
			t.Error("chat mode should show 'Send message'")
		}
		if strings.Contains(output, "Navigate") {
			t.Error("chat mode should not show 'Navigate'")
		}
	})

	t.Run("list mode", func(t *testing.T) {
		output := h.Render(80, "list")
		if !strings.Contains(output, "Navigate") {
			t.Error("list mode should show 'Navigate'")
		}
	})

	t.Run("no mode shows all", func(t *testing.T) {
		output := h.Render(80, "")
		if !strings.Contains(output, "Navigate") {
			t.Error("no mode should show 'Navigate'")
		}
		if !strings.Contains(output, "Send message") {
			t.Error("no mode should show 'Send message'")
		}
	})
}

func TestHelpModalIntegration(t *testing.T) {
	h := shell.NewHelp(
		shell.Category{
			Title: "Test",
			Bindings: []shell.Binding{
				{Key: "?", Description: "Help"},
			},
		},
	)

	m := shell.NewHelpModal(shell.HelpModalConfig{
		ID:    "test",
		Title: "Test Help",
		Help:  h,
		Mode:  "chat",
	})

	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if !strings.Contains(output, "Help") {
		t.Error("modal should render help content")
	}
}

func TestModeFilteringAcrossCategories(t *testing.T) {
	h := shell.NewHelp(
		shell.Category{
			Title: "Chat Actions",
			Bindings: []shell.Binding{
				{Key: "enter", Description: "Send message", Modes: []string{"chat"}},
			},
		},
		shell.Category{
			Title: "List Actions",
			Bindings: []shell.Binding{
				{Key: "enter", Description: "Select item", Modes: []string{"list"}},
			},
		},
		shell.Category{
			Title: "Global",
			Bindings: []shell.Binding{
				{Key: "ctrl+c", Description: "Quit"},
			},
		},
	).WithTheme(theme.Get("dracula"))

	t.Run("chat mode shows chat + global", func(t *testing.T) {
		output := h.Render(80, "chat")
		if !strings.Contains(output, "Send message") {
			t.Error("chat mode should show 'Send message'")
		}
		if strings.Contains(output, "Select item") {
			t.Error("chat mode should NOT show 'Select item'")
		}
		if !strings.Contains(output, "Quit") {
			t.Error("chat mode should show global 'Quit'")
		}
	})

	t.Run("list mode shows list + global", func(t *testing.T) {
		output := h.Render(80, "list")
		if !strings.Contains(output, "Select item") {
			t.Error("list mode should show 'Select item'")
		}
		if !strings.Contains(output, "Quit") {
			t.Error("list mode should show global 'Quit'")
		}
	})
}

func TestGlobalBindingsAppearInAllModes(t *testing.T) {
	h := shell.NewHelp(
		shell.Category{
			Title: "Global",
			Bindings: []shell.Binding{
				{Key: "ctrl+c", Description: "Quit"},
			},
		},
	).WithTheme(theme.Get("dracula"))

	modes := []string{"chat", "list", "history"}
	for _, mode := range modes {
		t.Run("mode: "+mode, func(t *testing.T) {
			output := h.Render(80, mode)
			if !strings.Contains(output, "Quit") {
				t.Errorf("global binding 'Quit' should appear in mode %q", mode)
			}
		})
	}
}
