# Architecture

> Core concepts and component relationships for tux

## Overview

tux is a shared library for building multi-agent terminal interfaces. It provides:

- **Shell** - Top-level container with tabs, input, status bar
- **Modals** - Stacked overlays (wizards, approvals, forms)
- **Content** - Composable primitives (viewports, lists, progress)
- **Theming** - Consistent styling with user overrides
- **Agent Backend** - Interface for connecting to agent orchestrators

## Visual Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                           Shell                                  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ [Tab1]  [Tab2]  [Tab3]                          TabBar     │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                                                            │ │
│  │                   Active Tab Content                       │ │
│  │                   (Viewport, List, Timeline, etc.)         │ │
│  │                                                            │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ > input area                                    [Spinner]  │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ model │ ● status │ tokens │ [progress] │ hints   StatusBar │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌──────────────── ModalManager ────────────────┐               │
│  │  Modal Stack (version-tracked)               │               │
│  │  ┌─────────────────────────────────────────┐ │               │
│  │  │  Active Modal (captures all input)      │ │               │
│  │  └─────────────────────────────────────────┘ │               │
│  └──────────────────────────────────────────────┘               │
└─────────────────────────────────────────────────────────────────┘
         │
         │ Events
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AgentBackend                                │
│  (implements streaming, tool execution, approval requests)       │
└─────────────────────────────────────────────────────────────────┘
```

## Design Principles

### 1. Domain-Agnostic

Works for any multi-agent terminal app:
- Coding agents (like hex)
- Personal assistants (like jeff)
- DevOps tools
- Data analysis interfaces

### 2. Composable

Small primitives combine into complex UIs. Content can be placed in tabs or modals interchangeably.

### 3. Pluggable Backends

Agent implementations are swappable. The Shell doesn't know about specific LLMs or tools—it receives events from an abstract backend.

### 4. User-Configurable

Two-tier configuration:
1. App developers set defaults
2. Users override via `~/.config/{appname}/ui.toml`

### 5. Consistent UX

Same interaction patterns across all apps using tux:
- Same keyboard shortcuts (unless user remaps)
- Same modal behaviors
- Same status bar layout

## Component Hierarchy

```
Shell
├── TabBar
│   └── Tab[]
│       └── Content (Viewport, List, Timeline, etc.)
├── InputArea
│   └── Autocomplete
├── StatusBar
│   └── Section[]
├── Spinner
└── ModalManager
    └── Modal[] (stack)
        └── Content
```

## Data Flow

### User Input Flow

```
KeyPress/MouseEvent
    │
    ▼
┌─────────────┐
│ModalManager │ ── Has active modal? ──▶ Route to Modal
└─────────────┘
    │ no
    ▼
┌─────────────┐
│   Shell     │ ── Check keybindings ──▶ Handle action
└─────────────┘
    │ unhandled
    ▼
┌─────────────┐
│ Active Tab  │ ── Content-specific ──▶ Scroll, select, etc.
└─────────────┘
    │ unhandled
    ▼
┌─────────────┐
│ InputArea   │ ── Text input
└─────────────┘
```

### Agent Event Flow

```
AgentBackend
    │
    │ Stream events
    ▼
┌─────────────────────────────────────────┐
│ EventText      → Append to chat content │
│ EventToolCall  → Show in activity tab   │
│ EventApproval  → Push ApprovalModal     │
│ EventComplete  → Update status bar      │
│ EventError     → Show error state       │
└─────────────────────────────────────────┘
```

## Focus Management

Focus priority (highest to lowest):
1. **Modal** - Active modal captures all input
2. **Tab Content** - When explicitly focused (scrolling, selecting)
3. **Input Area** - Default focus for typing

## Non-Goals

These are explicitly NOT part of tux:

1. **Tool definitions** - Domain-specific, stay in app
2. **Persistence** - Each app brings its own storage
3. **Business logic** - Stays in app code
4. **LLM integration** - Handled by backend (e.g., mux)
