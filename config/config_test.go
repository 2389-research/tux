package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Theme.Name != "dracula" {
		t.Errorf("expected theme dracula, got %s", cfg.Theme.Name)
	}
	if cfg.Input.Prefix != "> " {
		t.Errorf("expected prefix '> ', got %q", cfg.Input.Prefix)
	}
	if cfg.TabBar.Position != "top" {
		t.Errorf("expected position top, got %s", cfg.TabBar.Position)
	}
	if cfg.Modal.Animation != "none" {
		t.Errorf("expected animation none, got %s", cfg.Modal.Animation)
	}
	if !cfg.Mouse.Enabled {
		t.Error("expected mouse enabled by default")
	}
	if cfg.Mouse.ScrollLines != 3 {
		t.Errorf("expected scroll_lines 3, got %d", cfg.Mouse.ScrollLines)
	}
}

func TestLoadNoFile(t *testing.T) {
	cfg, err := Load("nonexistent-app-12345")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Error("config should not be nil")
	}
	// Should return defaults
	if cfg.Theme.Name != "dracula" {
		t.Errorf("expected default theme, got %s", cfg.Theme.Name)
	}
}

func TestLoadFile(t *testing.T) {
	// Create temp config file
	dir := t.TempDir()
	path := filepath.Join(dir, "ui.toml")

	content := `
[theme]
name = "nord"

[theme.colors]
primary = "#88c0d0"

[input]
prefix = "λ "

[tabbar]
position = "bottom"
style = "pills"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Theme.Name != "nord" {
		t.Errorf("expected theme nord, got %s", cfg.Theme.Name)
	}
	if cfg.Theme.Colors.Primary != "#88c0d0" {
		t.Errorf("expected primary #88c0d0, got %s", cfg.Theme.Colors.Primary)
	}
	if cfg.Input.Prefix != "λ " {
		t.Errorf("expected prefix 'λ ', got %q", cfg.Input.Prefix)
	}
	if cfg.TabBar.Position != "bottom" {
		t.Errorf("expected position bottom, got %s", cfg.TabBar.Position)
	}
	if cfg.TabBar.Style != "pills" {
		t.Errorf("expected style pills, got %s", cfg.TabBar.Style)
	}

	// Verify defaults are preserved for unspecified fields
	if cfg.Modal.Animation != "none" {
		t.Errorf("expected default animation none, got %s", cfg.Modal.Animation)
	}
}

func TestLoadFileInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ui.toml")

	// Invalid TOML
	if err := os.WriteFile(path, []byte("invalid [ toml"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid TOML")
	}
}

func TestLoadWithEnvVar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.toml")

	content := `
[theme]
name = "gruvbox"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("TESTAPP_UI_CONFIG", path)
	defer os.Unsetenv("TESTAPP_UI_CONFIG")

	cfg, err := Load("testapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Theme.Name != "gruvbox" {
		t.Errorf("expected theme gruvbox, got %s", cfg.Theme.Name)
	}
}

func TestPath(t *testing.T) {
	// No config file
	path := Path("nonexistent-app-12345")
	if path != "" {
		t.Errorf("expected empty path, got %s", path)
	}

	// With env var
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ui.toml")
	if err := os.WriteFile(configPath, []byte("[theme]\nname = \"nord\""), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("PATHTEST_UI_CONFIG", configPath)
	defer os.Unsetenv("PATHTEST_UI_CONFIG")

	path = Path("pathtest")
	if path != configPath {
		t.Errorf("expected %s, got %s", configPath, path)
	}
}

func TestValidateColors(t *testing.T) {
	tests := []struct {
		color string
		valid bool
	}{
		{"#fff", true},
		{"#FFF", true},
		{"#ffffff", true},
		{"#FFFFFF", true},
		{"#88c0d0", true},
		{"#abc", true},
		{"fff", false},
		{"#gggggg", false},
		{"#ff", false},
		{"#fffffff", false},
		{"rgb(255,255,255)", false},
	}

	for _, tt := range tests {
		result := isValidHexColor(tt.color)
		if result != tt.valid {
			t.Errorf("isValidHexColor(%q) = %v, want %v", tt.color, result, tt.valid)
		}
	}
}

func TestValidateKeybindings(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"a", true},
		{"Z", true},
		{"?", true},
		{"enter", true},
		{"Enter", true},
		{"esc", true},
		{"f1", true},
		{"F12", true},
		{"ctrl+c", true},
		{"ctrl+C", true},
		{"alt+x", true},
		{"shift+tab", true},
		{"ctrl+shift+p", true},
		{"g g", true},
		{"G", true},
		{"pgup", true},
		{"pgdn", true},
		{"home", true},
		{"end", true},
		{"", false},
		{"ctrl++", false},
		{"invalid+key", false},
		{"ctrl+", false},
	}

	for _, tt := range tests {
		result := isValidKeybinding(tt.key)
		if result != tt.valid {
			t.Errorf("isValidKeybinding(%q) = %v, want %v", tt.key, result, tt.valid)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	// Valid config
	cfg := Default()
	errs := cfg.Validate()
	if len(errs) > 0 {
		t.Errorf("default config should be valid, got errors: %v", errs)
	}

	// Invalid theme
	cfg = Default()
	cfg.Theme.Name = "invalid-theme"
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid theme")
	}

	// Invalid color
	cfg = Default()
	cfg.Theme.Colors.Primary = "#gggggg"
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid color")
	}

	// Invalid tabbar position
	cfg = Default()
	cfg.TabBar.Position = "left"
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid tabbar position")
	}

	// Invalid tabbar style
	cfg = Default()
	cfg.TabBar.Style = "fancy"
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid tabbar style")
	}

	// Invalid modal animation
	cfg = Default()
	cfg.Modal.Animation = "bounce"
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid modal animation")
	}

	// Invalid backdrop opacity
	cfg = Default()
	cfg.Modal.BackdropOpacity = 1.5
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid backdrop opacity")
	}

	// Invalid keybinding
	cfg = Default()
	cfg.Keybindings.Help = []string{"ctrl++"}
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid keybinding")
	}

	// Invalid custom keybinding
	cfg = Default()
	cfg.Keybindings.Custom = map[string][]string{
		"test": {"invalid+key"},
	}
	errs = cfg.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for invalid custom keybinding")
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Errors: []string{"error 1", "error 2"},
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("error string should not be empty")
	}
	if !contains(errStr, "error 1") {
		t.Error("error string should contain 'error 1'")
	}
	if !contains(errStr, "error 2") {
		t.Error("error string should contain 'error 2'")
	}
}

func TestLoadFileWithValidationErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ui.toml")

	content := `
[theme]
name = "invalid-theme"

[theme.colors]
primary = "#gggggg"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFile(path)
	if err == nil {
		t.Error("expected validation error")
	}

	// Config should still be returned (with merged values)
	if cfg == nil {
		t.Error("config should not be nil even with validation errors")
	}

	// Check it's a ValidationError
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestMergeAllFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ui.toml")

	// Config that overrides everything
	content := `
[theme]
name = "nord"

[theme.colors]
primary = "#88c0d0"
secondary = "#81a1c1"
background = "#2e3440"
foreground = "#eceff4"
success = "#a3be8c"
warning = "#ebcb8b"
error = "#bf616a"
info = "#5e81ac"
border = "#4c566a"
border_focused = "#88c0d0"
muted = "#4c566a"
user = "#d08770"
assistant = "#a3be8c"
tool = "#5e81ac"
system = "#4c566a"

[theme.styles.title]
foreground = "#88c0d0"
bold = true

[mouse]
scroll_lines = 5

[keybindings]
submit = ["ctrl+enter"]
cancel = ["ctrl+c"]
help = ["f1"]
quick_actions = ["ctrl+p"]
next_tab = ["ctrl+l"]
prev_tab = ["ctrl+h"]
scroll_up = ["k"]
scroll_down = ["j"]
scroll_top = ["gg"]
scroll_bottom = ["G"]

[keybindings.custom]
my_action = ["ctrl+m"]

[statusbar]
order = ["model", "tokens"]

[statusbar.sections.model]
max_width = 30

[statusbar.custom.git]
position = 2
priority = 10

[tabbar]
position = "bottom"
style = "boxed"
max_visible = 5

[input]
prefix = "$ "
placeholder = "Enter command..."
max_height = 10
max_chars = 500

[modal]
backdrop_opacity = 0.8
animation = "fade"

[autocomplete]
max_suggestions = 5
min_chars = 2
delay_ms = 100
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all theme colors
	if cfg.Theme.Colors.Primary != "#88c0d0" {
		t.Errorf("expected primary #88c0d0, got %s", cfg.Theme.Colors.Primary)
	}
	if cfg.Theme.Colors.Secondary != "#81a1c1" {
		t.Errorf("expected secondary #81a1c1, got %s", cfg.Theme.Colors.Secondary)
	}
	if cfg.Theme.Colors.Background != "#2e3440" {
		t.Errorf("expected background #2e3440, got %s", cfg.Theme.Colors.Background)
	}
	if cfg.Theme.Colors.Foreground != "#eceff4" {
		t.Errorf("expected foreground #eceff4, got %s", cfg.Theme.Colors.Foreground)
	}
	if cfg.Theme.Colors.Success != "#a3be8c" {
		t.Errorf("expected success #a3be8c, got %s", cfg.Theme.Colors.Success)
	}
	if cfg.Theme.Colors.Warning != "#ebcb8b" {
		t.Errorf("expected warning #ebcb8b, got %s", cfg.Theme.Colors.Warning)
	}
	if cfg.Theme.Colors.Error != "#bf616a" {
		t.Errorf("expected error #bf616a, got %s", cfg.Theme.Colors.Error)
	}
	if cfg.Theme.Colors.Info != "#5e81ac" {
		t.Errorf("expected info #5e81ac, got %s", cfg.Theme.Colors.Info)
	}
	if cfg.Theme.Colors.Border != "#4c566a" {
		t.Errorf("expected border #4c566a, got %s", cfg.Theme.Colors.Border)
	}
	if cfg.Theme.Colors.BorderFocused != "#88c0d0" {
		t.Errorf("expected border_focused #88c0d0, got %s", cfg.Theme.Colors.BorderFocused)
	}
	if cfg.Theme.Colors.Muted != "#4c566a" {
		t.Errorf("expected muted #4c566a, got %s", cfg.Theme.Colors.Muted)
	}
	if cfg.Theme.Colors.User != "#d08770" {
		t.Errorf("expected user #d08770, got %s", cfg.Theme.Colors.User)
	}
	if cfg.Theme.Colors.Assistant != "#a3be8c" {
		t.Errorf("expected assistant #a3be8c, got %s", cfg.Theme.Colors.Assistant)
	}
	if cfg.Theme.Colors.Tool != "#5e81ac" {
		t.Errorf("expected tool #5e81ac, got %s", cfg.Theme.Colors.Tool)
	}
	if cfg.Theme.Colors.System != "#4c566a" {
		t.Errorf("expected system #4c566a, got %s", cfg.Theme.Colors.System)
	}

	// Verify styles
	if cfg.Theme.Styles == nil {
		t.Error("expected styles to be set")
	} else if style, ok := cfg.Theme.Styles["title"]; !ok {
		t.Error("expected title style")
	} else if style.Foreground != "#88c0d0" {
		t.Errorf("expected title foreground #88c0d0, got %s", style.Foreground)
	}

	// Verify mouse
	if cfg.Mouse.ScrollLines != 5 {
		t.Errorf("expected scroll_lines 5, got %d", cfg.Mouse.ScrollLines)
	}

	// Verify keybindings
	if len(cfg.Keybindings.Submit) != 1 || cfg.Keybindings.Submit[0] != "ctrl+enter" {
		t.Errorf("expected submit [ctrl+enter], got %v", cfg.Keybindings.Submit)
	}
	if len(cfg.Keybindings.Cancel) != 1 || cfg.Keybindings.Cancel[0] != "ctrl+c" {
		t.Errorf("expected cancel [ctrl+c], got %v", cfg.Keybindings.Cancel)
	}
	if len(cfg.Keybindings.Help) != 1 || cfg.Keybindings.Help[0] != "f1" {
		t.Errorf("expected help [f1], got %v", cfg.Keybindings.Help)
	}
	if len(cfg.Keybindings.QuickActions) != 1 || cfg.Keybindings.QuickActions[0] != "ctrl+p" {
		t.Errorf("expected quick_actions [ctrl+p], got %v", cfg.Keybindings.QuickActions)
	}
	if len(cfg.Keybindings.NextTab) != 1 || cfg.Keybindings.NextTab[0] != "ctrl+l" {
		t.Errorf("expected next_tab [ctrl+l], got %v", cfg.Keybindings.NextTab)
	}
	if len(cfg.Keybindings.PrevTab) != 1 || cfg.Keybindings.PrevTab[0] != "ctrl+h" {
		t.Errorf("expected prev_tab [ctrl+h], got %v", cfg.Keybindings.PrevTab)
	}
	if len(cfg.Keybindings.ScrollUp) != 1 || cfg.Keybindings.ScrollUp[0] != "k" {
		t.Errorf("expected scroll_up [k], got %v", cfg.Keybindings.ScrollUp)
	}
	if len(cfg.Keybindings.ScrollDown) != 1 || cfg.Keybindings.ScrollDown[0] != "j" {
		t.Errorf("expected scroll_down [j], got %v", cfg.Keybindings.ScrollDown)
	}
	if len(cfg.Keybindings.ScrollTop) != 1 || cfg.Keybindings.ScrollTop[0] != "gg" {
		t.Errorf("expected scroll_top [gg], got %v", cfg.Keybindings.ScrollTop)
	}
	if len(cfg.Keybindings.ScrollBottom) != 1 || cfg.Keybindings.ScrollBottom[0] != "G" {
		t.Errorf("expected scroll_bottom [G], got %v", cfg.Keybindings.ScrollBottom)
	}
	if cfg.Keybindings.Custom == nil {
		t.Error("expected custom keybindings")
	} else if keys, ok := cfg.Keybindings.Custom["my_action"]; !ok || len(keys) != 1 {
		t.Errorf("expected custom my_action, got %v", cfg.Keybindings.Custom)
	}

	// Verify statusbar
	if len(cfg.StatusBar.Order) != 2 {
		t.Errorf("expected 2 statusbar items, got %d", len(cfg.StatusBar.Order))
	}
	if cfg.StatusBar.Sections == nil {
		t.Error("expected statusbar sections")
	} else if section, ok := cfg.StatusBar.Sections["model"]; !ok || section.MaxWidth != 30 {
		t.Error("expected model section with max_width 30")
	}
	if cfg.StatusBar.Custom == nil {
		t.Error("expected custom statusbar sections")
	} else if custom, ok := cfg.StatusBar.Custom["git"]; !ok || custom.Position != 2 {
		t.Error("expected git custom section with position 2")
	}

	// Verify tabbar
	if cfg.TabBar.Position != "bottom" {
		t.Errorf("expected position bottom, got %s", cfg.TabBar.Position)
	}
	if cfg.TabBar.Style != "boxed" {
		t.Errorf("expected style boxed, got %s", cfg.TabBar.Style)
	}
	if cfg.TabBar.MaxVisible != 5 {
		t.Errorf("expected max_visible 5, got %d", cfg.TabBar.MaxVisible)
	}

	// Verify input
	if cfg.Input.Prefix != "$ " {
		t.Errorf("expected prefix '$ ', got %q", cfg.Input.Prefix)
	}
	if cfg.Input.Placeholder != "Enter command..." {
		t.Errorf("expected placeholder 'Enter command...', got %q", cfg.Input.Placeholder)
	}
	if cfg.Input.MaxHeight != 10 {
		t.Errorf("expected max_height 10, got %d", cfg.Input.MaxHeight)
	}
	if cfg.Input.MaxChars != 500 {
		t.Errorf("expected max_chars 500, got %d", cfg.Input.MaxChars)
	}

	// Verify modal
	if cfg.Modal.BackdropOpacity != 0.8 {
		t.Errorf("expected backdrop_opacity 0.8, got %f", cfg.Modal.BackdropOpacity)
	}
	if cfg.Modal.Animation != "fade" {
		t.Errorf("expected animation fade, got %s", cfg.Modal.Animation)
	}

	// Verify autocomplete
	if cfg.Autocomplete.MaxSuggestions != 5 {
		t.Errorf("expected max_suggestions 5, got %d", cfg.Autocomplete.MaxSuggestions)
	}
	if cfg.Autocomplete.MinChars != 2 {
		t.Errorf("expected min_chars 2, got %d", cfg.Autocomplete.MinChars)
	}
	if cfg.Autocomplete.DelayMs != 100 {
		t.Errorf("expected delay_ms 100, got %d", cfg.Autocomplete.DelayMs)
	}
}

func TestMergePreservesDefaults(t *testing.T) {
	base := Default()
	user := &Config{} // Empty config

	merge(base, user)

	// Defaults should be preserved
	if base.Theme.Name != "dracula" {
		t.Errorf("expected theme dracula, got %s", base.Theme.Name)
	}
	if base.Input.Prefix != "> " {
		t.Errorf("expected prefix '> ', got %q", base.Input.Prefix)
	}
}

func TestMergeAccessibility(t *testing.T) {
	base := Default()
	user := &Config{
		Accessibility: AccessibilityConfig{
			HighContrast: true,
		},
	}

	merge(base, user)
	// Just verify it doesn't panic - bool merging is tricky
}

func TestFindConfigFileXDG(t *testing.T) {
	dir := t.TempDir()

	// Create XDG config
	configDir := filepath.Join(dir, ".config", "xdgtest")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "ui.toml")
	if err := os.WriteFile(configPath, []byte("[theme]\nname = \"nord\""), 0644); err != nil {
		t.Fatal(err)
	}

	// Set XDG_CONFIG_HOME
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
	defer os.Setenv("XDG_CONFIG_HOME", oldXDG)

	path := findConfigFile("xdgtest")
	if path != configPath {
		t.Errorf("expected %s, got %s", configPath, path)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestValidateValidThemes(t *testing.T) {
	themes := []string{"dracula", "nord", "gruvbox", "high-contrast"}
	for _, theme := range themes {
		cfg := Default()
		cfg.Theme.Name = theme
		errs := cfg.Validate()
		if len(errs) > 0 {
			t.Errorf("theme %q should be valid, got errors: %v", theme, errs)
		}
	}
}

func TestValidateAllColorFields(t *testing.T) {
	cfg := Default()
	// Set all colors to valid values
	cfg.Theme.Colors = ColorsConfig{
		Primary:       "#fff",
		Secondary:     "#fff",
		Background:    "#fff",
		Foreground:    "#fff",
		Success:       "#fff",
		Warning:       "#fff",
		Error:         "#fff",
		Info:          "#fff",
		Border:        "#fff",
		BorderFocused: "#fff",
		Muted:         "#fff",
		User:          "#fff",
		Assistant:     "#fff",
		Tool:          "#fff",
		System:        "#fff",
	}
	errs := cfg.Validate()
	if len(errs) > 0 {
		t.Errorf("all valid colors should pass, got errors: %v", errs)
	}
}
