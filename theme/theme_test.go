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
