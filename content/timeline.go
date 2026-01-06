package content

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TimelineStatus represents the status of a timeline item.
type TimelineStatus int

const (
	// TimelinePending indicates the item is waiting.
	TimelinePending TimelineStatus = iota
	// TimelineRunning indicates the item is in progress.
	TimelineRunning
	// TimelineSuccess indicates the item completed successfully.
	TimelineSuccess
	// TimelineError indicates the item failed.
	TimelineError
)

// TimelineItem represents an item in the timeline.
type TimelineItem struct {
	ID        string
	Timestamp time.Time
	Icon      string
	Title     string
	Content   string
	Status    TimelineStatus
	Expanded  bool
}

// Timeline displays a chronological list of items.
type Timeline struct {
	items    []TimelineItem
	viewport viewport.Model
	width    int
	height   int

	// Styles
	pendingStyle  lipgloss.Style
	runningStyle  lipgloss.Style
	successStyle  lipgloss.Style
	errorStyle    lipgloss.Style
	titleStyle    lipgloss.Style
	contentStyle  lipgloss.Style
	timeStyle     lipgloss.Style
}

// NewTimeline creates a new timeline.
func NewTimeline() *Timeline {
	return &Timeline{
		items:    make([]TimelineItem, 0),
		viewport: viewport.New(0, 0),
		pendingStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		runningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f1fa8c")),
		successStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")).
			Bold(true),
		contentStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		timeStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
	}
}

// Init implements Content.
func (t *Timeline) Init() tea.Cmd {
	return nil
}

// Update implements Content.
func (t *Timeline) Update(msg tea.Msg) (Content, tea.Cmd) {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return t, cmd
}

// View implements Content.
func (t *Timeline) View() string {
	return t.viewport.View()
}

// Value implements Content.
func (t *Timeline) Value() any {
	return nil
}

// SetSize implements Content.
func (t *Timeline) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.viewport.Width = width
	t.viewport.Height = height
	t.updateViewport()
}

// AddItem adds an item to the timeline.
func (t *Timeline) AddItem(item TimelineItem) {
	if item.Timestamp.IsZero() {
		item.Timestamp = time.Now()
	}
	t.items = append(t.items, item)
	t.updateViewport()
}

// UpdateItem updates an existing item by ID.
func (t *Timeline) UpdateItem(id string, updates TimelineItem) {
	for i := range t.items {
		if t.items[i].ID == id {
			if updates.Title != "" {
				t.items[i].Title = updates.Title
			}
			if updates.Content != "" {
				t.items[i].Content = updates.Content
			}
			if updates.Icon != "" {
				t.items[i].Icon = updates.Icon
			}
			t.items[i].Status = updates.Status
			t.items[i].Expanded = updates.Expanded
			t.updateViewport()
			return
		}
	}
}

// Clear removes all items from the timeline.
func (t *Timeline) Clear() {
	t.items = make([]TimelineItem, 0)
	t.updateViewport()
}

// updateViewport rebuilds the viewport content.
func (t *Timeline) updateViewport() {
	t.viewport.SetContent(t.renderItems())
}

// renderItems renders all timeline items.
func (t *Timeline) renderItems() string {
	if len(t.items) == 0 {
		return t.contentStyle.Render("No activity")
	}

	var b strings.Builder

	for i, item := range t.items {
		// Status icon
		var statusStyle lipgloss.Style
		icon := item.Icon
		if icon == "" {
			switch item.Status {
			case TimelinePending:
				icon = "○"
				statusStyle = t.pendingStyle
			case TimelineRunning:
				icon = "●"
				statusStyle = t.runningStyle
			case TimelineSuccess:
				icon = "✓"
				statusStyle = t.successStyle
			case TimelineError:
				icon = "✗"
				statusStyle = t.errorStyle
			}
		} else {
			statusStyle = t.getStatusStyle(item.Status)
		}

		// Timestamp
		timeStr := item.Timestamp.Format("15:04:05")

		// Build line
		b.WriteString(statusStyle.Render(icon))
		b.WriteString(" ")
		b.WriteString(t.timeStyle.Render(timeStr))
		b.WriteString(" ")
		b.WriteString(t.titleStyle.Render(item.Title))

		// Content (if expanded or has content)
		if item.Content != "" && item.Expanded {
			b.WriteString("\n    ")
			b.WriteString(t.contentStyle.Render(item.Content))
		}

		if i < len(t.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// getStatusStyle returns the style for a status.
func (t *Timeline) getStatusStyle(status TimelineStatus) lipgloss.Style {
	switch status {
	case TimelinePending:
		return t.pendingStyle
	case TimelineRunning:
		return t.runningStyle
	case TimelineSuccess:
		return t.successStyle
	case TimelineError:
		return t.errorStyle
	default:
		return t.pendingStyle
	}
}

// Count returns the number of items.
func (t *Timeline) Count() int {
	return len(t.items)
}

// GetItem returns an item by ID.
func (t *Timeline) GetItem(id string) *TimelineItem {
	for i := range t.items {
		if t.items[i].ID == id {
			return &t.items[i]
		}
	}
	return nil
}

// ScrollToBottom scrolls to the bottom of the timeline.
func (t *Timeline) ScrollToBottom() {
	t.viewport.GotoBottom()
}

// ScrollToTop scrolls to the top of the timeline.
func (t *Timeline) ScrollToTop() {
	t.viewport.GotoTop()
}
