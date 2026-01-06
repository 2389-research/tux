package config

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/lipgloss"
)

// BuildTheme creates a Theme from the configuration.
// It loads the base theme by name and applies any color overrides.
func (c *Config) BuildTheme() theme.Theme {
	base := theme.Get(c.Theme.Name)

	// If no color overrides, return base theme directly
	if c.Theme.Colors == (ColorsConfig{}) {
		return base
	}

	// Create custom theme with overrides
	return &customTheme{
		base:    base,
		colors:  c.Theme.Colors,
	}
}

// customTheme wraps a base theme and applies color overrides.
type customTheme struct {
	base   theme.Theme
	colors ColorsConfig
	styles *theme.Styles
}

func (t *customTheme) Name() string { return t.base.Name() + "-custom" }

func (t *customTheme) Background() lipgloss.Color {
	if t.colors.Background != "" {
		return lipgloss.Color(t.colors.Background)
	}
	return t.base.Background()
}

func (t *customTheme) Foreground() lipgloss.Color {
	if t.colors.Foreground != "" {
		return lipgloss.Color(t.colors.Foreground)
	}
	return t.base.Foreground()
}

func (t *customTheme) Primary() lipgloss.Color {
	if t.colors.Primary != "" {
		return lipgloss.Color(t.colors.Primary)
	}
	return t.base.Primary()
}

func (t *customTheme) Secondary() lipgloss.Color {
	if t.colors.Secondary != "" {
		return lipgloss.Color(t.colors.Secondary)
	}
	return t.base.Secondary()
}

func (t *customTheme) Success() lipgloss.Color {
	if t.colors.Success != "" {
		return lipgloss.Color(t.colors.Success)
	}
	return t.base.Success()
}

func (t *customTheme) Warning() lipgloss.Color {
	if t.colors.Warning != "" {
		return lipgloss.Color(t.colors.Warning)
	}
	return t.base.Warning()
}

func (t *customTheme) Error() lipgloss.Color {
	if t.colors.Error != "" {
		return lipgloss.Color(t.colors.Error)
	}
	return t.base.Error()
}

func (t *customTheme) Info() lipgloss.Color {
	if t.colors.Info != "" {
		return lipgloss.Color(t.colors.Info)
	}
	return t.base.Info()
}

func (t *customTheme) Border() lipgloss.Color {
	if t.colors.Border != "" {
		return lipgloss.Color(t.colors.Border)
	}
	return t.base.Border()
}

func (t *customTheme) BorderFocused() lipgloss.Color {
	if t.colors.BorderFocused != "" {
		return lipgloss.Color(t.colors.BorderFocused)
	}
	return t.base.BorderFocused()
}

func (t *customTheme) Muted() lipgloss.Color {
	if t.colors.Muted != "" {
		return lipgloss.Color(t.colors.Muted)
	}
	return t.base.Muted()
}

func (t *customTheme) UserColor() lipgloss.Color {
	if t.colors.User != "" {
		return lipgloss.Color(t.colors.User)
	}
	return t.base.UserColor()
}

func (t *customTheme) AssistantColor() lipgloss.Color {
	if t.colors.Assistant != "" {
		return lipgloss.Color(t.colors.Assistant)
	}
	return t.base.AssistantColor()
}

func (t *customTheme) ToolColor() lipgloss.Color {
	if t.colors.Tool != "" {
		return lipgloss.Color(t.colors.Tool)
	}
	return t.base.ToolColor()
}

func (t *customTheme) SystemColor() lipgloss.Color {
	if t.colors.System != "" {
		return lipgloss.Color(t.colors.System)
	}
	return t.base.SystemColor()
}

func (t *customTheme) Styles() theme.Styles {
	// Return base styles - could be enhanced to rebuild with custom colors
	return t.base.Styles()
}
