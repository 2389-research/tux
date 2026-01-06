package content

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Viewport is a scrollable text content area.
type Viewport struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
	style    lipgloss.Style
}

// NewViewport creates a new viewport.
func NewViewport() *Viewport {
	return &Viewport{
		viewport: viewport.New(0, 0),
		style:    lipgloss.NewStyle(),
	}
}

// Init implements Content.
func (v *Viewport) Init() tea.Cmd {
	return nil
}

// Update implements Content.
func (v *Viewport) Update(msg tea.Msg) (Content, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

// View implements Content.
func (v *Viewport) View() string {
	return v.viewport.View()
}

// Value implements Content. Returns nil as viewport doesn't produce values.
func (v *Viewport) Value() any {
	return nil
}

// SetSize implements Content.
func (v *Viewport) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.viewport.Width = width
	v.viewport.Height = height
}

// SetContent sets the viewport content.
func (v *Viewport) SetContent(content string) {
	v.content = content
	v.viewport.SetContent(content)
}

// AppendContent appends text to the viewport content.
func (v *Viewport) AppendContent(content string) {
	v.content += content
	v.viewport.SetContent(v.content)
}

// ScrollToTop scrolls to the top.
func (v *Viewport) ScrollToTop() {
	v.viewport.GotoTop()
}

// ScrollToBottom scrolls to the bottom.
func (v *Viewport) ScrollToBottom() {
	v.viewport.GotoBottom()
}

// ScrollUp scrolls up by n lines.
func (v *Viewport) ScrollUp(n int) {
	v.viewport.LineUp(n)
}

// ScrollDown scrolls down by n lines.
func (v *Viewport) ScrollDown(n int) {
	v.viewport.LineDown(n)
}

// AtTop returns true if scrolled to the top.
func (v *Viewport) AtTop() bool {
	return v.viewport.AtTop()
}

// AtBottom returns true if scrolled to the bottom.
func (v *Viewport) AtBottom() bool {
	return v.viewport.AtBottom()
}

// ScrollPercent returns the scroll position as a percentage.
func (v *Viewport) ScrollPercent() float64 {
	return v.viewport.ScrollPercent()
}

// LineCount returns the number of lines in the content.
func (v *Viewport) LineCount() int {
	return strings.Count(v.content, "\n") + 1
}

// Content returns the current content.
func (v *Viewport) Content() string {
	return v.content
}
