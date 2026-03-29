package preview

import (
	"strings"

	"github.com/anomaly/ghr/internal/types"
	"github.com/anomaly/ghr/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	visible  bool
	pr       *types.PR
	width    int
	height   int
	styles   *styles.Palette
	viewport viewport.Model
}

func New(s *styles.Palette) Model {
	return Model{
		styles:   s,
		viewport: viewport.New(0, 0),
	}
}

func (m *Model) SetPR(pr *types.PR) {
	m.pr = pr
	if pr == nil {
		m.viewport.SetContent("")
	}
	m.viewport.GotoTop()
}

func (m *Model) SetDiffContent(content string) {
	m.viewport.SetContent(m.colorizeDiff(content))
	m.viewport.GotoTop()
}

func (m *Model) SetVisible(visible bool) {
	m.visible = visible
	if visible {
		m.viewport.GotoTop()
	}
}

func (m *Model) Toggle() {
	m.visible = !m.visible
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.viewport.Width = width - 4
}

func (m *Model) SetHeight(height int) {
	m.height = height
	m.viewport.Height = height - 2
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) View() string {
	if !m.visible || m.pr == nil {
		return ""
	}

	content := m.viewport.View()

	border := m.styles.PreviewBorder.
		Width(m.width - 2).
		Height(m.height - 1)

	inner := m.styles.PreviewContent.Padding(0, 1).Render(content)

	return border.Render(inner)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *Model) ScrollUp(lines int) {
	m.viewport.ScrollUp(lines)
}

func (m *Model) ScrollDown(lines int) {
	m.viewport.ScrollDown(lines)
}

func (m *Model) colorizeDiff(content string) string {
	lines := strings.Split(content, "\n")
	var b strings.Builder

	for i, line := range lines {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(m.colorizeLine(line))
	}

	return b.String()
}

func (m *Model) colorizeLine(line string) string {
	if len(line) == 0 {
		return line
	}

	switch {
	case strings.HasPrefix(line, "diff --git"),
		strings.HasPrefix(line, "index "),
		strings.HasPrefix(line, "--- "),
		strings.HasPrefix(line, "+++ "):
		return m.styles.DiffFileHeader.Render(line)
	case strings.HasPrefix(line, "@@"):
		return m.styles.DiffHeader.Render(line)
	case strings.HasPrefix(line, "+"):
		return m.styles.DiffAdd.Render(line)
	case strings.HasPrefix(line, "-"):
		return m.styles.DiffDelete.Render(line)
	default:
		return line
	}
}
