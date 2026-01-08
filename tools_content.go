// tools_content.go
package tux

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// Compile-time check that ToolsContent implements content.Content.
var _ content.Content = (*ToolsContent)(nil)

// ToolsContent displays the tool call timeline in the Tools tab.
type ToolsContent struct {
	mu     sync.Mutex
	theme  theme.Theme
	items  []toolItem
	width  int
	height int
}

type toolItem struct {
	id        string
	name      string
	params    map[string]any
	output    string
	success   bool
	completed bool
	timestamp time.Time
}

// NewToolsContent creates a new ToolsContent.
func NewToolsContent(th theme.Theme) *ToolsContent {
	return &ToolsContent{
		theme: th,
		items: make([]toolItem, 0),
	}
}

// Init implements content.Content.
func (c *ToolsContent) Init() tea.Cmd {
	return nil
}

// Update implements content.Content.
func (c *ToolsContent) Update(msg tea.Msg) (content.Content, tea.Cmd) {
	return c, nil
}

// View implements content.Content.
func (c *ToolsContent) View() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) == 0 {
		return "No tool calls yet"
	}

	var parts []string

	for _, item := range c.items {
		var status string
		if item.completed {
			if item.success {
				status = "✓"
			} else {
				status = "✗"
			}
		} else {
			status = "⋯"
		}

		line := fmt.Sprintf("%s %s", status, item.name)
		if item.completed && item.output != "" {
			// Truncate output for display
			output := item.output
			if len(output) > 50 {
				output = output[:47] + "..."
			}
			line += fmt.Sprintf(" → %s", output)
		}

		parts = append(parts, line)
	}

	return strings.Join(parts, "\n")
}

// Value implements content.Content.
func (c *ToolsContent) Value() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.items
}

// SetSize implements content.Content.
func (c *ToolsContent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// AddToolCall adds a tool call to the timeline.
func (c *ToolsContent) AddToolCall(id, name string, params map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = append(c.items, toolItem{
		id:        id,
		name:      name,
		params:    params,
		timestamp: time.Now(),
	})
}

// AddToolResult adds a result to an existing tool call.
func (c *ToolsContent) AddToolResult(id, output string, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.items {
		if c.items[i].id == id {
			c.items[i].output = output
			c.items[i].success = success
			c.items[i].completed = true
			return
		}
	}
}
