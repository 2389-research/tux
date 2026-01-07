// form/huh_theme_test.go
package shell

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestToHuhTheme(t *testing.T) {
	th := theme.Get("dracula")
	huhTheme := ToHuhTheme(th)

	if huhTheme == nil {
		t.Fatal("expected non-nil huh theme")
	}

	// Verify focused title uses primary color
	// (huh theme is configured, just verify it returns something)
}

func TestToHuhThemeAllThemes(t *testing.T) {
	themes := []string{"dracula", "nord", "gruvbox", "high-contrast", "neo-terminal"}

	for _, name := range themes {
		t.Run(name, func(t *testing.T) {
			th := theme.Get(name)
			huhTheme := ToHuhTheme(th)
			if huhTheme == nil {
				t.Errorf("ToHuhTheme returned nil for %s", name)
			}
		})
	}
}
