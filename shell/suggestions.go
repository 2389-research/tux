// Package shell provides smart suggestions based on input analysis.
package shell

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Suggestion represents a tool suggestion with confidence score.
type Suggestion struct {
	ToolName   string  // Name of the suggested tool
	Confidence float64 // Confidence score (0.0 - 1.0)
	Reason     string  // Human-readable reason for the suggestion
	Action     string  // Suggested action text (e.g., ":read /path/to/file")
}

// SuggestionProvider analyzes input and returns suggestions.
type SuggestionProvider interface {
	// Analyze analyzes the input and returns suggestions.
	// Returns nil if no suggestions are applicable.
	Analyze(input string) []Suggestion
}

// Suggestions manages the suggestion display state.
type Suggestions struct {
	provider    SuggestionProvider
	suggestions []Suggestion
	active      bool

	// Styles
	boxStyle    lipgloss.Style
	actionStyle lipgloss.Style
	reasonStyle lipgloss.Style
}

// NewSuggestions creates a new suggestions component.
func NewSuggestions() *Suggestions {
	return &Suggestions{
		boxStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Italic(true),
		actionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")),
		reasonStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// SetProvider sets the suggestion provider.
func (s *Suggestions) SetProvider(p SuggestionProvider) {
	s.provider = p
}

// Update analyzes the current input and updates suggestions.
func (s *Suggestions) Update(input string) {
	if s.provider == nil {
		s.suggestions = nil
		s.active = false
		return
	}

	s.suggestions = s.provider.Analyze(input)
	s.active = len(s.suggestions) > 0
}

// Active returns whether suggestions are visible.
func (s *Suggestions) Active() bool {
	return s.active
}

// Hide hides the suggestions.
func (s *Suggestions) Hide() {
	s.active = false
	s.suggestions = nil
}

// Top returns the top suggestion, or nil if none.
func (s *Suggestions) Top() *Suggestion {
	if len(s.suggestions) == 0 {
		return nil
	}
	return &s.suggestions[0]
}

// All returns all suggestions.
func (s *Suggestions) All() []Suggestion {
	return s.suggestions
}

// View renders the suggestions.
func (s *Suggestions) View() string {
	if !s.active || len(s.suggestions) == 0 {
		return ""
	}

	var parts []string
	// Only show top suggestion inline
	top := s.suggestions[0]

	hint := s.boxStyle.Render("💡 ") +
		s.actionStyle.Render(top.Action) +
		s.reasonStyle.Render(" ("+top.Reason+")")
	parts = append(parts, hint)

	if len(s.suggestions) > 1 {
		more := s.reasonStyle.Render(
			strings.Repeat(" ", 3) +
				"+" + string(rune('0'+len(s.suggestions)-1)) + " more suggestions")
		parts = append(parts, more)
	}

	return strings.Join(parts, "\n")
}
