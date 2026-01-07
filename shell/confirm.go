package shell

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmOption represents an option in a confirm modal.
type ConfirmOption struct {
	Label string
	Value any
}

// ConfirmModal presents a message with selectable options.
type ConfirmModal struct {
	id       string
	title    string
	message  string
	options  []ConfirmOption
	selected int
	onResult func(value any)
	width    int
	height   int

	// Styles
	boxStyle      lipgloss.Style
	titleStyle    lipgloss.Style
	messageStyle  lipgloss.Style
	optionStyle   lipgloss.Style
	selectedStyle lipgloss.Style
}

// ConfirmModalConfig configures a ConfirmModal.
type ConfirmModalConfig struct {
	ID       string
	Title    string
	Message  string
	Options  []ConfirmOption
	OnResult func(value any)
}

// NewConfirmModal creates a new confirm modal.
func NewConfirmModal(cfg ConfirmModalConfig) *ConfirmModal {
	return &ConfirmModal{
		id:       cfg.ID,
		title:    cfg.Title,
		message:  cfg.Message,
		options:  cfg.Options,
		selected: 0,
		onResult: cfg.OnResult,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")).
			Bold(true),
		messageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		optionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")).
			Bold(true),
	}
}

// NewYesNoModal creates a Yes/No confirm modal.
func NewYesNoModal(title, message string, onResult func(bool)) *ConfirmModal {
	return NewConfirmModal(ConfirmModalConfig{
		ID:      "yes-no",
		Title:   title,
		Message: message,
		Options: []ConfirmOption{
			{Label: "Yes", Value: true},
			{Label: "No", Value: false},
		},
		OnResult: func(v any) {
			if onResult != nil {
				onResult(v.(bool))
			}
		},
	})
}

// NewOKCancelModal creates an OK/Cancel confirm modal.
func NewOKCancelModal(title, message string, onResult func(bool)) *ConfirmModal {
	return NewConfirmModal(ConfirmModalConfig{
		ID:      "ok-cancel",
		Title:   title,
		Message: message,
		Options: []ConfirmOption{
			{Label: "OK", Value: true},
			{Label: "Cancel", Value: false},
		},
		OnResult: func(v any) {
			if onResult != nil {
				onResult(v.(bool))
			}
		},
	})
}

// ID implements Modal.
func (m *ConfirmModal) ID() string { return m.id }

// Title implements Modal.
func (m *ConfirmModal) Title() string { return m.title }

// Size implements Modal.
func (m *ConfirmModal) Size() Size { return SizeSmall }

// OnPush implements Modal.
func (m *ConfirmModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop implements Modal.
func (m *ConfirmModal) OnPop() {}

// HandleKey implements Modal.
func (m *ConfirmModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type {
	case tea.KeyUp, tea.KeyShiftTab:
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case tea.KeyDown, tea.KeyTab:
		if m.selected < len(m.options)-1 {
			m.selected++
		}
		return true, nil
	case tea.KeyEnter:
		if m.onResult != nil && m.selected >= 0 && m.selected < len(m.options) {
			m.onResult(m.options[m.selected].Value)
		}
		return true, func() tea.Msg { return PopMsg{} }
	}

	switch key.String() {
	case "k":
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case "j":
		if m.selected < len(m.options)-1 {
			m.selected++
		}
		return true, nil
	case "y":
		// Quick select yes
		for i, opt := range m.options {
			if opt.Value == true {
				m.selected = i
				if m.onResult != nil {
					m.onResult(true)
				}
				return true, func() tea.Msg { return PopMsg{} }
			}
		}
	case "n":
		// Quick select no
		for i, opt := range m.options {
			if opt.Value == false {
				m.selected = i
				if m.onResult != nil {
					m.onResult(false)
				}
				return true, func() tea.Msg { return PopMsg{} }
			}
		}
	}

	return false, nil
}

// Render implements Modal.
func (m *ConfirmModal) Render(width, height int) string {
	var parts []string

	// Title
	if m.title != "" {
		parts = append(parts, m.titleStyle.Render(m.title))
		parts = append(parts, "")
	}

	// Message
	if m.message != "" {
		parts = append(parts, m.messageStyle.Render(m.message))
		parts = append(parts, "")
	}

	// Options
	var optionParts []string
	for i, opt := range m.options {
		style := m.optionStyle
		prefix := "  "
		if i == m.selected {
			style = m.selectedStyle
			prefix = "▸ "
		}
		optionParts = append(optionParts, style.Render(prefix+opt.Label))
	}
	parts = append(parts, strings.Join(optionParts, "\n"))

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.boxStyle.Width(width - 4).Render(content)
}

// Selected returns the currently selected option index.
func (m *ConfirmModal) Selected() int {
	return m.selected
}

// SetSelected sets the selected option index.
func (m *ConfirmModal) SetSelected(index int) {
	if index >= 0 && index < len(m.options) {
		m.selected = index
	}
}
