package shell

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem represents an item in a list modal.
type ListItem struct {
	ID          string
	Title       string
	Description string
	Value       any
}

// ListModal presents a filterable list for selection.
type ListModal struct {
	id         string
	title      string
	items      []ListItem
	filtered   []ListItem
	selected   int
	filter     string
	filterable bool
	onSelect   func(item ListItem)
	onCancel   func()
	width      int
	height     int
	maxVisible int

	// Styles
	boxStyle      lipgloss.Style
	titleStyle    lipgloss.Style
	filterStyle   lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
	descStyle     lipgloss.Style
}

// ListModalConfig configures a ListModal.
type ListModalConfig struct {
	ID         string
	Title      string
	Items      []ListItem
	Filterable bool
	OnSelect   func(item ListItem)
	OnCancel   func()
}

// NewListModal creates a new list modal.
func NewListModal(cfg ListModalConfig) *ListModal {
	m := &ListModal{
		id:         cfg.ID,
		title:      cfg.Title,
		items:      cfg.Items,
		filtered:   cfg.Items,
		filterable: cfg.Filterable,
		onSelect:   cfg.OnSelect,
		onCancel:   cfg.OnCancel,
		maxVisible: 10,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")).
			Bold(true),
		filterStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")),
		itemStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")).
			Bold(true),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
	return m
}

// ID implements Modal.
func (m *ListModal) ID() string { return m.id }

// Title implements Modal.
func (m *ListModal) Title() string { return m.title }

// Size implements Modal.
func (m *ListModal) Size() Size { return SizeMedium }

// OnPush implements Modal.
func (m *ListModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop implements Modal.
func (m *ListModal) OnPop() {
	if m.onCancel != nil {
		m.onCancel()
	}
}

// HandleKey implements Modal.
func (m *ListModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type {
	case tea.KeyUp:
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case tea.KeyDown:
		if m.selected < len(m.filtered)-1 {
			m.selected++
		}
		return true, nil
	case tea.KeyEnter:
		if m.selected >= 0 && m.selected < len(m.filtered) {
			if m.onSelect != nil {
				m.onSelect(m.filtered[m.selected])
			}
			return true, func() tea.Msg { return PopMsg{} }
		}
		return true, nil
	case tea.KeyBackspace:
		if m.filterable && len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			m.applyFilter()
			return true, nil
		}
	case tea.KeyRunes:
		if m.filterable {
			m.filter += string(key.Runes)
			m.applyFilter()
			return true, nil
		}
	}

	switch key.String() {
	case "k":
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case "j":
		if m.selected < len(m.filtered)-1 {
			m.selected++
		}
		return true, nil
	}

	return false, nil
}

// applyFilter filters items based on the current filter string.
func (m *ListModal) applyFilter() {
	if m.filter == "" {
		m.filtered = m.items
		m.selected = 0
		return
	}

	filter := strings.ToLower(m.filter)
	m.filtered = make([]ListItem, 0)
	for _, item := range m.items {
		if strings.Contains(strings.ToLower(item.Title), filter) ||
			strings.Contains(strings.ToLower(item.Description), filter) {
			m.filtered = append(m.filtered, item)
		}
	}
	m.selected = 0
}

// Render implements Modal.
func (m *ListModal) Render(width, height int) string {
	var parts []string

	// Title
	if m.title != "" {
		parts = append(parts, m.titleStyle.Render(m.title))
	}

	// Filter
	if m.filterable {
		filterDisplay := m.filter
		if filterDisplay == "" {
			filterDisplay = "Type to filter..."
		}
		parts = append(parts, m.filterStyle.Render("> "+filterDisplay))
	}

	parts = append(parts, "")

	// Items
	if len(m.filtered) == 0 {
		parts = append(parts, m.descStyle.Render("No items match"))
	} else {
		start := 0
		if m.selected >= m.maxVisible {
			start = m.selected - m.maxVisible + 1
		}
		end := start + m.maxVisible
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		for i := start; i < end; i++ {
			item := m.filtered[i]
			prefix := "  "
			style := m.itemStyle
			if i == m.selected {
				prefix = "▸ "
				style = m.selectedStyle
			}

			line := prefix + style.Render(item.Title)
			if item.Description != "" {
				line += "\n    " + m.descStyle.Render(item.Description)
			}
			parts = append(parts, line)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.boxStyle.Width(width - 4).Render(content)
}

// SetItems replaces the items.
func (m *ListModal) SetItems(items []ListItem) {
	m.items = items
	m.applyFilter()
}

// SelectedItem returns the currently selected item.
func (m *ListModal) SelectedItem() *ListItem {
	if m.selected >= 0 && m.selected < len(m.filtered) {
		return &m.filtered[m.selected]
	}
	return nil
}
