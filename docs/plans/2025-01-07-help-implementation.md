# Help Package Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a static help overlay component with categories, mode filtering, and modal integration.

**Architecture:** Categories contain bindings with optional mode tags. Help component filters and renders. HelpModal adapter for modal stack integration.

**Tech Stack:** Go, lipgloss (styling), bubbletea (modal integration)

---

### Task 1: Binding and Category Types

**Files:**
- Create: `help/binding.go`
- Create: `help/binding_test.go`

**Step 1: Write the failing test**

```go
// help/binding_test.go
package help

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./help/... -v`
Expected: FAIL - package doesn't exist

**Step 3: Write minimal implementation**

```go
// help/binding.go
package help

// Binding represents a keyboard shortcut or command.
type Binding struct {
	Key         string   // Display text: "ctrl+c", "?", "/feedback"
	Description string   // What it does: "Quit", "Toggle help"
	Modes       []string // Optional: empty = all modes
}

// MatchesMode returns true if this binding should be shown for the given mode.
// Empty mode parameter means show all bindings.
// Empty Modes field means binding is shown in all modes.
func (b Binding) MatchesMode(mode string) bool {
	if mode == "" || len(b.Modes) == 0 {
		return true
	}
	for _, m := range b.Modes {
		if m == mode {
			return true
		}
	}
	return false
}

// Category groups related bindings.
type Category struct {
	Title    string
	Bindings []Binding
}

// FilterByMode returns bindings that match the given mode.
func (c Category) FilterByMode(mode string) []Binding {
	var result []Binding
	for _, b := range c.Bindings {
		if b.MatchesMode(mode) {
			result = append(result, b)
		}
	}
	return result
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./help/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add help/binding.go help/binding_test.go
git commit -m "feat(help): add Binding and Category types with mode filtering"
```

---

### Task 2: Help Component

**Files:**
- Create: `help/help.go`
- Create: `help/help_test.go`

**Step 1: Write the failing test**

```go
// help/help_test.go
package help

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestHelpNew(t *testing.T) {
	h := New(
		Category{Title: "General", Bindings: []Binding{{Key: "?", Description: "Help"}}},
	)
	if h == nil {
		t.Fatal("New should return non-nil Help")
	}
}

func TestHelpWithTheme(t *testing.T) {
	h := New().WithTheme(theme.Default())
	if h.theme == nil {
		t.Error("WithTheme should set theme")
	}
}

func TestHelpRender(t *testing.T) {
	h := New(
		Category{
			Title: "Actions",
			Bindings: []Binding{
				{Key: "ctrl+c", Description: "Quit"},
				{Key: "?", Description: "Help"},
			},
		},
	).WithTheme(theme.Default())

	output := h.Render(60, "")

	if !strings.Contains(output, "Actions") {
		t.Error("should contain category title")
	}
	if !strings.Contains(output, "ctrl+c") {
		t.Error("should contain key")
	}
	if !strings.Contains(output, "Quit") {
		t.Error("should contain description")
	}
}

func TestHelpRenderModeFilter(t *testing.T) {
	h := New(
		Category{
			Title: "Input",
			Bindings: []Binding{
				{Key: "enter", Description: "Send", Modes: []string{"chat"}},
				{Key: "enter", Description: "Select", Modes: []string{"list"}},
			},
		},
	).WithTheme(theme.Default())

	// Chat mode
	output := h.Render(60, "chat")
	if !strings.Contains(output, "Send") {
		t.Error("chat mode should show Send")
	}
	if strings.Contains(output, "Select") {
		t.Error("chat mode should not show Select")
	}

	// List mode
	output = h.Render(60, "list")
	if strings.Contains(output, "Send") {
		t.Error("list mode should not show Send")
	}
	if !strings.Contains(output, "Select") {
		t.Error("list mode should show Select")
	}
}

func TestHelpRenderEmptyCategoryHidden(t *testing.T) {
	h := New(
		Category{
			Title: "Chat Only",
			Bindings: []Binding{
				{Key: "enter", Description: "Send", Modes: []string{"chat"}},
			},
		},
		Category{
			Title: "Always",
			Bindings: []Binding{
				{Key: "?", Description: "Help"},
			},
		},
	).WithTheme(theme.Default())

	// In list mode, "Chat Only" category should be hidden
	output := h.Render(60, "list")
	if strings.Contains(output, "Chat Only") {
		t.Error("empty category should be hidden")
	}
	if !strings.Contains(output, "Always") {
		t.Error("non-empty category should be shown")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./help/... -v`
Expected: FAIL - Help type doesn't exist

**Step 3: Write minimal implementation**

```go
// help/help.go
package help

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

// Help displays keyboard shortcuts and commands.
type Help struct {
	categories []Category
	theme      theme.Theme
}

// New creates a new Help component with the given categories.
func New(categories ...Category) *Help {
	return &Help{
		categories: categories,
		theme:      theme.Default(),
	}
}

// WithTheme sets the theme for rendering.
func (h *Help) WithTheme(th theme.Theme) *Help {
	h.theme = th
	return h
}

// Render returns the help overlay content.
// If mode is empty, shows all bindings.
// If mode is set, shows only bindings that match.
func (h *Help) Render(width int, mode string) string {
	if h.theme == nil {
		h.theme = theme.Default()
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(h.theme.Primary()).
		Bold(true).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(h.theme.Secondary()).
		Bold(true).
		Width(12)

	descStyle := lipgloss.NewStyle().
		Foreground(h.theme.Foreground())

	footerStyle := lipgloss.NewStyle().
		Foreground(h.theme.Muted()).
		MarginTop(1)

	var sections []string

	for _, cat := range h.categories {
		filtered := cat.FilterByMode(mode)
		if len(filtered) == 0 {
			continue
		}

		lines := []string{titleStyle.Render(cat.Title)}
		for _, b := range filtered {
			key := keyStyle.Render(b.Key)
			desc := descStyle.Render(b.Description)
			lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, key, desc))
		}
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, lines...))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	footer := footerStyle.Render("Press ? or esc to close")
	content = lipgloss.JoinVertical(lipgloss.Left, content, footer)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.theme.Border()).
		Padding(1, 2).
		Width(width)

	return boxStyle.Render(content)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./help/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add help/help.go help/help_test.go
git commit -m "feat(help): add Help component with Render and mode filtering"
```

---

### Task 3: HelpModal Adapter

**Files:**
- Create: `modal/help.go`
- Create: `modal/help_test.go`

**Step 1: Write the failing test**

```go
// modal/help_test.go
package modal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/help"
)

func TestHelpModal(t *testing.T) {
	h := help.New(
		help.Category{
			Title: "Test",
			Bindings: []help.Binding{{Key: "?", Description: "Help"}},
		},
	)

	m := NewHelpModal(HelpModalConfig{
		ID:    "test-help",
		Title: "Test Help",
		Help:  h,
	})

	if m.ID() != "test-help" {
		t.Errorf("expected ID 'test-help', got %s", m.ID())
	}
	if m.Title() != "Test Help" {
		t.Errorf("expected title 'Test Help', got %s", m.Title())
	}
}

func TestHelpModalDefaultID(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	if m.ID() != "help-modal" {
		t.Errorf("expected default 'help-modal', got %s", m.ID())
	}
}

func TestHelpModalDefaultTitle(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	if m.Title() != "Help" {
		t.Errorf("expected default 'Help', got %s", m.Title())
	}
}

func TestHelpModalCloseOnEscape(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)

	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})
	if !handled {
		t.Error("escape should be handled")
	}
	if cmd == nil {
		t.Error("escape should return PopMsg")
	}
}

func TestHelpModalCloseOnQuestionMark(t *testing.T) {
	h := help.New()
	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)

	handled, cmd := m.HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !handled {
		t.Error("? should be handled")
	}
	if cmd == nil {
		t.Error("? should return PopMsg")
	}
}

func TestHelpModalRender(t *testing.T) {
	h := help.New(
		help.Category{
			Title: "Actions",
			Bindings: []help.Binding{{Key: "ctrl+c", Description: "Quit"}},
		},
	)

	m := NewHelpModal(HelpModalConfig{Help: h})
	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestHelpModalWithMode(t *testing.T) {
	h := help.New(
		help.Category{
			Title: "Input",
			Bindings: []help.Binding{
				{Key: "enter", Description: "Send", Modes: []string{"chat"}},
			},
		},
	)

	m := NewHelpModal(HelpModalConfig{
		Help: h,
		Mode: "chat",
	})

	if m.mode != "chat" {
		t.Errorf("expected mode 'chat', got %s", m.mode)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./modal/... -v -run HelpModal`
Expected: FAIL - HelpModal doesn't exist

**Step 3: Write minimal implementation**

```go
// modal/help.go
package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/help"
	"github.com/2389-research/tux/theme"
)

// HelpModalConfig configures a HelpModal.
type HelpModalConfig struct {
	ID    string
	Title string
	Help  *help.Help
	Size  Size
	Theme theme.Theme
	Mode  string // Current mode for filtering
}

// HelpModal adapts Help to work in the modal stack.
type HelpModal struct {
	id     string
	title  string
	help   *help.Help
	size   Size
	theme  theme.Theme
	mode   string
	width  int
	height int
}

// NewHelpModal creates a new HelpModal.
func NewHelpModal(cfg HelpModalConfig) *HelpModal {
	id := cfg.ID
	if id == "" {
		id = "help-modal"
	}
	title := cfg.Title
	if title == "" {
		title = "Help"
	}
	size := cfg.Size
	if size == 0 {
		size = SizeMedium
	}
	th := cfg.Theme
	if th == nil {
		th = theme.Default()
	}

	h := cfg.Help
	if h != nil {
		h = h.WithTheme(th)
	}

	return &HelpModal{
		id:    id,
		title: title,
		help:  h,
		size:  size,
		theme: th,
		mode:  cfg.Mode,
	}
}

// ID returns the modal identifier.
func (m *HelpModal) ID() string { return m.id }

// Title returns the modal title.
func (m *HelpModal) Title() string { return m.title }

// Size returns the modal size.
func (m *HelpModal) Size() Size { return m.size }

// OnPush is called when the modal is pushed onto the stack.
func (m *HelpModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop is called when the modal is popped from the stack.
func (m *HelpModal) OnPop() {}

// HandleKey processes key events.
func (m *HelpModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type {
	case tea.KeyEscape:
		return true, func() tea.Msg { return PopMsg{} }
	case tea.KeyRunes:
		if len(key.Runes) == 1 && key.Runes[0] == '?' {
			return true, func() tea.Msg { return PopMsg{} }
		}
	}
	return false, nil
}

// Render renders the help modal.
func (m *HelpModal) Render(width, height int) string {
	if m.help == nil {
		return ""
	}
	return m.help.Render(width, m.mode)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./modal/... -v -run HelpModal`
Expected: PASS

**Step 5: Commit**

```bash
git add modal/help.go modal/help_test.go
git commit -m "feat(modal): add HelpModal adapter"
```

---

### Task 4: Integration Test

**Files:**
- Create: `help/integration_test.go`

**Step 1: Write integration test**

```go
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
	).WithTheme(theme.Default())

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
```

**Step 2: Run tests**

Run: `go test ./help/... -v`
Expected: PASS

**Step 3: Commit**

```bash
git add help/integration_test.go
git commit -m "test(help): add integration tests"
```

---

### Task 5: Coverage and Cleanup

**Step 1: Run coverage**

Run: `go test ./help/... ./modal/... -coverprofile=coverage.out -covermode=atomic`
Run: `go tool cover -func=coverage.out | grep -E "(help|modal)"`

**Step 2: Add any missing tests to reach 95%+**

Check uncovered lines and add tests as needed.

**Step 3: Final commit**

```bash
git add -A
git commit -m "test(help): increase coverage to 95%+"
```

---

### Task 6: Update go.mod (if needed)

**Step 1: Tidy dependencies**

Run: `go mod tidy`

**Step 2: Verify build**

Run: `go build ./...`

**Step 3: Commit if changes**

```bash
git add go.mod go.sum
git commit -m "chore: update dependencies"
```
