# Near-Term Improvements Design

Error display, history navigation, and approval handling for tux.App.

## 1. Error Display

### Behavior
- Status bar shows `⚠` indicator + truncated error message (~10 chars)
- `ctrl+e` opens modal with full error list
- Errors accumulate until cleared
- Errors clear on next successful agent run (EventComplete with no preceding EventError)

### Implementation
- Add `errors []error` field to tux.App
- On `EventError`: append to errors, update status bar
- On `EventComplete`: if no errors in this run, clear errors
- Add `ErrorModal` that lists all errors with full text
- Shell handles `ctrl+e` to push ErrorModal

### Status Bar Format
```
claude-3.5 │ ● connected │ ⚠ "connection t..." │ 1k/100k
```

## 2. History Navigation

### Behavior
- Up arrow in input: previous user prompt
- Down arrow: next user prompt (or empty if at end)
- Session only (no persistence)

### Implementation
- Add `historyIndex int` to tux.App (or Input component)
- Source: `ChatContent.messages` filtered to role="user"
- Up: decrement index, set input to that message
- Down: increment index, set input (or clear if past end)
- On submit: reset index to end

### Why session-only
Chat messages are session-only today. When we add chat persistence later, history comes free.

## 3. Approval Handling

### Architecture Context
```
mux (agent loop) ←→ app (hex/jeff) ←→ tux (UI)
```

- mux orchestrates tool execution
- mux has `ApprovalFunc` callback that blocks until decision
- app wires mux to tux
- tux shows approval UI

### Event Flow
```
mux                         tux.App                     UI
 │                           │                           │
 │ ── EventApproval ───────► │                           │
 │    (with Response chan)   ├─ show ApprovalModal ────► │
 │                           │                           │ user decides
 │                           │ ◄── decision ─────────────┤
 │ ◄── decision on chan ─────┤                           │
 │                           │                           │
 │ ── EventToolResult ─────► │ (denied: Success=false)   │
```

### Event Type Addition
```go
type Event struct {
    // ... existing fields ...

    // For EventApproval
    Response chan ApprovalDecision
}

type ApprovalDecision int

const (
    DecisionApprove ApprovalDecision = iota
    DecisionDeny
    DecisionAlwaysAllow
    DecisionNeverAllow
)
```

### New Event Type
```go
const (
    // ... existing ...
    EventApproval EventType = "approval"
)
```

### tux.App.processEvent
```go
case EventApproval:
    // Show approval modal
    modal := shell.NewApprovalModal(shell.ApprovalModalConfig{
        Tool: shell.ToolInfo{
            ID:     event.ToolID,
            Name:   event.ToolName,
            Params: event.ToolParams,
        },
        OnDecision: func(decision shell.ApprovalDecision) {
            // Map shell decision to tux decision and send
            event.Response <- mapDecision(decision)
        },
    })
    a.shell.PushModal(modal)
```

### Denial UX
When user denies:
1. Modal closes
2. tux sends `DecisionDeny` on Response channel
3. mux sends `EventToolResult{Success: false, ToolOutput: "denied by user"}`
4. Tools tab shows ✗ like any failed tool

### Modal Options
Using existing `shell.ApprovalModal`:
- Approve (run this time)
- Deny (skip this time)
- Always Allow (never ask again)
- Never Allow (block permanently)

"Always/Never" decisions go to mux which handles rule persistence.

## Summary

| Feature | Complexity | Notes |
|---------|------------|-------|
| Error display | Low | New field, status bar indicator, modal |
| History nav | Low | Filter chat messages, index tracking |
| Approval | Medium | New event type, response channel, modal wiring |
