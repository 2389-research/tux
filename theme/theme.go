package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color and style interface for tux components.
type Theme interface {
	Name() string

	// Base colors
	Background() lipgloss.Color
	Foreground() lipgloss.Color

	// Accent colors
	Primary() lipgloss.Color
	Secondary() lipgloss.Color

	// Status colors
	Success() lipgloss.Color
	Warning() lipgloss.Color
	Error() lipgloss.Color
	Info() lipgloss.Color

	// UI colors
	Border() lipgloss.Color
	BorderFocused() lipgloss.Color
	Muted() lipgloss.Color

	// Role colors (for messages)
	UserColor() lipgloss.Color
	AssistantColor() lipgloss.Color
	ToolColor() lipgloss.Color
	SystemColor() lipgloss.Color

	// Composed styles
	Styles() Styles
}

var themes = map[string]func() Theme{
	"dracula": NewDraculaTheme,
}

// Register adds a theme constructor to the registry.
func Register(name string, constructor func() Theme) {
	themes[name] = constructor
}

// Get returns a theme by name, defaulting to dracula if not found.
func Get(name string) Theme {
	if constructor, ok := themes[name]; ok {
		return constructor()
	}
	return NewDraculaTheme()
}

// Available returns the names of all registered themes.
func Available() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}
