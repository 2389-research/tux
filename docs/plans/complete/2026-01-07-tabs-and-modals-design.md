# Tabs and Modals Design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Establish an opinionated default for view management in tux - tabs for persistent content views, modals for temporary overlays.

**Architecture:** Shell owns TabBar (persistent views) and ModalManager (temporary overlays). Apps use tabs to organize content and modals for transient UI. No separate ViewMode concept needed.

**Tech Stack:** Go, Bubble Tea, existing tux shell components

---

## Background

### Problem

Hex and Jeff both have ad-hoc view management:

- **Hex:** Has `ViewMode` enum (Intro/Chat/History/Tools) but History and Tools are actually implemented as fullscreen overlays (Ctrl+R, Ctrl+O). The enum is vestigial.
- **Jeff:** Has `ViewMode` enum but History/Tools are unimplemented placeholders. Has split-pane layout with focus cycling that's separate from view modes.

Both apps reinvent view state management. Tux should provide an opinionated default.

### Solution

Two concepts, clear separation:

| Concept | Purpose | Lifecycle | Examples |
|---------|---------|-----------|----------|
| **Tabs** | Persistent content views | Always available, switch between | Chat, History, Tools |
| **Modals** | Temporary overlays | Push/pop, dismiss to return | Help, Approval, Forms |

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│ [Chat] [History] [Tools]              Tab Bar   │
├─────────────────────────────────────────────────┤
│                                                 │
│              Active Tab Content                 │
│              (app-provided content.Content)     │
│                                                 │
├─────────────────────────────────────────────────┤
│ > input                                         │
├─────────────────────────────────────────────────┤
│ status bar                                      │
└─────────────────────────────────────────────────┘

         Modal pushed (captures input)
                    ↓
┌─────────────────────────────────────────────────┐
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │           Modal Content                   │  │
│  │           (Help, Approval, Form)          │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Tabs

Tabs are persistent content views. The tab bar shows available tabs; users switch between them.

**Characteristics:**
- Always exist once added
- Content preserved when switching away
- Visible in tab bar (unless hidden)
- Keyboard shortcuts for quick access

**API (existing, may need extension):**
```go
// Existing
shell.AddTab(tab Tab)
shell.RemoveTab(id string)
shell.SetActiveTab(id string)

// May need
shell.AddTab(tab Tab, opts ...TabOption)  // WithHidden(), WithShortcut()
```

### Modals

Modals are temporary overlays that capture input until dismissed.

**Characteristics:**
- Stack-based (can push multiple)
- Capture all input when active
- Dismissed via Esc or explicit action
- Return to previous state on dismiss

**API (existing):**
```go
shell.PushModal(modal Modal)
shell.PopModal()
shell.HasModal() bool
```

---

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Ctrl+Tab` | Next tab |
| `Ctrl+Shift+Tab` | Previous tab |
| `Ctrl+1/2/3...` | Jump to tab by index |
| `Esc` | Close modal (if open) |

Apps can add custom shortcuts for specific tabs:
- `Ctrl+R` → Switch to History tab
- `Ctrl+O` → Switch to Tools tab

---

## Migration: Hex

### Current State

| Component | Implementation |
|-----------|---------------|
| Intro | ViewMode, renders welcome screen |
| Chat | ViewMode, main conversation |
| History | Fullscreen overlay (Ctrl+R) |
| Tool Timeline | Fullscreen overlay (Ctrl+O) |
| Help | Fullscreen overlay (Ctrl+H) |
| Tool Approval | Non-fullscreen overlay |

### Target State

| Component | Implementation |
|-----------|---------------|
| Intro | Initial state of Chat tab (no messages) |
| Chat | Tab |
| History | Tab (Ctrl+R shortcut) |
| Tools | Tab (Ctrl+O shortcut) |
| Help | Modal (Ctrl+H) |
| Tool Approval | Modal |

### Changes Required

1. Remove `ViewMode` enum and `NextView()` cycling
2. Convert History overlay to History tab
3. Convert Tool Timeline overlay to Tools tab
4. Keep Help as modal (temporary reference)
5. Keep Tool Approval as modal (transient decision)
6. Tab key no longer cycles views (use Ctrl+Tab or shortcuts)

---

## Migration: Jeff

### Current State

| Component | Implementation |
|-----------|---------------|
| Chat | ViewMode + split-pane (Chat/Display) |
| History | ViewMode placeholder (not implemented) |
| Tools | ViewMode placeholder (not implemented) |
| Display Cards | Right pane in split layout |
| Help | Modal |
| Tool Approval | Modal (huh form) |
| Quick Actions | Modal (`:` menu) |

### Target State

| Component | Implementation |
|-----------|---------------|
| Chat | Tab (split-pane is internal rendering) |
| History | Tab (needs implementation) |
| Tools | Tab (needs implementation) |
| Display | Future: specialized tab |
| Help | Modal |
| Tool Approval | Modal |
| Quick Actions | Modal |

### Changes Required

1. Remove `ViewMode` enum
2. Use tux tabs (Chat, History, Tools)
3. Split-pane becomes internal layout of Chat tab content
4. Pane focus cycling becomes internal to Chat tab
5. Display cards → future work (specialized tab)

---

## Tux Implementation

### Phase 1: Tab Shortcuts

Add keyboard shortcuts for tab navigation:

```go
// shell/shell.go - Update() method

case tea.KeyMsg:
    switch {
    case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+tab"))):
        s.tabs.Next()
    case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+shift+tab"))):
        s.tabs.Prev()
    case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+1"))):
        s.tabs.SetActiveByIndex(0)
    // ... etc
    }
```

### Phase 2: Hidden Tabs

Allow tabs that are accessible via shortcut but not shown in tab bar:

```go
type TabOption func(*Tab)

func WithHidden() TabOption {
    return func(t *Tab) { t.hidden = true }
}

func WithShortcut(key string) TabOption {
    return func(t *Tab) { t.shortcut = key }
}

// Usage
shell.AddTab(historyTab, WithShortcut("ctrl+r"))
```

### Phase 3: Tab Content Interface

Tabs may need lifecycle hooks:

```go
type TabContent interface {
    content.Content

    // Called when tab becomes active
    OnActivate() tea.Cmd

    // Called when tab becomes inactive
    OnDeactivate()
}
```

---

## Decision Log

| Question | Decision | Rationale |
|----------|----------|-----------|
| ViewMode vs Tabs? | Tabs | Single concept, already in tux |
| History: tab or modal? | Tab | Persistent content you browse |
| Tool Timeline: tab or modal? | Tab | Persistent content you reference |
| Help: tab or modal? | Modal | Temporary reference, dismiss to return |
| Tool Approval: tab or modal? | Modal | Transient decision, must resolve |
| Intro screen? | Initial tab state | Not a separate view, just empty state |
| Jeff's Display pane? | Future tab | Separate problem, defer |

---

## Out of Scope

- Jeff's Display cards redesign (separate effort)
- Glamour/markdown rendering (app layer)
- Session/conversation management (app layer)
- Onboarding wizards (run outside shell)

---

## Success Criteria

1. Hex can remove `ViewMode` enum and use tux tabs
2. Jeff can remove `ViewMode` enum and use tux tabs
3. Keyboard shortcuts work consistently across apps
4. Clear mental model: tabs persist, modals dismiss
