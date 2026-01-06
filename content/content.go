// Package content provides composable content primitives for tabs and modals.
package content

import tea "github.com/charmbracelet/bubbletea"

// Content is the interface for content primitives that can be used
// in Tabs or Modals.
type Content interface {
	// Init initializes the content, returning any initial command.
	Init() tea.Cmd

	// Update handles messages and returns updated content and command.
	Update(msg tea.Msg) (Content, tea.Cmd)

	// View renders the content as a string.
	View() string

	// Value returns the current value for forms/wizards.
	// Returns nil if the content doesn't produce a value.
	Value() any

	// SetSize updates the available width and height.
	SetSize(width, height int)
}
