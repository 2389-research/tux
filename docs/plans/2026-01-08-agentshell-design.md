# AgentShell Design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make `shell.New(agent)` produce a fully-wired agent TUI with sensible defaults, so spinning up a new agent app is just: define tools, create mux agent, pass to shell, done.

**Architecture:** Shell accepts a mux agent directly. Default tabs (Chat, Tools) and approval flow are auto-wired. Apps can customize via functional options.

**Tech Stack:** Go, Bubble Tea, tux shell components, mux orchestrator

---

## Background

### Problem

Today, creating a new agent TUI (like hex or jeff) requires significant boilerplate:
- Manual tab setup (~5 lines per tab)
- Manual event routing (subscribe to agent, update tabs)
- Manual approval wiring (bridge blocking ApprovalFunc to event-driven TUI)
- ~25+ lines of config/theme/shell setup, identical across apps

### Solution

`shell.New(agent)` gives you a working agent TUI. Tools are already in the agent's registry. Events auto-route to tabs. Approval modal auto-wires. Customize only what's different.

---

## API

### Minimal Usage

```go
// App defines tools
registry := tool.NewRegistry()
registry.Register(myTool1, myTool2, myTool3)

// App creates mux agent
agent := agent.New(agent.Config{
    Name:         "myagent",
    Registry:     registry,
    LLMClient:    llmClient,
    SystemPrompt: "You are a helpful assistant.",
})

// Shell does the rest
s := shell.New(agent)
s.Run()
```

### With Customizations

```go
s := shell.New(agent,
    // UI customizations
    shell.WithTheme(myTheme),
    shell.WithTab(shell.TabDef{
        ID:       "diff",
        Label:    "Diff Viewer",
        Shortcut: "ctrl+d",
        Content:  diffViewerContent,
    }),

    // Disable a default tab
    shell.WithoutTab("tools"),

    // Custom tool classifier
    shell.WithClassifier(myClassifier),

    // Additional shortcuts
    shell.WithShortcut("ctrl+d", showDiffAction),
)
s.Run()
```

---

## What You Get By Default

### Tabs

| Tab | Shortcut | Purpose |
|-----|----------|---------|
| Chat | Alt+1 | Streaming conversation, markdown rendering |
| Tools | Alt+2, Ctrl+O | Tool calls/results timeline |

### Modals

| Modal | Trigger | Purpose |
|-------|---------|---------|
| Tool Approval | Auto (when tool needs approval) | Show tool, params, risk; get user decision |
| Help | Ctrl+H | Keyboard shortcuts, available commands |
| Onboarding | First run | Initial setup wizard |

### Chrome

- **Input area** - User prompt entry
- **Status bar** - Agent state, streaming indicator
- **Keyboard shortcuts** - Alt+1-9 (tabs), Ctrl+O (tools), Ctrl+H (help), Escape (cancel/close)
- **Mouse support** - Click tabs, buttons, scroll

---

## Mux Integration

### Interface Requirements

Shell needs mux agent to provide:

```go
// mux's *agent.Agent already satisfies this
type Agent interface {
    Run(ctx context.Context, prompt string) error
    Subscribe() <-chan orchestrator.Event
    Config() agent.Config  // for ApprovalFunc, Registry access
}
```

### Event Routing

Shell subscribes to agent events and routes them:

| Event Type | Routed To |
|------------|-----------|
| `EventText` | Chat tab (streaming text) |
| `EventToolCall` | Tools tab (show tool starting) |
| `EventToolResult` | Tools tab (show result) |
| `EventComplete` | Status bar, re-enable input |
| `EventError` | Status bar, show error, re-enable input |
| `EventStateChange` | Status bar (streaming/executing/etc) |

### Approval Flow

Shell provides ApprovalFunc to mux that bridges blocking call to event-driven TUI:

```go
// Shell creates this internally
approvalFunc := func(ctx context.Context, t tool.Tool, params map[string]any) (bool, error) {
    responseChan := make(chan bool)

    // Send approval request to TUI (via tea.Cmd/message)
    sendToTUI(ApprovalRequestMsg{
        Tool:     t,
        Params:   params,
        Response: responseChan,
    })

    // Block until user decides
    select {
    case approved := <-responseChan:
        return approved, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// Shell injects this into agent config before Run()
```

### Cancellation

User hits Escape mid-stream:

```go
// Shell manages context
ctx, cancel := context.WithCancel(context.Background())

// Start agent in goroutine
go agent.Run(ctx, prompt)

// On Escape (when streaming):
cancel()
```

---

## Customization Options

### Adding Custom Tabs

```go
shell.WithTab(shell.TabDef{
    ID:       "diff",
    Label:    "Diff Viewer",
    Shortcut: "ctrl+d",      // optional
    Hidden:   false,          // optional
    Content:  diffContent,    // implements content.Content
})
```

### Removing Default Tabs

```go
shell.WithoutTab("tools")  // remove tools tab
```

### Custom Tool Classifier

```go
// Override which tools need approval
shell.WithClassifier(func(t tool.ToolInfo) (ToolAction, string) {
    if t.Name == "read_file" {
        return ActionAutoApprove, "read-only"
    }
    return ActionNeedsApproval, ""
})
```

### Theme

```go
shell.WithTheme(theme.NewNeoTerminalTheme())
```

### Additional Shortcuts

```go
shell.WithShortcut("ctrl+d", func() tea.Cmd {
    return showDiffModal()
})
```

---

## Implementation Layers

```
┌─────────────────────────────────────────┐
│         Your App (hex, jeff)             │
│   - defines tools                        │
│   - creates mux agent                    │
│   - calls shell.New(agent)               │
├─────────────────────────────────────────┤
│        shell.New(agent) [NEW]            │
│   - creates default tabs (Chat, Tools)   │
│   - wires event routing                  │
│   - wires approval flow                  │
│   - applies options/customizations       │
├─────────────────────────────────────────┤
│       tux/shell (existing)               │
│   - tabs, modals, input, status          │
│   - keyboard shortcuts, mouse            │
│   - themes, configuration                │
├─────────────────────────────────────────┤
│          mux (orchestrator)              │
│   - agent loop, tool execution           │
│   - event streaming, approval hooks      │
└─────────────────────────────────────────┘
```

---

## What Apps Provide

| Component | Required | Notes |
|-----------|----------|-------|
| Tools | Yes | Registered in mux Registry |
| LLM Client | Yes | Passed to mux agent.Config |
| System Prompt | Yes | Passed to mux agent.Config |
| Custom Tabs | No | Via shell.WithTab() |
| Custom Theme | No | Via shell.WithTheme() |
| Custom Classifier | No | Via shell.WithClassifier() |

---

## Open Questions

1. **History tab** - Is this a default tab, or customization? Current thinking: not default, since Chat shows conversation and Tools shows activity.

2. **Message persistence** - Where do conversations save? App layer, mux layer, or shell layer? Current thinking: app layer provides this if needed.

3. **Multi-conversation** - Current design assumes single conversation. Supporting multiple conversations (load/save/switch) is future work.

---

## Success Criteria

1. New agent app is <20 lines to get working TUI
2. Hex can migrate to shell.New(agent) with minimal changes
3. Jeff can migrate similarly
4. Custom tabs/modals work without forking shell
5. No mux changes required (current interface sufficient)
