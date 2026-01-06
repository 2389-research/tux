package content

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressStatus represents the status of a progress item.
type ProgressStatus int

const (
	// ProgressPending indicates the item hasn't started.
	ProgressPending ProgressStatus = iota
	// ProgressRunning indicates the item is in progress.
	ProgressRunning
	// ProgressComplete indicates the item finished successfully.
	ProgressComplete
	// ProgressError indicates the item failed.
	ProgressError
)

// ProgressItem represents a single item in the progress display.
type ProgressItem struct {
	Label  string
	Status ProgressStatus
}

// Progress displays progress with optional item list.
type Progress struct {
	items      []ProgressItem
	total      int
	current    int
	message    string
	showBar    bool
	showItems  bool
	maxVisible int
	width      int
	height     int

	// Styles
	barStyle      lipgloss.Style
	fillStyle     lipgloss.Style
	emptyStyle    lipgloss.Style
	pendingStyle  lipgloss.Style
	runningStyle  lipgloss.Style
	completeStyle lipgloss.Style
	errorStyle    lipgloss.Style
	messageStyle  lipgloss.Style
}

// ProgressConfig configures a Progress component.
type ProgressConfig struct {
	Total      int
	ShowBar    bool
	ShowItems  bool
	MaxVisible int
}

// NewProgress creates a new progress display.
func NewProgress(cfg ProgressConfig) *Progress {
	if cfg.MaxVisible == 0 {
		cfg.MaxVisible = 10
	}
	return &Progress{
		items:      make([]ProgressItem, 0),
		total:      cfg.Total,
		showBar:    cfg.ShowBar,
		showItems:  cfg.ShowItems,
		maxVisible: cfg.MaxVisible,
		barStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		fillStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
		emptyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#44475a")),
		pendingStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		runningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f1fa8c")),
		completeStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")),
		messageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
	}
}

// Init implements Content.
func (p *Progress) Init() tea.Cmd {
	return nil
}

// Update implements Content.
func (p *Progress) Update(msg tea.Msg) (Content, tea.Cmd) {
	return p, nil
}

// View implements Content.
func (p *Progress) View() string {
	var sections []string

	// Message
	if p.message != "" {
		sections = append(sections, p.messageStyle.Render(p.message))
	}

	// Progress bar
	if p.showBar && p.total > 0 {
		sections = append(sections, p.renderBar())
	}

	// Items
	if p.showItems && len(p.items) > 0 {
		sections = append(sections, p.renderItems())
	}

	return strings.Join(sections, "\n")
}

// renderBar renders the progress bar.
func (p *Progress) renderBar() string {
	barWidth := 20
	if p.width > 0 && p.width < 40 {
		barWidth = p.width / 2
	}

	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(barWidth))
	empty := barWidth - filled

	bar := p.fillStyle.Render(strings.Repeat("█", filled)) +
		p.emptyStyle.Render(strings.Repeat("░", empty))

	return fmt.Sprintf("[%s] %d/%d (%.0f%%)", bar, p.current, p.total, percent*100)
}

// renderItems renders the item list.
func (p *Progress) renderItems() string {
	var b strings.Builder

	start := 0
	if len(p.items) > p.maxVisible {
		start = len(p.items) - p.maxVisible
	}

	for i := start; i < len(p.items); i++ {
		item := p.items[i]
		var icon string
		var style lipgloss.Style

		switch item.Status {
		case ProgressPending:
			icon = "○"
			style = p.pendingStyle
		case ProgressRunning:
			icon = "●"
			style = p.runningStyle
		case ProgressComplete:
			icon = "✓"
			style = p.completeStyle
		case ProgressError:
			icon = "✗"
			style = p.errorStyle
		}

		b.WriteString(style.Render(icon + " " + item.Label))
		if i < len(p.items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Value implements Content.
func (p *Progress) Value() any {
	return p.current
}

// SetSize implements Content.
func (p *Progress) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// SetTotal sets the total count.
func (p *Progress) SetTotal(total int) {
	p.total = total
}

// SetCurrent sets the current count.
func (p *Progress) SetCurrent(current int) {
	p.current = current
}

// SetMessage sets the progress message.
func (p *Progress) SetMessage(message string) {
	p.message = message
}

// AddItem adds a progress item.
func (p *Progress) AddItem(item ProgressItem) {
	p.items = append(p.items, item)
}

// UpdateItem updates an item's status by index.
func (p *Progress) UpdateItem(index int, status ProgressStatus) {
	if index >= 0 && index < len(p.items) {
		p.items[index].Status = status
	}
}

// UpdateItemByLabel updates an item's status by label.
func (p *Progress) UpdateItemByLabel(label string, status ProgressStatus) {
	for i := range p.items {
		if p.items[i].Label == label {
			p.items[i].Status = status
			return
		}
	}
}

// Clear clears all items and resets progress.
func (p *Progress) Clear() {
	p.items = make([]ProgressItem, 0)
	p.current = 0
	p.message = ""
}

// Percent returns the current progress percentage.
func (p *Progress) Percent() float64 {
	if p.total == 0 {
		return 0
	}
	return float64(p.current) / float64(p.total)
}
