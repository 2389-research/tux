// form/confirm_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestConfirmField(t *testing.T) {
	f := NewConfirm().
		WithID("delete").
		WithLabel("Delete this file?").
		WithAffirmative("Yes, delete").
		WithNegative("Cancel")

	if f.ID() != "delete" {
		t.Errorf("expected ID 'delete', got %s", f.ID())
	}
}

func TestConfirmFieldValue(t *testing.T) {
	f := NewConfirm().WithID("test")

	// Default should be false
	if f.Value().(bool) != false {
		t.Error("default should be false")
	}

	f.SetValue(true)
	if f.Value().(bool) != true {
		t.Error("should be true after SetValue(true)")
	}
}

func TestConfirmFieldRender(t *testing.T) {
	f := NewConfirm().
		WithID("test").
		WithLabel("Confirm?")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}
