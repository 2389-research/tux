# tux for Go

> Go implementation using Bubble Tea

## Installation

```bash
go get github.com/2389-research/tux
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - Components (viewport, textinput, etc.)
- [Huh](https://github.com/charmbracelet/huh) - Forms
- [BurntSushi/toml](https://github.com/BurntSushi/toml) - Config parsing

## Quick Start

```go
package main

import (
    "github.com/2389-research/tux/shell"
    "github.com/2389-research/tux/config"
    "github.com/2389-research/tux/theme"
)

func main() {
    // Load config with app defaults
    cfg, _ := config.Load("myapp", config.Defaults{
        Theme: "dracula",
    })

    // Create shell
    s := shell.New(shell.Config{
        Theme:   theme.Get(cfg.Theme.Name),
        Backend: myBackend,  // Your AgentBackend implementation
    })

    // Add tabs
    s.AddTab(shell.Tab{
        ID:      "chat",
        Label:   "Chat",
        Content: content.NewViewport(),
    })

    // Run
    s.Run()
}
```

## Package Structure

```
github.com/2389-research/tux
├── shell/          # Shell, TabBar, StatusBar, Input
├── modal/          # ModalManager, built-in modals
├── content/        # Content primitives (Viewport, SelectList, etc.)
├── theme/          # Theme interface, built-in themes
├── config/         # Config loading, validation, hot reload
└── agent/          # AgentBackend interface, ToolQueue
```

## Implementing AgentBackend

```go
type AgentBackend interface {
    // Start streaming response
    Stream(ctx context.Context, messages []Message) (<-chan AgentEvent, error)

    // Handle tool approval response
    RespondToApproval(requestID string, decision ApprovalDecision) error

    // Get tool metadata for display
    DescribeTool(name string) ToolDescription
}

// Example: wrapping mux
type MuxBackend struct {
    agent *mux.Agent
}

func (b *MuxBackend) Stream(ctx context.Context, messages []Message) (<-chan AgentEvent, error) {
    // Convert messages to mux format
    // Subscribe to mux events
    // Transform to AgentEvent
}
```

## Custom Content

Implement the `Content` interface:

```go
type Content interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Content, tea.Cmd)
    View() string
    SetSize(width, height int)
    Value() any  // For forms/wizards
}

// Example: custom email list
type EmailList struct {
    emails   []Email
    selected int
    width    int
    height   int
}

func (e *EmailList) Update(msg tea.Msg) (Content, tea.Cmd) {
    // Handle keys, return self
}

func (e *EmailList) View() string {
    // Render email list
}
```

## See Also

- [Full API Documentation](./api.md)
- [Examples](./examples/)
- [Migration Guide](./migration.md) - Moving from hex/jeff internal UI
