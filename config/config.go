// Package config provides two-tier configuration for tux applications.
// App defaults are overridden by user preferences from config files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds all UI configuration settings.
type Config struct {
	Theme         ThemeConfig         `toml:"theme"`
	Mouse         MouseConfig         `toml:"mouse"`
	Keybindings   KeybindingsConfig   `toml:"keybindings"`
	StatusBar     StatusBarConfig     `toml:"statusbar"`
	TabBar        TabBarConfig        `toml:"tabbar"`
	Input         InputConfig         `toml:"input"`
	Modal         ModalConfig         `toml:"modal"`
	Autocomplete  AutocompleteConfig  `toml:"autocomplete"`
	Accessibility AccessibilityConfig `toml:"accessibility"`
}

// ThemeConfig holds theme settings.
type ThemeConfig struct {
	Name   string            `toml:"name"`
	Colors ColorsConfig      `toml:"colors"`
	Styles map[string]Style  `toml:"styles"`
}

// ColorsConfig holds color overrides.
type ColorsConfig struct {
	Primary       string `toml:"primary"`
	Secondary     string `toml:"secondary"`
	Background    string `toml:"background"`
	Foreground    string `toml:"foreground"`
	Success       string `toml:"success"`
	Warning       string `toml:"warning"`
	Error         string `toml:"error"`
	Info          string `toml:"info"`
	Border        string `toml:"border"`
	BorderFocused string `toml:"border_focused"`
	Muted         string `toml:"muted"`
	User          string `toml:"user"`
	Assistant     string `toml:"assistant"`
	Tool          string `toml:"tool"`
	System        string `toml:"system"`
}

// Style holds style settings for a component.
type Style struct {
	Foreground        string `toml:"foreground"`
	Background        string `toml:"background"`
	Bold              bool   `toml:"bold"`
	Italic            bool   `toml:"italic"`
	Border            string `toml:"border"`
	BorderForeground  string `toml:"border_foreground"`
	PaddingHorizontal int    `toml:"padding_horizontal"`
	PaddingVertical   int    `toml:"padding_vertical"`
}

// MouseConfig holds mouse settings.
type MouseConfig struct {
	Enabled          bool `toml:"enabled"`
	ScrollLines      int  `toml:"scroll_lines"`
	HoverEnabled     bool `toml:"hover_enabled"`
	ShiftPassthrough bool `toml:"shift_passthrough"`
}

// KeybindingsConfig holds keybinding settings.
type KeybindingsConfig struct {
	Submit       []string          `toml:"submit"`
	Cancel       []string          `toml:"cancel"`
	Help         []string          `toml:"help"`
	QuickActions []string          `toml:"quick_actions"`
	NextTab      []string          `toml:"next_tab"`
	PrevTab      []string          `toml:"prev_tab"`
	ScrollUp     []string          `toml:"scroll_up"`
	ScrollDown   []string          `toml:"scroll_down"`
	ScrollTop    []string          `toml:"scroll_top"`
	ScrollBottom []string          `toml:"scroll_bottom"`
	Custom       map[string][]string `toml:"custom"`
}

// StatusBarConfig holds status bar settings.
type StatusBarConfig struct {
	Order    []string                     `toml:"order"`
	Sections map[string]StatusBarSection  `toml:"sections"`
	Custom   map[string]CustomSection     `toml:"custom"`
}

// StatusBarSection holds settings for a status bar section.
type StatusBarSection struct {
	MaxWidth int    `toml:"max_width"`
	Format   string `toml:"format"`
}

// CustomSection holds settings for a custom status bar section.
type CustomSection struct {
	Position int `toml:"position"`
	Priority int `toml:"priority"`
}

// TabBarConfig holds tab bar settings.
type TabBarConfig struct {
	Position   string `toml:"position"`
	Style      string `toml:"style"`
	ShowBadges bool   `toml:"show_badges"`
	ShowClose  bool   `toml:"show_close"`
	MaxVisible int    `toml:"max_visible"`
}

// InputConfig holds input area settings.
type InputConfig struct {
	Prefix        string `toml:"prefix"`
	Placeholder   string `toml:"placeholder"`
	Multiline     bool   `toml:"multiline"`
	MaxHeight     int    `toml:"max_height"`
	ShowCharCount bool   `toml:"show_char_count"`
	MaxChars      int    `toml:"max_chars"`
}

// ModalConfig holds modal settings.
type ModalConfig struct {
	Backdrop            bool    `toml:"backdrop"`
	BackdropOpacity     float64 `toml:"backdrop_opacity"`
	Animation           string  `toml:"animation"`
	CloseOnEsc          bool    `toml:"close_on_esc"`
	CloseOnClickOutside bool    `toml:"close_on_click_outside"`
}

// AutocompleteConfig holds autocomplete settings.
type AutocompleteConfig struct {
	Enabled        bool `toml:"enabled"`
	MaxSuggestions int  `toml:"max_suggestions"`
	MinChars       int  `toml:"min_chars"`
	DelayMs        int  `toml:"delay_ms"`
}

// AccessibilityConfig holds accessibility settings.
type AccessibilityConfig struct {
	HighContrast      bool `toml:"high_contrast"`
	ReduceMotion      bool `toml:"reduce_motion"`
	ScreenReaderHints bool `toml:"screen_reader_hints"`
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Theme: ThemeConfig{
			Name: "dracula",
		},
		Mouse: MouseConfig{
			Enabled:          true,
			ScrollLines:      3,
			HoverEnabled:     true,
			ShiftPassthrough: true,
		},
		Keybindings: KeybindingsConfig{
			Submit:       []string{"enter"},
			Cancel:       []string{"esc"},
			Help:         []string{"ctrl+h", "?", "f1"},
			QuickActions: []string{"ctrl+k"},
			NextTab:      []string{"ctrl+tab", "ctrl+n"},
			PrevTab:      []string{"ctrl+shift+tab", "ctrl+p"},
			ScrollUp:     []string{"ctrl+u", "pgup"},
			ScrollDown:   []string{"ctrl+d", "pgdn"},
			ScrollTop:    []string{"g g", "home"},
			ScrollBottom: []string{"G", "end"},
		},
		StatusBar: StatusBarConfig{
			Order: []string{"model", "status", "tokens", "mode", "progress", "hints"},
		},
		TabBar: TabBarConfig{
			Position:   "top",
			Style:      "underline",
			ShowBadges: true,
			ShowClose:  true,
			MaxVisible: 8,
		},
		Input: InputConfig{
			Prefix:    "> ",
			MaxHeight: 5,
		},
		Modal: ModalConfig{
			Backdrop:        true,
			BackdropOpacity: 0.5,
			Animation:       "none",
			CloseOnEsc:      true,
		},
		Autocomplete: AutocompleteConfig{
			Enabled:        true,
			MaxSuggestions: 10,
			MinChars:       1,
			DelayMs:        50,
		},
		Accessibility: AccessibilityConfig{
			ScreenReaderHints: true,
		},
	}
}

// Load loads configuration for the given app name.
// It checks locations in order: XDG config, legacy rc file, env var.
// Returns default config merged with user overrides.
func Load(appName string) (*Config, error) {
	cfg := Default()

	path := findConfigFile(appName)
	if path == "" {
		return cfg, nil // No user config, return defaults
	}

	userCfg, err := loadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("loading config from %s: %w", path, err)
	}

	merge(cfg, userCfg)

	if errs := cfg.Validate(); len(errs) > 0 {
		return cfg, &ValidationError{Errors: errs}
	}

	return cfg, nil
}

// LoadFile loads configuration from a specific file path.
func LoadFile(path string) (*Config, error) {
	cfg := Default()

	userCfg, err := loadFile(path)
	if err != nil {
		return nil, err
	}

	merge(cfg, userCfg)

	if errs := cfg.Validate(); len(errs) > 0 {
		return cfg, &ValidationError{Errors: errs}
	}

	return cfg, nil
}

// findConfigFile finds the config file for the given app.
func findConfigFile(appName string) string {
	// Check env var first
	envVar := strings.ToUpper(appName) + "_UI_CONFIG"
	if path := os.Getenv(envVar); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check XDG config
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	xdgPath := filepath.Join(configDir, appName, "ui.toml")
	if _, err := os.Stat(xdgPath); err == nil {
		return xdgPath
	}

	// Check legacy rc file
	home, _ := os.UserHomeDir()
	rcPath := filepath.Join(home, "."+appName+"rc")
	if _, err := os.Stat(rcPath); err == nil {
		return rcPath
	}

	return ""
}

// loadFile loads config from a TOML file.
func loadFile(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Path returns the path to the config file for the given app, or empty if none exists.
func Path(appName string) string {
	return findConfigFile(appName)
}
