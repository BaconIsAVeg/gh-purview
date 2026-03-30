package statusbar

import (
	"fmt"

	"github.com/anomaly/ghr/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

const (
	ModeList = "list mode"
	ModeDiff = "diff mode"
)

type KeyBinding struct {
	Key  string
	Desc string
}

type Model struct {
	mode      string
	version   string
	width     int
	styles    *styles.Palette
	additions int
	deletions int
}

func New(s *styles.Palette, version string) Model {
	return Model{
		mode:    ModeList,
		styles:  s,
		version: version,
	}
}

func (m *Model) SetMode(mode string) {
	m.mode = mode
}

func (m *Model) SetStats(additions, deletions int) {
	m.additions = additions
	m.deletions = deletions
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) getKeybinds() []KeyBinding {
	switch m.mode {
	case ModeDiff:
		return []KeyBinding{
			{Key: "^n/^p", Desc: "scroll"},
			{Key: "a", Desc: "approve"},
			{Key: "o", Desc: "open on web"},
			{Key: "esc", Desc: "close"},
		}
	default:
		return []KeyBinding{
			{Key: "j/k", Desc: "navigate"},
			{Key: "p", Desc: "preview"},
			{Key: "r", Desc: "refresh"},
			{Key: "q", Desc: "quit"},
		}
	}
}

func (m Model) View() string {
	barBg := lipgloss.Color("234")

	modeContent := m.styles.StatusMode.Render(m.mode)

	keys := m.getKeybinds()
	keysText := ""
	for i, k := range keys {
		if i > 0 {
			keysText += m.styles.StatusSep.Render("  ")
		}
		keysText += m.styles.StatusKey.Render(k.Key) + m.styles.StatusDesc.Render(" "+k.Desc)
	}

	keysContent := m.styles.StatusBar.Render(keysText)

	leftWidth := lipgloss.Width(modeContent)
	rightWidth := lipgloss.Width(keysContent)
	middleWidth := max(m.width-leftWidth-rightWidth, 0)

	var middleText string
	if m.mode == ModeDiff && (m.additions > 0 || m.deletions > 0) {
		middleText = fmt.Sprintf(" +%d -%d ", m.additions, m.deletions)
	} else {
		middleText = " purview " + m.version
	}
	middleContent := lipgloss.NewStyle().
		Background(barBg).
		Width(middleWidth).
		Render(middleText)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Top, modeContent, middleContent, keysContent)

	if lipgloss.Width(statusBar) < m.width {
		statusBar = lipgloss.NewStyle().Width(m.width).Background(barBg).Render(statusBar)
	}

	return statusBar
}
