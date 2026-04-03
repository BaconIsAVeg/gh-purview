package prlist

import (
	"fmt"
	"strings"

	"github.com/BaconIsAVeg/gh-purview/internal/github"
	"github.com/BaconIsAVeg/gh-purview/internal/types"
	"github.com/BaconIsAVeg/github-tuis/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// linesPerPR is the number of lines each PR item occupies in the list
	linesPerPR = 2
	// titlePadding is the character width reserved for cursor, number, and spacing
	titlePadding = 14
)

type Model struct {
	prs    []types.PR
	cursor int
	offset int
	width  int
	height int
	styles *styles.Palette
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
		styles: s,
	}
}

func (m *Model) SetPRs(prs []types.PR) {
	m.prs = prs
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

// visiblePRCount returns the number of PRs that can fit in the current height.
// Each PR occupies linesPerPR lines.
func (m *Model) visiblePRCount() int {
	count := m.height / linesPerPR
	if count < 1 {
		return 1
	}
	return count
}

// clampOffset ensures the offset stays within valid bounds for cursor visibility.
func (m *Model) clampOffset() {
	visible := m.visiblePRCount()
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m *Model) EnsureCursorVisible() {
	m.clampOffset()
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
		m.clampOffset()
	}
}

func (m *Model) CursorDown() {
	if m.cursor < len(m.prs)-1 {
		m.cursor++
		m.clampOffset()
	}
}

func (m Model) View() string {
	if len(m.prs) == 0 {
		return m.renderEmptyState()
	}

	var b strings.Builder
	visibleCount := m.visiblePRCount()
	endIndex := min(m.offset+visibleCount, len(m.prs))

	for i := m.offset; i < endIndex; i++ {
		pr := m.prs[i]
		selected := i == m.cursor

		b.WriteString(m.renderLine1(pr, selected))
		b.WriteByte('\n')
		b.WriteString(m.renderLine2(pr, selected))

		if i < endIndex-1 {
			b.WriteByte('\n')
		}
	}

	// Pad remaining lines to fill height
	linesRendered := (endIndex - m.offset) * linesPerPR
	for i := linesRendered; i < m.height; i++ {
		b.WriteByte('\n')
	}

	return b.String()
}

// renderEmptyState renders the list when there are no PRs to display.
func (m Model) renderEmptyState() string {
	var b strings.Builder
	b.WriteString(m.styles.PRMeta.Render(""))
	for i := 1; i < m.height; i++ {
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Model) renderLine1(pr types.PR, selected bool) string {
	cursor := "  "
	if selected {
		cursor = m.styles.PRNumber.Render(" ┌")
	}
	num := m.styles.PRNumber.Render(fmt.Sprintf("#%d", pr.Number))
	title := m.styles.PRTitle.Render(truncate(pr.Title, m.width-titlePadding))
	return lipgloss.JoinHorizontal(lipgloss.Left, cursor, " ", num, " ", title)
}

func (m Model) renderLine2(pr types.PR, selected bool) string {
	cursor := "   "
	if selected {
		cursor = m.styles.PRNumber.Render(" └ ")
	}
	repo := m.styles.PRMeta.Render(pr.RepoPath())
	author := m.styles.PRMeta.Render(pr.Author)
	review := m.renderReviewDecision(pr.ReviewDecision)
	status := m.renderStatus(pr.Status)
	return lipgloss.JoinHorizontal(lipgloss.Left, cursor, repo, " ", author, " ", status, review)
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
