// help/integration_test.go
package help_test

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/help"
	"github.com/2389-research/tux/modal"
	"github.com/2389-research/tux/theme"
)

func TestHelpFullWorkflow(t *testing.T) {
	// Create help with multiple categories and modes
	h := help.New(
		help.Category{
			Title: "Navigation",
			Bindings: []help.Binding{
				{Key: "↑↓", Description: "Navigate", Modes: []string{"list", "history"}},
				{Key: "enter", Description: "Select", Modes: []string{"list"}},
				{Key: "enter", Description: "Send message", Modes: []string{"chat"}},
			},
		},
		help.Category{
			Title: "General",
			Bindings: []help.Binding{
				{Key: "ctrl+c", Description: "Quit"},
				{Key: "?", Description: "Toggle help"},
				{Key: "esc", Description: "Close/cancel"},
			},
		},
		help.Category{
			Title: "Commands",
			Bindings: []help.Binding{
				{Key: "/feedback", Description: "Send feedback"},
				{Key: "/clear", Description: "Clear screen"},
			},
		},
	).WithTheme(theme.Get("dracula"))

	// Test rendering in different modes
	t.Run("chat mode", func(t *testing.T) {
		output := h.Render(80, "chat")
		if !strings.Contains(output, "Send message") {
			t.Error("chat mode should show 'Send message'")
		}
		if strings.Contains(output, "Navigate") {
			t.Error("chat mode should not show 'Navigate'")
		}
		if !strings.Contains(output, "Quit") {
			t.Error("should always show 'Quit'")
		}
	})

	t.Run("list mode", func(t *testing.T) {
		output := h.Render(80, "list")
		if !strings.Contains(output, "Navigate") {
			t.Error("list mode should show 'Navigate'")
		}
		if !strings.Contains(output, "Select") {
			t.Error("list mode should show 'Select'")
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
	h := help.New(
		help.Category{
			Title: "Test",
			Bindings: []help.Binding{
				{Key: "?", Description: "Help"},
			},
		},
	)

	m := modal.NewHelpModal(modal.HelpModalConfig{
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
	// Test mode filtering with multiple categories
	h := help.New(
		help.Category{
			Title: "Chat Actions",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Send message", Modes: []string{"chat"}},
				{Key: "ctrl+u", Description: "Clear input", Modes: []string{"chat"}},
			},
		},
		help.Category{
			Title: "List Actions",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Select item", Modes: []string{"list"}},
				{Key: "j/k", Description: "Navigate", Modes: []string{"list"}},
			},
		},
		help.Category{
			Title: "Global",
			Bindings: []help.Binding{
				{Key: "ctrl+c", Description: "Quit"}, // No modes = global
				{Key: "?", Description: "Help"},      // No modes = global
			},
		},
	).WithTheme(theme.Get("dracula"))

	t.Run("chat mode shows chat + global", func(t *testing.T) {
		output := h.Render(80, "chat")

		// Chat bindings visible
		if !strings.Contains(output, "Send message") {
			t.Error("chat mode should show 'Send message'")
		}
		if !strings.Contains(output, "Clear input") {
			t.Error("chat mode should show 'Clear input'")
		}

		// List bindings hidden
		if strings.Contains(output, "Select item") {
			t.Error("chat mode should NOT show 'Select item'")
		}
		if strings.Contains(output, "Navigate") {
			t.Error("chat mode should NOT show 'Navigate'")
		}

		// Global bindings visible
		if !strings.Contains(output, "Quit") {
			t.Error("chat mode should show global 'Quit'")
		}
		if !strings.Contains(output, "Help") {
			t.Error("chat mode should show global 'Help'")
		}
	})

	t.Run("list mode shows list + global", func(t *testing.T) {
		output := h.Render(80, "list")

		// List bindings visible
		if !strings.Contains(output, "Select item") {
			t.Error("list mode should show 'Select item'")
		}
		if !strings.Contains(output, "Navigate") {
			t.Error("list mode should show 'Navigate'")
		}

		// Chat bindings hidden
		if strings.Contains(output, "Send message") {
			t.Error("list mode should NOT show 'Send message'")
		}

		// Global bindings visible
		if !strings.Contains(output, "Quit") {
			t.Error("list mode should show global 'Quit'")
		}
	})

	t.Run("empty mode shows all", func(t *testing.T) {
		output := h.Render(80, "")

		// All bindings visible
		if !strings.Contains(output, "Send message") {
			t.Error("empty mode should show 'Send message'")
		}
		if !strings.Contains(output, "Select item") {
			t.Error("empty mode should show 'Select item'")
		}
		if !strings.Contains(output, "Quit") {
			t.Error("empty mode should show 'Quit'")
		}
	})
}

func TestGlobalBindingsAppearInAllModes(t *testing.T) {
	// Global bindings (empty Modes) should appear regardless of mode filter
	h := help.New(
		help.Category{
			Title: "Global Shortcuts",
			Bindings: []help.Binding{
				{Key: "ctrl+c", Description: "Quit application"},
				{Key: "?", Description: "Toggle help overlay"},
				{Key: "ctrl+l", Description: "Clear screen"},
			},
		},
	).WithTheme(theme.Get("dracula"))

	modes := []string{"chat", "list", "history", "compose", "arbitrary-mode"}

	for _, mode := range modes {
		t.Run("mode: "+mode, func(t *testing.T) {
			output := h.Render(80, mode)

			if !strings.Contains(output, "Quit application") {
				t.Errorf("global binding 'Quit application' should appear in mode %q", mode)
			}
			if !strings.Contains(output, "Toggle help overlay") {
				t.Errorf("global binding 'Toggle help overlay' should appear in mode %q", mode)
			}
			if !strings.Contains(output, "Clear screen") {
				t.Errorf("global binding 'Clear screen' should appear in mode %q", mode)
			}
		})
	}
}

func TestHelpModalWithModeParameter(t *testing.T) {
	h := help.New(
		help.Category{
			Title: "Mode-Specific",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Chat send", Modes: []string{"chat"}},
				{Key: "enter", Description: "List select", Modes: []string{"list"}},
			},
		},
		help.Category{
			Title: "Always Visible",
			Bindings: []help.Binding{
				{Key: "?", Description: "Help"},
			},
		},
	)

	t.Run("modal with chat mode", func(t *testing.T) {
		m := modal.NewHelpModal(modal.HelpModalConfig{
			Help: h,
			Mode: "chat",
		})
		m.OnPush(80, 24)
		output := m.Render(60, 20)

		if !strings.Contains(output, "Chat send") {
			t.Error("modal with chat mode should show 'Chat send'")
		}
		if strings.Contains(output, "List select") {
			t.Error("modal with chat mode should NOT show 'List select'")
		}
		if !strings.Contains(output, "Help") {
			t.Error("modal should always show global 'Help'")
		}
	})

	t.Run("modal with list mode", func(t *testing.T) {
		m := modal.NewHelpModal(modal.HelpModalConfig{
			Help: h,
			Mode: "list",
		})
		m.OnPush(80, 24)
		output := m.Render(60, 20)

		if !strings.Contains(output, "List select") {
			t.Error("modal with list mode should show 'List select'")
		}
		if strings.Contains(output, "Chat send") {
			t.Error("modal with list mode should NOT show 'Chat send'")
		}
	})

	t.Run("modal without mode shows all", func(t *testing.T) {
		m := modal.NewHelpModal(modal.HelpModalConfig{
			Help: h,
			Mode: "", // No mode filter
		})
		m.OnPush(80, 24)
		output := m.Render(60, 20)

		if !strings.Contains(output, "Chat send") {
			t.Error("modal without mode should show 'Chat send'")
		}
		if !strings.Contains(output, "List select") {
			t.Error("modal without mode should show 'List select'")
		}
	})
}

func TestCategoryHidingWhenEmpty(t *testing.T) {
	// Categories with no visible bindings for the current mode should be hidden
	h := help.New(
		help.Category{
			Title: "Chat Only Category",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Send", Modes: []string{"chat"}},
			},
		},
		help.Category{
			Title: "List Only Category",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Select", Modes: []string{"list"}},
			},
		},
	).WithTheme(theme.Get("dracula"))

	t.Run("chat mode hides list category", func(t *testing.T) {
		output := h.Render(80, "chat")

		if !strings.Contains(output, "Chat Only Category") {
			t.Error("chat mode should show 'Chat Only Category'")
		}
		if strings.Contains(output, "List Only Category") {
			t.Error("chat mode should hide 'List Only Category'")
		}
	})

	t.Run("list mode hides chat category", func(t *testing.T) {
		output := h.Render(80, "list")

		if strings.Contains(output, "Chat Only Category") {
			t.Error("list mode should hide 'Chat Only Category'")
		}
		if !strings.Contains(output, "List Only Category") {
			t.Error("list mode should show 'List Only Category'")
		}
	})
}
