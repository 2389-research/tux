package config

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds multiple validation errors.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed:\n  %s", strings.Join(e.Errors, "\n  "))
}

// Validate checks the config for errors.
// Returns a slice of error messages, empty if valid.
func (c *Config) Validate() []string {
	var errs []string

	errs = append(errs, c.validateTheme()...)
	errs = append(errs, c.validateKeybindings()...)
	errs = append(errs, c.validateTabBar()...)
	errs = append(errs, c.validateModal()...)

	return errs
}

// validThemes is the list of built-in theme names.
var validThemes = map[string]bool{
	"dracula":       true,
	"nord":          true,
	"gruvbox":       true,
	"high-contrast": true,
}

func (c *Config) validateTheme() []string {
	var errs []string

	// Validate theme name
	if c.Theme.Name != "" && !validThemes[c.Theme.Name] {
		errs = append(errs, fmt.Sprintf("theme.name: %q is not a valid theme (valid: dracula, nord, gruvbox, high-contrast)", c.Theme.Name))
	}

	// Validate colors
	colorFields := map[string]string{
		"theme.colors.primary":        c.Theme.Colors.Primary,
		"theme.colors.secondary":      c.Theme.Colors.Secondary,
		"theme.colors.background":     c.Theme.Colors.Background,
		"theme.colors.foreground":     c.Theme.Colors.Foreground,
		"theme.colors.success":        c.Theme.Colors.Success,
		"theme.colors.warning":        c.Theme.Colors.Warning,
		"theme.colors.error":          c.Theme.Colors.Error,
		"theme.colors.info":           c.Theme.Colors.Info,
		"theme.colors.border":         c.Theme.Colors.Border,
		"theme.colors.border_focused": c.Theme.Colors.BorderFocused,
		"theme.colors.muted":          c.Theme.Colors.Muted,
		"theme.colors.user":           c.Theme.Colors.User,
		"theme.colors.assistant":      c.Theme.Colors.Assistant,
		"theme.colors.tool":           c.Theme.Colors.Tool,
		"theme.colors.system":         c.Theme.Colors.System,
	}

	for field, value := range colorFields {
		if value != "" && !isValidHexColor(value) {
			errs = append(errs, fmt.Sprintf("%s: %q is not a valid hex color", field, value))
		}
	}

	return errs
}

// hexColorRegex matches #RGB or #RRGGBB format.
var hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

func isValidHexColor(s string) bool {
	return hexColorRegex.MatchString(s)
}

func (c *Config) validateKeybindings() []string {
	var errs []string

	bindings := map[string][]string{
		"keybindings.submit":        c.Keybindings.Submit,
		"keybindings.cancel":        c.Keybindings.Cancel,
		"keybindings.help":          c.Keybindings.Help,
		"keybindings.quick_actions": c.Keybindings.QuickActions,
		"keybindings.next_tab":      c.Keybindings.NextTab,
		"keybindings.prev_tab":      c.Keybindings.PrevTab,
		"keybindings.scroll_up":     c.Keybindings.ScrollUp,
		"keybindings.scroll_down":   c.Keybindings.ScrollDown,
		"keybindings.scroll_top":    c.Keybindings.ScrollTop,
		"keybindings.scroll_bottom": c.Keybindings.ScrollBottom,
	}

	for field, keys := range bindings {
		for _, key := range keys {
			if !isValidKeybinding(key) {
				errs = append(errs, fmt.Sprintf("%s: %q is not a valid key binding", field, key))
			}
		}
	}

	for name, keys := range c.Keybindings.Custom {
		for _, key := range keys {
			if !isValidKeybinding(key) {
				errs = append(errs, fmt.Sprintf("keybindings.custom.%s: %q is not a valid key binding", name, key))
			}
		}
	}

	return errs
}

// isValidKeybinding checks if a key binding string is valid.
func isValidKeybinding(s string) bool {
	if s == "" {
		return false
	}

	// Single characters
	if len(s) == 1 {
		return true
	}

	// Special keys
	specialKeys := map[string]bool{
		"enter": true, "return": true, "esc": true, "escape": true,
		"tab": true, "space": true, "backspace": true, "delete": true,
		"up": true, "down": true, "left": true, "right": true,
		"home": true, "end": true, "pgup": true, "pgdn": true,
		"f1": true, "f2": true, "f3": true, "f4": true, "f5": true,
		"f6": true, "f7": true, "f8": true, "f9": true, "f10": true,
		"f11": true, "f12": true,
	}

	lower := strings.ToLower(s)

	// Check if it's a special key
	if specialKeys[lower] {
		return true
	}

	// Check for modifier combinations: ctrl+x, alt+x, shift+x
	parts := strings.Split(lower, "+")
	if len(parts) >= 2 {
		// All parts except last should be modifiers
		modifiers := map[string]bool{"ctrl": true, "alt": true, "shift": true}
		for i := 0; i < len(parts)-1; i++ {
			if !modifiers[parts[i]] {
				return false
			}
		}
		// Last part should be a valid key
		lastPart := parts[len(parts)-1]
		if len(lastPart) == 1 || specialKeys[lastPart] {
			return true
		}
	}

	// Vim-style double keys like "g g" or "gg"
	if strings.Contains(s, " ") {
		parts := strings.Split(s, " ")
		if len(parts) == 2 && len(parts[0]) == 1 && len(parts[1]) == 1 {
			return true
		}
	}

	// Double character sequences like "gg"
	if len(s) == 2 && s[0] == s[1] {
		return true
	}

	return false
}

func (c *Config) validateTabBar() []string {
	var errs []string

	validPositions := map[string]bool{"top": true, "bottom": true}
	if c.TabBar.Position != "" && !validPositions[c.TabBar.Position] {
		errs = append(errs, fmt.Sprintf("tabbar.position: %q is not valid (must be \"top\" or \"bottom\")", c.TabBar.Position))
	}

	validStyles := map[string]bool{"underline": true, "boxed": true, "pills": true}
	if c.TabBar.Style != "" && !validStyles[c.TabBar.Style] {
		errs = append(errs, fmt.Sprintf("tabbar.style: %q is not valid (must be \"underline\", \"boxed\", or \"pills\")", c.TabBar.Style))
	}

	return errs
}

func (c *Config) validateModal() []string {
	var errs []string

	validAnimations := map[string]bool{"none": true, "fade": true, "slide": true}
	if c.Modal.Animation != "" && !validAnimations[c.Modal.Animation] {
		errs = append(errs, fmt.Sprintf("modal.animation: %q is not valid (must be \"none\", \"fade\", or \"slide\")", c.Modal.Animation))
	}

	if c.Modal.BackdropOpacity < 0 || c.Modal.BackdropOpacity > 1 {
		errs = append(errs, fmt.Sprintf("modal.backdrop_opacity: %v is not valid (must be between 0.0 and 1.0)", c.Modal.BackdropOpacity))
	}

	return errs
}
