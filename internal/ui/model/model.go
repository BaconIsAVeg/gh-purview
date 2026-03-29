package model

import (
	"context"
	"fmt"

	"github.com/anomaly/ghr/internal/github"
	"github.com/anomaly/ghr/internal/types"
	"github.com/anomaly/ghr/internal/ui/header"
	"github.com/anomaly/ghr/internal/ui/helpers"
	"github.com/anomaly/ghr/internal/ui/notification"
	"github.com/anomaly/ghr/internal/ui/preview"
	"github.com/anomaly/ghr/internal/ui/prlist"
	"github.com/anomaly/ghr/internal/ui/statusbar"
	"github.com/anomaly/ghr/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	header       header.Model
	prlist       prlist.Model
	preview      preview.Model
	statusbar    statusbar.Model
	notification notification.Model
	styles       *styles.Palette
	ghClient     *github.Client
	version      string
	width        int
	height       int
	ready        bool
	loading      bool
}

func New(ghClient *github.Client, version string) Model {
	s := styles.NewPalette()
	notif := notification.New(s)
	notif.Set("Please wait...", notification.TypeInfo)
	return Model{
		header:       header.New(s, version),
		prlist:       prlist.New(s),
		preview:      preview.New(s),
		statusbar:    statusbar.New(s, version),
		notification: notif,
		styles:       s,
		ghClient:     ghClient,
		version:      version,
		loading:      true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadPRs()
}

func (m Model) loadPRs() tea.Cmd {
	return func() tea.Msg {
		if m.ghClient == nil {
			return prsLoadedMsg{err: fmt.Errorf("GitHub client not initialized")}
		}

		filter := m.ghClient.Query()
		prs, total, err := m.ghClient.FetchPRs(context.Background())
		if err != nil {
			return prsLoadedMsg{err: err, filter: filter}
		}
		return prsLoadedMsg{prs: prs, total: total, filter: filter}
	}
}

type prsLoadedMsg struct {
	prs    []types.PR
	total  int
	err    error
	filter string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.preview.Visible() {
				m.preview.Toggle()
				m.statusbar.SetMode("pr list")
				m.updateLayout()
			} else {
				return m, tea.Quit
			}
		case "j", "down":
			m.prlist.CursorDown()
			if m.preview.Visible() {
				m.preview.SetPR(m.prlist.SelectedPR())
			}
		case "k", "up":
			m.prlist.CursorUp()
			if m.preview.Visible() {
				m.preview.SetPR(m.prlist.SelectedPR())
			}
		case "p", "escape":
			m.preview.Toggle()
			if m.preview.Visible() {
				m.preview.SetPR(m.prlist.SelectedPR())
				m.statusbar.SetMode("preview")
			} else {
				m.statusbar.SetMode("pr list")
			}
			m.updateLayout()
		case "r":
			if !m.preview.Visible() {
				m.loading = true
				cmds = append(cmds, m.notification.ShowInfo("Please wait..."))
				cmds = append(cmds, m.loadPRs())
			}
		case "a":
			if m.preview.Visible() {
				pr := m.prlist.SelectedPR()
				if pr != nil {
					cmds = append(cmds, m.approvePR(pr))
				}
			}
		case "o":
			if m.preview.Visible() {
				pr := m.prlist.SelectedPR()
				if pr != nil {
					cmds = append(cmds, m.openOnWeb(pr))
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true
	case prsLoadedMsg:
		m.loading = false
		m.notification.Hide()
		m.header.SetFilter(msg.filter)
		if msg.err != nil {
			cmds = append(cmds, m.notification.ShowWarning(fmt.Sprintf("Error: %v", msg.err)))
		} else {
			m.prlist.SetPRs(msg.prs)
			m.header.SetCount(len(msg.prs), msg.total)
		}
		m.updateLayout()
	case approvePRMsg:
		if msg.pr != nil {
			msgText := fmt.Sprintf("PR #%d approved", msg.pr.Number)
			cmds = append(cmds, m.notification.Show(msgText))
			m.updateLayout()
		}
	case notification.HideMsg:
		m.updateLayout()
	}

	var cmd tea.Cmd
	m.notification, cmd = m.notification.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) approvePR(pr *types.PR) tea.Cmd {
	return func() tea.Msg {
		return approvePRMsg{pr: pr}
	}
}

func (m Model) openOnWeb(pr *types.PR) tea.Cmd {
	return func() tea.Msg {
		return openOnWebMsg{pr: pr}
	}
}

type approvePRMsg struct {
	pr *types.PR
}

type openOnWebMsg struct {
	pr *types.PR
}

func (m *Model) updateLayout() {
	headerHeight := 1
	statusbarHeight := 1

	m.header.SetWidth(m.width)
	m.statusbar.SetWidth(m.width)

	listHeight := m.height - headerHeight - statusbarHeight

	if m.preview.Visible() {
		previewHeight := int(float64(listHeight) * 0.75)
		listHeight = listHeight - previewHeight
		m.preview.SetWidth(m.width)
		m.preview.SetHeight(previewHeight)
	}

	m.prlist.SetWidth(m.width)
	m.prlist.SetHeight(listHeight)
	m.prlist.EnsureCursorVisible()
}

func (m Model) View() string {
	if !m.ready {
		notifView := m.notification.View()
		if notifView != "" {
			return notifView
		}
		return "Loading..."
	}

	headerView := m.header.View()
	listView := m.prlist.View()
	statusView := m.statusbar.View()

	var mainContent string
	if m.preview.Visible() {
		previewView := m.preview.View()
		mainContent = lipgloss.JoinVertical(lipgloss.Left, listView, previewView)
	} else {
		mainContent = listView
	}

	base := lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		mainContent,
		statusView,
	)

	if m.notification.Visible() {
		notifView := m.notification.View()
		notifWidth := lipgloss.Width(notifView)
		x := m.width - notifWidth - 1
		y := m.height - 2
		return helpers.PlaceOverlay(x, y, notifView, base, false)
	}

	return base
}
