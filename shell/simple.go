package shell

import (
	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SimpleModal is a basic modal with content and footer.
type SimpleModal struct {
	id      string
	title   string
	content content.Content
	footer  string
	size    Size
	onClose func()
	width   int
	height  int

	// Styles
	boxStyle    lipgloss.Style
	titleStyle  lipgloss.Style
	footerStyle lipgloss.Style
}

// SimpleModalConfig configures a SimpleModal.
type SimpleModalConfig struct {
	ID      string
	Title   string
	Content content.Content
	Footer  string
	Size    Size
	OnClose func()
}

// NewSimpleModal creates a new simple modal.
func NewSimpleModal(cfg SimpleModalConfig) *SimpleModal {
	if cfg.Size == 0 {
		cfg.Size = SizeMedium
	}
	return &SimpleModal{
		id:      cfg.ID,
		title:   cfg.Title,
		content: cfg.Content,
		footer:  cfg.Footer,
		size:    cfg.Size,
		onClose: cfg.OnClose,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")).
			Bold(true),
		footerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// ID implements Modal.
func (m *SimpleModal) ID() string { return m.id }

// Title implements Modal.
func (m *SimpleModal) Title() string { return m.title }

// Size implements Modal.
func (m *SimpleModal) Size() Size { return m.size }

// OnPush implements Modal.
func (m *SimpleModal) OnPush(width, height int) {
	m.width = width
	m.height = height
	if m.content != nil {
		contentWidth := int(float64(width)*m.size.WidthPercent()) - 6
		contentHeight := int(float64(height)*m.size.HeightPercent()) - 6
		m.content.SetSize(contentWidth, contentHeight)
	}
}

// OnPop implements Modal.
func (m *SimpleModal) OnPop() {
	if m.onClose != nil {
		m.onClose()
	}
}

// HandleKey implements Modal.
func (m *SimpleModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	if m.content != nil {
		m.content.Update(key)
		return true, nil
	}
	return false, nil
}

// Render implements Modal.
func (m *SimpleModal) Render(width, height int) string {
	var parts []string

	// Title
	if m.title != "" {
		parts = append(parts, m.titleStyle.Render(m.title))
		parts = append(parts, "")
	}

	// Content
	if m.content != nil {
		parts = append(parts, m.content.View())
	}

	// Footer
	if m.footer != "" {
		parts = append(parts, "")
		parts = append(parts, m.footerStyle.Render(m.footer))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.boxStyle.Width(width - 4).Render(content)
}

// SetContent sets the modal content.
func (m *SimpleModal) SetContent(c content.Content) {
	m.content = c
}

// Content returns the modal content.
func (m *SimpleModal) Content() content.Content {
	return m.content
}
