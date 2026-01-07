// modal/go
package shell

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/2389-research/tux/theme"
)

// HelpModalConfig configures a HelpModal.
type HelpModalConfig struct {
	ID    string
	Title string // Default: "Help"
	Help  *Help
	Size  Size        // Default: SizeMedium
	Theme theme.Theme
	Mode  string // Current mode for filtering
}

// HelpModal adapts a Help component to work as a Modal.
type HelpModal struct {
	id     string
	title  string
	help   *Help
	size   Size
	theme  theme.Theme
	mode   string
	width  int
	height int
}

// NewHelpModal creates a new help modal.
func NewHelpModal(cfg HelpModalConfig) *HelpModal {
	id := cfg.ID
	if id == "" {
		id = "help-modal"
	}

	title := cfg.Title
	if title == "" {
		title = "Help"
	}

	size := cfg.Size
	if size == 0 {
		size = SizeMedium
	}

	th := cfg.Theme
	if th == nil {
		th = theme.Get("dracula")
	}

	h := cfg.Help
	if h == nil {
		h = NewHelp()
	}
	h.WithTheme(th)

	return &HelpModal{
		id:    id,
		title: title,
		help:  h,
		size:  size,
		theme: th,
		mode:  cfg.Mode,
	}
}

// ID implements Modal.
func (m *HelpModal) ID() string {
	return m.id
}

// Title implements Modal.
func (m *HelpModal) Title() string {
	return m.title
}

// Size implements Modal.
func (m *HelpModal) Size() Size {
	return m.size
}

// OnPush implements Modal.
func (m *HelpModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop implements Modal.
func (m *HelpModal) OnPop() {}

// HandleKey implements Modal.
// Escape and ? both close the modal (toggle behavior).
func (m *HelpModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	// Check for Escape
	if key.Type == tea.KeyEscape {
		return true, func() tea.Msg { return PopMsg{} }
	}

	// Check for ? (toggle close)
	if key.Type == tea.KeyRunes && len(key.Runes) == 1 && key.Runes[0] == '?' {
		return true, func() tea.Msg { return PopMsg{} }
	}

	return false, nil
}

// Render implements Modal.
func (m *HelpModal) Render(width, height int) string {
	return m.help.Render(width, m.mode)
}
