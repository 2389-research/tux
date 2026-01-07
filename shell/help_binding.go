// help/binding.go
package shell

// Binding represents a keyboard shortcut or command.
type Binding struct {
	Key         string   // Display text: "ctrl+c", "?", "/feedback"
	Description string   // What it does: "Quit", "Toggle help"
	Modes       []string // Optional: empty = all modes
}

// MatchesMode returns true if this binding should be shown for the given mode.
// Empty mode parameter means show all bindings.
// Empty Modes field means binding is shown in all modes.
func (b Binding) MatchesMode(mode string) bool {
	if mode == "" || len(b.Modes) == 0 {
		return true
	}
	for _, m := range b.Modes {
		if m == mode {
			return true
		}
	}
	return false
}

// Category groups related bindings.
type Category struct {
	Title    string
	Bindings []Binding
}

// FilterByMode returns bindings that match the given mode.
func (c Category) FilterByMode(mode string) []Binding {
	var result []Binding
	for _, b := range c.Bindings {
		if b.MatchesMode(mode) {
			result = append(result, b)
		}
	}
	return result
}
