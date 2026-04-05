package statusbar

import (
	"fmt"

	"github.com/BaconIsAVeg/github-tuis/buildinfo"
	"github.com/BaconIsAVeg/github-tuis/ui/statusbar"
	"github.com/BaconIsAVeg/github-tuis/ui/styles"
)

const (
	ModeList       = "L"
	ModeDiff       = "D"
	ModeFilterEdit = "F"
)

type Model struct {
	inner          statusbar.Model
	additions      int
	deletions      int
	scrollPosition string
	version        string
}

func New(s *styles.Palette) Model {
	return Model{
		inner:   statusbar.New(s),
		version: buildinfo.GetVersion(),
	}
}

func (m *Model) SetMode(mode string) {
	m.inner.SetMode(mode)
	m.updateKeybinds(mode)
	m.updateMiddleContent()
}

func (m *Model) SetStats(additions, deletions int) {
	m.additions = additions
	m.deletions = deletions
	m.updateMiddleContent()
}

func (m *Model) SetScrollPosition(pos string) {
	m.scrollPosition = pos
	m.updateMiddleContent()
}

func (m *Model) SetWidth(width int) {
	m.inner.SetWidth(width)
}

func (m *Model) updateKeybinds(mode string) {
	switch mode {
	case ModeDiff:
		m.inner.SetKeybindings([]statusbar.KeyBinding{
			{Key: "^n/^p", Desc: "scroll"},
			{Key: "g/G", Desc: "top/bot"},
			{Key: "^a", Desc: "approve"},
			{Key: "o", Desc: "open on web"},
			{Key: "esc", Desc: "close"},
		})
	case ModeFilterEdit:
		m.inner.SetKeybindings([]statusbar.KeyBinding{
			{Key: "enter", Desc: "apply"},
			{Key: "esc", Desc: "cancel"},
		})
	default:
		m.inner.SetKeybindings([]statusbar.KeyBinding{
			{Key: "j/k", Desc: "navigate"},
			{Key: "p", Desc: "preview"},
			{Key: "f", Desc: "filter"},
			{Key: "r", Desc: "refresh"},
			{Key: "q", Desc: "quit"},
		})
	}
}

func (m *Model) updateMiddleContent() {
	mode := m.inner.Mode()
	if mode == ModeDiff && (m.additions > 0 || m.deletions > 0) {
		text := fmt.Sprintf("+%d -%d ", m.additions, m.deletions)
		if m.scrollPosition != "" {
			text += fmt.Sprintf("[%s] ", m.scrollPosition)
		}
		m.inner.SetMiddleContent(text)
	} else {
		m.inner.SetMiddleContent(m.version)
	}
}

func (m Model) View() string {
	return m.inner.View()
}
