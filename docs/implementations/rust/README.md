# tux for Rust

> Rust implementation using ratatui

## Installation

```toml
# Cargo.toml
[dependencies]
tux = "0.1"
```

## Dependencies

- [ratatui](https://github.com/ratatui-org/ratatui) - TUI framework
- [crossterm](https://github.com/crossterm-rs/crossterm) - Terminal backend
- [toml](https://github.com/toml-rs/toml) - Config parsing

## Quick Start

```rust
use tux::{Shell, Config, Theme};

fn main() -> Result<()> {
    // Load config
    let cfg = Config::load("myapp")?;

    // Create shell
    let mut shell = Shell::new(ShellConfig {
        theme: Theme::get(&cfg.theme.name),
        backend: Box::new(MyBackend::new()),
    });

    // Add tabs
    shell.add_tab(Tab {
        id: "chat".into(),
        label: "Chat".into(),
        content: Box::new(Viewport::new()),
    });

    // Run
    shell.run()
}
```

## Crate Structure

```
tux/
├── src/
│   ├── shell/      # Shell, TabBar, StatusBar, Input
│   ├── modal/      # ModalManager, built-in modals
│   ├── content/    # Content traits and primitives
│   ├── theme/      # Theme trait, built-in themes
│   ├── config/     # Config loading
│   └── agent/      # AgentBackend trait
```

## Implementing AgentBackend

```rust
pub trait AgentBackend {
    fn stream(&self, messages: Vec<Message>) -> impl Stream<Item = AgentEvent>;
    fn respond_to_approval(&self, request_id: &str, decision: ApprovalDecision);
    fn describe_tool(&self, name: &str) -> ToolDescription;
}
```

## Custom Content

Implement the `Content` trait:

```rust
pub trait Content {
    fn update(&mut self, event: Event) -> Option<Command>;
    fn render(&self, area: Rect, buf: &mut Buffer);
    fn set_size(&mut self, width: u16, height: u16);
    fn value(&self) -> Option<Value>;
}
```

## Status

🚧 **Planned** - Not yet implemented

Tracking issue: #{issue_number}
