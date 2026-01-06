# tux for TypeScript

> TypeScript implementation using ink

## Installation

```bash
npm install @2389-research/tux
# or
yarn add @2389-research/tux
```

## Dependencies

- [ink](https://github.com/vadimdemedes/ink) - React for CLI
- [ink-text-input](https://github.com/vadimdemedes/ink-text-input) - Text input
- [@toml-tools/parser](https://github.com/nicolo-ribaudo/toml-tools) - Config parsing

## Quick Start

```tsx
import { Shell, Tab, Viewport, loadConfig } from '@2389-research/tux';

const config = await loadConfig('myapp');

const App = () => (
  <Shell
    theme={config.theme.name}
    backend={myBackend}
  >
    <Tab id="chat" label="Chat">
      <Viewport />
    </Tab>
  </Shell>
);

render(<App />);
```

## Package Structure

```
@2389-research/tux/
├── src/
│   ├── shell/      # Shell, TabBar, StatusBar, Input components
│   ├── modal/      # ModalManager, built-in modals
│   ├── content/    # Content components (Viewport, SelectList, etc.)
│   ├── theme/      # Theme context, built-in themes
│   ├── config/     # Config loading
│   └── agent/      # AgentBackend interface
```

## Implementing AgentBackend

```typescript
interface AgentBackend {
  stream(messages: Message[]): AsyncIterable<AgentEvent>;
  respondToApproval(requestId: string, decision: ApprovalDecision): void;
  describeTool(name: string): ToolDescription;
}
```

## Custom Content

Create a React component:

```tsx
interface ContentProps {
  width: number;
  height: number;
}

const EmailList: React.FC<ContentProps> = ({ width, height }) => {
  const [selected, setSelected] = useState(0);

  useInput((input, key) => {
    if (key.downArrow) setSelected(s => s + 1);
    if (key.upArrow) setSelected(s => s - 1);
  });

  return (
    <Box flexDirection="column">
      {emails.map((email, i) => (
        <EmailRow
          key={email.id}
          email={email}
          selected={i === selected}
        />
      ))}
    </Box>
  );
};
```

## Status

🚧 **Planned** - Not yet implemented

Tracking issue: #{issue_number}
