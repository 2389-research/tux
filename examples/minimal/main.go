// Package main demonstrates a minimal tux application.
package main

import (
	"log"

	"github.com/2389-research/tux/content"
	"github.com/2389-research/tux/shell"
	"github.com/2389-research/tux/theme"
)

func main() {
	// Create shell with dracula theme
	th := theme.NewDraculaTheme()
	cfg := shell.DefaultConfig()
	cfg.InputPrefix = "λ "
	cfg.InputPlaceholder = "Type a message..."

	s := shell.New(th, cfg)

	// Add a chat tab with some content
	chatContent := content.NewSelectList([]content.SelectItem{
		{Label: "Welcome to tux!", Description: "A multi-agent terminal interface library"},
		{Label: "Use ↑↓ to navigate", Description: "Or j/k for vim-style navigation"},
		{Label: "Press Ctrl+C to quit", Description: "Or Ctrl+Q"},
	})

	s.AddTab(shell.Tab{
		ID:      "chat",
		Label:   "Chat",
		Content: chatContent,
	})

	// Add an activity tab
	activityContent := content.NewSelectList([]content.SelectItem{
		{Label: "No activity yet", Description: "Tool calls will appear here"},
	})

	s.AddTab(shell.Tab{
		ID:      "activity",
		Label:   "Activity",
		Content: activityContent,
	})

	// Set initial status
	s.SetStatus(shell.Status{
		Model:      "claude-3.5",
		Connected:  true,
		TokensUsed: 1250,
		TokensMax:  100000,
		Hints:      "Tab: switch tabs │ Ctrl+C: quit",
	})

	// Run the shell
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
