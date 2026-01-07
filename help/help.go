package help

import (
	"strings"

	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

// Help displays categorized keyboard shortcuts and commands.
type Help struct {
	categories []Category
	theme      theme.Theme
}

// New creates a new Help component with the given categories.
func New(categories ...Category) *Help {
	return &Help{
		categories: categories,
	}
}

// WithTheme sets the theme and returns the Help for chaining.
func (h *Help) WithTheme(th theme.Theme) *Help {
	h.theme = th
	return h
}

// Render returns the help overlay as a styled string.
// If mode is empty, all bindings are shown.
// If mode is set, each category is filtered to show only matching bindings.
// Categories with no visible bindings are hidden.
func (h *Help) Render(width int, mode string) string {
	th := h.theme
	if th == nil {
		th = theme.Get("dracula")
	}

	var sections []string

	for _, cat := range h.categories {
		bindings := cat.FilterByMode(mode)
		if len(bindings) == 0 {
			continue
		}

		// Category title: Primary color, bold
		titleStyle := lipgloss.NewStyle().
			Foreground(th.Primary()).
			Bold(true)
		title := titleStyle.Render(cat.Title)

		// Bindings: Secondary (bold) for keys, Foreground for descriptions
		keyStyle := lipgloss.NewStyle().
			Foreground(th.Secondary()).
			Bold(true)
		descStyle := lipgloss.NewStyle().
			Foreground(th.Foreground())

		var bindingLines []string
		for _, b := range bindings {
			line := keyStyle.Render(b.Key) + "  " + descStyle.Render(b.Description)
			bindingLines = append(bindingLines, line)
		}

		section := title + "\n" + strings.Join(bindingLines, "\n")
		sections = append(sections, section)
	}

	// Footer: Muted color
	footerStyle := lipgloss.NewStyle().
		Foreground(th.Muted())
	footer := footerStyle.Render("Press ? or esc to close")

	// Combine all sections with spacing
	var content string
	if len(sections) > 0 {
		content = strings.Join(sections, "\n\n") + "\n\n" + footer
	} else {
		content = footer
	}

	// Wrap in border box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Border()).
		Padding(1, 2).
		Width(width - 4) // Account for border width

	return boxStyle.Render(content)
}
