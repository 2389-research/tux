// form/select_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestSelectField(t *testing.T) {
	f := NewSelect[string]().
		WithID("theme").
		WithLabel("Choose Theme").
		WithOptions(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		)

	if f.ID() != "theme" {
		t.Errorf("expected ID 'theme', got %s", f.ID())
	}
}

func TestSelectFieldDefault(t *testing.T) {
	f := NewSelect[string]().
		WithID("theme").
		WithOptions(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		).
		WithDefault("nord")

	if f.Value().(string) != "nord" {
		t.Errorf("expected default 'nord', got %v", f.Value())
	}
}

func TestSelectFieldValue(t *testing.T) {
	f := NewSelect[string]().WithID("test").WithOptions(
		Option("A", "a"),
		Option("B", "b"),
	)

	f.SetValue("b")
	if f.Value().(string) != "b" {
		t.Errorf("expected 'b', got %v", f.Value())
	}
}

func TestSelectFieldRender(t *testing.T) {
	f := NewSelect[string]().
		WithID("test").
		WithLabel("Pick One").
		WithOptions(Option("A", "a"))

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}
