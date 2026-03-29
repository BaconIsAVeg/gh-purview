package statusbar

import (
	"github.com/anomaly/ghr/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type KeyBinding struct {
	Key  string
	Desc string
}

type Model struct {
	mode    string
	version string
	width   int
	styles  *styles.Palette
}

func New(s *styles.Palette, version string) Model {
	return Model{
		mode:    "pr list",
		styles:  s,
		version: version,
	}
}

func (m *Model) SetMode(mode string) {
	m.mode = mode
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) getKeybinds() []KeyBinding {
	switch m.mode {
	case "preview":
		return []KeyBinding{
			{Key: "a", Desc: "approve"},
			{Key: "o", Desc: "open on web"},
			{Key: "p", Desc: "close"},
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

	// Placeholder for future information
	middleContent := lipgloss.NewStyle().
		Background(barBg).
		Width(middleWidth).
		Render(" " + m.version)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Top, modeContent, middleContent, keysContent)

	if lipgloss.Width(statusBar) < m.width {
		statusBar = lipgloss.NewStyle().Width(m.width).Background(barBg).Render(statusBar)
	}

	return statusBar
}
