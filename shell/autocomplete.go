package shell

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CompletionProvider generates completions for input text.
type CompletionProvider interface {
	// GetCompletions returns suggestions for the given input.
	GetCompletions(input string) []Completion
}

// Completion represents a single autocomplete suggestion.
type Completion struct {
	Value       string // What gets inserted
	Display     string // What shows in the dropdown
	Description string // Context about the completion
	Score       int    // For ranking (higher = better)
}

// Autocomplete provides input completion with pluggable providers.
type Autocomplete struct {
	providers      map[string]CompletionProvider
	completions    []Completion
	selectedIndex  int
	maxCompletions int
	active         bool
	input          string
	activeProvider string

	// Styles
	dropdownStyle lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
	descStyle     lipgloss.Style
}

// NewAutocomplete creates a new autocomplete component.
func NewAutocomplete() *Autocomplete {
	return &Autocomplete{
		providers:      make(map[string]CompletionProvider),
		maxCompletions: 10,
		dropdownStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6272a4")).
			Padding(0, 1),
		itemStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282a36")).
			Background(lipgloss.Color("#bd93f9")).
			Bold(true),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// RegisterProvider adds a completion provider with the given name.
func (ac *Autocomplete) RegisterProvider(name string, provider CompletionProvider) {
	ac.providers[name] = provider
}

// UnregisterProvider removes a completion provider.
func (ac *Autocomplete) UnregisterProvider(name string) {
	delete(ac.providers, name)
}

// SetMaxCompletions sets the maximum number of completions to show.
func (ac *Autocomplete) SetMaxCompletions(max int) {
	if max > 0 {
		ac.maxCompletions = max
	}
}

// SetStyles sets the styles for the autocomplete dropdown.
func (ac *Autocomplete) SetStyles(dropdown, item, selected, desc lipgloss.Style) {
	ac.dropdownStyle = dropdown
	ac.itemStyle = item
	ac.selectedStyle = selected
	ac.descStyle = desc
}

// Show activates autocomplete for the given input using a specific provider.
func (ac *Autocomplete) Show(input string, providerName string) {
	provider, ok := ac.providers[providerName]
	if !ok {
		return
	}

	ac.input = input
	ac.activeProvider = providerName
	ac.completions = provider.GetCompletions(input)

	// Sort by score (highest first)
	sort.Slice(ac.completions, func(i, j int) bool {
		return ac.completions[i].Score > ac.completions[j].Score
	})

	// Limit completions
	if len(ac.completions) > ac.maxCompletions {
		ac.completions = ac.completions[:ac.maxCompletions]
	}

	ac.selectedIndex = 0
	ac.active = len(ac.completions) > 0
}

// ShowAuto activates autocomplete with automatic provider detection.
func (ac *Autocomplete) ShowAuto(input string) {
	providerName := ac.DetectProvider(input)
	if providerName != "" {
		ac.Show(input, providerName)
	}
}

// Hide deactivates autocomplete.
func (ac *Autocomplete) Hide() {
	ac.active = false
	ac.completions = nil
	ac.selectedIndex = 0
}

// Active returns whether autocomplete is currently showing.
func (ac *Autocomplete) Active() bool {
	return ac.active
}

// Next moves selection to the next completion.
func (ac *Autocomplete) Next() {
	if !ac.active || len(ac.completions) == 0 {
		return
	}
	ac.selectedIndex = (ac.selectedIndex + 1) % len(ac.completions)
}

// Previous moves selection to the previous completion.
func (ac *Autocomplete) Previous() {
	if !ac.active || len(ac.completions) == 0 {
		return
	}
	ac.selectedIndex--
	if ac.selectedIndex < 0 {
		ac.selectedIndex = len(ac.completions) - 1
	}
}

// GetSelected returns the currently selected completion, or nil if none.
func (ac *Autocomplete) GetSelected() *Completion {
	if !ac.active || len(ac.completions) == 0 {
		return nil
	}
	return &ac.completions[ac.selectedIndex]
}

// Completions returns all current completions.
func (ac *Autocomplete) Completions() []Completion {
	return ac.completions
}

// SelectedIndex returns the current selection index.
func (ac *Autocomplete) SelectedIndex() int {
	return ac.selectedIndex
}

// DetectProvider determines the appropriate provider based on input context.
// "/" prefix → "command"
// "./" or "~/" or "/" at start → "file"
// Otherwise → "history"
func (ac *Autocomplete) DetectProvider(input string) string {
	if strings.HasPrefix(input, "/") {
		// Check if "command" provider exists, otherwise try "file"
		if _, ok := ac.providers["command"]; ok {
			return "command"
		}
	}
	if strings.HasPrefix(input, "./") || strings.HasPrefix(input, "~/") || strings.HasPrefix(input, "/") {
		if _, ok := ac.providers["file"]; ok {
			return "file"
		}
	}
	if _, ok := ac.providers["history"]; ok {
		return "history"
	}
	// Return first available provider
	for name := range ac.providers {
		return name
	}
	return ""
}

// Update handles keyboard events for autocomplete navigation.
func (ac *Autocomplete) Update(msg tea.Msg) (*Autocomplete, tea.Cmd) {
	if !ac.active {
		return ac, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			ac.Next()
		case "shift+tab", "up":
			ac.Previous()
		case "esc":
			ac.Hide()
		}
	}

	return ac, nil
}

// View renders the autocomplete dropdown.
func (ac *Autocomplete) View() string {
	if !ac.active || len(ac.completions) == 0 {
		return ""
	}

	var lines []string
	for i, c := range ac.completions {
		var line string
		if c.Description != "" {
			line = c.Display + " " + ac.descStyle.Render(c.Description)
		} else {
			line = c.Display
		}

		if i == ac.selectedIndex {
			lines = append(lines, ac.selectedStyle.Render(line))
		} else {
			lines = append(lines, ac.itemStyle.Render(line))
		}
	}

	return ac.dropdownStyle.Render(strings.Join(lines, "\n"))
}

// CommandProvider provides completions for slash commands.
type CommandProvider struct {
	commands []Completion
}

// NewCommandProvider creates a command completion provider.
func NewCommandProvider(commands []Completion) *CommandProvider {
	return &CommandProvider{commands: commands}
}

// GetCompletions returns commands matching the input.
func (p *CommandProvider) GetCompletions(input string) []Completion {
	if !strings.HasPrefix(input, "/") {
		return nil
	}

	query := strings.ToLower(strings.TrimPrefix(input, "/"))
	var matches []Completion

	for _, cmd := range p.commands {
		cmdName := strings.ToLower(strings.TrimPrefix(cmd.Value, "/"))
		if strings.HasPrefix(cmdName, query) {
			// Higher score for exact prefix match
			score := cmd.Score
			if cmdName == query {
				score += 100
			} else if strings.HasPrefix(cmdName, query) {
				score += 50
			}
			matches = append(matches, Completion{
				Value:       cmd.Value,
				Display:     cmd.Display,
				Description: cmd.Description,
				Score:       score,
			})
		}
	}

	return matches
}

// HistoryProvider provides completions from command history.
type HistoryProvider struct {
	history []string
}

// NewHistoryProvider creates a history completion provider.
func NewHistoryProvider(history []string) *HistoryProvider {
	return &HistoryProvider{history: history}
}

// AddHistory adds an entry to history.
func (p *HistoryProvider) AddHistory(entry string) {
	// Avoid duplicates at the end
	if len(p.history) > 0 && p.history[len(p.history)-1] == entry {
		return
	}
	p.history = append(p.history, entry)
}

// GetCompletions returns history entries matching the input.
func (p *HistoryProvider) GetCompletions(input string) []Completion {
	if input == "" {
		return nil
	}

	query := strings.ToLower(input)
	var matches []Completion

	// Search from most recent
	for i := len(p.history) - 1; i >= 0; i-- {
		entry := p.history[i]
		if strings.Contains(strings.ToLower(entry), query) {
			// Score based on recency and match quality
			score := i // More recent = higher base index
			if strings.HasPrefix(strings.ToLower(entry), query) {
				score += 100
			}
			matches = append(matches, Completion{
				Value:   entry,
				Display: entry,
				Score:   score,
			})
		}
	}

	return matches
}
