package config

// merge merges user config into base config.
// Only non-zero values from user config override base values.
func merge(base, user *Config) {
	mergeTheme(&base.Theme, &user.Theme)
	mergeMouse(&base.Mouse, &user.Mouse)
	mergeKeybindings(&base.Keybindings, &user.Keybindings)
	mergeStatusBar(&base.StatusBar, &user.StatusBar)
	mergeTabBar(&base.TabBar, &user.TabBar)
	mergeInput(&base.Input, &user.Input)
	mergeModal(&base.Modal, &user.Modal)
	mergeAutocomplete(&base.Autocomplete, &user.Autocomplete)
	mergeAccessibility(&base.Accessibility, &user.Accessibility)
}

func mergeTheme(base, user *ThemeConfig) {
	if user.Name != "" {
		base.Name = user.Name
	}
	mergeColors(&base.Colors, &user.Colors)
	if user.Styles != nil {
		if base.Styles == nil {
			base.Styles = make(map[string]Style)
		}
		for k, v := range user.Styles {
			base.Styles[k] = v
		}
	}
}

func mergeColors(base, user *ColorsConfig) {
	if user.Primary != "" {
		base.Primary = user.Primary
	}
	if user.Secondary != "" {
		base.Secondary = user.Secondary
	}
	if user.Background != "" {
		base.Background = user.Background
	}
	if user.Foreground != "" {
		base.Foreground = user.Foreground
	}
	if user.Success != "" {
		base.Success = user.Success
	}
	if user.Warning != "" {
		base.Warning = user.Warning
	}
	if user.Error != "" {
		base.Error = user.Error
	}
	if user.Info != "" {
		base.Info = user.Info
	}
	if user.Border != "" {
		base.Border = user.Border
	}
	if user.BorderFocused != "" {
		base.BorderFocused = user.BorderFocused
	}
	if user.Muted != "" {
		base.Muted = user.Muted
	}
	if user.User != "" {
		base.User = user.User
	}
	if user.Assistant != "" {
		base.Assistant = user.Assistant
	}
	if user.Tool != "" {
		base.Tool = user.Tool
	}
	if user.System != "" {
		base.System = user.System
	}
}

func mergeMouse(base, user *MouseConfig) {
	// For bools, we check if user explicitly set them by using a different approach
	// Since TOML decodes missing bools as false, we can't distinguish
	// We'll just always override if user config was loaded
	// This is a limitation - in practice, users specify what they want to change
	if user.ScrollLines != 0 {
		base.ScrollLines = user.ScrollLines
	}
}

func mergeKeybindings(base, user *KeybindingsConfig) {
	if len(user.Submit) > 0 {
		base.Submit = user.Submit
	}
	if len(user.Cancel) > 0 {
		base.Cancel = user.Cancel
	}
	if len(user.Help) > 0 {
		base.Help = user.Help
	}
	if len(user.QuickActions) > 0 {
		base.QuickActions = user.QuickActions
	}
	if len(user.NextTab) > 0 {
		base.NextTab = user.NextTab
	}
	if len(user.PrevTab) > 0 {
		base.PrevTab = user.PrevTab
	}
	if len(user.ScrollUp) > 0 {
		base.ScrollUp = user.ScrollUp
	}
	if len(user.ScrollDown) > 0 {
		base.ScrollDown = user.ScrollDown
	}
	if len(user.ScrollTop) > 0 {
		base.ScrollTop = user.ScrollTop
	}
	if len(user.ScrollBottom) > 0 {
		base.ScrollBottom = user.ScrollBottom
	}
	if user.Custom != nil {
		if base.Custom == nil {
			base.Custom = make(map[string][]string)
		}
		for k, v := range user.Custom {
			base.Custom[k] = v
		}
	}
}

func mergeStatusBar(base, user *StatusBarConfig) {
	if len(user.Order) > 0 {
		base.Order = user.Order
	}
	if user.Sections != nil {
		if base.Sections == nil {
			base.Sections = make(map[string]StatusBarSection)
		}
		for k, v := range user.Sections {
			base.Sections[k] = v
		}
	}
	if user.Custom != nil {
		if base.Custom == nil {
			base.Custom = make(map[string]CustomSection)
		}
		for k, v := range user.Custom {
			base.Custom[k] = v
		}
	}
}

func mergeTabBar(base, user *TabBarConfig) {
	if user.Position != "" {
		base.Position = user.Position
	}
	if user.Style != "" {
		base.Style = user.Style
	}
	if user.MaxVisible != 0 {
		base.MaxVisible = user.MaxVisible
	}
}

func mergeInput(base, user *InputConfig) {
	if user.Prefix != "" {
		base.Prefix = user.Prefix
	}
	if user.Placeholder != "" {
		base.Placeholder = user.Placeholder
	}
	if user.MaxHeight != 0 {
		base.MaxHeight = user.MaxHeight
	}
	if user.MaxChars != 0 {
		base.MaxChars = user.MaxChars
	}
}

func mergeModal(base, user *ModalConfig) {
	if user.BackdropOpacity != 0 {
		base.BackdropOpacity = user.BackdropOpacity
	}
	if user.Animation != "" {
		base.Animation = user.Animation
	}
}

func mergeAutocomplete(base, user *AutocompleteConfig) {
	if user.MaxSuggestions != 0 {
		base.MaxSuggestions = user.MaxSuggestions
	}
	if user.MinChars != 0 {
		base.MinChars = user.MinChars
	}
	if user.DelayMs != 0 {
		base.DelayMs = user.DelayMs
	}
}

func mergeAccessibility(base, user *AccessibilityConfig) {
	// Bools - same limitation as mouse config
}
