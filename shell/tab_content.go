package shell

import (
	"github.com/2389-research/tux/content"
	tea "github.com/charmbracelet/bubbletea"
)

// TabContent extends content.Content with lifecycle hooks.
// Content implementations can optionally implement this interface
// to receive notifications when the tab becomes active or inactive.
type TabContent interface {
	content.Content

	// OnActivate is called when the tab becomes active.
	// Return a command to run on activation (e.g., start a timer).
	OnActivate() tea.Cmd

	// OnDeactivate is called when the tab becomes inactive.
	OnDeactivate()
}
