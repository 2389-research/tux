# tux

> Multi-language library for building multi-agent terminal interfaces

## What is tux?

tux provides the UI primitives for building terminal-based agent interfaces. It handles:

- **Shell** - Tabs, input, status bar
- **Modals** - Wizards, approvals, forms
- **Theming** - Consistent styling with user overrides
- **Configuration** - Two-tier config (app defaults + user preferences)

You bring the agent backend—tux handles the terminal UI.

```
mux (backend) → tux (frontend)
```

## Implementations

| Language | Package | Status | Framework |
|----------|---------|--------|-----------|
| **Go** | `github.com/2389-research/tux` | 🚧 In Progress | Bubble Tea |
| **Rust** | `github.com/2389-research/tux-rs` | 📋 Planned | ratatui |
| **TypeScript** | `@2389-research/tux` | 📋 Planned | ink |
| **Python** | `pip install tux-ui` | 📋 Planned | Textual |

## Documentation

- **[Specification](./docs/spec/)** - Language-agnostic interface definitions
  - [Architecture](./docs/spec/architecture.md)
  - [Configuration](./docs/spec/configuration.md)
- **[Implementations](./docs/implementations/)**
  - [Go](./docs/implementations/go/)
  - [Rust](./docs/implementations/rust/)
  - [TypeScript](./docs/implementations/typescript/)
  - [Python](./docs/implementations/python/)

## Quick Example (Go)

```go
package main

import (
    "github.com/2389-research/tux/shell"
    "github.com/2389-research/tux/config"
)

func main() {
    cfg, _ := config.Load("myapp", defaults)

    s := shell.New(shell.Config{
        Theme:   cfg.BuildTheme(),
        Backend: myMuxBackend,
    })

    s.AddTab("chat", "Chat", chatContent)
    s.Run()
}
```

## User Configuration

Users customize via `~/.config/{appname}/ui.toml`:

```toml
[theme]
name = "dracula"

[theme.colors]
primary = "#ff79c6"

[keybindings]
help = ["f1", "?"]

[input]
prefix = "λ "
```

Same config format works across all language implementations.

## Origin

Extracted from [hex](https://github.com/2389-research/hex) (coding agent) and [jeff](https://github.com/2389-research/jeff) (personal agent) to provide a shared foundation for multi-agent terminal interfaces.

## License

MIT
