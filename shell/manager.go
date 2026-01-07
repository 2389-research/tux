package shell

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Manager manages a stack of modals with version tracking.
type Manager struct {
	stack   []Modal
	version int
	width   int
	height  int

	// Styles
	backdropStyle lipgloss.Style
}

// NewManager creates a new modal manager.
func NewManager() *Manager {
	return &Manager{
		stack:   make([]Modal, 0),
		version: 0,
		backdropStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#44475a")),
	}
}

// Push adds a modal to the top of the stack.
func (m *Manager) Push(modal Modal) {
	modal.OnPush(m.width, m.height)
	m.stack = append(m.stack, modal)
	m.version++
}

// Pop removes and returns the top modal from the stack.
// Returns nil if the stack is empty.
func (m *Manager) Pop() Modal {
	if len(m.stack) == 0 {
		return nil
	}
	modal := m.stack[len(m.stack)-1]
	modal.OnPop()
	m.stack = m.stack[:len(m.stack)-1]
	m.version++
	return modal
}

// Peek returns the top modal without removing it.
// Returns nil if the stack is empty.
func (m *Manager) Peek() Modal {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

// HasActive returns true if there's at least one modal on the stack.
func (m *Manager) HasActive() bool {
	return len(m.stack) > 0
}

// Version returns the current version number.
// Increments on every Push/Pop for change detection.
func (m *Manager) Version() int {
	return m.version
}

// Clear removes all modals from the stack.
func (m *Manager) Clear() {
	for len(m.stack) > 0 {
		m.Pop()
	}
}

// Count returns the number of modals on the stack.
func (m *Manager) Count() int {
	return len(m.stack)
}

// SetSize updates the available dimensions for modals.
func (m *Manager) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// HandleKey routes a key event to the active modal.
// Returns true if the key was handled.
func (m *Manager) HandleKey(key tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	if !m.HasActive() {
		return false, nil
	}
	return m.Peek().HandleKey(key)
}

// Render renders the active modal (if any) with backdrop.
func (m *Manager) Render(width, height int) string {
	if !m.HasActive() {
		return ""
	}

	modal := m.Peek()
	size := modal.Size()

	// Calculate modal dimensions
	modalWidth := int(float64(width) * size.WidthPercent())
	modalHeight := int(float64(height) * size.HeightPercent())

	// Render modal content
	content := modal.Render(modalWidth, modalHeight)

	// Center the modal
	return m.centerContent(content, width, height)
}

// centerContent centers content within the given dimensions.
func (m *Manager) centerContent(content string, width, height int) string {
	lines := strings.Split(content, "\n")
	contentHeight := len(lines)
	contentWidth := 0
	for _, line := range lines {
		w := lipgloss.Width(line)
		if w > contentWidth {
			contentWidth = w
		}
	}

	// Calculate padding
	topPad := (height - contentHeight) / 2
	leftPad := (width - contentWidth) / 2

	if topPad < 0 {
		topPad = 0
	}
	if leftPad < 0 {
		leftPad = 0
	}

	// Build centered output
	var b strings.Builder

	// Top padding
	for i := 0; i < topPad; i++ {
		b.WriteString(strings.Repeat(" ", width))
		b.WriteString("\n")
	}

	// Content with left padding
	for i, line := range lines {
		b.WriteString(strings.Repeat(" ", leftPad))
		b.WriteString(line)
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// PopMsg is a tea.Msg that requests popping the current modal.
type PopMsg struct{}

// PushMsg is a tea.Msg that requests pushing a modal.
type PushMsg struct {
	Modal Modal
}
