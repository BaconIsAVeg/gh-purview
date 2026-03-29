package header

import (
	"fmt"

	"github.com/anomaly/ghr/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	filter     string
	count      int
	totalCount int
	width      int
	styles     *styles.Palette
}

func New(s *styles.Palette, version string) Model {
	return Model{
		filter: "",
		styles: s,
	}
}

func (m *Model) SetFilter(filter string) {
	m.filter = filter
}

func (m *Model) SetCount(count, total int) {
	m.count = count
	m.totalCount = total
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) View() string {
	leftContent := m.styles.HeaderTitle.Render("purview")
	rightContent := m.styles.HeaderText.Padding(0, 1).Render(fmt.Sprintf("%d/%d", m.count, m.totalCount))

	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	middleWidth := m.width - leftWidth - rightWidth

	if middleWidth < 0 {
		middleWidth = 0
	}

	middleContent := m.styles.StatusBar.Width(middleWidth).Render(m.filter)

	header := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, middleContent, rightContent)

	if lipgloss.Width(header) < m.width {
		header = m.styles.Header.Width(m.width).Render(header)
	}

	return header
}
