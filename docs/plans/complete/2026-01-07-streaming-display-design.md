# Streaming Display Design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add streaming display capabilities to tux for LLM orchestration apps (hex, jeff) - status indicators in statusbar and optional typewriter effect for content.

**Architecture:** Two components - a StreamingController for state management exposed via shell, and a StreamingContent wrapper for typewriter effects. Statusbar auto-displays streaming status when active.

**Tech Stack:** Go, Bubble Tea (internal), tux theme system

---

## Components

### 1. StreamingController

Accessed via `shell.Streaming()`. Apps call methods to update streaming state; shell handles display internally.

```go
// shell/streaming.go

type StreamingController struct {
    text           string
    tokenCount     int
    tokenRate      float64
    lastTokenTime  time.Time
    startTime      time.Time

    streaming      bool
    thinking       bool
    waiting        bool  // streaming started, no tokens yet

    toolCalls      []ToolCall

    // Spinner animation
    spinnerFrames  []string
    spinnerFrame   int
    lastSpinTime   time.Time
}

type ToolCall struct {
    ID         string
    Name       string
    InProgress bool
}

// Lifecycle
func (s *StreamingController) Start()
func (s *StreamingController) End()
func (s *StreamingController) Reset()

// Tokens
func (s *StreamingController) AppendToken(text string)
func (s *StreamingController) GetText() string
func (s *StreamingController) TokenRate() float64
func (s *StreamingController) TokenCount() int

// Status
func (s *StreamingController) SetThinking(active bool)
func (s *StreamingController) StartToolCall(id, name string)
func (s *StreamingController) EndToolCall(id string)

// State queries
func (s *StreamingController) IsStreaming() bool
func (s *StreamingController) IsThinking() bool
func (s *StreamingController) IsWaiting() bool
func (s *StreamingController) ActiveToolCalls() []ToolCall
```

### 2. StreamingContent Wrapper

Wraps any content type, adds typewriter effect with cursor.

```go
// shell/streaming_content.go

type StreamingContent struct {
    inner           content.Content
    typewriter      bool
    typewriterSpeed time.Duration  // Default: 30ms
    position        int
    text            string
}

func NewStreamingContent(inner content.Content) *StreamingContent

// Configuration
func (s *StreamingContent) WithTypewriter(enabled bool) *StreamingContent
func (s *StreamingContent) WithSpeed(d time.Duration) *StreamingContent

// Update text
func (s *StreamingContent) SetText(text string)

// content.Content interface
func (s *StreamingContent) Init() tea.Cmd
func (s *StreamingContent) Update(msg tea.Msg) (content.Content, tea.Cmd)
func (s *StreamingContent) View() string
func (s *StreamingContent) SetSize(width, height int)
func (s *StreamingContent) Value() any
```

### 3. Statusbar Auto-Display

Shell automatically shows streaming status in statusbar when streaming is active.

**Display states:**

| State | Display | Symbol |
|-------|---------|--------|
| Waiting | `Waiting...` | (none, italic) |
| Thinking | `⠋ Thinking` | Animated braille spinner |
| Tool call | `▍ Bash` | Box drawing character |
| Streaming | `▸ 42 tok/s` | Triangle + rate |

Multiple states combine: `⠋ Thinking  ▍ Bash  ▸ 42 tok/s`

Apps can disable via `shell.SetStreamingStatusVisible(false)`.

---

## Data Flow

```
App receives API event
    ↓
App calls shell.Streaming().AppendToken("hello")
    ↓
StreamingController updates internal state (text, token rate EMA)
    ↓
Shell's statusbar reads from controller on render
    ↓
Statusbar displays: ▸ 42 tok/s
```

For typewriter:
```
App calls streamingContent.SetText(shell.Streaming().GetText())
    ↓
StreamingContent stores text, position advances on tick
    ↓
View() renders text[:position] + cursor
```

---

## Theming

All symbols use theme colors (no color emoji):

| Element | Theme Color |
|---------|-------------|
| Spinner | `theme.Primary()` |
| Tool calls | `theme.Secondary()` |
| Token rate | `theme.Muted()` |
| Waiting text | `theme.Muted()` + italic |
| Typewriter cursor | `theme.Primary()` |

---

## App Usage Example

```go
type Model struct {
    shell *shell.Shell
    view  *shell.StreamingContent
}

func NewModel() Model {
    sh := shell.New()
    return Model{
        shell: sh,
        view:  shell.NewStreamingContent(myMarkdownView).WithTypewriter(true),
    }
}

func (m Model) HandleAPIEvent(event api.Event) {
    s := m.shell.Streaming()

    switch e := event.(type) {
    case api.StreamStart:
        s.Start()
    case api.Token:
        s.AppendToken(e.Text)
        m.view.SetText(s.GetText())
    case api.ThinkingStart:
        s.SetThinking(true)
    case api.ThinkingEnd:
        s.SetThinking(false)
    case api.ToolStart:
        s.StartToolCall(e.ID, e.Name)
    case api.ToolEnd:
        s.EndToolCall(e.ID)
    case api.StreamEnd:
        s.End()
    }
}
```

---

## Token Rate Calculation

Uses exponential moving average (EMA) for smooth display:

```go
func (s *StreamingController) AppendToken(text string) {
    now := time.Now()
    elapsed := now.Sub(s.lastTokenTime).Seconds()

    if elapsed > 0 {
        instantRate := 1.0 / elapsed
        if s.tokenRate == 0 {
            s.tokenRate = instantRate
        } else {
            // EMA with alpha = 0.3
            s.tokenRate = 0.3*instantRate + 0.7*s.tokenRate
        }
    }

    s.text += text
    s.tokenCount++
    s.lastTokenTime = now
    s.waiting = false
}
```

---

## Spinner Animation

Braille dots spinner (same as hex), 80ms frame interval:

```go
var defaultSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
```

Frame advances on each render if 80ms has passed since last advance.

---

## Out of Scope (YAGNI)

- Configurable spinner frames (can add later)
- Multiple concurrent streaming sessions
- Progressive markdown rendering hints (apps handle markdown)
- Typewriter speed per-character adjustment
