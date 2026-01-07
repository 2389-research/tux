// Package shell provides modal overlays, modal overlay management for tux.
package shell

import tea "github.com/charmbracelet/bubbletea"

// Modal is the interface for overlay modals that capture input focus.
type Modal interface {
	// ID returns a unique identifier for this modal.
	ID() string

	// Title returns the modal title for display.
	Title() string

	// Size returns the preferred size of this modal.
	Size() Size

	// Render renders the modal content at the given dimensions.
	Render(width, height int) string

	// OnPush is called when the modal is pushed onto the stack.
	OnPush(width, height int)

	// OnPop is called when the modal is popped from the stack.
	OnPop()

	// HandleKey processes key input. Returns true if the key was handled.
	HandleKey(key tea.KeyMsg) (handled bool, cmd tea.Cmd)
}

// Size represents the preferred size of a modal.
type Size int

const (
	// SizeSmall is approximately 30% of screen height.
	SizeSmall Size = iota
	// SizeMedium is approximately 50% of screen height.
	SizeMedium
	// SizeLarge is approximately 80% of screen height.
	SizeLarge
	// SizeFullscreen is 100% of screen.
	SizeFullscreen
)

// HeightPercent returns the percentage of screen height for this size.
func (s Size) HeightPercent() float64 {
	switch s {
	case SizeSmall:
		return 0.30
	case SizeMedium:
		return 0.50
	case SizeLarge:
		return 0.80
	case SizeFullscreen:
		return 1.0
	default:
		return 0.50
	}
}

// WidthPercent returns the percentage of screen width for this size.
func (s Size) WidthPercent() float64 {
	switch s {
	case SizeSmall:
		return 0.50
	case SizeMedium:
		return 0.60
	case SizeLarge:
		return 0.80
	case SizeFullscreen:
		return 1.0
	default:
		return 0.60
	}
}
