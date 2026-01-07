// shell/modal_form.go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/tux/theme"
)

// FormModalConfig configures a FormModal.
type FormModalConfig struct {
	ID       string
	Title    string
	Form     *Form
	Size     Size
	Theme    theme.Theme
	OnSubmit func(Values)
	OnCancel func()
}

// FormModal adapts a Form to work as a Modal.
type FormModal struct {
	id       string
	title    string
	form     *Form
	size     Size
	theme    theme.Theme
	onSubmit func(Values)
	onCancel func()
	width    int
	height   int

	boxStyle   lipgloss.Style
	titleStyle lipgloss.Style
}

// NewFormModal creates a new form modal.
func NewFormModal(cfg FormModalConfig) *FormModal {
	size := cfg.Size
	if size == 0 {
		size = SizeMedium
	}

	th := cfg.Theme
	if th == nil {
		th = theme.Get("dracula")
	}

	m := &FormModal{
		id:       cfg.ID,
		title:    cfg.Title,
		form:     cfg.Form,
		size:     size,
		theme:    th,
		onSubmit: cfg.OnSubmit,
		onCancel: cfg.OnCancel,
	}

	// Set up styles
	m.boxStyle = th.Styles().ModalBox
	m.titleStyle = th.Styles().ModalTitle

	// Wire form callbacks
	if m.form != nil {
		m.form.OnSubmit(func(v Values) {
			if m.onSubmit != nil {
				m.onSubmit(v)
			}
		})
		m.form.OnCancel(func() {
			if m.onCancel != nil {
				m.onCancel()
			}
		})
		m.form.WithTheme(th)
	}

	return m
}

// ID implements Modal.
func (m *FormModal) ID() string {
	if m.id != "" {
		return m.id
	}
	return "form-modal"
}

// Title implements Modal.
func (m *FormModal) Title() string {
	return m.title
}

// Size implements Modal.
func (m *FormModal) Size() Size {
	return m.size
}

// OnPush implements Modal.
func (m *FormModal) OnPush(width, height int) {
	m.width = width
	m.height = height
	if m.form != nil {
		m.form.Init()
	}
}

// OnPop implements Modal.
func (m *FormModal) OnPop() {}

// HandleKey implements Modal.
func (m *FormModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	if m.form == nil {
		return false, nil
	}

	handled := m.form.HandleKey(key)

	// Check if form completed
	switch m.form.State() {
	case StateSubmitted:
		return true, func() tea.Msg { return PopMsg{} }
	case StateCancelled:
		return true, func() tea.Msg { return PopMsg{} }
	}

	return handled, nil
}

// Render implements Modal.
func (m *FormModal) Render(width, height int) string {
	var content string

	// Title
	if m.title != "" {
		content = m.titleStyle.Render(m.title) + "\n\n"
	}

	// Form content
	if m.form != nil {
		content += m.form.Render(width-6, height-6)
	}

	return m.boxStyle.Width(width - 4).Render(content)
}
