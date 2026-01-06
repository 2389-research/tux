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

func TestDraculaAllColors(t *testing.T) {
	theme := NewDraculaTheme()

	// Test all color getters return non-empty values
	tests := []struct {
		name  string
		color func() string
	}{
		{"Background", func() string { return string(theme.Background()) }},
		{"Foreground", func() string { return string(theme.Foreground()) }},
		{"Primary", func() string { return string(theme.Primary()) }},
		{"Secondary", func() string { return string(theme.Secondary()) }},
		{"Success", func() string { return string(theme.Success()) }},
		{"Warning", func() string { return string(theme.Warning()) }},
		{"Error", func() string { return string(theme.Error()) }},
		{"Info", func() string { return string(theme.Info()) }},
		{"Border", func() string { return string(theme.Border()) }},
		{"BorderFocused", func() string { return string(theme.BorderFocused()) }},
		{"Muted", func() string { return string(theme.Muted()) }},
		{"UserColor", func() string { return string(theme.UserColor()) }},
		{"AssistantColor", func() string { return string(theme.AssistantColor()) }},
		{"ToolColor", func() string { return string(theme.ToolColor()) }},
		{"SystemColor", func() string { return string(theme.SystemColor()) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := tt.color()
			if color == "" {
				t.Errorf("%s color should not be empty", tt.name)
			}
		})
	}
}

func TestDraculaColorValues(t *testing.T) {
	theme := NewDraculaTheme()

	// Verify specific Dracula colors
	if theme.Background() != "#282a36" {
		t.Errorf("expected background #282a36, got %s", theme.Background())
	}
	if theme.Primary() != "#bd93f9" {
		t.Errorf("expected primary #bd93f9 (purple), got %s", theme.Primary())
	}
	if theme.Success() != "#50fa7b" {
		t.Errorf("expected success #50fa7b (green), got %s", theme.Success())
	}
	if theme.Error() != "#ff5555" {
		t.Errorf("expected error #ff5555 (red), got %s", theme.Error())
	}
}

func TestNordTheme(t *testing.T) {
	theme := NewNordTheme()

	if theme.Name() != "nord" {
		t.Errorf("expected name 'nord', got %s", theme.Name())
	}

	// Verify Nord colors
	if theme.Background() != "#2e3440" {
		t.Errorf("expected background #2e3440, got %s", theme.Background())
	}
	if theme.Primary() != "#88c0d0" {
		t.Errorf("expected primary #88c0d0, got %s", theme.Primary())
	}

	// Verify styles
	styles := theme.Styles()
	if !styles.Title.GetBold() {
		t.Error("title style should be bold")
	}
}

func TestGruvboxTheme(t *testing.T) {
	theme := NewGruvboxTheme()

	if theme.Name() != "gruvbox" {
		t.Errorf("expected name 'gruvbox', got %s", theme.Name())
	}

	// Verify Gruvbox colors
	if theme.Background() != "#282828" {
		t.Errorf("expected background #282828, got %s", theme.Background())
	}
	if theme.Primary() != "#fabd2f" {
		t.Errorf("expected primary #fabd2f, got %s", theme.Primary())
	}

	// Verify styles
	styles := theme.Styles()
	if !styles.Title.GetBold() {
		t.Error("title style should be bold")
	}
}

func TestHighContrastTheme(t *testing.T) {
	theme := NewHighContrastTheme()

	if theme.Name() != "high-contrast" {
		t.Errorf("expected name 'high-contrast', got %s", theme.Name())
	}

	// Verify high contrast colors
	if theme.Background() != "#000000" {
		t.Errorf("expected background #000000, got %s", theme.Background())
	}
	if theme.Foreground() != "#ffffff" {
		t.Errorf("expected foreground #ffffff, got %s", theme.Foreground())
	}
	if theme.Primary() != "#ffff00" {
		t.Errorf("expected primary #ffff00, got %s", theme.Primary())
	}

	// Verify styles
	styles := theme.Styles()
	if !styles.Title.GetBold() {
		t.Error("title style should be bold")
	}
}

func TestGetAllThemes(t *testing.T) {
	themes := []string{"dracula", "nord", "gruvbox", "high-contrast", "neo-terminal"}

	for _, name := range themes {
		theme := Get(name)
		if theme.Name() != name {
			t.Errorf("Get(%q) returned theme with name %q", name, theme.Name())
		}
	}
}

func TestNeoTerminalTheme(t *testing.T) {
	theme := NewNeoTerminalTheme()

	if theme.Name() != "neo-terminal" {
		t.Errorf("expected name 'neo-terminal', got %s", theme.Name())
	}

	// Verify Neo-Terminal colors
	if theme.Background() != "#1a1b26" {
		t.Errorf("expected background #1a1b26, got %s", theme.Background())
	}
	if theme.Primary() != "#7aa2f7" {
		t.Errorf("expected primary #7aa2f7, got %s", theme.Primary())
	}
	if theme.UserColor() != "#ff9e64" {
		t.Errorf("expected user color #ff9e64, got %s", theme.UserColor())
	}
	if theme.AssistantColor() != "#9ece6a" {
		t.Errorf("expected assistant color #9ece6a, got %s", theme.AssistantColor())
	}

	// Verify styles
	styles := theme.Styles()
	if !styles.Title.GetBold() {
		t.Error("title style should be bold")
	}
}

func TestAllThemesHaveAllColors(t *testing.T) {
	themes := []Theme{
		NewDraculaTheme(),
		NewNordTheme(),
		NewGruvboxTheme(),
		NewHighContrastTheme(),
		NewNeoTerminalTheme(),
	}

	for _, theme := range themes {
		t.Run(theme.Name(), func(t *testing.T) {
			if theme.Background() == "" {
				t.Error("Background not set")
			}
			if theme.Foreground() == "" {
				t.Error("Foreground not set")
			}
			if theme.Primary() == "" {
				t.Error("Primary not set")
			}
			if theme.Secondary() == "" {
				t.Error("Secondary not set")
			}
			if theme.Success() == "" {
				t.Error("Success not set")
			}
			if theme.Warning() == "" {
				t.Error("Warning not set")
			}
			if theme.Error() == "" {
				t.Error("Error not set")
			}
			if theme.Info() == "" {
				t.Error("Info not set")
			}
			if theme.Border() == "" {
				t.Error("Border not set")
			}
			if theme.BorderFocused() == "" {
				t.Error("BorderFocused not set")
			}
			if theme.Muted() == "" {
				t.Error("Muted not set")
			}
			if theme.UserColor() == "" {
				t.Error("UserColor not set")
			}
			if theme.AssistantColor() == "" {
				t.Error("AssistantColor not set")
			}
			if theme.ToolColor() == "" {
				t.Error("ToolColor not set")
			}
			if theme.SystemColor() == "" {
				t.Error("SystemColor not set")
			}
		})
	}
}
