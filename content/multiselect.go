package content

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelectItem represents an item in a multi-select list.
type MultiSelectItem struct {
	Label    string
	Key      string
	Selected bool
}

// MultiSelect is a multiple-selection list content primitive.
type MultiSelect struct {
	items  []MultiSelectItem
	cursor int
	width  int
	height int

	// Styles
	cursorStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	checkStyle      lipgloss.Style
}

// NewMultiSelect creates a new multi-select list.
func NewMultiSelect(items []MultiSelectItem) *MultiSelect {
	return &MultiSelect{
		items:  items,
		cursor: 0,
		cursorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
		unselectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		checkStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
	}
}

// Init implements Content.
func (m *MultiSelect) Init() tea.Cmd {
	return nil
}

// Update implements Content.
func (m *MultiSelect) Update(msg tea.Msg) (Content, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case tea.KeySpace:
			m.Toggle()
		}
		switch msg.String() {
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "x":
			m.Toggle()
		case "a":
			m.SelectAll()
		case "n":
			m.SelectNone()
		}
	}
	return m, nil
}

// View implements Content.
func (m *MultiSelect) View() string {
	if len(m.items) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")).Render("No items")
	}

	var b strings.Builder

	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = m.cursorStyle.Render("▸ ")
		}

		check := "[ ]"
		style := m.unselectedStyle
		if item.Selected {
			check = m.checkStyle.Render("[✓]")
			style = m.selectedStyle
		}

		b.WriteString(cursor)
		b.WriteString(check)
		b.WriteString(" ")
		b.WriteString(style.Render(item.Label))

		if i < len(m.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Value implements Content. Returns slice of selected keys.
func (m *MultiSelect) Value() any {
	var selected []string
	for _, item := range m.items {
		if item.Selected {
			selected = append(selected, item.Key)
		}
	}
	return selected
}

// SetSize implements Content.
func (m *MultiSelect) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Toggle toggles the selection of the current item.
func (m *MultiSelect) Toggle() {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		m.items[m.cursor].Selected = !m.items[m.cursor].Selected
	}
}

// SelectAll selects all items.
func (m *MultiSelect) SelectAll() {
	for i := range m.items {
		m.items[i].Selected = true
	}
}

// SelectNone deselects all items.
func (m *MultiSelect) SelectNone() {
	for i := range m.items {
		m.items[i].Selected = false
	}
}

// SelectedCount returns the number of selected items.
func (m *MultiSelect) SelectedCount() int {
	count := 0
	for _, item := range m.items {
		if item.Selected {
			count++
		}
	}
	return count
}

// SetItems replaces all items.
func (m *MultiSelect) SetItems(items []MultiSelectItem) {
	m.items = items
	if m.cursor >= len(items) {
		m.cursor = len(items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// Items returns all items.
func (m *MultiSelect) Items() []MultiSelectItem {
	return m.items
}
