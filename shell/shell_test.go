package shell

import (
	"testing"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/modal"
	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
)

func TestShellNew(t *testing.T) {
	th := theme.NewDraculaTheme()
	cfg := DefaultConfig()

	s := New(th, cfg)

	if s.theme.Name() != "dracula" {
		t.Errorf("expected dracula theme, got %s", s.theme.Name())
	}
	if s.focused != FocusInput {
		t.Error("initial focus should be input")
	}
}

func TestShellTabs(t *testing.T) {
	s := New(nil, DefaultConfig())

	// Add tabs
	s.AddTab(Tab{ID: "tab1", Label: "Tab 1"})
	s.AddTab(Tab{ID: "tab2", Label: "Tab 2"})

	if s.tabs.Count() != 2 {
		t.Errorf("expected 2 tabs, got %d", s.tabs.Count())
	}

	// First tab should be active
	if s.tabs.ActiveTab().ID != "tab1" {
		t.Errorf("expected tab1 active, got %s", s.tabs.ActiveTab().ID)
	}

	// Switch tabs
	s.SetActiveTab("tab2")
	if s.tabs.ActiveTab().ID != "tab2" {
		t.Errorf("expected tab2 active, got %s", s.tabs.ActiveTab().ID)
	}

	// Remove tab
	s.RemoveTab("tab1")
	if s.tabs.Count() != 1 {
		t.Errorf("expected 1 tab after remove, got %d", s.tabs.Count())
	}
}

func TestShellModal(t *testing.T) {
	s := New(nil, DefaultConfig())

	if s.HasModal() {
		t.Error("should not have modal initially")
	}

	// Push modal
	m := &testModal{id: "test"}
	s.PushModal(m)

	if !s.HasModal() {
		t.Error("should have modal after push")
	}
	if s.focused != FocusModal {
		t.Error("focus should be modal after push")
	}

	// Pop modal
	popped := s.PopModal()
	if popped.ID() != "test" {
		t.Errorf("expected to pop test modal, got %s", popped.ID())
	}
	if s.HasModal() {
		t.Error("should not have modal after pop")
	}
	if s.focused != FocusInput {
		t.Error("focus should return to input after pop")
	}
}

func TestShellInput(t *testing.T) {
	s := New(nil, DefaultConfig())

	s.SetInputValue("hello")
	if s.InputValue() != "hello" {
		t.Errorf("expected 'hello', got %s", s.InputValue())
	}

	s.ClearInput()
	if s.InputValue() != "" {
		t.Error("input should be empty after clear")
	}
}

func TestShellUpdate(t *testing.T) {
	s := New(nil, DefaultConfig())

	// Window size
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if !s.ready {
		t.Error("shell should be ready after window size")
	}
	if s.width != 80 || s.height != 24 {
		t.Errorf("expected 80x24, got %dx%d", s.width, s.height)
	}
}

func TestShellView(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	s.AddTab(Tab{
		ID:      "chat",
		Label:   "Chat",
		Content: content.NewSelectList(nil),
	})

	view := s.View()
	if view == "" {
		t.Error("view should not be empty")
	}
	if len(view) < 10 {
		t.Error("view seems too short")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.ShowTabBar {
		t.Error("ShowTabBar should default to true")
	}
	if !cfg.ShowStatusBar {
		t.Error("ShowStatusBar should default to true")
	}
	if !cfg.ShowInput {
		t.Error("ShowInput should default to true")
	}
	if cfg.InputPrefix != "> " {
		t.Errorf("InputPrefix should default to '> ', got %q", cfg.InputPrefix)
	}
}

// testModal for testing
type testModal struct {
	id string
}

func (m *testModal) ID() string                                       { return m.id }
func (m *testModal) Title() string                                    { return "Test" }
func (m *testModal) Size() modal.Size                                 { return modal.SizeMedium }
func (m *testModal) Render(width, height int) string                  { return "test modal" }
func (m *testModal) OnPush(width, height int)                         {}
func (m *testModal) OnPop()                                           {}
func (m *testModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) { return false, nil }
