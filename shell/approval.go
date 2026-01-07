package shell

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ApprovalDecision represents the user's decision on a tool approval.
type ApprovalDecision int

const (
	// DecisionApprove allows the tool to run this time.
	DecisionApprove ApprovalDecision = iota
	// DecisionDeny skips the tool this time.
	DecisionDeny
	// DecisionAlwaysAllow remembers to always allow this tool.
	DecisionAlwaysAllow
	// DecisionNeverAllow remembers to never allow this tool.
	DecisionNeverAllow
)

// RiskLevel indicates the risk level of a tool.
type RiskLevel int

const (
	// RiskLow indicates minimal risk.
	RiskLow RiskLevel = iota
	// RiskMedium indicates moderate risk.
	RiskMedium
	// RiskHigh indicates significant risk.
	RiskHigh
)

// ToolInfo contains information about a tool for approval.
type ToolInfo struct {
	ID      string
	Name    string
	Params  map[string]any
	Preview string
	Risk    RiskLevel
}

// ApprovalOption represents an option in the approval modal.
type ApprovalOption struct {
	Label    string
	Decision ApprovalDecision
	Hint     string
}

// DefaultApprovalOptions are the standard approval options.
var DefaultApprovalOptions = []ApprovalOption{
	{Label: "Approve", Decision: DecisionApprove, Hint: "Run this time"},
	{Label: "Deny", Decision: DecisionDeny, Hint: "Skip this time"},
	{Label: "Always Allow", Decision: DecisionAlwaysAllow, Hint: "Never ask again"},
	{Label: "Never Allow", Decision: DecisionNeverAllow, Hint: "Block permanently"},
}

// ApprovalModal presents a tool approval request to the user.
type ApprovalModal struct {
	id         string
	tool       ToolInfo
	options    []ApprovalOption
	selected   int
	queueHint  string
	onDecision func(decision ApprovalDecision)
	width      int
	height     int

	// Styles
	boxStyle      lipgloss.Style
	titleStyle    lipgloss.Style
	toolStyle     lipgloss.Style
	paramStyle    lipgloss.Style
	previewStyle  lipgloss.Style
	optionStyle   lipgloss.Style
	selectedStyle lipgloss.Style
	hintStyle     lipgloss.Style
	riskLowStyle  lipgloss.Style
	riskMedStyle  lipgloss.Style
	riskHighStyle lipgloss.Style
}

// ApprovalModalConfig configures an ApprovalModal.
type ApprovalModalConfig struct {
	Tool       ToolInfo
	Options    []ApprovalOption
	QueueHint  string
	OnDecision func(decision ApprovalDecision)
}

// NewApprovalModal creates a new approval modal.
func NewApprovalModal(cfg ApprovalModalConfig) *ApprovalModal {
	options := cfg.Options
	if len(options) == 0 {
		options = DefaultApprovalOptions
	}
	return &ApprovalModal{
		id:         "approval-" + cfg.Tool.ID,
		tool:       cfg.Tool,
		options:    options,
		queueHint:  cfg.QueueHint,
		onDecision: cfg.OnDecision,
		boxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffb86c")).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffb86c")).
			Bold(true),
		toolStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8be9fd")).
			Bold(true),
		paramStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		previewStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")),
		optionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")).
			Bold(true),
		hintStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")),
		riskLowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")),
		riskMedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f1fa8c")),
		riskHighStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")),
	}
}

// ID implements Modal.
func (m *ApprovalModal) ID() string { return m.id }

// Title implements Modal.
func (m *ApprovalModal) Title() string { return "Tool Approval" }

// Size implements Modal.
func (m *ApprovalModal) Size() Size { return SizeMedium }

// OnPush implements Modal.
func (m *ApprovalModal) OnPush(width, height int) {
	m.width = width
	m.height = height
}

// OnPop implements Modal.
func (m *ApprovalModal) OnPop() {}

// HandleKey implements Modal.
func (m *ApprovalModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type {
	case tea.KeyUp:
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case tea.KeyDown:
		if m.selected < len(m.options)-1 {
			m.selected++
		}
		return true, nil
	case tea.KeyEnter:
		if m.onDecision != nil && m.selected >= 0 && m.selected < len(m.options) {
			m.onDecision(m.options[m.selected].Decision)
		}
		return true, func() tea.Msg { return PopMsg{} }
	}

	switch key.String() {
	case "k":
		if m.selected > 0 {
			m.selected--
		}
		return true, nil
	case "j":
		if m.selected < len(m.options)-1 {
			m.selected++
		}
		return true, nil
	case "y", "a":
		m.selected = 0 // Approve
		if m.onDecision != nil {
			m.onDecision(DecisionApprove)
		}
		return true, func() tea.Msg { return PopMsg{} }
	case "n", "d":
		m.selected = 1 // Deny
		if m.onDecision != nil {
			m.onDecision(DecisionDeny)
		}
		return true, func() tea.Msg { return PopMsg{} }
	}

	return false, nil
}

// Render implements Modal.
func (m *ApprovalModal) Render(width, height int) string {
	var parts []string

	// Title with risk indicator
	title := "Tool Approval Required"
	if m.queueHint != "" {
		title += " " + m.queueHint
	}
	parts = append(parts, m.titleStyle.Render(title))
	parts = append(parts, "")

	// Tool name with risk
	riskStyle := m.riskLowStyle
	riskText := "low risk"
	switch m.tool.Risk {
	case RiskMedium:
		riskStyle = m.riskMedStyle
		riskText = "medium risk"
	case RiskHigh:
		riskStyle = m.riskHighStyle
		riskText = "HIGH RISK"
	}
	parts = append(parts, m.toolStyle.Render(m.tool.Name)+" "+riskStyle.Render("("+riskText+")"))

	// Parameters
	if len(m.tool.Params) > 0 {
		parts = append(parts, "")
		for k, v := range m.tool.Params {
			paramStr := fmt.Sprintf("  %s: %v", k, v)
			// Truncate long values
			if len(paramStr) > width-10 {
				paramStr = paramStr[:width-13] + "..."
			}
			parts = append(parts, m.paramStyle.Render(paramStr))
		}
	}

	// Preview
	if m.tool.Preview != "" {
		parts = append(parts, "")
		parts = append(parts, m.previewStyle.Render(m.tool.Preview))
	}

	parts = append(parts, "")

	// Options
	for i, opt := range m.options {
		prefix := "  "
		style := m.optionStyle
		if i == m.selected {
			prefix = "▸ "
			style = m.selectedStyle
		}
		line := prefix + style.Render(opt.Label)
		if opt.Hint != "" {
			line += " " + m.hintStyle.Render("("+opt.Hint+")")
		}
		parts = append(parts, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.boxStyle.Width(width - 4).Render(content)
}

// Tool returns the tool info.
func (m *ApprovalModal) Tool() ToolInfo {
	return m.tool
}

// SetQueueHint sets the queue progress hint.
func (m *ApprovalModal) SetQueueHint(hint string) {
	m.queueHint = hint
}

// Selected returns the currently selected option index.
func (m *ApprovalModal) Selected() int {
	return m.selected
}
