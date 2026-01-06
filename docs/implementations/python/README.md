# tux for Python

> Python implementation using Textual

## Installation

```bash
pip install tux-ui
# or
poetry add tux-ui
```

## Dependencies

- [Textual](https://github.com/Textualize/textual) - TUI framework
- [Rich](https://github.com/Textualize/rich) - Terminal rendering
- [toml](https://github.com/uiri/toml) - Config parsing

## Quick Start

```python
from tux import Shell, Tab, Viewport, load_config

config = load_config("myapp")

class MyApp(Shell):
    def compose(self):
        yield Tab("chat", "Chat", Viewport())

    def on_mount(self):
        self.backend = MyBackend()

if __name__ == "__main__":
    app = MyApp(theme=config.theme.name)
    app.run()
```

## Package Structure

```
tux/
├── shell/          # Shell, TabBar, StatusBar, Input
├── modal/          # ModalManager, built-in modals
├── content/        # Content widgets (Viewport, SelectList, etc.)
├── theme/          # Theme class, built-in themes
├── config/         # Config loading
└── agent/          # AgentBackend protocol
```

## Implementing AgentBackend

```python
from typing import Protocol, AsyncIterator

class AgentBackend(Protocol):
    async def stream(self, messages: list[Message]) -> AsyncIterator[AgentEvent]:
        ...

    def respond_to_approval(self, request_id: str, decision: ApprovalDecision) -> None:
        ...

    def describe_tool(self, name: str) -> ToolDescription:
        ...
```

## Custom Content

Subclass `Content`:

```python
from tux.content import Content

class EmailList(Content):
    def __init__(self):
        super().__init__()
        self.emails = []
        self.selected = 0

    def on_key(self, event):
        if event.key == "down":
            self.selected += 1
        elif event.key == "up":
            self.selected -= 1
        self.refresh()

    def render(self):
        # Return Rich renderables
        ...
```

## Why Textual?

- Built on Rich (excellent terminal rendering)
- CSS-like styling
- Async-first design
- Active development
- Good for rapid prototyping

## Status

🚧 **Planned** - Not yet implemented

Tracking issue: #{issue_number}
