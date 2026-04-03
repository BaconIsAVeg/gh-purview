package model

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/BaconIsAVeg/gh-purview/internal/github"
	"github.com/BaconIsAVeg/gh-purview/internal/types"
	"github.com/BaconIsAVeg/gh-purview/internal/ui/preview"
	"github.com/BaconIsAVeg/gh-purview/internal/ui/prlist"
	"github.com/BaconIsAVeg/gh-purview/internal/ui/statusbar"
	"github.com/BaconIsAVeg/github-tuis/ui/diffstyles"
	"github.com/BaconIsAVeg/github-tuis/ui/header"
	"github.com/BaconIsAVeg/github-tuis/ui/helpers"
	"github.com/BaconIsAVeg/github-tuis/ui/notification"
	"github.com/BaconIsAVeg/github-tuis/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

type Model struct {
	header       header.Model
	prlist       prlist.Model
	preview      preview.Model
	statusbar    statusbar.Model
	notification notification.Model
	styles       *styles.Palette
	ghClient     *github.Client
	width        int
	height       int
	ready        bool
	loading      bool
}

func New(ghClient *github.Client) Model {
	hasDarkBg := lipgloss.HasDarkBackground()
	s := styles.NewPalette(hasDarkBg)
	ds := diffstyles.NewDiffPalette(hasDarkBg)
	h := header.New(s)
	h.SetLeft("filter")
	sb := statusbar.New(s)
	sb.SetMode(statusbar.ModeList)
	notif := notification.New(s)
	notif.Set("Please wait...", notification.TypeInfo)
	return Model{
		header:       h,
		prlist:       prlist.New(s),
		preview:      preview.New(s, ds),
		statusbar:    sb,
		notification: notif,
		styles:       s,
		ghClient:     ghClient,
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

func (m Model) loadDiff(pr *types.PR) tea.Cmd {
	return func() tea.Msg {
		if m.ghClient == nil || pr == nil {
			return diffLoadedMsg{err: fmt.Errorf("cannot load diff")}
		}
		result, err := m.ghClient.FetchPRDiff(context.Background(), pr)
		if err != nil {
			return diffLoadedMsg{err: err}
		}
		return diffLoadedMsg{
			content:   result.Content,
			truncated: result.Truncated,
			additions: result.Additions,
			deletions: result.Deletions,
		}
	}
}

func (m *Model) openPreview() []tea.Cmd {
	m.preview.SetVisible(true)
	pr := m.prlist.SelectedPR()
	m.preview.SetPR(pr)
	m.statusbar.SetMode(statusbar.ModeDiff)
	m.statusbar.SetStats(0, 0)
	m.updateLayout()
	return []tea.Cmd{
		m.notification.ShowInfo("Loading diff..."),
		m.loadDiff(pr),
	}
}

func (m *Model) closePreview() {
	m.preview.SetVisible(false)
	m.statusbar.SetMode(statusbar.ModeList)
	m.statusbar.SetStats(0, 0)
	m.statusbar.SetScrollPosition("")
	m.updateLayout()
}

func (m *Model) loadDiffForSelectedPR() []tea.Cmd {
	pr := m.prlist.SelectedPR()
	m.preview.SetPR(pr)
	m.statusbar.SetStats(0, 0)
	return []tea.Cmd{
		m.notification.ShowInfo("Loading diff..."),
		m.loadDiff(pr),
	}
}

type prsLoadedMsg struct {
	prs    []types.PR
	total  int
	err    error
	filter string
}

type diffLoadedMsg struct {
	content   string
	truncated bool
	additions int
	deletions int
	err       error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	wasEditing := m.header.IsEditing()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true
	case prsLoadedMsg:
		m.loading = false
		m.notification.Hide()
		m.header.SetMiddle(msg.filter)
		if msg.err != nil {
			cmds = append(cmds, m.notification.ShowWarning(fmt.Sprintf("Error: %v", msg.err)))
		} else {
			m.prlist.SetPRs(msg.prs)
			m.header.SetRight(fmt.Sprintf("%d/%d", len(msg.prs), msg.total))
			m.statusbar.SetMode(statusbar.ModeList)
		}
		m.updateLayout()
	case diffLoadedMsg:
		m.notification.Hide()
		if msg.err != nil {
			m.preview.SetDiffContent(fmt.Sprintf("Error loading diff: %v", msg.err))
			m.statusbar.SetStats(0, 0)
			m.statusbar.SetScrollPosition("")
		} else {
			m.preview.SetDiffContent(msg.content)
			m.statusbar.SetStats(msg.additions, msg.deletions)
			m.statusbar.SetScrollPosition(m.preview.ScrollPosition())
		}
	case approvePRMsg:
		if msg.err != nil {
			cmds = append(cmds, m.notification.ShowWithTimeout(fmt.Sprintf("Error: %v", msg.err), notification.TypeWarning, 2*time.Second))
		} else if msg.pr != nil {
			cmds = append(cmds, m.notification.ShowWithTimeout(fmt.Sprintf("PR #%d approved", msg.pr.Number), notification.TypeSuccess, 2*time.Second))
		}
		m.updateLayout()
	case notification.HideMsg:
		m.updateLayout()
	}

	if m.header.IsEditing() && wasEditing {
		var cmd tea.Cmd
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.notification, cmd = m.notification.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKey(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	if m.header.IsEditing() {
		switch msg.String() {
		case "enter":
			newFilter := m.header.StopEditing()
			m.ghClient.SetQuery(newFilter)
			m.header.SetMiddle(m.ghClient.Query())
			m.loading = true
			cmds = append(cmds, m.notification.ShowInfo("Please wait..."))
			cmds = append(cmds, m.loadPRs())
			m.statusbar.SetMode(statusbar.ModeList)
			return cmds
		case "esc":
			m.header.CancelEditing()
			m.statusbar.SetMode(statusbar.ModeList)
			return cmds
		}
		return cmds
	}

	switch msg.String() {
	case "ctrl+c":
		return []tea.Cmd{tea.Quit}
	case "q":
		if m.preview.Visible() {
			m.closePreview()
		} else {
			return []tea.Cmd{tea.Quit}
		}
	case "esc":
		if m.preview.Visible() {
			m.closePreview()
		}
	case "j", "down":
		m.prlist.CursorDown()
		if m.preview.Visible() {
			cmds = append(cmds, m.loadDiffForSelectedPR()...)
		}
	case "k", "up":
		m.prlist.CursorUp()
		if m.preview.Visible() {
			cmds = append(cmds, m.loadDiffForSelectedPR()...)
		}
	case "p", "enter":
		if !m.preview.Visible() {
			cmds = append(cmds, m.openPreview()...)
		}
	case "r":
		if !m.preview.Visible() {
			m.loading = true
			cmds = append(cmds, m.notification.ShowInfo("Please wait..."))
			cmds = append(cmds, m.loadPRs())
		}
	case "ctrl+a":
		if m.preview.Visible() {
			if pr := m.prlist.SelectedPR(); pr != nil {
				cmds = append(cmds, m.notification.ShowInfo("Approving PR..."))
				cmds = append(cmds, m.approvePR(pr))
			}
		}
	case "o":
		if m.preview.Visible() {
			if pr := m.prlist.SelectedPR(); pr != nil {
				cmds = append(cmds, m.notification.Show("Opening on GitHub.com..."))
				cmds = append(cmds, m.openOnWeb(pr))
			}
		}
	case "ctrl+n":
		if m.preview.Visible() {
			m.preview.ScrollDown(1)
			m.statusbar.SetScrollPosition(m.preview.ScrollPosition())
		}
	case "ctrl+p":
		if m.preview.Visible() {
			m.preview.ScrollUp(1)
			m.statusbar.SetScrollPosition(m.preview.ScrollPosition())
		}
	case "g":
		if m.preview.Visible() {
			m.preview.ScrollToTop()
			m.statusbar.SetScrollPosition(m.preview.ScrollPosition())
		}
	case "G":
		if m.preview.Visible() {
			m.preview.ScrollToBottom()
			m.statusbar.SetScrollPosition(m.preview.ScrollPosition())
		}
	case "f":
		if !m.preview.Visible() {
			currentFilter := m.ghClient.Query()
			m.header.StartEditing(currentFilter)
			m.statusbar.SetMode(statusbar.ModeFilterEdit)
		}
	}

	return cmds
}

func (m Model) approvePR(pr *types.PR) tea.Cmd {
	return func() tea.Msg {
		err := m.ghClient.ApprovePR(context.Background(), pr)
		return approvePRMsg{pr: pr, err: err}
	}
}

func transformURLForMDCA(url string) string {
	if os.Getenv("GH_MDCA") == "" {
		return url
	}
	return strings.Replace(url, "github.com", "github.com.mcas.ms", 1)
}

func (m Model) openOnWeb(pr *types.PR) tea.Cmd {
	return func() tea.Msg {
		if pr != nil && pr.URL != "" {
			browser.Stdout = io.Discard
			browser.Stderr = io.Discard
			browser.OpenURL(transformURLForMDCA(pr.URL))
		}
		return nil
	}
}

type approvePRMsg struct {
	pr  *types.PR
	err error
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
		return helpers.PlaceOverlay(x, y, notifView, base, false, m.styles.ShadowFg)
	}

	return base
}
