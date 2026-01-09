# Roadmap

## Current State (v0.1.0)

The Go implementation is functional:

- **tux.App** - High-level agent shell with Chat/Tools tabs, event routing
- **shell package** - Tabs, modals, forms, status bar, input
- **theme package** - Dracula and NeoTerminal themes
- **config package** - TOML config loading with theme building

## Near Term

### Error Display
Show agent errors in status bar or modal (currently a TODO in event routing).

### History Navigation
Up/down arrow to cycle through previous prompts.

### Approval Modal
Built-in modal for tool approval flows (the shell supports modals, but no approval-specific one yet).

## Future

### Rust Implementation
Port core tux to Rust using ratatui. Same config format, same UX patterns.

**Why Rust:**
- Performance-critical agent applications
- Native binaries without Go runtime
- Ecosystem compatibility (Rust AI/ML tools)

### Python Implementation
Port to Python using Textual. Lower priority than Rust.

**Why Python:**
- Rapid prototyping
- Data science / notebook integration
- Accessibility for non-systems programmers

### Cross-Language Config
All implementations will read the same `~/.config/{appname}/config.toml` format, so user preferences carry across apps regardless of implementation language.

## Non-Goals

- **TypeScript/ink** - Deprioritized; terminal UIs in JS have limited use cases
- **Tool definitions** - Stay in app layer, not tux
- **LLM integration** - Handled by agent backends (e.g., mux)
- **Persistence** - Each app brings its own storage
