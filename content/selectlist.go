package content

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectItem represents a single item in a SelectList.
type SelectItem struct {
	Label       string
	Description string
	Value       any
}

// SelectList is a single-selection list content primitive.
type SelectList struct {
	items    []SelectItem
	selected int
	width    int
	height   int

	// Styles
	cursorStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	descStyle       lipgloss.Style
}

// NewSelectList creates a new SelectList with the given items.
func NewSelectList(items []SelectItem) *SelectList {
	return &SelectList{
		items:    items,
		selected: 0,
		cursorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")).
			Bold(true),
		unselectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// Init implements Content.
func (s *SelectList) Init() tea.Cmd {
	return nil
}

// Update implements Content.
func (s *SelectList) Update(msg tea.Msg) (Content, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftTab:
			if s.selected > 0 {
				s.selected--
			}
		case tea.KeyDown, tea.KeyTab:
			if s.selected < len(s.items)-1 {
				s.selected++
			}
		}
		switch msg.String() {
		case "k":
			if s.selected > 0 {
				s.selected--
			}
		case "j":
			if s.selected < len(s.items)-1 {
				s.selected++
			}
		case "g":
			s.selected = 0
		case "G":
			s.selected = len(s.items) - 1
		}
	}
	return s, nil
}

// View implements Content.
func (s *SelectList) View() string {
	if len(s.items) == 0 {
		return s.descStyle.Render("No items")
	}

	var b strings.Builder

	for i, item := range s.items {
		cursor := "  "
		style := s.unselectedStyle

		if i == s.selected {
			cursor = s.cursorStyle.Render("▸ ")
			style = s.selectedStyle
		}

		b.WriteString(cursor)
		b.WriteString(style.Render(item.Label))

		if item.Description != "" {
			b.WriteString("\n    ")
			b.WriteString(s.descStyle.Render(item.Description))
		}

		if i < len(s.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Value implements Content. Returns the value of the selected item.
func (s *SelectList) Value() any {
	if s.selected >= 0 && s.selected < len(s.items) {
		return s.items[s.selected].Value
	}
	return nil
}

// SetSize implements Content.
func (s *SelectList) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Selected returns the index of the currently selected item.
func (s *SelectList) Selected() int {
	return s.selected
}

// SetSelected sets the selected index.
func (s *SelectList) SetSelected(index int) {
	if index >= 0 && index < len(s.items) {
		s.selected = index
	}
}

// SelectedItem returns the currently selected item.
func (s *SelectList) SelectedItem() *SelectItem {
	if s.selected >= 0 && s.selected < len(s.items) {
		return &s.items[s.selected]
	}
	return nil
}

// SetItems replaces all items in the list.
func (s *SelectList) SetItems(items []SelectItem) {
	s.items = items
	if s.selected >= len(items) {
		s.selected = len(items) - 1
	}
	if s.selected < 0 {
		s.selected = 0
	}
}

// Items returns all items in the list.
func (s *SelectList) Items() []SelectItem {
	return s.items
}
