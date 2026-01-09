# Tabs and Modals Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add tab keyboard shortcuts, hidden tabs, and lifecycle hooks to enable the tabs-and-modals pattern.

**Architecture:** Extend existing TabBar with index-based navigation, hidden tab support, custom shortcuts, and lifecycle hooks. Shell handles global tab shortcuts; TabBar manages tab state.

**Tech Stack:** Go, Bubble Tea, existing tux shell components

---

### Task 1: Add SetActiveByIndex to TabBar

**Files:**
- Modify: `shell/tabbar.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestTabBarSetActiveByIndex -v`
Expected: FAIL with "tb.SetActiveByIndex undefined"

**Step 3: Write minimal implementation**

```go
// Add to shell/tabbar.go after SetActive method

// SetActiveByIndex sets the active tab by index (0-based).
// Does nothing if index is out of range.
func (t *TabBar) SetActiveByIndex(index int) {
	if index >= 0 && index < len(t.tabs) {
		t.active = index
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestTabBarSetActiveByIndex -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/tabbar.go shell/shell_test.go
git commit -m "feat(shell): add SetActiveByIndex to TabBar"
```

---

### Task 2: Add Tab Index Shortcuts to Shell

**Files:**
- Modify: `shell/shell.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShellTabIndexShortcuts -v`
Expected: FAIL (Alt+N doesn't switch tabs)

**Step 3: Write minimal implementation**

```go
// In shell/shell.go Update() method, add after global keys section:

		// Tab index shortcuts (Alt+1 through Alt+9)
		if msg.Alt && len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r >= '1' && r <= '9' {
				index := int(r - '1') // '1' -> 0, '2' -> 1, etc.
				s.tabs.SetActiveByIndex(index)
				return s, nil
			}
		}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShellTabIndexShortcuts -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): add Alt+1-9 tab index shortcuts"
```

---

### Task 3: Add Hidden Tabs Support

**Files:**
- Modify: `shell/tabbar.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestTabBarHiddenTabs -v`
Expected: FAIL with "unknown field 'Hidden'"

**Step 3: Write minimal implementation**

```go
// Update Tab struct in shell/tabbar.go

// Tab represents a single tab in the tab bar.
type Tab struct {
	ID       string
	Label    string
	Badge    string
	Content  content.Content
	Closable bool
	Hidden   bool // Hidden tabs are accessible but not shown in tab bar
}

// Update View() method in shell/tabbar.go

// View renders the tab bar.
func (t *TabBar) View() string {
	if len(t.tabs) == 0 {
		return ""
	}

	styles := t.theme.Styles()
	var tabs []string

	for i, tab := range t.tabs {
		if tab.Hidden {
			continue
		}

		label := tab.Label
		if tab.Badge != "" {
			label += " " + tab.Badge
		}

		var style lipgloss.Style
		if i == t.active {
			style = styles.TabActive
		} else {
			style = styles.TabInactive
		}

		tabs = append(tabs, style.Render(label))
	}

	if len(tabs) == 0 {
		return ""
	}

	return strings.Join(tabs, "  ")
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestTabBarHiddenTabs -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/tabbar.go shell/shell_test.go
git commit -m "feat(shell): add hidden tabs support"
```

---

### Task 4: Add Custom Tab Shortcuts

**Files:**
- Modify: `shell/tabbar.go`
- Modify: `shell/shell.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShellCustomTabShortcuts -v`
Expected: FAIL with "unknown field 'Shortcut'"

**Step 3: Write minimal implementation**

```go
// Update Tab struct in shell/tabbar.go

// Tab represents a single tab in the tab bar.
type Tab struct {
	ID       string
	Label    string
	Badge    string
	Content  content.Content
	Closable bool
	Hidden   bool   // Hidden tabs are accessible but not shown in tab bar
	Shortcut string // Keyboard shortcut to activate this tab (e.g., "ctrl+r")
}

// Add to TabBar in shell/tabbar.go

// FindByShortcut returns the tab ID matching the given shortcut, or empty string.
func (t *TabBar) FindByShortcut(shortcut string) string {
	for _, tab := range t.tabs {
		if tab.Shortcut == shortcut {
			return tab.ID
		}
	}
	return ""
}

// In shell/shell.go Update(), add after tab index shortcuts:

		// Custom tab shortcuts
		shortcut := keyMsgToShortcut(msg)
		if shortcut != "" {
			if tabID := s.tabs.FindByShortcut(shortcut); tabID != "" {
				s.tabs.SetActive(tabID)
				return s, nil
			}
		}

// Add helper function in shell/shell.go

// keyMsgToShortcut converts a tea.KeyMsg to a shortcut string.
func keyMsgToShortcut(msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyCtrlA:
		return "ctrl+a"
	case tea.KeyCtrlB:
		return "ctrl+b"
	case tea.KeyCtrlD:
		return "ctrl+d"
	case tea.KeyCtrlE:
		return "ctrl+e"
	case tea.KeyCtrlF:
		return "ctrl+f"
	case tea.KeyCtrlG:
		return "ctrl+g"
	case tea.KeyCtrlH:
		return "ctrl+h"
	case tea.KeyCtrlI:
		return "ctrl+i"
	case tea.KeyCtrlJ:
		return "ctrl+j"
	case tea.KeyCtrlK:
		return "ctrl+k"
	case tea.KeyCtrlL:
		return "ctrl+l"
	case tea.KeyCtrlN:
		return "ctrl+n"
	case tea.KeyCtrlO:
		return "ctrl+o"
	case tea.KeyCtrlP:
		return "ctrl+p"
	case tea.KeyCtrlR:
		return "ctrl+r"
	case tea.KeyCtrlS:
		return "ctrl+s"
	case tea.KeyCtrlT:
		return "ctrl+t"
	case tea.KeyCtrlU:
		return "ctrl+u"
	case tea.KeyCtrlV:
		return "ctrl+v"
	case tea.KeyCtrlW:
		return "ctrl+w"
	case tea.KeyCtrlX:
		return "ctrl+x"
	case tea.KeyCtrlY:
		return "ctrl+y"
	case tea.KeyCtrlZ:
		return "ctrl+z"
	}
	return ""
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShellCustomTabShortcuts -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/tabbar.go shell/shell.go shell/shell_test.go
git commit -m "feat(shell): add custom tab shortcuts"
```

---

### Task 5: Add Tab Lifecycle Hooks

**Files:**
- Create: `shell/tab_content.go`
- Modify: `shell/tabbar.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestTabLifecycleHooks -v`
Expected: FAIL with undefined types/methods

**Step 3: Write minimal implementation**

```go
// Create shell/tab_content.go

package shell

import (
	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

// TabContent extends content.Content with lifecycle hooks.
// Content implementations can optionally implement this interface
// to receive notifications when the tab becomes active or inactive.
type TabContent interface {
	content.Content

	// OnActivate is called when the tab becomes active.
	// Return a command to run on activation (e.g., start a timer).
	OnActivate() tea.Cmd

	// OnDeactivate is called when the tab becomes inactive.
	OnDeactivate()
}

// Add to shell/tabbar.go

// Track previous active tab for lifecycle hooks
// Add field to TabBar struct:
//     lastActive int

// Add method to TabBar:

// ActivateCurrentTab calls lifecycle hooks when switching tabs.
// Call this after changing the active tab.
func (t *TabBar) ActivateCurrentTab() tea.Cmd {
	// Deactivate previous tab
	if t.lastActive >= 0 && t.lastActive < len(t.tabs) && t.lastActive != t.active {
		if tc, ok := t.tabs[t.lastActive].Content.(TabContent); ok {
			tc.OnDeactivate()
		}
	}

	// Activate current tab
	var cmd tea.Cmd
	if t.active >= 0 && t.active < len(t.tabs) {
		if tc, ok := t.tabs[t.active].Content.(TabContent); ok {
			cmd = tc.OnActivate()
		}
	}

	t.lastActive = t.active
	return cmd
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestTabLifecycleHooks -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/tab_content.go shell/tabbar.go shell/shell_test.go
git commit -m "feat(shell): add tab lifecycle hooks (OnActivate/OnDeactivate)"
```

---

### Task 6: Wire Lifecycle Hooks in Shell

**Files:**
- Modify: `shell/shell.go`
- Modify: `shell/shell_test.go`

**Step 1: Write the failing test**

```go
// Add to shell/shell_test.go

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./shell -run TestShellTabLifecycleOnSwitch -v`
Expected: FAIL (lifecycle hooks not called by shell)

**Step 3: Write minimal implementation**

```go
// Update shell/shell.go

// In Init() method, add:
func (s *Shell) Init() tea.Cmd {
	return tea.Batch(
		s.input.Init(),
		s.tabs.ActivateCurrentTab(), // Activate initial tab
	)
}

// Update SetActiveTab:
// SetActiveTab switches to the tab with the given ID.
func (s *Shell) SetActiveTab(id string) tea.Cmd {
	s.tabs.SetActive(id)
	return s.tabs.ActivateCurrentTab()
}

// Note: The return type changes from void to tea.Cmd
// This is a breaking change but necessary for lifecycle hooks
```

**Step 4: Run test to verify it passes**

Run: `go test ./shell -run TestShellTabLifecycleOnSwitch -v`
Expected: PASS

**Step 5: Commit**

```bash
git add shell/shell.go shell/shell_test.go
git commit -m "feat(shell): wire tab lifecycle hooks in shell"
```

---

### Task 7: Final Verification

**Files:**
- None (verification only)

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: All tests pass

**Step 2: Run go vet**

Run: `go vet ./...`
Expected: No issues

**Step 3: Verify feature summary**

Verify the following features work:
- [ ] `SetActiveByIndex(n)` - switch tab by 0-based index
- [ ] Alt+1-9 shortcuts - switch to tab by number
- [ ] Hidden tabs - tabs not shown in bar but accessible
- [ ] Custom shortcuts - `Shortcut: "ctrl+r"` activates tab
- [ ] Lifecycle hooks - `OnActivate()` and `OnDeactivate()` called on switch

**Step 4: Done**

All tabs-and-modals features implemented.
