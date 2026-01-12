// chat_content.go
package tux

import (
	"strings"
	"sync"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
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
	}
}

// Init implements content.Content.
func (c *ChatContent) Init() tea.Cmd {
	return nil
}

// Update implements content.Content.
func (c *ChatContent) Update(msg tea.Msg) (content.Content, tea.Cmd) {
	return c, nil
}

// View implements content.Content.
func (c *ChatContent) View() string {
	c.mu.Lock()
	defer c.mu.Unlock()

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
}

// AppendText appends streaming text to the current assistant message.
func (c *ChatContent) AppendText(text string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current.WriteString(text)
}

// AddUserMessage adds a user message to the conversation.
func (c *ChatContent) AddUserMessage(content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, chatMessage{
		role:    "user",
		content: content,
	})
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
}
