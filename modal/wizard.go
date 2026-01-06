package modal

import (
	"fmt"
	"strings"

	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WizardStep represents a single step in a wizard.
type WizardStep struct {
	ID          string
	Title       string
	Description string
	Content     content.Content
	Validate    func(value any) error
	Optional    bool
}

// WizardModal guides users through multi-step processes.
type WizardModal struct {
	id         string
	title      string
	steps      []WizardStep
	current    int
	results    map[string]any
	onComplete func(results map[string]any)
	onCancel   func()
	width      int
	height     int
	err        string

	// Styles
	boxStyle      lipgloss.Style
	titleStyle    lipgloss.Style
	stepStyle     lipgloss.Style
	progressStyle lipgloss.Style
	descStyle     lipgloss.Style
	errorStyle    lipgloss.Style
	hintStyle     lipgloss.Style
}

// WizardModalConfig configures a WizardModal.
type WizardModalConfig struct {
	ID         string
	Title      string
	Steps      []WizardStep
	OnComplete func(results map[string]any)
	OnCancel   func()
}

// NewWizardModal creates a new wizard modal.
func NewWizardModal(cfg WizardModalConfig) *WizardModal {
	return &WizardModal{
		id:         cfg.ID,
		title:      cfg.Title,
		steps:      cfg.Steps,
		current:    0,
		results:    make(map[string]any),
		onComplete: cfg.OnComplete,
		onCancel:   cfg.OnCancel,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")).
			Bold(true),
		stepStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")),
		progressStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")),
		hintStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// ID implements Modal.
func (m *WizardModal) ID() string { return m.id }

// Title implements Modal.
func (m *WizardModal) Title() string { return m.title }

// Size implements Modal.
func (m *WizardModal) Size() Size { return SizeLarge }

// OnPush implements Modal.
func (m *WizardModal) OnPush(width, height int) {
	m.width = width
	m.height = height
	m.updateContentSize()
}

// OnPop implements Modal.
func (m *WizardModal) OnPop() {
	if m.onCancel != nil {
		m.onCancel()
	}
}

// updateContentSize updates the current step's content size.
func (m *WizardModal) updateContentSize() {
	if step := m.CurrentStep(); step != nil && step.Content != nil {
		contentWidth := int(float64(m.width)*m.Size().WidthPercent()) - 8
		contentHeight := int(float64(m.height)*m.Size().HeightPercent()) - 12
		step.Content.SetSize(contentWidth, contentHeight)
	}
}

// HandleKey implements Modal.
func (m *WizardModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	// Pass to current content first
	if step := m.CurrentStep(); step != nil && step.Content != nil {
		step.Content.Update(key)
	}

	switch key.Type {
	case tea.KeyEnter:
		return true, m.Next()
	}

	switch key.String() {
	case "ctrl+n":
		return true, m.Next()
	case "ctrl+p":
		m.Previous()
		return true, nil
	}

	return true, nil // Wizard captures all input
}

// Next advances to the next step or completes the wizard.
func (m *WizardModal) Next() tea.Cmd {
	step := m.CurrentStep()
	if step == nil {
		return nil
	}

	// Get value from current step
	var value any
	if step.Content != nil {
		value = step.Content.Value()
	}

	// Validate if needed
	if step.Validate != nil {
		if err := step.Validate(value); err != nil {
			m.err = err.Error()
			return nil
		}
	}

	// Store result
	m.results[step.ID] = value
	m.err = ""

	// Advance or complete
	if m.current < len(m.steps)-1 {
		m.current++
		m.updateContentSize()
		return nil
	}

	// Complete
	if m.onComplete != nil {
		m.onComplete(m.results)
	}
	return func() tea.Msg { return PopMsg{} }
}

// Previous goes back to the previous step.
func (m *WizardModal) Previous() {
	if m.current > 0 {
		m.current--
		m.err = ""
		m.updateContentSize()
	}
}

// CanGoNext returns true if we can advance.
func (m *WizardModal) CanGoNext() bool {
	return m.current < len(m.steps)
}

// CanGoPrevious returns true if we can go back.
func (m *WizardModal) CanGoPrevious() bool {
	return m.current > 0
}

// CurrentStep returns the current step.
func (m *WizardModal) CurrentStep() *WizardStep {
	if m.current >= 0 && m.current < len(m.steps) {
		return &m.steps[m.current]
	}
	return nil
}

// Progress returns current and total step counts.
func (m *WizardModal) Progress() (current, total int) {
	return m.current + 1, len(m.steps)
}

// Render implements Modal.
func (m *WizardModal) Render(width, height int) string {
	var parts []string

	// Title
	if m.title != "" {
		parts = append(parts, m.titleStyle.Render(m.title))
	}

	// Progress
	current, total := m.Progress()
	progress := m.progressStyle.Render(fmt.Sprintf("Step %d of %d", current, total))
	parts = append(parts, progress)

	// Current step
	step := m.CurrentStep()
	if step != nil {
		parts = append(parts, "")
		parts = append(parts, m.stepStyle.Render(step.Title))

		if step.Description != "" {
			parts = append(parts, m.descStyle.Render(step.Description))
		}

		parts = append(parts, "")

		// Content
		if step.Content != nil {
			parts = append(parts, step.Content.View())
		}
	}

	// Error
	if m.err != "" {
		parts = append(parts, "")
		parts = append(parts, m.errorStyle.Render("Error: "+m.err))
	}

	// Hints
	parts = append(parts, "")
	var hints []string
	if m.CanGoPrevious() {
		hints = append(hints, "Ctrl+P: Back")
	}
	hints = append(hints, "Enter: Continue")
	hints = append(hints, "Esc: Cancel")
	parts = append(parts, m.hintStyle.Render(strings.Join(hints, " │ ")))

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.boxStyle.Width(width - 4).Render(content)
}

// Results returns the collected results.
func (m *WizardModal) Results() map[string]any {
	return m.results
}

// SetResults sets pre-filled results.
func (m *WizardModal) SetResults(results map[string]any) {
	m.results = results
}

// GoToStep jumps to a specific step by index.
func (m *WizardModal) GoToStep(index int) {
	if index >= 0 && index < len(m.steps) {
		m.current = index
		m.err = ""
		m.updateContentSize()
	}
}
