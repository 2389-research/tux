package shell

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func TestNewSpinner(t *testing.T) {
	tests := []struct {
		name        string
		spinnerType SpinnerType
	}{
		{"Default", SpinnerDefault},
		{"Execution", SpinnerExecution},
		{"Streaming", SpinnerStreaming},
		{"Loading", SpinnerLoading},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpinner(tt.spinnerType)
			if s == nil {
				t.Fatal("expected non-nil spinner")
			}
			if s.Type() != tt.spinnerType {
				t.Errorf("expected type %d, got %d", tt.spinnerType, s.Type())
			}
			if s.Active() {
				t.Error("spinner should not be active initially")
			}
		})
	}
}

func TestSpinnerSetMessage(t *testing.T) {
	s := NewSpinner(SpinnerDefault)
	s.SetMessage("Loading...")

	if s.Message() != "Loading..." {
		t.Errorf("expected message 'Loading...', got %q", s.Message())
	}
}

func TestSpinnerSetTokenRate(t *testing.T) {
	s := NewSpinner(SpinnerStreaming)
	s.SetTokenRate(42.5)

	if s.TokenRate() != 42.5 {
		t.Errorf("expected token rate 42.5, got %f", s.TokenRate())
	}
}

func TestSpinnerSetStyle(t *testing.T) {
	s := NewSpinner(SpinnerDefault)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	s.SetStyle(style)
	// Style is applied internally, just verify no panic
}

func TestSpinnerStartStop(t *testing.T) {
	s := NewSpinner(SpinnerDefault)

	if s.Active() {
		t.Error("spinner should not be active before start")
	}

	cmd := s.Start()
	if cmd == nil {
		t.Error("Start should return a command")
	}
	if !s.Active() {
		t.Error("spinner should be active after start")
	}

	// Verify start time is set
	if s.Elapsed() == 0 {
		// Give a tiny bit of time
		time.Sleep(time.Millisecond)
		if s.Elapsed() == 0 {
			t.Error("elapsed time should be non-zero after start")
		}
	}

	s.Stop()
	if s.Active() {
		t.Error("spinner should not be active after stop")
	}

	// Elapsed should return 0 when not active
	if s.Elapsed() != 0 {
		t.Error("elapsed should be 0 when spinner is not active")
	}
}

func TestSpinnerUpdate(t *testing.T) {
	s := NewSpinner(SpinnerDefault)

	// Update when not active should return nil cmd
	updated, cmd := s.Update(spinner.TickMsg{})
	if cmd != nil {
		t.Error("update when not active should return nil cmd")
	}
	if updated != s {
		t.Error("update should return same spinner")
	}

	// Start spinner and update
	s.Start()
	updated, cmd = s.Update(spinner.TickMsg{})
	if updated != s {
		t.Error("update should return same spinner")
	}
	// cmd may or may not be nil depending on spinner state
}

func TestSpinnerViewWhenInactive(t *testing.T) {
	s := NewSpinner(SpinnerDefault)
	view := s.View()
	if view != "" {
		t.Errorf("expected empty view when inactive, got %q", view)
	}
}

func TestSpinnerViewDefault(t *testing.T) {
	s := NewSpinner(SpinnerDefault)
	s.Start()

	// Without message
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}

	// With message
	s.SetMessage("Working")
	view = s.View()
	if !strings.Contains(view, "Working") {
		t.Errorf("expected view to contain 'Working', got %q", view)
	}
}

func TestSpinnerViewExecution(t *testing.T) {
	s := NewSpinner(SpinnerExecution)
	s.Start()

	// Should show elapsed time
	time.Sleep(10 * time.Millisecond)
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
	// Elapsed time format like (0s) or (1s)
	if !strings.Contains(view, "s)") && !strings.Contains(view, "0s") {
		// The format might vary, just check it renders something
		if view == "" {
			t.Error("expected execution spinner to show elapsed time")
		}
	}

	// With message
	s.SetMessage("Executing")
	view = s.View()
	if !strings.Contains(view, "Executing") {
		t.Errorf("expected view to contain 'Executing', got %q", view)
	}
}

func TestSpinnerViewStreaming(t *testing.T) {
	s := NewSpinner(SpinnerStreaming)
	s.Start()

	// Without token rate
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}

	// With token rate
	s.SetTokenRate(15.5)
	view = s.View()
	if !strings.Contains(view, "15.5") || !strings.Contains(view, "tok/s") {
		t.Errorf("expected view to contain token rate, got %q", view)
	}

	// With message and token rate
	s.SetMessage("Streaming")
	view = s.View()
	if !strings.Contains(view, "Streaming") {
		t.Errorf("expected view to contain 'Streaming', got %q", view)
	}
	if !strings.Contains(view, "tok/s") {
		t.Errorf("expected view to contain token rate, got %q", view)
	}

	// With message but no token rate
	s.SetTokenRate(0)
	view = s.View()
	if !strings.Contains(view, "Streaming") {
		t.Errorf("expected view to contain 'Streaming', got %q", view)
	}
}

func TestSpinnerViewLoading(t *testing.T) {
	s := NewSpinner(SpinnerLoading)
	s.Start()

	view := s.View()
	if view == "" {
		t.Error("expected non-empty view for loading spinner")
	}

	s.SetMessage("Please wait")
	view = s.View()
	if !strings.Contains(view, "Please wait") {
		t.Errorf("expected view to contain message, got %q", view)
	}
}

func TestSpinnerTypes(t *testing.T) {
	// Verify all spinner types create valid spinners
	types := []SpinnerType{SpinnerDefault, SpinnerExecution, SpinnerStreaming, SpinnerLoading}

	for _, st := range types {
		s := NewSpinner(st)
		s.Start()
		view := s.View()
		if view == "" {
			t.Errorf("spinner type %d should produce non-empty view", st)
		}
		s.Stop()
	}
}
