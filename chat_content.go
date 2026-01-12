// chat_content.go
package tux

import (
	"strings"
	"sync"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Compile-time check that ChatContent implements content.Content.
var _ content.Content = (*ChatContent)(nil)

// ChatContent displays the conversation in the Chat tab.
type ChatContent struct {
	mu             sync.Mutex
	theme          theme.Theme
	messages       []chatMessage
	current        strings.Builder // Current streaming message
	width          int
	height         int
	userStyle      lipgloss.Style
	assistantStyle lipgloss.Style

	// Viewport for scrolling
	viewport   viewport.Model
	autoScroll bool // Whether to auto-scroll to bottom on new content
	ready      bool // Whether viewport has been sized
}

type chatMessage struct {
	role    string // "user" or "assistant"
	content string
}

// NewChatContent creates a new ChatContent.
func NewChatContent(th theme.Theme) *ChatContent {
	if th == nil {
		panic("NewChatContent: nil theme")
	}
	return &ChatContent{
		theme:          th,
		messages:       make([]chatMessage, 0),
		userStyle:      lipgloss.NewStyle().Foreground(th.UserColor()),
		assistantStyle: lipgloss.NewStyle().Foreground(th.AssistantColor()),
		viewport:       viewport.New(0, 0),
		autoScroll:     true, // Auto-scroll by default
	}
}

// Init implements content.Content.
func (c *ChatContent) Init() tea.Cmd {
	return nil
}

// Update implements content.Content.
// Handles keyboard navigation for scrolling.
func (c *ChatContent) Update(msg tea.Msg) (content.Content, tea.Cmd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Track if user manually scrolled (disable auto-scroll)
		wasAtBottom := c.viewport.AtBottom()

		switch msg.String() {
		case "k", "up":
			c.viewport.LineUp(1)
			c.autoScroll = false
		case "j", "down":
			c.viewport.LineDown(1)
			if !c.viewport.AtBottom() {
				c.autoScroll = false
			}
		case "g":
			// For "gg", go to top (single g is handled, double-g needs state)
			// Simplified: just go to top on single g
			c.viewport.GotoTop()
			c.autoScroll = false
		case "G":
			c.viewport.GotoBottom()
			c.autoScroll = true
		case "ctrl+u":
			c.viewport.HalfViewUp()
			c.autoScroll = false
		case "ctrl+d":
			c.viewport.HalfViewDown()
			if !c.viewport.AtBottom() {
				c.autoScroll = false
			}
		case "pgup":
			c.viewport.ViewUp()
			c.autoScroll = false
		case "pgdown":
			c.viewport.ViewDown()
			if !c.viewport.AtBottom() {
				c.autoScroll = false
			}
		default:
			// Let viewport handle other keys (arrows, etc.)
			var cmd tea.Cmd
			c.viewport, cmd = c.viewport.Update(msg)
			// Check if user scrolled away from bottom
			if wasAtBottom && !c.viewport.AtBottom() {
				c.autoScroll = false
			} else if c.viewport.AtBottom() {
				c.autoScroll = true
			}
			return c, cmd
		}
	}

	return c, nil
}

// View implements content.Content.
func (c *ChatContent) View() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Before viewport is sized, fallback to direct render
	if !c.ready {
		return c.renderContent()
	}

	return c.viewport.View()
}

// renderContent builds the styled content string from messages.
// Must be called with mutex held.
func (c *ChatContent) renderContent() string {
	var parts []string

	for _, msg := range c.messages {
		var style lipgloss.Style
		if msg.role == "user" {
			style = c.userStyle
		} else {
			style = c.assistantStyle
		}
		parts = append(parts, style.Render(msg.content))
	}

	// Add current streaming message if any
	if c.current.Len() > 0 {
		parts = append(parts, c.assistantStyle.Render(c.current.String()))
	}

	return strings.Join(parts, "\n\n")
}

// updateViewport rebuilds the viewport content and optionally scrolls to bottom.
// Must be called with mutex held.
func (c *ChatContent) updateViewport() {
	if !c.ready {
		return
	}
	content := c.renderContent()
	c.viewport.SetContent(content)
	if c.autoScroll {
		c.viewport.GotoBottom()
	}
}

// Value implements content.Content.
func (c *ChatContent) Value() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Return a copy to prevent callers from mutating internal slice
	result := make([]chatMessage, len(c.messages))
	copy(result, c.messages)
	return result
}

// SetSize implements content.Content.
func (c *ChatContent) SetSize(width, height int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.width = width
	c.height = height
	c.viewport.Width = width
	c.viewport.Height = height
	if !c.ready {
		c.ready = true
		c.updateViewport()
	}
}

// AppendText appends streaming text to the current assistant message.
func (c *ChatContent) AppendText(text string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current.WriteString(text)
	c.updateViewport()
}

// AddUserMessage adds a user message to the conversation.
func (c *ChatContent) AddUserMessage(content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, chatMessage{
		role:    "user",
		content: content,
	})
	c.updateViewport()
}

// FinishAssistantMessage completes the current streaming message.
func (c *ChatContent) FinishAssistantMessage() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.current.Len() > 0 {
		c.messages = append(c.messages, chatMessage{
			role:    "assistant",
			content: c.current.String(),
		})
		c.current.Reset()
		c.updateViewport()
	}
}

// UserMessages returns all user message contents in order (oldest to newest).
func (c *ChatContent) UserMessages() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var result []string
	for _, msg := range c.messages {
		if msg.role == "user" {
			result = append(result, msg.content)
		}
	}
	return result
}

// Clear removes all messages and resets the current streaming state.
func (c *ChatContent) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = make([]chatMessage, 0)
	c.current.Reset()
	c.autoScroll = true // Reset auto-scroll on clear
	c.updateViewport()
}

// AddAssistantMessage adds a completed assistant message to the conversation.
// Use this when restoring a previous session, not for streaming.
func (c *ChatContent) AddAssistantMessage(content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, chatMessage{
		role:    "assistant",
		content: content,
	})
	c.updateViewport()
}
