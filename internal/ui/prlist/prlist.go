package prlist

import (
	"fmt"
	"strings"

	"github.com/anomaly/ghr/internal/github"
	"github.com/anomaly/ghr/internal/types"
	"github.com/anomaly/ghr/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	prs      []types.PR
	cursor   int
	offset   int
	width    int
	height   int
	styles   *styles.Palette
	viewport viewport.Model
}

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Preview key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:      key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		Down:    key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Preview: key.NewBinding(key.WithKeys("P"), key.WithHelp("P", "preview")),
		Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func New(s *styles.Palette) Model {
	return Model{
		styles:   s,
		viewport: viewport.New(0, 0),
	}
}

func (m *Model) SetPRs(prs []types.PR) {
	m.prs = prs
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.viewport.Width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
	m.viewport.Height = height
}

func (m *Model) EnsureCursorVisible() {
	visibleItems := m.height / 2
	if visibleItems < 1 {
		visibleItems = 1
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visibleItems {
		m.offset = m.cursor - visibleItems + 1
	}
}

func (m *Model) SelectedPR() *types.PR {
	if len(m.prs) == 0 || m.cursor < 0 || m.cursor >= len(m.prs) {
		return nil
	}
	return &m.prs[m.cursor]
}

func (m *Model) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
	}
}

func (m *Model) CursorDown() {
	if m.cursor < len(m.prs)-1 {
		m.cursor++
		visibleItems := m.height / 2
		if m.cursor >= m.offset+visibleItems {
			m.offset = m.cursor - visibleItems + 1
		}
	}
}

func (m Model) View() string {
	if len(m.prs) == 0 {
		var b strings.Builder
		b.WriteString(m.styles.PRMeta.Render(""))
		for i := 1; i < m.height; i++ {
			b.WriteString("\n")
		}
		return b.String()
	}

	var b strings.Builder
	visibleItems := m.height / 2
	if visibleItems < 1 {
		visibleItems = 1
	}

	itemCount := 0
	for i := m.offset; i < len(m.prs) && i < m.offset+visibleItems; i++ {
		pr := m.prs[i]
		isSelected := i == m.cursor

		line1 := m.renderLine1(pr, isSelected)
		line2 := m.renderLine2(pr, isSelected)

		b.WriteString(line1 + "\n" + line2)
		itemCount++
		if i < len(m.prs)-1 && i < m.offset+visibleItems-1 {
			b.WriteString("\n")
		}
	}

	linesRendered := itemCount * 2
	for i := linesRendered; i < m.height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderLine1(pr types.PR, selected bool) string {
	cursor := "  "
	if selected {
		cursor = m.styles.PRNumber.Render(" ▶")
	}
	num := m.styles.PRNumber.Render(fmt.Sprintf("#%d", pr.Number))
	title := m.styles.PRTitle.Render(truncate(pr.Title, m.width-14))
	return lipgloss.JoinHorizontal(lipgloss.Left, cursor, " ", num, " ", title)
}

func (m Model) renderLine2(pr types.PR, selected bool) string {
	repo := m.styles.PRMeta.Render(pr.RepoPath())
	author := m.styles.PRMeta.Render(pr.Author)
	review := m.renderReviewDecision(pr.ReviewDecision)
	status := m.renderStatus(pr.Status)
	return lipgloss.JoinHorizontal(lipgloss.Left, "    ", repo, " ", author, " ", status, review)
}

func (m Model) renderReviewDecision(decision string) string {
	switch decision {
	case string(github.ReviewApproved):
		return m.styles.ReviewApproved.Render(" ✓")
	case string(github.ReviewChangesRequested):
		return m.styles.ReviewChanges.Render(" ✗")
	case string(github.ReviewRequired):
		return m.styles.ReviewRequired.Render(" ~")
	default:
		return ""
	}
}

func (m Model) renderStatus(status types.PRStatus) string {
	switch status {
	case types.StatusOpen:
		return m.styles.StatusOpen.Render("OPEN")
	case types.StatusClosed:
		return m.styles.StatusClosed.Render("CLOSED")
	case types.StatusMerged:
		return m.styles.StatusMerged.Render("MERGED")
	default:
		return string(status)
	}
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) Height() int {
	return m.height
}
