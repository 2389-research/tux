# tux Specification

> Language-agnostic specification for multi-agent terminal interfaces

This spec defines the interfaces, behaviors, and configuration format for tux implementations across languages.

## Documents

1. **[Architecture](./architecture.md)** - Core concepts, components, data flow
2. **[Shell](./shell.md)** - Top-level container spec
3. **[Modals](./modals.md)** - Modal system and built-in types
4. **[Content](./content.md)** - Content primitives
5. **[Theming](./theming.md)** - Theme interface and built-in themes
6. **[Configuration](./configuration.md)** - User config format (TOML)
7. **[Agent Backend](./agent-backend.md)** - Backend interface for agent integration
8. **[Accessibility](./accessibility.md)** - Keyboard nav, screen reader hooks

## Implementations

| Language | Package | Framework |
|----------|---------|-----------|
| Go | `github.com/2389-research/tux` | Bubble Tea |
| Rust | `github.com/2389-research/tux-rs` | ratatui |
| TypeScript | `@2389-research/tux` | ink |
| Python | `tux-ui` (PyPI) | textual |

## Versioning

Spec and implementations are versioned together. A `v1.2.0` implementation conforms to spec `v1.2.0`.

## Config Compatibility

All implementations read the same config format:
- Location: `~/.config/{appname}/ui.toml`
- Format: TOML (see [Configuration](./configuration.md))

A user's config works regardless of which language their app is built in.
