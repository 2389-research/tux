package shell

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestStatusBarErrorIndicator(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	// Set an error
	sb.SetError("connection timeout", 1)

	view := sb.View(80)
	if !strings.Contains(view, "⚠") {
		t.Error("status bar should show error indicator")
	}
	if !strings.Contains(view, "connection...") {
		t.Error("status bar should show truncated error message")
	}
}

func TestStatusBarNoError(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	view := sb.View(80)
	if strings.Contains(view, "⚠") {
		t.Error("status bar should not show error indicator when no errors")
	}
}

func TestStatusBarClearError(t *testing.T) {
	th := theme.NewDraculaTheme()
	sb := NewStatusBar(th)

	sb.SetError("some error", 1)
	sb.ClearError()

	view := sb.View(80)
	if strings.Contains(view, "⚠") {
		t.Error("status bar should not show error after clearing")
	}
}
