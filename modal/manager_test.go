package modal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// testModal is a simple modal for testing.
type testModal struct {
	id       string
	title    string
	size     Size
	pushed   bool
	popped   bool
	rendered bool
}

func newTestModal(id string) *testModal {
	return &testModal{
		id:    id,
		title: "Test " + id,
		size:  SizeMedium,
	}
}

func (m *testModal) ID() string              { return m.id }
func (m *testModal) Title() string           { return m.title }
func (m *testModal) Size() Size              { return m.size }
func (m *testModal) OnPush(width, height int) { m.pushed = true }
func (m *testModal) OnPop()                  { m.popped = true }

func (m *testModal) Render(width, height int) string {
	m.rendered = true
	return "Modal: " + m.id
}

func (m *testModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	if key.String() == "enter" {
		return true, nil
	}
	return false, nil
}

func TestManagerPushPop(t *testing.T) {
	mgr := NewManager()

	if mgr.HasActive() {
		t.Error("new manager should not have active modal")
	}
	if mgr.Count() != 0 {
		t.Errorf("expected count 0, got %d", mgr.Count())
	}

	// Push first modal
	modal1 := newTestModal("modal1")
	mgr.Push(modal1)

	if !mgr.HasActive() {
		t.Error("manager should have active modal after push")
	}
	if mgr.Count() != 1 {
		t.Errorf("expected count 1, got %d", mgr.Count())
	}
	if !modal1.pushed {
		t.Error("OnPush should have been called")
	}

	// Push second modal
	modal2 := newTestModal("modal2")
	mgr.Push(modal2)

	if mgr.Count() != 2 {
		t.Errorf("expected count 2, got %d", mgr.Count())
	}
	if mgr.Peek().ID() != "modal2" {
		t.Errorf("expected peek to return modal2, got %s", mgr.Peek().ID())
	}

	// Pop second modal
	popped := mgr.Pop()
	if popped.ID() != "modal2" {
		t.Errorf("expected to pop modal2, got %s", popped.ID())
	}
	if !modal2.popped {
		t.Error("OnPop should have been called")
	}
	if mgr.Count() != 1 {
		t.Errorf("expected count 1 after pop, got %d", mgr.Count())
	}

	// Pop first modal
	popped = mgr.Pop()
	if popped.ID() != "modal1" {
		t.Errorf("expected to pop modal1, got %s", popped.ID())
	}
	if mgr.HasActive() {
		t.Error("manager should not have active modal after all pops")
	}

	// Pop empty stack
	popped = mgr.Pop()
	if popped != nil {
		t.Error("popping empty stack should return nil")
	}
}

func TestManagerVersion(t *testing.T) {
	mgr := NewManager()
	initialVersion := mgr.Version()

	modal := newTestModal("test")
	mgr.Push(modal)

	if mgr.Version() != initialVersion+1 {
		t.Errorf("version should increment on push")
	}

	v := mgr.Version()
	mgr.Pop()

	if mgr.Version() != v+1 {
		t.Errorf("version should increment on pop")
	}
}

func TestManagerClear(t *testing.T) {
	mgr := NewManager()

	mgr.Push(newTestModal("m1"))
	mgr.Push(newTestModal("m2"))
	mgr.Push(newTestModal("m3"))

	if mgr.Count() != 3 {
		t.Errorf("expected count 3, got %d", mgr.Count())
	}

	mgr.Clear()

	if mgr.Count() != 0 {
		t.Errorf("expected count 0 after clear, got %d", mgr.Count())
	}
	if mgr.HasActive() {
		t.Error("should not have active modal after clear")
	}
}

func TestManagerHandleKey(t *testing.T) {
	mgr := NewManager()

	// No modal - should not handle
	handled, _ := mgr.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if handled {
		t.Error("should not handle key when no active modal")
	}

	modal := newTestModal("test")
	mgr.Push(modal)

	// Enter should be handled
	handled, _ = mgr.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if !handled {
		t.Error("enter key should be handled")
	}

	// Other keys should not be handled
	handled, _ = mgr.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if handled {
		t.Error("tab key should not be handled")
	}
}

func TestManagerRender(t *testing.T) {
	mgr := NewManager()
	mgr.SetSize(80, 24)

	// No modal - empty string
	output := mgr.Render(80, 24)
	if output != "" {
		t.Error("render should return empty string when no modal")
	}

	modal := newTestModal("test")
	mgr.Push(modal)

	output = mgr.Render(80, 24)
	if output == "" {
		t.Error("render should return content when modal active")
	}
	if !modal.rendered {
		t.Error("modal Render should have been called")
	}
}

func TestSizePercent(t *testing.T) {
	tests := []struct {
		size          Size
		heightPercent float64
		widthPercent  float64
	}{
		{SizeSmall, 0.30, 0.50},
		{SizeMedium, 0.50, 0.60},
		{SizeLarge, 0.80, 0.80},
		{SizeFullscreen, 1.0, 1.0},
	}

	for _, tt := range tests {
		if tt.size.HeightPercent() != tt.heightPercent {
			t.Errorf("size %d: expected height %v, got %v", tt.size, tt.heightPercent, tt.size.HeightPercent())
		}
		if tt.size.WidthPercent() != tt.widthPercent {
			t.Errorf("size %d: expected width %v, got %v", tt.size, tt.widthPercent, tt.size.WidthPercent())
		}
	}
}
