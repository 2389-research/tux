// Package shell provides the top-level container for tux applications.
package shell

import (
	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Shell is the top-level container that manages tabs, modals, input, and status.
type Shell struct {
	// Components
	tabs         *TabBar
	input        *Input
	statusBar    *StatusBar
	modalManager *Manager

	// State
	width   int
	height  int
	focused FocusTarget
	ready   bool

	// Configuration
	theme  theme.Theme
	config Config
}

// Config holds shell configuration options.
type Config struct {
	// ShowTabBar controls whether the tab bar is visible.
	ShowTabBar bool
	// ShowStatusBar controls whether the status bar is visible.
	ShowStatusBar bool
	// ShowInput controls whether the input area is visible.
	ShowInput bool
	// InputPrefix is the prefix shown before user input.
	InputPrefix string
	// InputPlaceholder is shown when input is empty.
	InputPlaceholder string
}

// DefaultConfig returns the default shell configuration.
func DefaultConfig() Config {
	return Config{
		ShowTabBar:       true,
		ShowStatusBar:    true,
		ShowInput:        true,
		InputPrefix:      "> ",
		InputPlaceholder: "",
	}
}

// FocusTarget represents which component has focus.
type FocusTarget int

const (
	// FocusInput means the input area has focus (default).
	FocusInput FocusTarget = iota
	// FocusTab means the active tab content has focus.
	FocusTab
	// FocusModal means a modal has focus.
	FocusModal
)

// New creates a new Shell with the given theme and config.
func New(th theme.Theme, cfg Config) *Shell {
	if th == nil {
		th = theme.NewDraculaTheme()
	}

	s := &Shell{
		theme:        th,
		config:       cfg,
		tabs:         NewTabBar(th),
		input:        NewInput(th, cfg.InputPrefix, cfg.InputPlaceholder),
		statusBar:    NewStatusBar(th),
		modalManager: NewManager(),
		focused:      FocusInput,
	}

	return s
}

// Init implements tea.Model.
func (s *Shell) Init() tea.Cmd {
	return tea.Batch(
		s.input.Init(),
	)
}

// Update implements tea.Model.
func (s *Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.modalManager.SetSize(msg.Width, msg.Height)
		s.ready = true
		s.updateSizes()

	case tea.KeyMsg:
		// Modal captures all input when active
		if s.modalManager.HasActive() {
			handled, cmd := s.modalManager.HandleKey(msg)
			if handled {
				return s, cmd
			}
			// Esc closes modal
			if msg.Type == tea.KeyEsc {
				s.modalManager.Pop()
				return s, nil
			}
		}

		// Global keys
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return s, tea.Quit
		}

		// Route to focused component
		switch s.focused {
		case FocusInput:
			var cmd tea.Cmd
			s.input, cmd = s.input.Update(msg)
			cmds = append(cmds, cmd)
		case FocusTab:
			cmd := s.tabs.HandleKey(msg)
			cmds = append(cmds, cmd)
		}

	case PopMsg:
		s.modalManager.Pop()

	case PushMsg:
		s.modalManager.Push(msg.Modal)
	}

	return s, tea.Batch(cmds...)
}

// View implements tea.Model.
func (s *Shell) View() string {
	if !s.ready {
		return "Loading..."
	}

	var sections []string

	// Tab bar
	if s.config.ShowTabBar {
		sections = append(sections, s.tabs.View())
	}

	// Content area
	contentHeight := s.contentHeight()
	content := s.tabs.RenderActiveContent(s.width, contentHeight)
	sections = append(sections, content)

	// Input
	if s.config.ShowInput {
		sections = append(sections, s.input.View())
	}

	// Status bar
	if s.config.ShowStatusBar {
		sections = append(sections, s.statusBar.View(s.width))
	}

	// Join sections
	output := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Overlay modal if active
	if s.modalManager.HasActive() {
		modalView := s.modalManager.Render(s.width, s.height)
		output = s.overlayModal(output, modalView)
	}

	return output
}

// contentHeight calculates available height for tab content.
func (s *Shell) contentHeight() int {
	h := s.height
	if s.config.ShowTabBar {
		h -= 1
	}
	if s.config.ShowInput {
		h -= 3 // Input with border
	}
	if s.config.ShowStatusBar {
		h -= 1
	}
	if h < 1 {
		h = 1
	}
	return h
}

// updateSizes updates component sizes after window resize.
func (s *Shell) updateSizes() {
	s.input.SetWidth(s.width)
	s.tabs.SetSize(s.width, s.contentHeight())
}

// overlayModal composites the modal view over the base view.
func (s *Shell) overlayModal(base, modal string) string {
	// For now, just return the modal centered
	// A more sophisticated implementation would blend them
	return modal
}

// AddTab adds a tab to the shell.
func (s *Shell) AddTab(tab Tab) {
	s.tabs.AddTab(tab)
}

// RemoveTab removes a tab by ID.
func (s *Shell) RemoveTab(id string) {
	s.tabs.RemoveTab(id)
}

// SetActiveTab switches to the tab with the given ID.
func (s *Shell) SetActiveTab(id string) {
	s.tabs.SetActive(id)
}

// PushModal pushes a modal onto the stack.
func (s *Shell) PushModal(m Modal) {
	s.modalManager.Push(m)
	s.focused = FocusModal
}

// PopModal pops the top modal from the stack.
func (s *Shell) PopModal() Modal {
	m := s.modalManager.Pop()
	if !s.modalManager.HasActive() {
		s.focused = FocusInput
	}
	return m
}

// HasModal returns true if there's an active 
func (s *Shell) HasModal() bool {
	return s.modalManager.HasActive()
}

// Focus sets the focus target.
func (s *Shell) Focus(target FocusTarget) {
	s.focused = target
}

// InputValue returns the current input text.
func (s *Shell) InputValue() string {
	return s.input.Value()
}

// SetInputValue sets the input text.
func (s *Shell) SetInputValue(value string) {
	s.input.SetValue(value)
}

// ClearInput clears the input text.
func (s *Shell) ClearInput() {
	s.input.SetValue("")
}

// SetStatus updates the status bar.
func (s *Shell) SetStatus(status Status) {
	s.statusBar.SetStatus(status)
}

// Theme returns the current theme.
func (s *Shell) Theme() theme.Theme {
	return s.theme
}

// Run starts the shell as a Bubble Tea program.
func (s *Shell) Run() error {
	p := tea.NewProgram(s, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
