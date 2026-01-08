package shell

import (
	"strings"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab in the tab bar.
type Tab struct {
	ID       string
	Label    string
	Badge    string
	Content  content.Content
	Closable bool
	Hidden   bool // Hidden tabs are accessible but not shown in tab bar
}

// TabBar manages tabs and renders the tab bar.
type TabBar struct {
	tabs     []Tab
	active   int
	width    int
	height   int
	theme    theme.Theme
}

// NewTabBar creates a new tab bar.
func NewTabBar(th theme.Theme) *TabBar {
	return &TabBar{
		tabs:  make([]Tab, 0),
		theme: th,
	}
}

// AddTab adds a tab.
func (t *TabBar) AddTab(tab Tab) {
	t.tabs = append(t.tabs, tab)
	if len(t.tabs) == 1 {
		t.active = 0
	}
}

// RemoveTab removes a tab by ID.
func (t *TabBar) RemoveTab(id string) {
	for i, tab := range t.tabs {
		if tab.ID == id {
			t.tabs = append(t.tabs[:i], t.tabs[i+1:]...)
			if t.active >= len(t.tabs) {
				t.active = len(t.tabs) - 1
			}
			if t.active < 0 {
				t.active = 0
			}
			return
		}
	}
}

// SetActive sets the active tab by ID.
func (t *TabBar) SetActive(id string) {
	for i, tab := range t.tabs {
		if tab.ID == id {
			t.active = i
			return
		}
	}
}

// SetActiveByIndex sets the active tab by index (0-based).
// Does nothing if index is out of range.
func (t *TabBar) SetActiveByIndex(index int) {
	if index >= 0 && index < len(t.tabs) {
		t.active = index
	}
}

// ActiveTab returns the currently active tab.
func (t *TabBar) ActiveTab() *Tab {
	if t.active >= 0 && t.active < len(t.tabs) {
		return &t.tabs[t.active]
	}
	return nil
}

// SetSize sets the available size.
func (t *TabBar) SetSize(width, height int) {
	t.width = width
	t.height = height
	// Update content sizes
	for i := range t.tabs {
		if t.tabs[i].Content != nil {
			t.tabs[i].Content.SetSize(width, height)
		}
	}
}

// HandleKey handles keyboard input for tab navigation.
func (t *TabBar) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "tab", "ctrl+tab":
		t.NextTab()
	case "shift+tab", "ctrl+shift+tab":
		t.PrevTab()
	}

	// Pass to active content
	if tab := t.ActiveTab(); tab != nil && tab.Content != nil {
		tab.Content.Update(msg)
	}

	return nil
}

// NextTab switches to the next tab.
func (t *TabBar) NextTab() {
	if len(t.tabs) > 0 {
		t.active = (t.active + 1) % len(t.tabs)
	}
}

// PrevTab switches to the previous tab.
func (t *TabBar) PrevTab() {
	if len(t.tabs) > 0 {
		t.active--
		if t.active < 0 {
			t.active = len(t.tabs) - 1
		}
	}
}

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

// RenderActiveContent renders the content of the active tab.
func (t *TabBar) RenderActiveContent(width, height int) string {
	tab := t.ActiveTab()
	if tab == nil || tab.Content == nil {
		return strings.Repeat("\n", height-1)
	}

	tab.Content.SetSize(width, height)
	view := tab.Content.View()

	// Pad to fill height
	lines := strings.Split(view, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// Count returns the number of tabs.
func (t *TabBar) Count() int {
	return len(t.tabs)
}

// SetBadge sets the badge on a tab.
func (t *TabBar) SetBadge(id string, badge string) {
	for i := range t.tabs {
		if t.tabs[i].ID == id {
			t.tabs[i].Badge = badge
			return
		}
	}
}
