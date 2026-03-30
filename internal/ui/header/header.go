package header

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/anomaly/ghr/internal/ui/styles"
)

type Model struct {
	filter     string
	count      int
	totalCount int
	width      int
	styles     *styles.Palette
	textInput  textinput.Model
	editing    bool
}

func New(s *styles.Palette) Model {
	ti := textinput.New()
	ti.CharLimit = 500
	ti.Prompt = ""

	return Model{
		filter:    "",
		styles:    s,
		textInput: ti,
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

func (m *Model) StartEditing(currentFilter string) {
	m.textInput.SetValue(currentFilter)
	m.textInput.Focus()
	m.editing = true
}

func (m *Model) StopEditing() string {
	m.editing = false
	m.textInput.Blur()
	return m.textInput.Value()
}

func (m *Model) CancelEditing() {
	m.editing = false
	m.textInput.Blur()
}

func (m Model) IsEditing() bool {
	return m.editing
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.editing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	leftContent := m.styles.HeaderTitle.Render("filter")
	rightContent := m.styles.HeaderText.Padding(0, 1).Render(fmt.Sprintf("%d/%d", m.count, m.totalCount))

	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	middleWidth := m.width - leftWidth - rightWidth

	if middleWidth < 1 {
		middleWidth = 1
	}

	var middleContent string
	if m.editing {
		m.textInput.Width = middleWidth - 1
		middleContent = lipgloss.NewStyle().Width(middleWidth).Render(" " + m.textInput.View())
	} else {
		middleContent = m.styles.StatusBar.Width(middleWidth).Render(m.filter)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, middleContent, rightContent)

	if lipgloss.Width(header) < m.width {
		header = m.styles.Header.Width(m.width).Render(header)
	}

	return header
}
