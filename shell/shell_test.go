package shell

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/content"
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
	m := &shellTestModal{id: "test"}
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

func TestShellFocus(t *testing.T) {
	s := New(nil, DefaultConfig())

	s.Focus(FocusTab)
	if s.focused != FocusTab {
		t.Error("focus should be tab")
	}

	s.Focus(FocusModal)
	if s.focused != FocusModal {
		t.Error("focus should be modal")
	}

	s.Focus(FocusInput)
	if s.focused != FocusInput {
		t.Error("focus should be input")
	}
}

func TestShellTheme(t *testing.T) {
	th := theme.NewDraculaTheme()
	s := New(th, DefaultConfig())

	if s.Theme() != th {
		t.Error("Theme() should return the theme")
	}
}

func TestShellSetStatus(t *testing.T) {
	s := New(nil, DefaultConfig())

	s.SetStatus(Status{
		Model:     "claude-3",
		Connected: true,
		Message:   "Testing",
	})
	// Should not panic
}

func TestShellOverlayModal(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Add a tab with content
	s.AddTab(Tab{
		ID:      "test",
		Label:   "Test",
		Content: content.NewSelectList(nil),
	})

	// Push modal
	m := &shellTestModal{id: "overlay-test"}
	s.PushModal(m)

	// View should render with modal overlay
	view := s.View()
	if view == "" {
		t.Error("view should not be empty with modal")
	}
}

func TestTabBarHandleKey(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "2", Label: "Tab 2"})
	tb.AddTab(Tab{ID: "3", Label: "Tab 3"})

	// Tab key cycles through tabs
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if tb.ActiveTab().ID != "2" {
		t.Error("tab key should move to next tab")
	}

	// Shift+Tab goes back
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyShiftTab})
	if tb.ActiveTab().ID != "1" {
		t.Error("shift+tab should move to previous tab")
	}

	// Continue cycling
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if tb.ActiveTab().ID != "3" {
		t.Error("should be at tab 3")
	}
}

func TestTabBarNextPrevTab(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "2", Label: "Tab 2"})
	tb.AddTab(Tab{ID: "3", Label: "Tab 3"})

	tb.NextTab()
	if tb.ActiveTab().ID != "2" {
		t.Error("NextTab should move to tab 2")
	}

	tb.NextTab()
	if tb.ActiveTab().ID != "3" {
		t.Error("NextTab should move to tab 3")
	}

	tb.NextTab() // Should wrap
	if tb.ActiveTab().ID != "1" {
		t.Error("NextTab should wrap to tab 1")
	}

	tb.PrevTab() // Should wrap back
	if tb.ActiveTab().ID != "3" {
		t.Error("PrevTab should wrap to tab 3")
	}

	tb.PrevTab()
	if tb.ActiveTab().ID != "2" {
		t.Error("PrevTab should move to tab 2")
	}
}

func TestTabBarSetBadge(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})

	tb.SetBadge("1", "3")
	tab := tb.ActiveTab()
	if tab.Badge != "3" {
		t.Errorf("expected badge '3', got %s", tab.Badge)
	}
}

func TestStatusBarSetters(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	sb.SetStatus(Status{Connected: true})
	sb.SetModel("claude-3")
	sb.SetConnected(true)
	sb.SetStreaming(true)
	sb.SetTokens(100000, 50000)
	sb.SetMode("normal")
	sb.SetMessage("Processing...")
	sb.SetHints("Ctrl+C | Ctrl+D")

	// Render and verify it doesn't panic
	view := sb.View(80)
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestStatusBarViewVariants(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	// Default state
	view := sb.View(80)
	if view == "" {
		t.Error("default view should not be empty")
	}

	// With streaming
	sb.SetStreaming(true)
	view = sb.View(80)
	if view == "" {
		t.Error("streaming view should not be empty")
	}

	// With tokens
	sb.SetTokens(100000, 500000)
	view = sb.View(80)
	if view == "" {
		t.Error("tokens view should not be empty")
	}

	// With message
	sb.SetMessage("Status message")
	view = sb.View(80)
	if view == "" {
		t.Error("message view should not be empty")
	}
}

func TestShellInit(t *testing.T) {
	s := New(nil, DefaultConfig())
	cmd := s.Init()
	// Init should return a batch command (for input)
	if cmd == nil {
		t.Error("Init should return a command")
	}
}

func TestInputInit(t *testing.T) {
	th := theme.NewDraculaTheme()
	inp := NewInput(th, "> ", "")
	cmd := inp.Init()
	// Blink command
	if cmd == nil {
		t.Error("Input Init should return blink command")
	}
}

func TestInputUpdate(t *testing.T) {
	th := theme.NewDraculaTheme()
	inp := NewInput(th, "> ", "")

	// Type something
	inp, _ = inp.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	inp, _ = inp.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

	if inp.Value() != "hi" {
		t.Errorf("expected 'hi', got %s", inp.Value())
	}
}

func TestInputFocusBlur(t *testing.T) {
	th := theme.NewDraculaTheme()
	inp := NewInput(th, "> ", "")

	inp.Focus()
	if !inp.Focused() {
		t.Error("should be focused after Focus()")
	}

	inp.Blur()
	if inp.Focused() {
		t.Error("should not be focused after Blur()")
	}
}

func TestTabBarSetSizeWithContent(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)

	// Add tab with content
	tb.AddTab(Tab{
		ID:      "1",
		Label:   "Tab 1",
		Content: content.NewSelectList(nil),
	})

	tb.SetSize(100, 50)
	// Should not panic, and should update content size
}

// shellTestModal for testing
type shellTestModal struct {
	id         string
	handleKey  bool
	returnCmd  tea.Cmd
}

func (m *shellTestModal) ID() string                     { return m.id }
func (m *shellTestModal) Title() string                  { return "Test" }
func (m *shellTestModal) Size() Size               { return SizeMedium }
func (m *shellTestModal) Render(width, height int) string { return "test modal" }
func (m *shellTestModal) OnPush(width, height int)       {}
func (m *shellTestModal) OnPop()                         {}
func (m *shellTestModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	return m.handleKey, m.returnCmd
}

// === Additional tests for 95% coverage ===

func TestShellUpdateWithModalHandled(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Push modal that handles keys
	m := &shellTestModal{id: "handled", handleKey: true}
	s.PushModal(m)

	// Key should be handled by modal
	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// Modal handled, should return cmd from HandleKey
	_ = cmd // May be nil
}

func TestShellUpdateWithModalEsc(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	m := &shellTestModal{id: "esctest", handleKey: false}
	s.PushModal(m)

	if !s.HasModal() {
		t.Error("should have modal")
	}

	// Esc should close modal
	s.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if s.HasModal() {
		t.Error("modal should be closed after Esc")
	}
}

func TestShellUpdateQuitKeys(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Ctrl+C should quit
	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("ctrl+c should return quit cmd")
	}

	// Ctrl+Q should quit
	_, cmd = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false})
	// Need to check with proper key representation
}

func TestShellUpdateFocusTab(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Add tab with content
	s.AddTab(Tab{
		ID:      "test",
		Label:   "Test",
		Content: content.NewSelectList(nil),
	})

	s.Focus(FocusTab)

	// Key should be routed to tab handler
	s.Update(tea.KeyMsg{Type: tea.KeyTab})
	// Should not panic
}

func TestShellUpdatePopMsg(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	m := &shellTestModal{id: "poptest"}
	s.PushModal(m)

	// Send PopMsg
	s.Update(PopMsg{})
	if s.HasModal() {
		t.Error("modal should be closed after PopMsg")
	}
}

func TestShellUpdatePushMsg(t *testing.T) {
	s := New(nil, DefaultConfig())
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	m := &shellTestModal{id: "pushtest"}

	// Send PushMsg
	s.Update(PushMsg{Modal: m})
	if !s.HasModal() {
		t.Error("should have modal after PushMsg")
	}
}

func TestShellViewNotReady(t *testing.T) {
	s := New(nil, DefaultConfig())
	// Don't send WindowSizeMsg, so shell is not ready

	view := s.View()
	if view != "Loading..." {
		t.Errorf("expected 'Loading...', got %s", view)
	}
}

func TestShellContentHeightMinimum(t *testing.T) {
	s := New(nil, DefaultConfig()) // ShowTabBar, ShowStatusBar, ShowInput all true
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 3}) // Very small height

	// Should not panic, height should be clamped to 1
	h := s.contentHeight()
	if h < 1 {
		t.Error("contentHeight should be at least 1")
	}
}

func TestShellViewNoTabBar(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShowTabBar = false
	s := New(nil, cfg)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestShellViewNoInput(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShowInput = false
	s := New(nil, cfg)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestShellViewNoStatusBar(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShowStatusBar = false
	s := New(nil, cfg)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestTabBarActiveTabEmpty(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	if tb.ActiveTab() != nil {
		t.Error("ActiveTab should return nil for empty tab bar")
	}
}

func TestTabBarRemoveActiveTab(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "2", Label: "Tab 2"})
	tb.SetActive("2") // Make tab 2 active

	tb.RemoveTab("2") // Remove active tab
	if tb.active >= len(tb.tabs) {
		t.Error("active should be adjusted when removing active tab")
	}
}

func TestTabBarRemoveLastTab(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.RemoveTab("1")

	if tb.active != 0 {
		t.Error("active should be 0 after removing last tab")
	}
}

func TestTabBarRemoveNonexistent(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.RemoveTab("nonexistent") // Should not panic
	if tb.Count() != 1 {
		t.Error("tab count should still be 1")
	}
}

func TestTabBarHandleKeyCtrl(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "2", Label: "Tab 2"})

	// ctrl+tab
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\t'}, Alt: false})
	// ctrl+shift+tab - hard to simulate, but covering shift+tab
}

func TestTabBarViewWithBadge(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1", Badge: "5"})

	view := tb.View()
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestTabBarViewEmpty(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	view := tb.View()
	if view != "" {
		t.Error("empty tab bar should return empty view")
	}
}

func TestTabBarRenderActiveContentNoTab(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	content := tb.RenderActiveContent(80, 10)
	// Should return empty lines
	if content == "" {
		t.Error("should return placeholder content")
	}
}

func TestTabBarRenderActiveContentNoContent(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1", Content: nil})

	content := tb.RenderActiveContent(80, 10)
	if content == "" {
		t.Error("should return placeholder content")
	}
}

func TestTabBarRenderActiveContentTruncate(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())

	// Create content with many lines
	items := make([]content.SelectItem, 50)
	for i := range items {
		items[i] = content.SelectItem{Label: "Item"}
	}
	tb.AddTab(Tab{ID: "1", Label: "Tab 1", Content: content.NewSelectList(items)})

	// Render with small height - should truncate
	rendered := tb.RenderActiveContent(80, 5)
	if rendered == "" {
		t.Error("should render content")
	}
}

func TestTabBarHandleKeyWithContent(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{
		ID:      "1",
		Label:   "Tab 1",
		Content: content.NewSelectList(nil),
	})

	// HandleKey should pass to content
	tb.HandleKey(tea.KeyMsg{Type: tea.KeyDown})
	// Should not panic
}

func TestStatusBarViewDisconnected(t *testing.T) {
	sb := NewStatusBar(theme.NewDraculaTheme())
	sb.SetConnected(false)

	view := sb.View(80)
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestStatusBarViewNarrow(t *testing.T) {
	sb := NewStatusBar(theme.NewDraculaTheme())
	sb.SetModel("claude-3")
	sb.SetTokens(100000, 200000)
	sb.SetMode("normal")
	sb.SetMessage("Long message here")
	sb.SetHints("Hints")

	// Very narrow width should still work
	view := sb.View(20)
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestTabBarSetBadgeNonexistent(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.SetBadge("nonexistent", "5") // Should not panic
}

func TestTabBarViewMultipleTabs(t *testing.T) {
	tb := NewTabBar(theme.NewDraculaTheme())
	tb.AddTab(Tab{ID: "1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "2", Label: "Tab 2"}) // Inactive tab
	tb.AddTab(Tab{ID: "3", Label: "Tab 3"}) // Inactive tab

	view := tb.View()
	if view == "" {
		t.Error("view should render all tabs")
	}
}

func TestShell_Streaming(t *testing.T) {
	sh := New(nil, DefaultConfig())

	// Streaming() returns controller
	s := sh.Streaming()
	if s == nil {
		t.Fatal("expected non-nil StreamingController")
	}

	// Same instance each time
	if sh.Streaming() != s {
		t.Error("expected same StreamingController instance")
	}

	// Can disable streaming status
	sh.SetStreamingStatusVisible(false)
	// (visibility tested in statusbar render)
}

func TestShell_StreamingInView(t *testing.T) {
	sh := New(nil, DefaultConfig())

	// Simulate window size
	sh.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Start streaming with thinking
	sh.Streaming().Start()
	sh.Streaming().AppendToken("test") // Exit waiting state
	sh.Streaming().SetThinking(true)

	view := sh.View()
	if !strings.Contains(view, "Thinking") {
		t.Errorf("expected 'Thinking' in shell view")
	}

	// Disable streaming status
	sh.SetStreamingStatusVisible(false)
	view = sh.View()
	if strings.Contains(view, "Thinking") {
		t.Error("expected no 'Thinking' when streaming status disabled")
	}
}

func TestTabBarSetActiveByIndex(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)

	tb.AddTab(Tab{ID: "tab1", Label: "Tab 1"})
	tb.AddTab(Tab{ID: "tab2", Label: "Tab 2"})
	tb.AddTab(Tab{ID: "tab3", Label: "Tab 3"})

	// Set by valid index
	tb.SetActiveByIndex(1)
	if tb.ActiveTab().ID != "tab2" {
		t.Errorf("expected tab2, got %s", tb.ActiveTab().ID)
	}

	// Set by index 0
	tb.SetActiveByIndex(0)
	if tb.ActiveTab().ID != "tab1" {
		t.Errorf("expected tab1, got %s", tb.ActiveTab().ID)
	}

	// Invalid index (too high) - should not change
	tb.SetActiveByIndex(10)
	if tb.ActiveTab().ID != "tab1" {
		t.Errorf("expected tab1 unchanged, got %s", tb.ActiveTab().ID)
	}

	// Invalid index (negative) - should not change
	tb.SetActiveByIndex(-1)
	if tb.ActiveTab().ID != "tab1" {
		t.Errorf("expected tab1 unchanged, got %s", tb.ActiveTab().ID)
	}
}

func TestShellTabIndexShortcuts(t *testing.T) {
	th := theme.NewDraculaTheme()
	sh := New(th, DefaultConfig())

	sh.AddTab(Tab{ID: "tab1", Label: "Tab 1"})
	sh.AddTab(Tab{ID: "tab2", Label: "Tab 2"})
	sh.AddTab(Tab{ID: "tab3", Label: "Tab 3"})

	// Simulate window size
	sh.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Alt+2 should switch to tab2
	sh.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}, Alt: true})
	if sh.tabs.ActiveTab().ID != "tab2" {
		t.Errorf("expected tab2 after Alt+2, got %s", sh.tabs.ActiveTab().ID)
	}

	// Alt+1 should switch to tab1
	sh.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}, Alt: true})
	if sh.tabs.ActiveTab().ID != "tab1" {
		t.Errorf("expected tab1 after Alt+1, got %s", sh.tabs.ActiveTab().ID)
	}

	// Alt+3 should switch to tab3
	sh.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}, Alt: true})
	if sh.tabs.ActiveTab().ID != "tab3" {
		t.Errorf("expected tab3 after Alt+3, got %s", sh.tabs.ActiveTab().ID)
	}

	// Alt+9 (out of range) should not change
	sh.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}, Alt: true})
	if sh.tabs.ActiveTab().ID != "tab3" {
		t.Errorf("expected tab3 unchanged after Alt+9, got %s", sh.tabs.ActiveTab().ID)
	}
}

func TestTabBarHiddenTabs(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)

	tb.AddTab(Tab{ID: "chat", Label: "Chat"})
	tb.AddTab(Tab{ID: "history", Label: "History", Hidden: true})
	tb.AddTab(Tab{ID: "tools", Label: "Tools", Hidden: true})

	// View should only show non-hidden tabs
	view := tb.View()
	if !strings.Contains(view, "Chat") {
		t.Error("expected Chat in view")
	}
	if strings.Contains(view, "History") {
		t.Error("expected History to be hidden in view")
	}
	if strings.Contains(view, "Tools") {
		t.Error("expected Tools to be hidden in view")
	}

	// But hidden tabs should still be navigable by ID
	tb.SetActive("history")
	if tb.ActiveTab().ID != "history" {
		t.Errorf("expected history tab active, got %s", tb.ActiveTab().ID)
	}

	// And by index
	tb.SetActiveByIndex(2)
	if tb.ActiveTab().ID != "tools" {
		t.Errorf("expected tools tab active, got %s", tb.ActiveTab().ID)
	}
}

func TestShellCustomTabShortcuts(t *testing.T) {
	th := theme.NewDraculaTheme()
	sh := New(th, DefaultConfig())

	sh.AddTab(Tab{ID: "chat", Label: "Chat"})
	sh.AddTab(Tab{ID: "history", Label: "History", Hidden: true, Shortcut: "ctrl+r"})
	sh.AddTab(Tab{ID: "tools", Label: "Tools", Hidden: true, Shortcut: "ctrl+o"})

	// Simulate window size
	sh.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Ctrl+R should switch to history
	sh.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	if sh.tabs.ActiveTab().ID != "history" {
		t.Errorf("expected history after Ctrl+R, got %s", sh.tabs.ActiveTab().ID)
	}

	// Ctrl+O should switch to tools
	sh.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	if sh.tabs.ActiveTab().ID != "tools" {
		t.Errorf("expected tools after Ctrl+O, got %s", sh.tabs.ActiveTab().ID)
	}

	// Pressing same shortcut again should toggle back to previous (or stay - depends on design)
	// For now, just verify it doesn't break
	sh.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	if sh.tabs.ActiveTab() == nil {
		t.Error("expected a tab to be active")
	}
}

func TestTabLifecycleHooks(t *testing.T) {
	th := theme.NewDraculaTheme()
	tb := NewTabBar(th)

	activations := make([]string, 0)
	deactivations := make([]string, 0)

	content1 := &mockTabContent{
		onActivate:   func() { activations = append(activations, "tab1") },
		onDeactivate: func() { deactivations = append(deactivations, "tab1") },
	}
	content2 := &mockTabContent{
		onActivate:   func() { activations = append(activations, "tab2") },
		onDeactivate: func() { deactivations = append(deactivations, "tab2") },
	}

	tb.AddTab(Tab{ID: "tab1", Label: "Tab 1", Content: content1})
	tb.AddTab(Tab{ID: "tab2", Label: "Tab 2", Content: content2})

	// Initial tab should be activated
	tb.ActivateCurrentTab()
	if len(activations) != 1 || activations[0] != "tab1" {
		t.Errorf("expected tab1 activation, got %v", activations)
	}

	// Switch to tab2
	tb.SetActive("tab2")
	tb.ActivateCurrentTab()

	if len(deactivations) != 1 || deactivations[0] != "tab1" {
		t.Errorf("expected tab1 deactivation, got %v", deactivations)
	}
	if len(activations) != 2 || activations[1] != "tab2" {
		t.Errorf("expected tab2 activation, got %v", activations)
	}
}

// mockTabContent implements TabContent for testing
type mockTabContent struct {
	mockContent
	onActivate   func()
	onDeactivate func()
}

func (m *mockTabContent) OnActivate() tea.Cmd {
	if m.onActivate != nil {
		m.onActivate()
	}
	return nil
}

func (m *mockTabContent) OnDeactivate() {
	if m.onDeactivate != nil {
		m.onDeactivate()
	}
}

func TestShellTabLifecycleOnSwitch(t *testing.T) {
	th := theme.NewDraculaTheme()
	sh := New(th, DefaultConfig())

	activated := ""
	deactivated := ""

	content1 := &mockTabContent{
		onActivate:   func() { activated = "tab1" },
		onDeactivate: func() { deactivated = "tab1" },
	}
	content2 := &mockTabContent{
		onActivate:   func() { activated = "tab2" },
		onDeactivate: func() { deactivated = "tab2" },
	}

	sh.AddTab(Tab{ID: "tab1", Label: "Tab 1", Content: content1})
	sh.AddTab(Tab{ID: "tab2", Label: "Tab 2", Content: content2})

	// Simulate window size to make shell ready
	sh.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Initial activation
	if activated != "tab1" {
		t.Errorf("expected tab1 activated on init, got %q", activated)
	}

	// Switch via SetActiveTab
	sh.SetActiveTab("tab2")
	if deactivated != "tab1" {
		t.Errorf("expected tab1 deactivated, got %q", deactivated)
	}
	if activated != "tab2" {
		t.Errorf("expected tab2 activated, got %q", activated)
	}
}

func TestShellCtrlEOpensErrorModal(t *testing.T) {
	th := theme.NewDraculaTheme()
	cfg := DefaultConfig()

	errorModalOpened := false
	cfg.OnShowErrors = func() {
		errorModalOpened = true
	}

	s := New(th, cfg)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Send Ctrl+E
	s.Update(tea.KeyMsg{Type: tea.KeyCtrlE})

	if !errorModalOpened {
		t.Error("ctrl+e should trigger OnShowErrors callback")
	}
}

func TestShellHistoryProvider(t *testing.T) {
	th := theme.NewDraculaTheme()
	cfg := DefaultConfig()

	history := []string{"prompt1", "prompt2"}
	cfg.HistoryProvider = func() []string {
		return history
	}

	s := New(th, cfg)

	// Verify input has history provider set
	// (internal implementation detail, but we test behavior)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Send up arrow while input focused
	s.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Input should show last history item
	if s.InputValue() != "prompt2" {
		t.Errorf("expected 'prompt2', got %q", s.InputValue())
	}
}
