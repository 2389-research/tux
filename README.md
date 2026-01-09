# tux

> Terminal UI for multi-agent applications

## What is tux?

tux provides a ready-to-use terminal interface for agent applications. Bring your agent backend, tux handles the UI:

- **Chat tab** - Streaming conversation display
- **Tools tab** - Tool call timeline with status
- **Theming** - Dracula, NeoTerminal, or custom themes
- **Tabs & Modals** - Extensible UI primitives

```
your agent → tux.App → terminal UI
```

## Quick Start

```go
package main

import "github.com/2389-research/tux"

func main() {
    app := tux.New(myAgent)
    app.Run()
}
```

That's it. tux wires up:
- Input submission → `agent.Run(ctx, prompt)`
- Agent events → Chat/Tools display
- Escape → Cancel running agent

## Agent Interface

Your agent implements three methods:

```go
type Agent interface {
    // Run executes with the user's prompt
    Run(ctx context.Context, prompt string) error

    // Subscribe returns events as they happen
    Subscribe() <-chan tux.Event

    // Cancel stops the current run
    Cancel()
}
```

Events flow from your agent to the UI:

```go
// Streaming text → Chat tab
Event{Type: EventText, Text: "Hello..."}

// Tool calls → Tools tab
Event{Type: EventToolCall, ToolID: "1", ToolName: "read_file", ToolParams: params}
Event{Type: EventToolResult, ToolID: "1", ToolOutput: "contents", Success: true}

// Completion → Finalize message
Event{Type: EventComplete}
```

## Customization

```go
app := tux.New(agent,
    tux.WithTheme(theme.NewNeoTerminalTheme()),
    tux.WithTab(tux.TabDef{ID: "logs", Label: "Logs", Content: logsContent}),
    tux.WithoutTab("tools"),  // Remove default tab
)
```

## Low-Level API

For full control, use the shell package directly:

```go
import "github.com/2389-research/tux/shell"

s := shell.New(theme, shell.DefaultConfig())
s.AddTab(shell.Tab{ID: "chat", Label: "Chat", Content: myContent})
s.Run()
```

The shell package provides tabs, modals, forms, status bar, and input handling as composable primitives.

## User Configuration

Users can customize via `~/.config/{appname}/config.toml`:

```toml
[theme]
name = "dracula"

[theme.colors]
primary = "#ff79c6"

[input]
prefix = "→ "
```

## Status

| Component | Status |
|-----------|--------|
| Agent shell (tux.App) | ✅ Working |
| Shell primitives | ✅ Working |
| Themes | ✅ Dracula, NeoTerminal |
| Forms & Modals | ✅ Working |
| Config loading | ✅ Working |

## License

MIT
