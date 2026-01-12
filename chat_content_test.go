// chat_content_test.go
package tux

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNewChatContent(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)
	if chat == nil {
		t.Error("NewChatContent should return a ChatContent")
	}
}

func TestChatContentAppendText(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AppendText("Hello ")
	chat.AppendText("World")

	view := chat.View()
	if !strings.Contains(view, "Hello") || !strings.Contains(view, "World") {
		t.Errorf("View should contain 'Hello' and 'World', got: %s", view)
	}
}

func TestChatContentAddUserMessage(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AddUserMessage("What is 2+2?")

	view := chat.View()
	if !strings.Contains(view, "What is 2+2?") {
		t.Errorf("View should contain 'What is 2+2?', got: %s", view)
	}
}

func TestChatContentFinishAssistant(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	chat.AppendText("The answer is 4")
	chat.FinishAssistantMessage()

	// Verify message was added to messages list
	messages, ok := chat.Value().([]chatMessage)
	if !ok {
		t.Error("Value() should return []chatMessage")
		return
	}
	if len(messages) != 1 {
		t.Errorf("Expected 1 message after FinishAssistantMessage, got %d", len(messages))
	}
	if len(messages) > 0 && messages[0].role != "assistant" {
		t.Errorf("Expected assistant role, got %s", messages[0].role)
	}

	// Should be able to append new message (current was reset)
	chat.AppendText("New response")
}

func TestChatContentValue(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	// Add a user message
	chat.AddUserMessage("Hello")

	// Add an assistant message
	chat.AppendText("Hi there")
	chat.FinishAssistantMessage()

	messages, ok := chat.Value().([]chatMessage)
	if !ok {
		t.Error("Value() should return []chatMessage")
		return
	}
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	// Verify user message
	if len(messages) > 0 {
		if messages[0].role != "user" {
			t.Errorf("Expected first message role 'user', got '%s'", messages[0].role)
		}
		if messages[0].content != "Hello" {
			t.Errorf("Expected first message content 'Hello', got '%s'", messages[0].content)
		}
	}

	// Verify assistant message
	if len(messages) > 1 {
		if messages[1].role != "assistant" {
			t.Errorf("Expected second message role 'assistant', got '%s'", messages[1].role)
		}
		if messages[1].content != "Hi there" {
			t.Errorf("Expected second message content 'Hi there', got '%s'", messages[1].content)
		}
	}
}

func TestChatContentViewport(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	// Set size to initialize viewport
	chat.SetSize(80, 10)

	// Add some content
	chat.AddUserMessage("Hello")
	chat.AppendText("Response")
	chat.FinishAssistantMessage()

	// View should work through viewport
	view := chat.View()
	if !strings.Contains(view, "Hello") {
		t.Errorf("View should contain 'Hello', got: %s", view)
	}
	if !strings.Contains(view, "Response") {
		t.Errorf("View should contain 'Response', got: %s", view)
	}
}

func TestChatContentAutoScroll(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	// Set size to initialize viewport
	chat.SetSize(80, 5)

	// Add many messages to exceed viewport height
	for i := 0; i < 20; i++ {
		chat.AddUserMessage("Message " + string(rune('A'+i)))
	}

	// autoScroll should be true by default
	chat.mu.Lock()
	autoScroll := chat.autoScroll
	chat.mu.Unlock()

	if !autoScroll {
		t.Error("autoScroll should be true by default")
	}
}

func TestChatContentClearResetsAutoScroll(t *testing.T) {
	th := theme.NewDraculaTheme()
	chat := NewChatContent(th)

	// Set size and disable auto-scroll
	chat.SetSize(80, 10)
	chat.mu.Lock()
	chat.autoScroll = false
	chat.mu.Unlock()

	// Clear should reset auto-scroll
	chat.Clear()

	chat.mu.Lock()
	autoScroll := chat.autoScroll
	chat.mu.Unlock()

	if !autoScroll {
		t.Error("Clear should reset autoScroll to true")
	}
}
