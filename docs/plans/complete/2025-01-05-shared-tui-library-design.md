# Shared TUI Library Design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Extract a shared, domain-agnostic TUI library from hex and jeff that provides primitives for building multi-agent terminal interfaces.

**Architecture:** Shell-based container with pluggable Tabs, stacked Modals, and composable Content primitives. Agent backends implement a simple interface; the TUI handles all rendering, input, and state management. Follows the mux pattern for backend abstraction.

**Tech Stack:** Go, Bubble Tea (charmbracelet/bubbletea), Lip Gloss (charmbracelet/lipgloss), Huh (charmbracelet/huh)

---

## Table of Contents

1. [Design Goals](#design-goals)
2. [Architecture Overview](#architecture-overview)
3. [Shell](#shell)
4. [Tabs](#tabs)
5. [Modals](#modals)
6. [Content Primitives](#content-primitives)
7. [Agent Integration](#agent-integration)
8. [Theming](#theming)
9. [Implementation Plan](#implementation-plan)
10. [Open Questions (Resolved)](#open-questions-resolved)
    - [Huh Integration](#1-huh-integration)
    - [Autocomplete System](#2-autocomplete-system)
    - [Message Persistence](#3-message-persistence)
    - [Plugin/Extensibility System](#4-pluginextensibility-system)
    - [Accessibility](#5-accessibility)
    - [Mouse Support](#6-mouse-support)
    - [Theming (Extended)](#7-theming-extended)
    - [Customizable Status Bar](#8-customizable-status-bar)
    - [Other Customization](#9-other-customization)
    - [User Configuration System](#10-user-configuration-system)

---

## Design Goals

### Primary Goals

1. **Domain-agnostic** - Works for coding agents (hex), personal agents (jeff), or any future agent
2. **Composable** - Small primitives that combine into complex UIs
3. **Pluggable backends** - Agent implementations are swappable via interface
4. **Consistent UX** - Same interaction patterns across all consumers
5. **Hardened** - Built-in safety patterns from hex (approval queue, context management)

### Non-Goals

1. Tool definitions (domain-specific)
2. Persistence layer (each app brings its own)
3. Business logic (stays in hex/jeff)

### Design Principles

- **Mux as model** - Follow how hex/jeff integrate with mux for backend abstraction
- **Hex patterns preferred** - Use hex's hardened patterns (OverlayManager, ToolQueue, etc.)
- **Progressive adoption** - Consumers can adopt incrementally
- **No magic** - Explicit over implicit

---

## Architecture Overview

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

---

## Shell

The Shell is the top-level container that owns all UI state.

### Interface

```go
package tui

type Shell struct {
    // Core state
    tabs         []Tab
    activeTab    string
    modalManager *ModalManager
    input        *Input
    statusBar    *StatusBar
    spinner      *Spinner

    // Configuration
    theme        Theme
    renderer     MessageRenderer

    // Backend
    backend      AgentBackend
    events       chan ShellEvent

    // Bubble Tea
    width, height int
    focused       FocusTarget
}

type ShellConfig struct {
    Theme    Theme
    Backend  AgentBackend
    Renderer MessageRenderer  // Optional, defaults to MarkdownRenderer
}

func NewShell(cfg ShellConfig) *Shell

// Bubble Tea interface
func (s *Shell) Init() tea.Cmd
func (s *Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (s *Shell) View() string

// Tab management
func (s *Shell) AddTab(tab Tab)
func (s *Shell) RemoveTab(id string)
func (s *Shell) SetActiveTab(id string)
func (s *Shell) GetTab(id string) *Tab
func (s *Shell) UpdateTabBadge(id string, badge string)

// Modal management
func (s *Shell) PushModal(modal Modal)
func (s *Shell) PopModal() Modal
func (s *Shell) HasModal() bool

// Focus management
func (s *Shell) Focus(target FocusTarget)
func (s *Shell) FocusInput()
func (s *Shell) FocusActiveTab()

// Input control
func (s *Shell) SetInputEnabled(enabled bool)
func (s *Shell) SetInputPlaceholder(text string)
func (s *Shell) GetInputValue() string
func (s *Shell) ClearInput()

// Status
func (s *Shell) SetStatus(cfg StatusConfig)
func (s *Shell) SetSpinner(spinnerType SpinnerType, message string)
func (s *Shell) StopSpinner()

// Events (for async pushes from backend)
func (s *Shell) Events() chan<- ShellEvent
```

### Focus Management

```go
type FocusTarget int

const (
    FocusInput FocusTarget = iota
    FocusTab
    FocusModal
)
```

Focus priority:
1. **Modal** (if present) - captures all input
2. **Tab** (if focused) - for scrolling, selection
3. **Input** (default) - for typing

### Event Bus

```go
type ShellEvent interface {
    shellEvent()  // marker method
}

// Backend pushes content
type AddToTabEvent struct {
    TabID   string
    Content any  // Card, TimelineItem, etc.
}

type ShowModalEvent struct {
    Modal Modal
}

type UpdateStatusEvent struct {
    Status StatusConfig
}

type StreamTextEvent struct {
    Text string
}

type StreamCompleteEvent struct {
    Message Message
}
```

---

## Tabs

Tabs are switchable content panes.

### Interface

```go
type Tab struct {
    ID       string
    Label    string
    Badge    string    // Optional notification count
    Content  Content   // What renders in this tab
    Closable bool
    OnClose  func()
}

func NewTab(id, label string, content Content) Tab
```

### Built-in Tab Contents

```go
// Chat tab - conversation viewport
type ChatContent struct {
    messages  []Message
    viewport  viewport.Model
    renderer  MessageRenderer
}

func NewChatContent(renderer MessageRenderer) *ChatContent
func (c *ChatContent) AddMessage(msg Message)
func (c *ChatContent) SetMessages(msgs []Message)
func (c *ChatContent) ScrollToBottom()

// Activity tab - tool timeline + display cards
type ActivityContent struct {
    items    []ActivityItem
    viewport viewport.Model
}

type ActivityItem struct {
    ID        string
    Type      ActivityType  // ToolCall, ToolResult, Card
    Timestamp time.Time
    Title     string
    Content   string
    Status    ActivityStatus  // Pending, Running, Success, Error
    Expanded  bool
}

func NewActivityContent() *ActivityContent
func (a *ActivityContent) AddItem(item ActivityItem)
func (a *ActivityContent) UpdateItem(id string, updates ActivityItem)
```

---

## Modals

Modals are overlays that capture input focus.

### Base Interface

```go
type Modal interface {
    // Identity
    ID() string

    // Rendering
    Title() string
    Size() ModalSize
    Render(width, height int) string

    // Lifecycle
    OnPush(width, height int)
    OnPop()

    // Input handling - returns true if key was handled
    HandleKey(key tea.KeyMsg) (handled bool, cmd tea.Cmd)
}

type ModalSize int

const (
    ModalSmall      ModalSize = iota  // ~30% height
    ModalMedium                        // ~50% height
    ModalLarge                         // ~80% height
    ModalFullscreen                    // 100%
)
```

### ModalManager

```go
type ModalManager struct {
    stack   []Modal
    version int  // Increments on every Push/Pop for change detection
}

func NewModalManager() *ModalManager

func (m *ModalManager) Push(modal Modal)
func (m *ModalManager) Pop() Modal
func (m *ModalManager) Peek() Modal
func (m *ModalManager) HasActive() bool
func (m *ModalManager) Version() int
func (m *ModalManager) Clear()

// Render the active modal (if any)
func (m *ModalManager) Render(width, height int) string

// Route key to active modal
func (m *ModalManager) HandleKey(key tea.KeyMsg) (handled bool, cmd tea.Cmd)
```

### Concrete Modal Types

#### SimpleModal

```go
type SimpleModal struct {
    id      string
    title   string
    content Content
    footer  string
    size    ModalSize
    onClose func()
}

func NewSimpleModal(cfg SimpleModalConfig) *SimpleModal

type SimpleModalConfig struct {
    ID      string
    Title   string
    Content Content
    Footer  string  // Key hints
    Size    ModalSize
    OnClose func()
}
```

#### ConfirmModal

```go
type ConfirmModal struct {
    id       string
    title    string
    message  string
    options  []ConfirmOption
    selected int
    onResult func(value any)
}

type ConfirmOption struct {
    Label string
    Value any
}

func NewConfirmModal(cfg ConfirmModalConfig) *ConfirmModal

// Preset constructors
func NewYesNoModal(title, message string, onResult func(bool)) *ConfirmModal
func NewOKCancelModal(title, message string, onResult func(bool)) *ConfirmModal
```

#### FormModal

```go
type FormModal struct {
    id       string
    title    string
    form     *huh.Form  // Wraps Huh for form handling
    onSubmit func(values map[string]any)
    onCancel func()
}

func NewFormModal(cfg FormModalConfig) *FormModal

type FormModalConfig struct {
    ID       string
    Title    string
    Fields   []FormField
    OnSubmit func(values map[string]any)
    OnCancel func()
}

type FormField struct {
    Key         string
    Label       string
    Type        FormFieldType  // Text, Password, Select, MultiSelect
    Placeholder string
    Options     []string       // For Select/MultiSelect
    Validate    func(string) error
}
```

#### ListModal

```go
type ListModal struct {
    id         string
    title      string
    items      []ListItem
    selected   int
    filter     string
    filterable bool
    onSelect   func(item ListItem)
    onCancel   func()
}

type ListItem struct {
    ID          string
    Title       string
    Description string
    Value       any
}

func NewListModal(cfg ListModalConfig) *ListModal
```

#### WizardModal

```go
type WizardModal struct {
    id         string
    title      string
    steps      []WizardStep
    current    int
    results    map[string]any
    onComplete func(results map[string]any)
    onCancel   func()
}

type WizardStep struct {
    ID          string
    Title       string
    Description string
    Content     Content   // SelectList, MultiSelect, TextInput, Progress, etc.
    Validate    func(value any) error
    Optional    bool
}

func NewWizardModal(cfg WizardModalConfig) *WizardModal

// Navigation
func (w *WizardModal) Next() tea.Cmd
func (w *WizardModal) Previous()
func (w *WizardModal) CanGoNext() bool
func (w *WizardModal) CanGoPrevious() bool
func (w *WizardModal) CurrentStep() *WizardStep
func (w *WizardModal) Progress() (current, total int)
```

#### ApprovalModal

```go
type ApprovalModal struct {
    id         string
    tool       ToolInfo
    options    []ApprovalOption
    selected   int
    queueHint  string  // "[●✓✗?]" from ToolQueue
    onDecision func(decision ApprovalDecision)
}

type ToolInfo struct {
    ID       string
    Name     string
    Params   map[string]any
    Preview  string     // Human-readable description
    Risk     RiskLevel
}

type RiskLevel int

const (
    RiskLow RiskLevel = iota
    RiskMedium
    RiskHigh
)

type ApprovalDecision int

const (
    DecisionApprove ApprovalDecision = iota
    DecisionDeny
    DecisionAlwaysAllow
    DecisionNeverAllow
)

type ApprovalOption struct {
    Label    string
    Decision ApprovalDecision
    Hint     string
}

func NewApprovalModal(cfg ApprovalModalConfig) *ApprovalModal

// Default options
var DefaultApprovalOptions = []ApprovalOption{
    {Label: "Approve (run this time)", Decision: DecisionApprove},
    {Label: "Deny (skip this time)", Decision: DecisionDeny},
    {Label: "Always Allow (never ask again)", Decision: DecisionAlwaysAllow},
    {Label: "Never Allow (block permanently)", Decision: DecisionNeverAllow},
}
```

#### QuickActionsModal

```go
type QuickActionsModal struct {
    id       string
    actions  []QuickAction
    filter   string
    filtered []QuickAction
    selected int
    onSelect func(action QuickAction)
    onCancel func()
}

type QuickAction struct {
    ID          string
    Label       string
    Description string
    Category    string
    Handler     func() tea.Cmd
}

func NewQuickActionsModal(actions []QuickAction, onSelect func(QuickAction)) *QuickActionsModal
```

---

## Content Primitives

Content primitives implement the Content interface and can be used in Tabs or Modals.

### Base Interface

```go
type Content interface {
    // Bubble Tea model
    Init() tea.Cmd
    Update(msg tea.Msg) (Content, tea.Cmd)
    View() string

    // Value extraction (for forms/wizards)
    Value() any

    // Size management
    SetSize(width, height int)
}
```

### Viewport

Scrollable text/markdown display.

```go
type Viewport struct {
    content  string
    viewport viewport.Model
    renderer MessageRenderer
}

func NewViewport() *Viewport
func (v *Viewport) SetContent(content string)
func (v *Viewport) AppendContent(content string)
func (v *Viewport) ScrollToTop()
func (v *Viewport) ScrollToBottom()
```

### SelectList

Single selection from options.

```go
type SelectList struct {
    items    []SelectItem
    selected int
    onSelect func(item SelectItem)
}

type SelectItem struct {
    Label       string
    Description string
    Value       any
}

func NewSelectList(items []SelectItem) *SelectList
func (s *SelectList) Value() any  // Returns selected item's value
```

### MultiSelect

Multiple selection with checkboxes.

```go
type MultiSelect struct {
    items   []MultiSelectItem
    cursor  int
}

type MultiSelectItem struct {
    Label    string
    Key      string
    Selected bool
}

func NewMultiSelect(items []MultiSelectItem) *MultiSelect
func (m *MultiSelect) Value() any  // Returns []string of selected keys
func (m *MultiSelect) Toggle()
func (m *MultiSelect) SelectAll()
func (m *MultiSelect) SelectNone()
```

### TextInput

Single-line text input.

```go
type TextInput struct {
    input       textinput.Model
    label       string
    placeholder string
    validate    func(string) error
    errMsg      string
}

func NewTextInput(cfg TextInputConfig) *TextInput
func (t *TextInput) Value() any  // Returns string
func (t *TextInput) SetValue(value string)
func (t *TextInput) Focus()
func (t *TextInput) Blur()
```

### TextArea

Multi-line text input.

```go
type TextArea struct {
    textarea    textarea.Model
    label       string
    placeholder string
}

func NewTextArea(cfg TextAreaConfig) *TextArea
func (t *TextArea) Value() any  // Returns string
```

### Progress

Live-updating progress display.

```go
type Progress struct {
    items      []ProgressItem
    total      int
    current    int
    message    string
    showBar    bool
    showItems  bool
    maxVisible int
}

type ProgressItem struct {
    Label  string
    Status ProgressStatus  // Pending, Running, Complete, Error
}

func NewProgress(cfg ProgressConfig) *Progress
func (p *Progress) SetTotal(total int)
func (p *Progress) SetCurrent(current int)
func (p *Progress) AddItem(item ProgressItem)
func (p *Progress) UpdateItem(index int, status ProgressStatus)
func (p *Progress) SetMessage(message string)
```

### Timeline

Chronological list of items (for tool history, activity).

```go
type Timeline struct {
    items    []TimelineItem
    viewport viewport.Model
}

type TimelineItem struct {
    ID        string
    Timestamp time.Time
    Icon      string
    Title     string
    Content   string
    Status    TimelineStatus  // Pending, Running, Success, Error
}

func NewTimeline() *Timeline
func (t *Timeline) AddItem(item TimelineItem)
func (t *Timeline) UpdateItem(id string, updates TimelineItem)
func (t *Timeline) Clear()
```

### Spinner

Loading indicator with different types.

```go
type SpinnerType int

const (
    SpinnerDefault    SpinnerType = iota  // Dot spinner
    SpinnerExecution                       // Points, shows elapsed time
    SpinnerStreaming                       // MiniDot, shows token rate
    SpinnerLoading                         // Line spinner
)

type Spinner struct {
    spinnerType SpinnerType
    spinner     spinner.Model
    message     string
    startTime   time.Time
    tokenRate   float64
}

func NewSpinner(spinnerType SpinnerType) *Spinner
func (s *Spinner) SetMessage(message string)
func (s *Spinner) SetTokenRate(rate float64)
func (s *Spinner) Start() tea.Cmd
func (s *Spinner) Stop()
```

---

## Agent Integration

### AgentBackend Interface

```go
type AgentBackend interface {
    // Start a conversation turn
    // Returns event channel for streaming responses
    Stream(ctx context.Context, messages []Message) (<-chan AgentEvent, error)

    // Execute a tool (after approval)
    ExecuteTool(ctx context.Context, tool ToolUse) (ToolResult, error)

    // For interactive approval flows
    ApprovalRequests() <-chan ApprovalRequest
    RespondToApproval(requestID string, decision ApprovalDecision) error

    // Tool metadata
    DescribeTool(name string) ToolDescription
}

type AgentEvent struct {
    Type       AgentEventType
    Text       string       // For EventText
    Tool       *ToolUse     // For EventToolCall
    Result     *ToolResult  // For EventToolResult
    Error      error        // For EventError
    Usage      *TokenUsage  // For EventComplete
}

type AgentEventType int

const (
    EventText AgentEventType = iota
    EventToolCall
    EventToolResult
    EventComplete
    EventError
)

type ApprovalRequest struct {
    ID       string
    Tool     ToolInfo
    Response chan<- ApprovalDecision
}
```

### ToolQueue

Sequential tool processing with pre-classification.

```go
type ToolQueue struct {
    items   []ToolDisposition
    current int
    results []ToolResult
}

type ToolDisposition struct {
    Tool    ToolInfo
    Action  ToolAction
    Outcome ToolOutcome
    Reason  string
}

type ToolAction int

const (
    ActionNeedsApproval ToolAction = iota
    ActionAutoApprove
    ActionAutoDeny
)

type ToolOutcome int

const (
    OutcomePending ToolOutcome = iota
    OutcomeApproved
    OutcomeDenied
)

func NewToolQueue(tools []ToolInfo, classifier ToolClassifier) *ToolQueue

type ToolClassifier func(tool ToolInfo) (ToolAction, string)

func (q *ToolQueue) Next() *ToolDisposition
func (q *ToolQueue) SetOutcome(outcome ToolOutcome, result *ToolResult)
func (q *ToolQueue) ProgressHint() string  // "[●✓✗?]"
func (q *ToolQueue) IsComplete() bool
func (q *ToolQueue) Results() []ToolResult
```

### Message Types

```go
type Message struct {
    ID           string
    Role         MessageRole
    Content      string
    ContentBlocks []ContentBlock
    Timestamp    time.Time
}

type MessageRole string

const (
    RoleUser      MessageRole = "user"
    RoleAssistant MessageRole = "assistant"
    RoleTool      MessageRole = "tool"
    RoleSystem    MessageRole = "system"
)

type ContentBlock struct {
    Type      string         // "text", "tool_use", "tool_result"
    Text      string
    ToolUse   *ToolUse
    ToolResult *ToolResult
}

type ToolUse struct {
    ID     string
    Name   string
    Input  map[string]any
}

type ToolResult struct {
    ToolUseID string
    Content   string
    IsError   bool
}
```

---

## Theming

### Theme Interface

```go
type Theme interface {
    // Identity
    Name() string

    // Base colors
    Background() lipgloss.Color
    Foreground() lipgloss.Color

    // Accent colors
    Primary() lipgloss.Color
    Secondary() lipgloss.Color

    // Status colors
    Success() lipgloss.Color
    Warning() lipgloss.Color
    Error() lipgloss.Color
    Info() lipgloss.Color

    // UI colors
    Border() lipgloss.Color
    BorderFocused() lipgloss.Color
    Muted() lipgloss.Color

    // Role colors (for messages)
    UserColor() lipgloss.Color
    AssistantColor() lipgloss.Color
    ToolColor() lipgloss.Color
    SystemColor() lipgloss.Color

    // Styles (composed from colors)
    Styles() ThemeStyles
}

type ThemeStyles struct {
    // Text styles
    Title      lipgloss.Style
    Subtitle   lipgloss.Style
    Body       lipgloss.Style
    Muted      lipgloss.Style

    // Component styles
    Input        lipgloss.Style
    InputFocused lipgloss.Style
    Button       lipgloss.Style
    ButtonActive lipgloss.Style

    // Status styles
    StatusBar    lipgloss.Style
    TabBar       lipgloss.Style
    TabActive    lipgloss.Style
    TabInactive  lipgloss.Style

    // Modal styles
    ModalBox     lipgloss.Style
    ModalTitle   lipgloss.Style
    ModalFooter  lipgloss.Style

    // Tool styles
    ToolApproval lipgloss.Style
    ToolSuccess  lipgloss.Style
    ToolError    lipgloss.Style
    ToolPending  lipgloss.Style
}
```

### Built-in Themes

```go
// Dracula theme (from hex)
func NewDraculaTheme() Theme

// Nord theme
func NewNordTheme() Theme

// Gruvbox theme
func NewGruvboxTheme() Theme
```

### MessageRenderer

```go
type MessageRenderer interface {
    RenderMessage(msg Message, width int, theme Theme) string
    RenderRole(role MessageRole, theme Theme) string
}

// Simple renderer - plain text
type SimpleRenderer struct{}

// Markdown renderer - uses glamour
type MarkdownRenderer struct {
    glamourRenderer *glamour.TermRenderer
}

// NeoTerminal renderer - colored borders by role (from hex)
type NeoTerminalRenderer struct{}
```

---

## Implementation Plan

### Phase 1: Core Primitives

#### Task 1: Project Setup

**Files:**
- Create: `go.mod`
- Create: `tui.go` (package entry point)
- Create: `README.md`

**Step 1: Initialize Go module**

```bash
cd /Users/dylanr/work/tui
go mod init github.com/2389-research/tui
```

**Step 2: Add dependencies**

```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/glamour
go get github.com/charmbracelet/huh
```

**Step 3: Create package entry point**

```go
// tui.go
package tui

// Package tui provides a shared TUI library for building
// multi-agent terminal interfaces.
//
// The library provides:
// - Shell: top-level container with tabs, modals, input, status
// - Tabs: switchable content panes
// - Modals: overlays including wizards, forms, approvals
// - Content: composable primitives (viewport, lists, forms)
// - Agent: backend interface for streaming, tool execution

const Version = "0.1.0"
```

**Step 4: Commit**

```bash
git init
git add .
git commit -m "feat: initialize tui library project"
```

---

#### Task 2: Theme System

**Files:**
- Create: `theme/theme.go`
- Create: `theme/dracula.go`
- Create: `theme/styles.go`
- Test: `theme/theme_test.go`

**Step 1: Write theme interface test**

```go
// theme/theme_test.go
package theme

import (
    "testing"
    "github.com/charmbracelet/lipgloss"
)

func TestDraculaTheme(t *testing.T) {
    theme := NewDraculaTheme()

    if theme.Name() != "dracula" {
        t.Errorf("expected name 'dracula', got %s", theme.Name())
    }

    // Verify colors are set
    if theme.Primary() == "" {
        t.Error("primary color not set")
    }

    // Verify styles are composed
    styles := theme.Styles()
    if styles.Title.GetForeground() == lipgloss.NoColor {
        t.Error("title style foreground not set")
    }
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./theme/... -v
```

Expected: FAIL - package doesn't exist

**Step 3: Implement theme interface**

```go
// theme/theme.go
package theme

import "github.com/charmbracelet/lipgloss"

type Theme interface {
    Name() string

    // Base colors
    Background() lipgloss.Color
    Foreground() lipgloss.Color

    // Accent colors
    Primary() lipgloss.Color
    Secondary() lipgloss.Color

    // Status colors
    Success() lipgloss.Color
    Warning() lipgloss.Color
    Error() lipgloss.Color
    Info() lipgloss.Color

    // UI colors
    Border() lipgloss.Color
    BorderFocused() lipgloss.Color
    Muted() lipgloss.Color

    // Role colors
    UserColor() lipgloss.Color
    AssistantColor() lipgloss.Color
    ToolColor() lipgloss.Color
    SystemColor() lipgloss.Color

    // Composed styles
    Styles() ThemeStyles
}
```

```go
// theme/styles.go
package theme

import "github.com/charmbracelet/lipgloss"

type ThemeStyles struct {
    // Text
    Title      lipgloss.Style
    Subtitle   lipgloss.Style
    Body       lipgloss.Style
    Muted      lipgloss.Style

    // Components
    Input        lipgloss.Style
    InputFocused lipgloss.Style
    Button       lipgloss.Style
    ButtonActive lipgloss.Style

    // Layout
    StatusBar   lipgloss.Style
    TabBar      lipgloss.Style
    TabActive   lipgloss.Style
    TabInactive lipgloss.Style

    // Modals
    ModalBox    lipgloss.Style
    ModalTitle  lipgloss.Style
    ModalFooter lipgloss.Style

    // Tools
    ToolApproval lipgloss.Style
    ToolSuccess  lipgloss.Style
    ToolError    lipgloss.Style
    ToolPending  lipgloss.Style
}
```

```go
// theme/dracula.go
package theme

import "github.com/charmbracelet/lipgloss"

// Dracula color palette
const (
    draculaBackground  = lipgloss.Color("#282a36")
    draculaCurrentLine = lipgloss.Color("#44475a")
    draculaForeground  = lipgloss.Color("#f8f8f2")
    draculaComment     = lipgloss.Color("#6272a4")
    draculaCyan        = lipgloss.Color("#8be9fd")
    draculaGreen       = lipgloss.Color("#50fa7b")
    draculaOrange      = lipgloss.Color("#ffb86c")
    draculaPink        = lipgloss.Color("#ff79c6")
    draculaPurple      = lipgloss.Color("#bd93f9")
    draculaRed         = lipgloss.Color("#ff5555")
    draculaYellow      = lipgloss.Color("#f1fa8c")
)

type draculaTheme struct {
    styles ThemeStyles
}

func NewDraculaTheme() Theme {
    t := &draculaTheme{}
    t.styles = t.buildStyles()
    return t
}

func (t *draculaTheme) Name() string { return "dracula" }

func (t *draculaTheme) Background() lipgloss.Color  { return draculaBackground }
func (t *draculaTheme) Foreground() lipgloss.Color  { return draculaForeground }
func (t *draculaTheme) Primary() lipgloss.Color     { return draculaPurple }
func (t *draculaTheme) Secondary() lipgloss.Color   { return draculaCyan }
func (t *draculaTheme) Success() lipgloss.Color     { return draculaGreen }
func (t *draculaTheme) Warning() lipgloss.Color     { return draculaOrange }
func (t *draculaTheme) Error() lipgloss.Color       { return draculaRed }
func (t *draculaTheme) Info() lipgloss.Color        { return draculaCyan }
func (t *draculaTheme) Border() lipgloss.Color      { return draculaComment }
func (t *draculaTheme) BorderFocused() lipgloss.Color { return draculaPurple }
func (t *draculaTheme) Muted() lipgloss.Color       { return draculaComment }
func (t *draculaTheme) UserColor() lipgloss.Color   { return draculaOrange }
func (t *draculaTheme) AssistantColor() lipgloss.Color { return draculaGreen }
func (t *draculaTheme) ToolColor() lipgloss.Color   { return draculaCyan }
func (t *draculaTheme) SystemColor() lipgloss.Color { return draculaComment }

func (t *draculaTheme) Styles() ThemeStyles { return t.styles }

func (t *draculaTheme) buildStyles() ThemeStyles {
    return ThemeStyles{
        Title: lipgloss.NewStyle().
            Foreground(draculaPurple).
            Bold(true),
        Subtitle: lipgloss.NewStyle().
            Foreground(draculaCyan),
        Body: lipgloss.NewStyle().
            Foreground(draculaForeground),
        Muted: lipgloss.NewStyle().
            Foreground(draculaComment),
        Input: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(draculaComment).
            Padding(0, 1),
        InputFocused: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(draculaPurple).
            Padding(0, 1),
        Button: lipgloss.NewStyle().
            Foreground(draculaForeground).
            Background(draculaCurrentLine).
            Padding(0, 2),
        ButtonActive: lipgloss.NewStyle().
            Foreground(draculaBackground).
            Background(draculaPurple).
            Padding(0, 2),
        StatusBar: lipgloss.NewStyle().
            Foreground(draculaForeground).
            Background(draculaCurrentLine),
        TabBar: lipgloss.NewStyle().
            Foreground(draculaForeground),
        TabActive: lipgloss.NewStyle().
            Foreground(draculaPurple).
            Bold(true).
            Underline(true),
        TabInactive: lipgloss.NewStyle().
            Foreground(draculaComment),
        ModalBox: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(draculaPurple).
            Padding(1, 2),
        ModalTitle: lipgloss.NewStyle().
            Foreground(draculaPurple).
            Bold(true),
        ModalFooter: lipgloss.NewStyle().
            Foreground(draculaComment),
        ToolApproval: lipgloss.NewStyle().
            Foreground(draculaOrange),
        ToolSuccess: lipgloss.NewStyle().
            Foreground(draculaGreen),
        ToolError: lipgloss.NewStyle().
            Foreground(draculaRed),
        ToolPending: lipgloss.NewStyle().
            Foreground(draculaYellow),
    }
}
```

**Step 4: Run tests**

```bash
go test ./theme/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add theme/
git commit -m "feat: add theme system with Dracula theme"
```

---

#### Task 3: Content Primitives - SelectList

**Files:**
- Create: `content/content.go`
- Create: `content/selectlist.go`
- Test: `content/selectlist_test.go`

**Step 1: Write test**

```go
// content/selectlist_test.go
package content

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestSelectList(t *testing.T) {
    items := []SelectItem{
        {Label: "Option A", Value: "a"},
        {Label: "Option B", Value: "b"},
        {Label: "Option C", Value: "c"},
    }

    list := NewSelectList(items)

    // Initial selection is first item
    if list.Value() != "a" {
        t.Errorf("expected initial value 'a', got %v", list.Value())
    }

    // Move down
    list, _ = list.Update(tea.KeyMsg{Type: tea.KeyDown})
    if list.Value() != "b" {
        t.Errorf("expected value 'b' after down, got %v", list.Value())
    }

    // Move down again
    list, _ = list.Update(tea.KeyMsg{Type: tea.KeyDown})
    if list.Value() != "c" {
        t.Errorf("expected value 'c' after down, got %v", list.Value())
    }

    // Can't go past end
    list, _ = list.Update(tea.KeyMsg{Type: tea.KeyDown})
    if list.Value() != "c" {
        t.Errorf("expected value to stay 'c', got %v", list.Value())
    }
}
```

**Step 2: Run test**

```bash
go test ./content/... -v
```

Expected: FAIL

**Step 3: Implement**

```go
// content/content.go
package content

import tea "github.com/charmbracelet/bubbletea"

// Content is the interface for content primitives that can
// be used in Tabs or Modals.
type Content interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Content, tea.Cmd)
    View() string
    Value() any
    SetSize(width, height int)
}
```

```go
// content/selectlist.go
package content

import (
    "strings"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type SelectItem struct {
    Label       string
    Description string
    Value       any
}

type SelectList struct {
    items    []SelectItem
    selected int
    width    int
    height   int

    // Styles
    selectedStyle   lipgloss.Style
    unselectedStyle lipgloss.Style
    descStyle       lipgloss.Style
}

func NewSelectList(items []SelectItem) *SelectList {
    return &SelectList{
        items:    items,
        selected: 0,
        selectedStyle: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#50fa7b")).
            Bold(true),
        unselectedStyle: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#f8f8f2")),
        descStyle: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#6272a4")),
    }
}

func (s *SelectList) Init() tea.Cmd {
    return nil
}

func (s *SelectList) Update(msg tea.Msg) (Content, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            if s.selected > 0 {
                s.selected--
            }
        case tea.KeyDown:
            if s.selected < len(s.items)-1 {
                s.selected++
            }
        }
    }
    return s, nil
}

func (s *SelectList) View() string {
    var b strings.Builder

    for i, item := range s.items {
        cursor := "  "
        style := s.unselectedStyle

        if i == s.selected {
            cursor = "▸ "
            style = s.selectedStyle
        }

        b.WriteString(cursor)
        b.WriteString(style.Render(item.Label))

        if item.Description != "" {
            b.WriteString("\n    ")
            b.WriteString(s.descStyle.Render(item.Description))
        }

        if i < len(s.items)-1 {
            b.WriteString("\n")
        }
    }

    return b.String()
}

func (s *SelectList) Value() any {
    if s.selected >= 0 && s.selected < len(s.items) {
        return s.items[s.selected].Value
    }
    return nil
}

func (s *SelectList) SetSize(width, height int) {
    s.width = width
    s.height = height
}

func (s *SelectList) Selected() int {
    return s.selected
}

func (s *SelectList) SetSelected(index int) {
    if index >= 0 && index < len(s.items) {
        s.selected = index
    }
}
```

**Step 4: Run tests**

```bash
go test ./content/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add content/
git commit -m "feat: add Content interface and SelectList primitive"
```

---

### Phase 2: Modal System

#### Task 4: ModalManager

**Files:**
- Create: `modal/modal.go`
- Create: `modal/manager.go`
- Test: `modal/manager_test.go`

_(Continue with similar TDD pattern for each component)_

---

### Phase 3: Shell

#### Task 5: Shell Core

**Files:**
- Create: `shell/shell.go`
- Create: `shell/input.go`
- Create: `shell/statusbar.go`
- Create: `shell/tabbar.go`
- Test: `shell/shell_test.go`

---

### Phase 4: Agent Integration

#### Task 6: AgentBackend Interface

**Files:**
- Create: `agent/backend.go`
- Create: `agent/types.go`
- Create: `agent/toolqueue.go`
- Test: `agent/toolqueue_test.go`

---

### Phase 5: Specialized Modals

#### Task 7: WizardModal

**Files:**
- Create: `modal/wizard.go`
- Test: `modal/wizard_test.go`

#### Task 8: ApprovalModal

**Files:**
- Create: `modal/approval.go`
- Test: `modal/approval_test.go`

---

### Phase 6: Integration Examples

#### Task 9: Example Consumer (minimal agent)

**Files:**
- Create: `examples/minimal/main.go`

This demonstrates how hex or jeff would consume the library.

---

## File Structure

```
tui/
├── go.mod
├── go.sum
├── tui.go                    # Package entry point
├── README.md
│
├── theme/
│   ├── theme.go              # Theme interface
│   ├── styles.go             # ThemeStyles struct
│   ├── dracula.go            # Dracula implementation
│   ├── nord.go               # Nord implementation
│   └── theme_test.go
│
├── content/
│   ├── content.go            # Content interface
│   ├── viewport.go           # Scrollable text
│   ├── selectlist.go         # Single selection
│   ├── multiselect.go        # Multi selection
│   ├── textinput.go          # Single line input
│   ├── textarea.go           # Multi line input
│   ├── progress.go           # Progress display
│   ├── timeline.go           # Chronological items
│   ├── spinner.go            # Loading indicator
│   └── *_test.go
│
├── modal/
│   ├── modal.go              # Modal interface
│   ├── manager.go            # ModalManager (stack)
│   ├── simple.go             # SimpleModal
│   ├── confirm.go            # ConfirmModal
│   ├── form.go               # FormModal (wraps Huh)
│   ├── list.go               # ListModal
│   ├── wizard.go             # WizardModal
│   ├── approval.go           # ApprovalModal
│   ├── quickactions.go       # QuickActionsModal
│   └── *_test.go
│
├── shell/
│   ├── shell.go              # Shell struct and methods
│   ├── input.go              # Input area
│   ├── statusbar.go          # Status bar
│   ├── tabbar.go             # Tab bar
│   ├── events.go             # ShellEvent types
│   └── *_test.go
│
├── agent/
│   ├── backend.go            # AgentBackend interface
│   ├── types.go              # Message, ToolUse, etc.
│   ├── toolqueue.go          # ToolQueue
│   └── *_test.go
│
├── render/
│   ├── renderer.go           # MessageRenderer interface
│   ├── simple.go             # Plain text renderer
│   ├── markdown.go           # Glamour-based renderer
│   ├── neoterminal.go        # Colored border renderer
│   └── *_test.go
│
├── docs/
│   └── plans/
│       └── 2025-01-05-shared-tui-library-design.md
│
└── examples/
    └── minimal/
        └── main.go           # Minimal consumer example
```

---

## Migration Path

### For Hex

1. Add `github.com/2389-research/tui` dependency
2. Replace `internal/ui/theme` with `tui/theme`
3. Replace `internal/ui/overlay*` with `tui/modal`
4. Replace `internal/ui/model.go` with `tui/shell`
5. Implement `AgentBackend` wrapping existing mux integration
6. Keep domain-specific tools in hex

### For Jeff

1. Add `github.com/2389-research/tui` dependency
2. Replace `internal/ui/themes` with `tui/theme`
3. Replace modal handling with `tui/modal`
4. Replace `internal/ui/model.go` with `tui/shell`
5. Implement `AgentBackend` wrapping existing mux backend
6. Keep domain-specific tools in jeff
7. Convert wizard.go/onboarding.go to use `WizardModal`

---

## Open Questions (Resolved)

### 1. Huh Integration

**Decision: Wrap Huh with themed adapters**

Both hex and jeff use Huh for forms (approval dialogs, settings, onboarding). The pattern is consistent:

1. Create a wrapper component that embeds `*huh.Form`
2. Provide a theme adapter function: `huhThemeFromTUITheme(theme Theme) *huh.Theme`
3. Handle Huh's state machine (`huh.StateCompleted`) internally
4. Expose simple callbacks: `OnSubmit`, `OnCancel`

**Implementation:**

```go
// modal/form.go
type FormModal struct {
    id       string
    title    string
    form     *huh.Form
    theme    Theme
    onSubmit func(values map[string]any)
    onCancel func()
}

// HuhThemeAdapter converts TUI theme to Huh theme
func HuhThemeAdapter(theme Theme) *huh.Theme {
    t := huh.ThemeBase()

    t.Focused.Title = t.Focused.Title.Foreground(theme.Primary())
    t.Focused.Description = t.Focused.Description.Foreground(theme.Muted())
    t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(theme.Primary())
    t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(theme.Primary())
    t.Focused.FocusedButton = t.Focused.FocusedButton.
        Foreground(theme.Background()).
        Background(theme.Primary())

    t.Blurred.Base = t.Blurred.Base.Foreground(theme.Muted())
    t.Blurred.Title = t.Blurred.Title.Foreground(theme.Muted())

    return t
}
```

**Key insight from hex:** Use `huh.NewSelect` instead of `huh.NewConfirm` because Select responds to Enter key properly.

---

### 2. Autocomplete System

**Decision: Include CompletionProvider interface in Shell with pluggable providers**

Hex has a sophisticated autocomplete system with fuzzy matching. This should be part of the shared library.

**Interface:**

```go
// shell/autocomplete.go
type CompletionProvider interface {
    // GetCompletions returns suggestions for the input
    GetCompletions(input string) []Completion
}

type Completion struct {
    Value       string  // The actual completion text
    Display     string  // What shows in the dropdown
    Description string  // Context about the completion
    Score       int     // For ranking (higher = better)
}

type Autocomplete struct {
    providers       map[string]CompletionProvider
    active          bool
    completions     []Completion
    selectedIndex   int
    maxCompletions  int
}

func NewAutocomplete() *Autocomplete

func (ac *Autocomplete) RegisterProvider(name string, provider CompletionProvider)
func (ac *Autocomplete) Show(input string, providerName string)
func (ac *Autocomplete) Hide()
func (ac *Autocomplete) Next()
func (ac *Autocomplete) Previous()
func (ac *Autocomplete) GetSelected() *Completion

// DetectProvider determines which provider based on input context
// "/" prefix → command provider
// "./" or "~/" prefix → file provider
// Otherwise → history provider
func DetectProvider(input string) string
```

**Built-in Providers:**

```go
// SlashCommandProvider - for /commands
type SlashCommandProvider struct {
    commands     []string
    descriptions map[string]string
}

func (sp *SlashCommandProvider) SetCommands(commands []string, descriptions map[string]string)

// FileProvider - for file paths
type FileProvider struct {
    basePath string
}

func (fp *FileProvider) SetBasePath(path string)

// HistoryProvider - for recent inputs
type HistoryProvider struct {
    history []string
}

func (hp *HistoryProvider) AddToHistory(command string)
```

**Integration with Shell:**

```go
type Shell struct {
    // ...
    autocomplete *Autocomplete
}

func (s *Shell) RegisterCompletionProvider(name string, provider CompletionProvider)
```

**Dependency:** Uses `github.com/sahilm/fuzzy` for fuzzy matching.

---

### 3. Message Persistence

**Decision: Fully backend responsibility, no TUI-level caching**

Both hex and jeff use SQLite with golang-migrate for persistence. This is domain-specific:
- Different schemas for different agents
- Different conversation structures
- Different metadata needs

The TUI library should:
- Accept messages via `AgentBackend.Stream()` events
- Store in-memory for rendering only
- Delegate all persistence to the consumer

**The contract:**

```go
// Shell holds messages in memory for rendering
type Shell struct {
    messages []Message  // In-memory only
}

// Consumer is responsible for persistence
type MyBackend struct {
    db *sql.DB  // Consumer's database
}

func (b *MyBackend) Stream(ctx context.Context, messages []Message) (<-chan AgentEvent, error) {
    // Consumer persists messages before/after streaming
}
```

**Rationale:** Storage schemas differ significantly:
- Hex: conversations, messages, todos, history, favorites
- Jeff: conversations, messages, feedback, oauth tokens

Embedding persistence would couple the TUI to specific schemas.

---

### 4. Plugin/Extensibility System

**Decision: Hooks and registration points, not a full plugin system**

The TUI provides extensibility through:

**a) Component Registration:**

```go
// Register custom content types
tui.RegisterContent("email-card", func() Content {
    return NewEmailCard()
})

// Register custom modal types
tui.RegisterModal("email-compose", func() Modal {
    return NewComposeModal()
})
```

**b) Event Hooks:**

```go
type ShellHooks struct {
    // Called before/after key handling
    BeforeKeyPress func(key tea.KeyMsg) (handled bool)
    AfterKeyPress  func(key tea.KeyMsg)

    // Called on lifecycle events
    OnTabSwitch    func(from, to string)
    OnModalPush    func(modal Modal)
    OnModalPop     func(modal Modal)
    OnInputSubmit  func(input string) (handled bool)

    // Called for custom rendering
    RenderHeader   func(width int) string
    RenderFooter   func(width int) string
}

func (s *Shell) SetHooks(hooks ShellHooks)
```

**c) Custom Tab Content:**

Consumers implement the `Content` interface for domain-specific displays:

```go
// In hex: tool timeline
type ToolTimelineContent struct { ... }

// In jeff: email inbox
type InboxContent struct { ... }
```

**Why not a full plugin system:**
- Plugins add complexity (discovery, versioning, isolation)
- Hooks provide 80% of the value with 20% of the complexity
- Consumers can build their own plugin system on top

---

### 5. Accessibility

**Decision: Minimal base support, document extensibility patterns**

Terminal accessibility is limited but we can provide:

**a) Keyboard-only navigation (required):**
- All actions reachable via keyboard
- Clear focus indicators
- Consistent key bindings documented

**b) High contrast support:**

```go
type Theme interface {
    // ...
    HighContrast() bool

    // Contrast-aware colors
    ContrastForeground() lipgloss.Color
    ContrastBackground() lipgloss.Color
}

// Built-in high contrast theme
func NewHighContrastTheme() Theme
```

**c) Screen reader hooks (optional):**

```go
type AccessibilityHooks struct {
    // Called when focus changes
    OnFocusChange func(description string)

    // Called on important state changes
    OnAnnounce func(message string, priority AnnouncePriority)
}

type AnnouncePriority int

const (
    AnnouncePolite AnnouncePriority = iota  // Queue behind current speech
    AnnounceAssertive                        // Interrupt current speech
)
```

**Note:** Actual screen reader integration depends on the terminal emulator and OS. The TUI provides hooks; consumers wire them to platform-specific APIs.

---

### 6. Mouse Support

**Decision: Optional, enabled via Shell configuration**

Both codebases use mouse:
- Jeff: `tea.WithMouseCellMotion()` for scrolling
- Hex: Mouse wheel scrolling + hover detection for messages

**Configuration:**

```go
type ShellConfig struct {
    // ...
    Mouse MouseConfig
}

type MouseConfig struct {
    Enabled      bool  // Enable mouse support
    ScrollLines  int   // Lines to scroll per wheel event (default: 3)
    HoverEnabled bool  // Enable hover detection
    ShiftPassthrough bool // Pass through when Shift held (for text selection)
}

// Default config
var DefaultMouseConfig = MouseConfig{
    Enabled:      true,
    ScrollLines:  3,
    HoverEnabled: true,
    ShiftPassthrough: true,
}
```

**Shell mouse handling:**

```go
func (s *Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Pass through for text selection when Shift held
        if msg.Shift && s.config.Mouse.ShiftPassthrough {
            return s, nil
        }

        switch msg.Button {
        case tea.MouseButtonWheelUp:
            s.activeContent.ScrollUp(s.config.Mouse.ScrollLines)
        case tea.MouseButtonWheelDown:
            s.activeContent.ScrollDown(s.config.Mouse.ScrollLines)
        case tea.MouseButtonNone:
            if msg.Action == tea.MouseActionMotion && s.config.Mouse.HoverEnabled {
                s.updateHover(msg.X, msg.Y)
            }
        }
    }
}
```

**Bubble Tea program setup:**

```go
func (s *Shell) Run() error {
    opts := []tea.ProgramOption{tea.WithAltScreen()}

    if s.config.Mouse.Enabled {
        opts = append(opts, tea.WithMouseCellMotion())
    }

    p := tea.NewProgram(s, opts...)
    return p.Start()
}
```

---

### 7. Theming (Extended)

**Decision: Comprehensive theme system with 50+ style fields**

The hex theme has ~50 style fields covering all UI states. The shared library should match this comprehensiveness.

**Extended Theme Interface:**

```go
type Theme interface {
    Name() string

    // Base colors (10)
    Background() lipgloss.Color
    Foreground() lipgloss.Color
    Primary() lipgloss.Color
    Secondary() lipgloss.Color
    Success() lipgloss.Color
    Warning() lipgloss.Color
    Error() lipgloss.Color
    Info() lipgloss.Color
    Border() lipgloss.Color
    Muted() lipgloss.Color

    // Role colors (4)
    UserColor() lipgloss.Color
    AssistantColor() lipgloss.Color
    ToolColor() lipgloss.Color
    SystemColor() lipgloss.Color

    // Gradients (for fancy effects)
    TitleGradient() []lipgloss.Color

    // Full style set
    Styles() ThemeStyles
}

type ThemeStyles struct {
    // Text (5)
    Title      lipgloss.Style
    Subtitle   lipgloss.Style
    Body       lipgloss.Style
    Muted      lipgloss.Style
    Emphasized lipgloss.Style

    // Status (4)
    Success lipgloss.Style
    Error   lipgloss.Style
    Warning lipgloss.Style
    Info    lipgloss.Style

    // Interactive (6)
    Border        lipgloss.Style
    BorderFocused lipgloss.Style
    Input         lipgloss.Style
    InputFocused  lipgloss.Style
    Button        lipgloss.Style
    ButtonActive  lipgloss.Style

    // Layout (6)
    StatusBar   lipgloss.Style
    TabBar      lipgloss.Style
    TabActive   lipgloss.Style
    TabInactive lipgloss.Style
    ViewMode    lipgloss.Style
    TokenCounter lipgloss.Style

    // Tool states (6)
    ToolApproval  lipgloss.Style
    ToolExecuting lipgloss.Style
    ToolSuccess   lipgloss.Style
    ToolError     lipgloss.Style
    ToolResult    lipgloss.Style
    ToolCall      lipgloss.Style

    // Autocomplete (4)
    AutocompleteDropdown lipgloss.Style
    AutocompleteItem     lipgloss.Style
    AutocompleteSelected lipgloss.Style
    AutocompleteHelp     lipgloss.Style

    // Modal (3)
    ModalBox    lipgloss.Style
    ModalTitle  lipgloss.Style
    ModalFooter lipgloss.Style

    // Help (3)
    HelpPanel lipgloss.Style
    HelpKey   lipgloss.Style
    HelpDesc  lipgloss.Style

    // List (3)
    ListItem         lipgloss.Style
    ListItemSelected lipgloss.Style
    ListItemActive   lipgloss.Style

    // Code (5)
    Code      lipgloss.Style
    CodeBlock lipgloss.Style
    Keyword   lipgloss.Style
    String    lipgloss.Style
    Number    lipgloss.Style

    // Links (2)
    Link      lipgloss.Style
    LinkHover lipgloss.Style

    // Messages (1)
    UserMessage lipgloss.Style
}
```

**Theme Registration:**

```go
var themes = map[string]func() Theme{
    "dracula": NewDraculaTheme,
    "nord":    NewNordTheme,
    "gruvbox": NewGruvboxTheme,
}

func RegisterTheme(name string, constructor func() Theme) {
    themes[name] = constructor
}

func GetTheme(name string) Theme {
    if constructor, ok := themes[name]; ok {
        return constructor()
    }
    return NewDraculaTheme()
}

func AvailableThemes() []string
```

---

### 8. Customizable Status Bar

**Decision: Segmented status bar with pluggable sections**

Hex's status bar has multiple independently updateable sections. The shared library should support this.

**Interface:**

```go
// shell/statusbar.go
type StatusBar struct {
    sections  map[string]StatusSection
    order     []string  // Display order
    width     int
    theme     Theme
}

type StatusSection struct {
    ID       string
    Content  string
    Style    lipgloss.Style
    MinWidth int
    MaxWidth int
    Priority int  // Higher priority sections get space first
}

type StatusConfig struct {
    Model       string  // "claude-3.5-sonnet"
    Status      ConnectionStatus
    TokenCount  string  // "12.5k / 100k"
    Mode        string  // "vim:NORMAL" or empty
    CustomMsg   string  // Temporary message
    Progress    string  // "[●●●○○]" or empty
    Hints       string  // "Ctrl+H help"
}

type ConnectionStatus int

const (
    StatusDisconnected ConnectionStatus = iota
    StatusConnected
    StatusStreaming
    StatusError
)

func NewStatusBar(theme Theme) *StatusBar

// Section management
func (s *StatusBar) SetSection(id string, section StatusSection)
func (s *StatusBar) RemoveSection(id string)
func (s *StatusBar) SetOrder(ids []string)

// Convenience setters (map to sections internally)
func (s *StatusBar) SetModel(model string)
func (s *StatusBar) SetConnectionStatus(status ConnectionStatus)
func (s *StatusBar) SetTokenCount(used, total int)
func (s *StatusBar) SetMode(mode string)
func (s *StatusBar) SetCustomMessage(msg string, duration time.Duration)
func (s *StatusBar) SetProgress(progress string)
func (s *StatusBar) SetHints(hints string)
```

**Default sections:**

```
┌─────────────────────────────────────────────────────────────────────┐
│ claude-3.5 │ ● Connected │ 12.5k/100k │ NORMAL │ [●●○] │ Ctrl+H    │
└─────────────────────────────────────────────────────────────────────┘
   model       status        tokens       mode    progress   hints
```

---

### 9. Other Customization

**a) Input Area Customization:**

```go
type InputConfig struct {
    Placeholder    string
    Prefix         string  // "> " or ">>> "
    MultiLine      bool
    MaxHeight      int     // For multi-line
    ShowCharCount  bool
    MaxChars       int     // 0 = unlimited
}

func (s *Shell) SetInputConfig(cfg InputConfig)
```

**b) Tab Bar Customization:**

```go
type TabBarConfig struct {
    Position    TabBarPosition  // Top, Bottom
    Style       TabStyle        // Underline, Boxed, Pills
    ShowBadges  bool
    ShowClose   bool
    MaxTabs     int  // Show "..." overflow
}

type TabBarPosition int

const (
    TabBarTop TabBarPosition = iota
    TabBarBottom
)

type TabStyle int

const (
    TabStyleUnderline TabStyle = iota
    TabStyleBoxed
    TabStylePills
)
```

**c) Modal Customization:**

```go
type ModalConfig struct {
    Backdrop     bool           // Dim background
    BackdropChar rune           // Character for backdrop
    Animation    ModalAnimation // None, FadeIn, SlideUp
    CloseOnEsc   bool
    CloseOnClick bool           // Click outside to close
}

type ModalAnimation int

const (
    ModalAnimationNone ModalAnimation = iota
    ModalAnimationFadeIn
    ModalAnimationSlideUp
)
```

**d) Key Binding Customization:**

```go
type KeyBindings struct {
    Submit       key.Binding  // Enter
    Cancel       key.Binding  // Esc
    NextTab      key.Binding  // Ctrl+Tab
    PrevTab      key.Binding  // Ctrl+Shift+Tab
    Help         key.Binding  // Ctrl+H or ?
    QuickActions key.Binding  // Ctrl+K
    ScrollUp     key.Binding  // Ctrl+U or PgUp
    ScrollDown   key.Binding  // Ctrl+D or PgDn
}

var DefaultKeyBindings = KeyBindings{
    Submit:       key.NewBinding(key.WithKeys("enter")),
    Cancel:       key.NewBinding(key.WithKeys("esc")),
    NextTab:      key.NewBinding(key.WithKeys("ctrl+tab")),
    PrevTab:      key.NewBinding(key.WithKeys("ctrl+shift+tab")),
    Help:         key.NewBinding(key.WithKeys("ctrl+h", "?")),
    QuickActions: key.NewBinding(key.WithKeys("ctrl+k")),
    ScrollUp:     key.NewBinding(key.WithKeys("ctrl+u", "pgup")),
    ScrollDown:   key.NewBinding(key.WithKeys("ctrl+d", "pgdn")),
}

func (s *Shell) SetKeyBindings(bindings KeyBindings)
```

---

### 10. User Configuration System

**Decision: Two-tier configuration with app defaults and user overrides**

App developers set sensible defaults; users can override via config file.

**Configuration Priority (lowest to highest):**

1. **Library defaults** - Built into the TUI library
2. **App defaults** - Set by hex/jeff at initialization
3. **User config file** - User's personal overrides

**Config File Locations (checked in order):**

```go
// For any app (hex, jeff, etc.):
// 1. ~/.config/{appname}/ui.toml
// 2. ~/.{appname}rc (legacy fallback)
// 3. ${APPNAME}_UI_CONFIG (env override)
//
// Examples:
//   hex:  ~/.config/hex/ui.toml,  ~/.hexrc,  $HEX_UI_CONFIG
//   jeff: ~/.config/jeff/ui.toml, ~/.jeffrc, $JEFF_UI_CONFIG

func DefaultConfigPaths(appName string) []string {
    home, _ := os.UserHomeDir()
    return []string{
        filepath.Join(home, ".config", appName, "ui.toml"),
        filepath.Join(home, "."+appName+"rc"),
    }
}
```

**User Config File Format (TOML):**

```toml
# ~/.config/{appname}/ui.toml
# e.g., ~/.config/hex/ui.toml or ~/.config/jeff/ui.toml

[theme]
name = "dracula"  # Use built-in theme as base

# Override specific colors
[theme.colors]
primary = "#ff79c6"      # Override primary to pink
background = "#1a1a2e"   # Darker background
success = "#00ff00"      # Brighter green

# Override specific styles (partial - only what you specify)
[theme.styles.title]
foreground = "#bd93f9"
bold = true
italic = false

[theme.styles.input]
border = "rounded"
border_foreground = "#6272a4"

[mouse]
enabled = true
scroll_lines = 5         # Override default of 3
hover_enabled = false    # Disable hover

[keybindings]
submit = ["enter", "ctrl+m"]
cancel = ["esc", "ctrl+c"]
help = ["f1", "?"]       # Remap help to F1
quick_actions = ["ctrl+p"]  # VS Code style

[statusbar]
# Hide sections by setting empty
mode = ""                # Hide mode section
# Reorder sections
order = ["model", "status", "tokens", "hints"]

[statusbar.custom_sections.git]
content = "main"         # Static or use hooks to update
priority = 50

[tabbar]
position = "bottom"
style = "pills"

[input]
prefix = "λ "            # Custom prompt
placeholder = "Ask me anything..."

[modal]
backdrop = true
animation = "fade"
```

**Go API:**

```go
// config/config.go
package config

// UserConfig represents user-overridable settings
type UserConfig struct {
    Theme     ThemeConfig     `toml:"theme"`
    Mouse     MouseConfig     `toml:"mouse"`
    Keys      KeyConfig       `toml:"keybindings"`
    StatusBar StatusBarConfig `toml:"statusbar"`
    TabBar    TabBarConfig    `toml:"tabbar"`
    Input     InputConfig     `toml:"input"`
    Modal     ModalConfig     `toml:"modal"`
}

type ThemeConfig struct {
    Name   string            `toml:"name"`   // Base theme name
    Colors map[string]string `toml:"colors"` // Color overrides
    Styles map[string]any    `toml:"styles"` // Style overrides
}

type KeyConfig struct {
    Submit       []string `toml:"submit"`
    Cancel       []string `toml:"cancel"`
    NextTab      []string `toml:"next_tab"`
    PrevTab      []string `toml:"prev_tab"`
    Help         []string `toml:"help"`
    QuickActions []string `toml:"quick_actions"`
    ScrollUp     []string `toml:"scroll_up"`
    ScrollDown   []string `toml:"scroll_down"`
    // App-specific keys can be added
    Custom       map[string][]string `toml:"custom"`
}

// Load loads user config, merging with defaults
func Load(appName string, appDefaults UserConfig) (*UserConfig, error) {
    // Start with app defaults
    cfg := appDefaults

    // Find config file
    paths := DefaultConfigPaths(appName)
    if envPath := os.Getenv(strings.ToUpper(appName) + "_UI_CONFIG"); envPath != "" {
        paths = append([]string{envPath}, paths...)
    }

    // Load and merge first found config
    for _, path := range paths {
        if data, err := os.ReadFile(path); err == nil {
            var userCfg UserConfig
            if err := toml.Unmarshal(data, &userCfg); err != nil {
                return nil, fmt.Errorf("parse %s: %w", path, err)
            }
            cfg = merge(cfg, userCfg)
            break
        }
    }

    return &cfg, nil
}

// merge combines base config with overrides (overrides win)
func merge(base, override UserConfig) UserConfig
```

**Shell Integration:**

```go
// In your app's main.go (hex, jeff, or any other):
func main() {
    appName := "myapp"  // "hex", "jeff", etc.

    // App sets its defaults
    appDefaults := config.UserConfig{
        Theme: config.ThemeConfig{Name: "dracula"},
        Mouse: config.MouseConfig{
            Enabled:     true,
            ScrollLines: 3,
        },
        Keys: config.KeyConfig{
            Help: []string{"ctrl+h"},
        },
        StatusBar: config.StatusBarConfig{
            Order: []string{"model", "status", "tokens", "mode", "hints"},
        },
    }

    // Load user overrides from ~/.config/{appName}/ui.toml
    cfg, err := config.Load(appName, appDefaults)
    if err != nil {
        log.Warn("Failed to load user config: %v", err)
        cfg = &appDefaults
    }

    // Create shell with merged config
    shell := tui.NewShell(tui.ShellConfig{
        Theme:  cfg.BuildTheme(),
        Mouse:  cfg.Mouse,
        Keys:   cfg.BuildKeyBindings(),
        // ...
    })
}
```

**Theme Override Resolution:**

```go
// theme/override.go

// WithOverrides creates a new theme by applying overrides to a base theme
func WithOverrides(base Theme, overrides ThemeConfig) Theme {
    return &overrideTheme{
        base:      base,
        colors:    overrides.Colors,
        styles:    overrides.Styles,
    }
}

type overrideTheme struct {
    base    Theme
    colors  map[string]string
    styles  map[string]any
}

func (t *overrideTheme) Primary() lipgloss.Color {
    if c, ok := t.colors["primary"]; ok {
        return lipgloss.Color(c)
    }
    return t.base.Primary()
}

// ... similar for other colors

func (t *overrideTheme) Styles() ThemeStyles {
    base := t.base.Styles()
    // Apply style overrides from t.styles
    return applyStyleOverrides(base, t.styles)
}
```

**Hot Reload (Optional):**

```go
// config/watcher.go

type Watcher struct {
    path     string
    onChange func(UserConfig)
    done     chan struct{}
}

func NewWatcher(path string, onChange func(UserConfig)) *Watcher

func (w *Watcher) Start() error {
    // Use fsnotify to watch config file
    // On change, reload and call onChange
}

func (w *Watcher) Stop()

// Shell integration
func (s *Shell) EnableConfigHotReload(configPath string) {
    watcher := config.NewWatcher(configPath, func(cfg UserConfig) {
        s.ApplyConfig(cfg)
    })
    watcher.Start()
}
```

**Validation & Errors:**

```go
// config/validate.go

type ConfigError struct {
    Field   string
    Value   any
    Message string
}

func Validate(cfg UserConfig) []ConfigError {
    var errors []ConfigError

    // Validate theme exists
    if cfg.Theme.Name != "" && !ThemeExists(cfg.Theme.Name) {
        errors = append(errors, ConfigError{
            Field:   "theme.name",
            Value:   cfg.Theme.Name,
            Message: fmt.Sprintf("unknown theme %q, available: %v", cfg.Theme.Name, AvailableThemes()),
        })
    }

    // Validate colors are valid hex
    for name, color := range cfg.Theme.Colors {
        if !isValidHexColor(color) {
            errors = append(errors, ConfigError{
                Field:   "theme.colors." + name,
                Value:   color,
                Message: "invalid hex color, expected format: #RRGGBB",
            })
        }
    }

    // Validate key bindings
    for name, keys := range cfg.Keys.Custom {
        for _, k := range keys {
            if _, err := key.ParseBinding(k); err != nil {
                errors = append(errors, ConfigError{
                    Field:   "keybindings.custom." + name,
                    Value:   k,
                    Message: fmt.Sprintf("invalid key binding: %v", err),
                })
            }
        }
    }

    return errors
}
```

**Config Generation Helper:**

```go
// CLI commands for UI config management
// $ {appname} ui init       - Generate default UI config
// $ {appname} ui validate   - Validate UI config
// $ {appname} ui path       - Show config file location
//
// e.g., $ hex ui init, $ jeff ui validate

func GenerateDefaultConfig(appName string, defaults UserConfig) string {
    var buf bytes.Buffer
    buf.WriteString("# " + appName + " UI Configuration\n")
    buf.WriteString("# Generated with: " + appName + " ui init\n\n")

    encoder := toml.NewEncoder(&buf)
    encoder.Encode(defaults)

    return buf.String()
}

// UIInitCommand generates a default config file
// $ {appname} ui init [--output PATH]
func UIInitCommand(appName string, outputPath string, defaults UserConfig) error {
    if outputPath == "" {
        outputPath = DefaultConfigPaths(appName)[0]  // ~/.config/{appname}/ui.toml
    }

    content := GenerateDefaultConfig(appName, defaults)

    // Ensure directory exists
    if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
        return fmt.Errorf("create config directory: %w", err)
    }

    if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("write config: %w", err)
    }

    fmt.Printf("Created UI config at %s\n", outputPath)
    return nil
}

// UIValidateCommand validates the current UI config
// $ {appname} ui validate
func UIValidateCommand(appName string) error {
    cfg, err := Load(appName, UserConfig{})
    if err != nil {
        return fmt.Errorf("load config: %w", err)
    }

    errors := Validate(*cfg)
    if len(errors) > 0 {
        for _, e := range errors {
            fmt.Printf("  %s: %v - %s\n", e.Field, e.Value, e.Message)
        }
        return fmt.Errorf("%d validation errors", len(errors))
    }

    fmt.Println("UI config is valid!")
    return nil
}

// UIPathCommand shows where the config file is located
// $ {appname} ui path
func UIPathCommand(appName string) {
    paths := DefaultConfigPaths(appName)
    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            fmt.Printf("Active: %s\n", path)
            return
        }
    }
    fmt.Printf("Not found. Expected at: %s\n", paths[0])
}
```

**Documentation for Users:**

Each app should ship a `ui.example.toml` (or users can generate one via `{appname} ui init`) that documents all options:

```toml
# Example TUI Configuration
# Copy to ~/.config/{appname}/ui.toml
# e.g., ~/.config/hex/ui.toml or ~/.config/jeff/ui.toml

# =============================================================================
# THEME
# =============================================================================

[theme]
# Base theme: "dracula", "nord", "gruvbox", "high-contrast"
name = "dracula"

# Override individual colors (hex format)
[theme.colors]
# primary = "#bd93f9"     # Main accent (buttons, active elements)
# secondary = "#8be9fd"   # Secondary accent
# background = "#282a36"  # Main background
# foreground = "#f8f8f2"  # Main text
# success = "#50fa7b"     # Success states
# warning = "#ffb86c"     # Warning states
# error = "#ff5555"       # Error states
# muted = "#6272a4"       # Dimmed text

# =============================================================================
# MOUSE
# =============================================================================

[mouse]
enabled = true
scroll_lines = 3          # Lines per scroll wheel tick
hover_enabled = true      # Show timestamps on hover
shift_passthrough = true  # Allow Shift+click for text selection

# =============================================================================
# KEY BINDINGS
# =============================================================================

[keybindings]
# Multiple keys can be bound to the same action
# submit = ["enter"]
# cancel = ["esc"]
# help = ["ctrl+h", "?", "f1"]
# quick_actions = ["ctrl+k"]
# scroll_up = ["ctrl+u", "pgup"]
# scroll_down = ["ctrl+d", "pgdn"]

# =============================================================================
# STATUS BAR
# =============================================================================

[statusbar]
# Section order (remove to hide)
order = ["model", "status", "tokens", "mode", "progress", "hints"]

# =============================================================================
# TAB BAR
# =============================================================================

[tabbar]
position = "top"          # "top" or "bottom"
style = "underline"       # "underline", "boxed", or "pills"
show_badges = true
show_close = true

# =============================================================================
# INPUT
# =============================================================================

[input]
prefix = "> "
# placeholder = "Type a message..."
# multiline = false

# =============================================================================
# MODALS
# =============================================================================

[modal]
backdrop = true
# animation = "none"      # "none", "fade", "slide"
close_on_esc = true
```
