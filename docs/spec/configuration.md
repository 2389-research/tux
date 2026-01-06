# Configuration Specification

> User-configurable settings via TOML

## Overview

tux uses a two-tier configuration system:

1. **App defaults** - Set by the application developer
2. **User overrides** - Personal customizations via config file

User settings always win when specified.

## File Locations

Config files are checked in order (first found wins):

1. `~/.config/{appname}/ui.toml` (XDG standard)
2. `~/.{appname}rc` (legacy/simple)
3. `${APPNAME}_UI_CONFIG` environment variable (explicit override)

Examples:
- hex: `~/.config/hex/ui.toml` or `~/.hexrc` or `$HEX_UI_CONFIG`
- jeff: `~/.config/jeff/ui.toml` or `~/.jeffrc` or `$JEFF_UI_CONFIG`

## Config Format

TOML format. All sections are optional—only specify what you want to change.

### Complete Reference

```toml
# {appname} UI Configuration
# Only include sections you want to customize

# =============================================================================
# THEME
# =============================================================================

[theme]
# Base theme to extend: "dracula", "nord", "gruvbox", "high-contrast"
name = "dracula"

# Override individual colors (hex format: #RRGGBB)
[theme.colors]
primary = "#bd93f9"       # Main accent (buttons, active elements)
secondary = "#8be9fd"     # Secondary accent
background = "#282a36"    # Main background
foreground = "#f8f8f2"    # Main text color
success = "#50fa7b"       # Success states
warning = "#ffb86c"       # Warning states
error = "#ff5555"         # Error states
info = "#8be9fd"          # Info states
border = "#6272a4"        # Border color
border_focused = "#bd93f9" # Focused border color
muted = "#6272a4"         # Dimmed/subtle text

# Role-specific colors (for message display)
user = "#ffb86c"          # User messages
assistant = "#50fa7b"     # Assistant messages
tool = "#8be9fd"          # Tool-related messages
system = "#6272a4"        # System messages

# Override specific component styles
[theme.styles.title]
foreground = "#bd93f9"
bold = true
italic = false

[theme.styles.input]
border = "rounded"        # "rounded", "square", "thick", "double", "none"
border_foreground = "#6272a4"
padding_horizontal = 1
padding_vertical = 0

# =============================================================================
# MOUSE
# =============================================================================

[mouse]
enabled = true            # Enable mouse support
scroll_lines = 3          # Lines per scroll wheel tick
hover_enabled = true      # Enable hover detection (timestamps, etc.)
shift_passthrough = true  # Pass through Shift+click for text selection

# =============================================================================
# KEY BINDINGS
# =============================================================================

[keybindings]
# Each action can have multiple keys (array)
# Key format: "ctrl+x", "alt+x", "shift+x", "f1", "enter", "esc", "tab", etc.

submit = ["enter"]
cancel = ["esc"]
help = ["ctrl+h", "?", "f1"]
quick_actions = ["ctrl+k"]
next_tab = ["ctrl+tab", "ctrl+n"]
prev_tab = ["ctrl+shift+tab", "ctrl+p"]
scroll_up = ["ctrl+u", "pgup"]
scroll_down = ["ctrl+d", "pgdn"]
scroll_top = ["g g", "home"]      # Vim-style double key
scroll_bottom = ["G", "end"]

# App-specific custom bindings
[keybindings.custom]
# Apps can define their own actions
# run_tests = ["ctrl+t"]
# toggle_sidebar = ["ctrl+b"]

# =============================================================================
# STATUS BAR
# =============================================================================

[statusbar]
# Section display order (omit to hide a section)
order = ["model", "status", "tokens", "mode", "progress", "hints"]

# Section-specific settings
[statusbar.sections.model]
max_width = 20

[statusbar.sections.tokens]
format = "{used}/{total}"   # or "{used}" or "{percent}%"

# Custom sections (app fills content via hooks)
[statusbar.custom.git]
position = 3                # Insert at position in order
priority = 50               # Space allocation priority

# =============================================================================
# TAB BAR
# =============================================================================

[tabbar]
position = "top"            # "top" or "bottom"
style = "underline"         # "underline", "boxed", "pills"
show_badges = true          # Show notification badges
show_close = true           # Show close button on tabs
max_visible = 8             # Max tabs before overflow menu

# =============================================================================
# INPUT AREA
# =============================================================================

[input]
prefix = "> "               # Prompt prefix
placeholder = ""            # Placeholder text when empty
multiline = false           # Allow multi-line input
max_height = 5              # Max lines when multiline
show_char_count = false     # Show character count
max_chars = 0               # Max characters (0 = unlimited)

# =============================================================================
# MODALS
# =============================================================================

[modal]
backdrop = true             # Dim background behind modals
backdrop_opacity = 0.5      # Backdrop dimming (0.0 - 1.0)
animation = "none"          # "none", "fade", "slide"
close_on_esc = true         # Esc closes modal
close_on_click_outside = false  # Click outside closes modal

# =============================================================================
# AUTOCOMPLETE
# =============================================================================

[autocomplete]
enabled = true
max_suggestions = 10        # Max items in dropdown
min_chars = 1               # Min chars before showing suggestions
delay_ms = 50               # Debounce delay

# =============================================================================
# ACCESSIBILITY
# =============================================================================

[accessibility]
high_contrast = false       # Force high contrast colors
reduce_motion = false       # Disable animations
screen_reader_hints = true  # Include screen reader announcements
```

## Validation Rules

Implementations MUST validate configs and report errors clearly:

| Field | Rule |
|-------|------|
| `theme.name` | Must be a registered theme name |
| `theme.colors.*` | Must be valid hex color (#RGB or #RRGGBB) |
| `keybindings.*` | Must be valid key binding strings |
| `statusbar.order` | Must contain only valid section names |
| `tabbar.position` | Must be "top" or "bottom" |
| `tabbar.style` | Must be "underline", "boxed", or "pills" |
| `modal.animation` | Must be "none", "fade", or "slide" |

## CLI Commands

Implementations SHOULD provide these commands:

```bash
# Generate default config with documentation
{appname} ui init [--output PATH]

# Validate current config
{appname} ui validate

# Show active config file location
{appname} ui path

# Show effective config (merged defaults + user)
{appname} ui show
```

## Hot Reload

Implementations MAY support hot reload:
- Watch config file for changes
- Re-apply config without restart
- Validate before applying (reject invalid configs)

## Errors

When config is invalid, implementations SHOULD:
1. Log a warning with specific error details
2. Fall back to app defaults
3. Continue running (don't crash)

Example error output:
```
Warning: Invalid UI config at ~/.config/hex/ui.toml
  theme.colors.primary: "#gggggg" is not a valid hex color
  keybindings.help: "ctrl++" is not a valid key binding
Using default configuration.
```
