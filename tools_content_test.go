// tools_content_test.go
package tux

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNewToolsContent(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)
	if tools == nil {
		t.Error("NewToolsContent should return a ToolsContent")
	}
}

func TestToolsContentAddCall(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})

	view := tools.View()
	if !strings.Contains(view, "read_file") {
		t.Errorf("View should contain 'read_file', got: %s", view)
	}
}

func TestToolsContentAddResult(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})
	tools.AddToolResult("tool-1", "file contents here", true)

	view := tools.View()
	if !strings.Contains(view, "✓") {
		t.Errorf("View should contain success marker, got: %s", view)
	}
}

func TestToolsContentAddResultFailure(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})
	tools.AddToolResult("tool-1", "error: file not found", false)

	view := tools.View()
	if !strings.Contains(view, "✗") {
		t.Errorf("View should contain failure marker, got: %s", view)
	}
}

func TestToolsContentPendingStatus(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/tmp/test.txt"})

	view := tools.View()
	if !strings.Contains(view, "⋯") {
		t.Errorf("View should contain pending marker, got: %s", view)
	}
}

func TestToolsContentValue(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/test"})
	tools.AddToolResult("tool-1", "contents", true)

	items, ok := tools.Value().([]toolItem)
	if !ok {
		t.Error("Value() should return []toolItem")
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
	if !items[0].completed {
		t.Error("Item should be completed after AddToolResult")
	}
}

func TestToolsContentEmptyView(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	view := tools.View()
	if !strings.Contains(view, "No tool calls yet") {
		t.Errorf("Empty view should show 'No tool calls yet', got: %s", view)
	}
}

func TestToolsContentMultipleCalls(t *testing.T) {
	th := theme.NewDraculaTheme()
	tools := NewToolsContent(th)

	tools.AddToolCall("tool-1", "read_file", map[string]any{"path": "/a"})
	tools.AddToolCall("tool-2", "write_file", map[string]any{"path": "/b"})
	tools.AddToolResult("tool-1", "ok", true)

	view := tools.View()
	if !strings.Contains(view, "read_file") || !strings.Contains(view, "write_file") {
		t.Errorf("View should contain both tools, got: %s", view)
	}
	if !strings.Contains(view, "✓") {
		t.Errorf("View should contain success marker for completed tool, got: %s", view)
	}
	if !strings.Contains(view, "⋯") {
		t.Errorf("View should contain pending marker for pending tool, got: %s", view)
	}
}
