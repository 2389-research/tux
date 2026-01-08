// chat_content.go
package tux

import (
	"strings"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatContent displays the conversation in the Chat tab.
type ChatContent struct {
	theme    theme.Theme
	messages []chatMessage
	current  strings.Builder // Current streaming message
	width    int
	height   int
}

type chatMessage struct {
	role    string // "user" or "assistant"
	content string
}

// NewChatContent creates a new ChatContent.
func NewChatContent(th theme.Theme) *ChatContent {
	return &ChatContent{
		theme:    th,
		messages: make([]chatMessage, 0),
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
	var parts []string

	userStyle := lipgloss.NewStyle().Foreground(c.theme.UserColor())
	assistantStyle := lipgloss.NewStyle().Foreground(c.theme.AssistantColor())

	for _, msg := range c.messages {
		var style lipgloss.Style
		if msg.role == "user" {
			style = userStyle
		} else {
			style = assistantStyle
		}
		parts = append(parts, style.Render(msg.content))
	}

	// Add current streaming message if any
	if c.current.Len() > 0 {
		parts = append(parts, assistantStyle.Render(c.current.String()))
	}

	return strings.Join(parts, "\n\n")
}

// Value implements content.Content.
func (c *ChatContent) Value() any {
	return c.messages
}

// SetSize implements content.Content.
func (c *ChatContent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// AppendText appends streaming text to the current assistant message.
func (c *ChatContent) AppendText(text string) {
	c.current.WriteString(text)
}

// AddUserMessage adds a user message to the conversation.
func (c *ChatContent) AddUserMessage(content string) {
	c.messages = append(c.messages, chatMessage{
		role:    "user",
		content: content,
	})
}

// FinishAssistantMessage completes the current streaming message.
func (c *ChatContent) FinishAssistantMessage() {
	if c.current.Len() > 0 {
		c.messages = append(c.messages, chatMessage{
			role:    "assistant",
			content: c.current.String(),
		})
		c.current.Reset()
	}
}
