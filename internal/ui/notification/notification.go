package notification

import (
	"time"

	"github.com/BaconIsAVeg/gh-purview/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type Type int

const (
	TypeSuccess Type = iota
	TypeInfo
	TypeWarning
)

type Model struct {
	message   string
	notifType Type
	visible   bool
	startTime time.Time
	duration  time.Duration
	styles    *styles.Palette
}

func New(s *styles.Palette) Model {
	return Model{
		duration: 3 * time.Second,
		styles:   s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Show(message string) tea.Cmd {
	return m.ShowWithType(message, TypeSuccess)
}

func (m *Model) ShowInfo(message string) tea.Cmd {
	return m.ShowWithType(message, TypeInfo)
}

func (m *Model) ShowWarning(message string) tea.Cmd {
	return m.ShowWithType(message, TypeWarning)
}

func (m *Model) ShowWithType(message string, notifType Type) tea.Cmd {
	m.message = message
	m.notifType = notifType
	m.visible = true
	m.startTime = time.Now()
	return tea.Tick(m.duration, func(time.Time) tea.Msg {
		return HideMsg{}
	})
}

func (m *Model) ShowWithTimeout(message string, notifType Type, timeout time.Duration) tea.Cmd {
	m.message = message
	m.notifType = notifType
	m.visible = true
	m.startTime = time.Now()
	return tea.Tick(timeout, func(time.Time) tea.Msg {
		return HideMsg{}
	})
}

func (m *Model) Set(message string, notifType Type) {
	m.message = message
	m.notifType = notifType
	m.visible = true
	m.startTime = time.Now()
}

func (m *Model) Hide() {
	m.visible = false
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case HideMsg:
		m.visible = false
	}
	return m, nil
}

func (m Model) View() string {
	if !m.visible {
		return ""
	}

	style := m.styles.Notification
	switch m.notifType {
	case TypeInfo:
		style = m.styles.NotificationInfo
	case TypeWarning:
		style = m.styles.NotificationWarn
	}

	return style.Render(" " + m.message + " ")
}

func (m Model) Visible() bool {
	return m.visible
}

type HideMsg struct{}
