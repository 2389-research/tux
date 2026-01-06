package theme

import "testing"

func TestDraculaTheme(t *testing.T) {
	theme := NewDraculaTheme()

	if theme.Name() != "dracula" {
		t.Errorf("expected name 'dracula', got %s", theme.Name())
	}

	// Verify colors are set
	if theme.Primary() == "" {
		t.Error("primary color not set")
	}
	if theme.Background() == "" {
		t.Error("background color not set")
	}
	if theme.Foreground() == "" {
		t.Error("foreground color not set")
	}

	// Verify styles are composed
	styles := theme.Styles()

	// Title style should be bold
	if !styles.Title.GetBold() {
		t.Error("title style should be bold")
	}

	// Verify other styles exist
	if styles.Input.GetPaddingLeft()+styles.Input.GetPaddingRight() == 0 {
		t.Error("input style should have padding")
	}
}

func TestGet(t *testing.T) {
	// Known theme
	theme := Get("dracula")
	if theme.Name() != "dracula" {
		t.Errorf("expected dracula theme, got %s", theme.Name())
	}

	// Unknown theme defaults to dracula
	theme = Get("nonexistent")
	if theme.Name() != "dracula" {
		t.Errorf("expected default dracula theme, got %s", theme.Name())
	}
}

func TestAvailable(t *testing.T) {
	themes := Available()
	if len(themes) == 0 {
		t.Error("expected at least one theme")
	}

	found := false
	for _, name := range themes {
		if name == "dracula" {
			found = true
			break
		}
	}
	if !found {
		t.Error("dracula theme not in available themes")
	}
}

func TestRegister(t *testing.T) {
	// Register a custom theme
	Register("test-theme", func() Theme {
		return NewDraculaTheme() // reuse dracula for simplicity
	})

	themes := Available()
	found := false
	for _, name := range themes {
		if name == "test-theme" {
			found = true
			break
		}
	}
	if !found {
		t.Error("registered theme not in available themes")
	}
}
