// form/huh_theme.go
package form

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ToHuhTheme converts a tux theme to a huh theme.
func ToHuhTheme(th theme.Theme) *huh.Theme {
	t := huh.ThemeBase()

	// Focused state styles
	t.Focused.Base = t.Focused.Base.
		BorderForeground(th.Primary())

	t.Focused.Title = lipgloss.NewStyle().
		Foreground(th.Primary()).
		Bold(true)

	t.Focused.Description = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(th.Primary())

	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(th.Success()).
		Bold(true)

	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(th.Background()).
		Background(th.Primary()).
		Bold(true).
		Padding(0, 1)

	t.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(th.Foreground()).
		Background(th.Border()).
		Padding(0, 1)

	t.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(th.Primary())

	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.TextInput.Prompt = lipgloss.NewStyle().
		Foreground(th.Primary())

	// Blurred state - more muted
	t.Blurred.Base = t.Blurred.Base.
		BorderForeground(th.Border())

	t.Blurred.Title = lipgloss.NewStyle().
		Foreground(th.Foreground())

	t.Blurred.Description = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Blurred.SelectSelector = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Blurred.SelectedOption = lipgloss.NewStyle().
		Foreground(th.Foreground())

	t.Blurred.UnselectedOption = lipgloss.NewStyle().
		Foreground(th.Muted())

	// Error styling
	t.Focused.ErrorMessage = lipgloss.NewStyle().
		Foreground(th.Error())

	t.Focused.ErrorIndicator = lipgloss.NewStyle().
		Foreground(th.Error())

	return t
}
