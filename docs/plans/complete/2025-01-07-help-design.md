# Help Package Design

> Static help overlay for tux with categories, mode filtering, and modal integration.

## Overview

Help provides a reference card overlay showing keyboard shortcuts and commands. Apps define bindings with optional mode tags; the help component filters and renders them. It's a cheat sheet, not an interactive menu.

## Architecture

```
┌─────────────────────────────────────────┐
│ Modal Stack (existing)                   │
│   └── Modal interface                    │
│         └── HelpModal (adapter)          │
│               └── Help                   │
│                     └── Category[]       │
│                           └── Binding[]  │
└─────────────────────────────────────────┘
```

**Package structure:**
```
tux/
├── help/           # NEW
│   ├── help.go     # Help component, Render
│   ├── binding.go  # Binding, Category types
│   └── help_test.go
├── modal/
│   └── help.go     # HelpModal adapter (new file)
```

## Core Types

### Binding

```go
type Binding struct {
    Key         string   // Display text: "ctrl+c", "?", "/feedback"
    Description string   // What it does: "Quit", "Toggle help"
    Modes       []string // Optional: ["chat", "history"], empty = all modes
}
```

### Category

```go
type Category struct {
    Title    string    // "Navigation", "Actions", "Commands"
    Bindings []Binding
}
```

### Help Component

```go
type Help struct {
    categories []Category
    theme      theme.Theme
}

func New(categories ...Category) *Help

func (h *Help) WithTheme(th theme.Theme) *Help

// Render returns the help overlay.
// If mode is empty, shows all bindings.
// If mode is set, shows only bindings that match (or have no mode restriction).
func (h *Help) Render(width int, mode string) string
```

## Usage Examples

### Simple app (no modes)

```go
h := help.New(
    help.Category{
        Title: "General",
        Bindings: []help.Binding{
            {Key: "ctrl+c", Description: "Quit"},
            {Key: "?", Description: "Toggle help"},
            {Key: "esc", Description: "Close/cancel"},
        },
    },
    help.Category{
        Title: "Commands",
        Bindings: []help.Binding{
            {Key: "/feedback", Description: "Send feedback"},
            {Key: "/clear", Description: "Clear screen"},
        },
    },
)

output := h.Render(80, "") // empty mode = show all
```

### Mode-aware app

```go
h := help.New(
    help.Category{
        Title: "Navigation",
        Bindings: []help.Binding{
            {Key: "↑↓", Description: "Navigate", Modes: []string{"list", "history"}},
            {Key: "enter", Description: "Select", Modes: []string{"list"}},
            {Key: "enter", Description: "Send message", Modes: []string{"chat"}},
        },
    },
    help.Category{
        Title: "Actions",
        Bindings: []help.Binding{
            {Key: "ctrl+c", Description: "Quit"}, // no modes = always shown
            {Key: "?", Description: "Toggle help"},
        },
    },
)

// In chat mode: shows "Send message" for enter, hides "Navigate"
output := h.Render(80, "chat")

// In list mode: shows "Navigate" and "Select"
output := h.Render(80, "list")
```

## Modal Integration

HelpModal adapts Help to work in the modal stack:

```go
type HelpModalConfig struct {
    ID         string
    Title      string       // Default: "Help"
    Help       *help.Help
    Size       Size         // Default: SizeMedium
    Theme      theme.Theme
    Mode       string       // Current mode for filtering
}

func NewHelpModal(cfg HelpModalConfig) *HelpModal

// Usage
modal := modal.NewHelpModal(modal.HelpModalConfig{
    Help: myHelp,
    Mode: "chat",
})
manager.Push(modal)
```

HelpModal:
- Renders the help overlay in a modal frame
- Closes on `esc` or `?` (toggle)
- Read-only, no form submission

## Rendering

### Layout

```
┌─ Help ─────────────────────────────┐
│                                    │
│  Navigation                        │
│  ↑↓        Navigate                │
│  enter     Send message            │
│                                    │
│  Actions                           │
│  ctrl+c    Quit                    │
│  ?         Toggle help             │
│                                    │
│  Commands                          │
│  /feedback Send feedback           │
│  /clear    Clear screen            │
│                                    │
│            Press ? or esc to close │
└────────────────────────────────────┘
```

### Styling

Uses theme colors:
- **Title:** Primary, bold (category headers)
- **Key:** Secondary, bold (keybinding)
- **Description:** Foreground (what it does)
- **Footer hint:** Muted

## Mode Filtering Logic

```go
func (h *Help) Render(width int, mode string) string {
    for _, cat := range h.categories {
        var visible []Binding
        for _, b := range cat.Bindings {
            if len(b.Modes) == 0 || contains(b.Modes, mode) {
                visible = append(visible, b)
            }
        }
        if len(visible) > 0 {
            // render category with visible bindings
        }
    }
}
```

- Empty `mode` parameter: show all bindings
- Binding with empty `Modes`: always shown
- Binding with `Modes` set: shown only if current mode matches

## Testing Strategy

1. **Binding tests** - mode matching logic
2. **Category tests** - filtering, empty categories hidden
3. **Render tests** - output contains expected content
4. **Modal tests** - open/close behavior

Coverage target: 95%+

## Comparison with References

| Feature | Jeff | Hex | tux |
|---------|------|-----|-----|
| Categories | No | Yes | Yes |
| Mode filtering | No | Yes | Yes |
| Compact view | No | Yes | No |
| Search | No | No | No |
| Modal integration | No | No | Yes |
| Caching | Yes | No | Optional |

## Implementation Tasks

1. **Binding and Category types** - `help/binding.go`
2. **Help component** - `help/help.go` with Render
3. **Mode filtering** - filter bindings by mode
4. **HelpModal adapter** - `modal/help.go`
5. **Tests** - unit tests for all components
