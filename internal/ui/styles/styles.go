package styles

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	Header           lipgloss.Style
	HeaderTitle      lipgloss.Style
	HeaderText       lipgloss.Style
	StatusBar        lipgloss.Style
	StatusMode       lipgloss.Style
	StatusKey        lipgloss.Style
	StatusDesc       lipgloss.Style
	StatusSep        lipgloss.Style
	PRNumber         lipgloss.Style
	PRTitle          lipgloss.Style
	PRMeta           lipgloss.Style
	PreviewBorder    lipgloss.Style
	PreviewContent   lipgloss.Style
	StatusOpen       lipgloss.Style
	StatusClosed     lipgloss.Style
	StatusMerged     lipgloss.Style
	ReviewApproved   lipgloss.Style
	ReviewChanges    lipgloss.Style
	ReviewRequired   lipgloss.Style
	Notification     lipgloss.Style
	NotificationInfo lipgloss.Style
	NotificationWarn lipgloss.Style
	DiffAdd          lipgloss.Style
	DiffDelete       lipgloss.Style
	DiffHeader       lipgloss.Style
	DiffFileHeader   lipgloss.Style
	SecondaryBg      lipgloss.Color
	ShadowFg         lipgloss.Color
}

func NewPalette() *Palette {
	primaryBg := lipgloss.Color("97")
	secondaryBg := lipgloss.Color("234")
	primaryFg := lipgloss.Color("15")
	dimFg := lipgloss.Color("244")
	popFg := lipgloss.Color("178")
	notificationBg := lipgloss.Color("28")
	notificationFg := lipgloss.Color("15")
	infoBg := lipgloss.Color("25")
	warnBg := lipgloss.Color("124")
	shadowFg := lipgloss.Color("#333333")

	return &Palette{
		Header: lipgloss.NewStyle().
			Background(primaryBg).
			Foreground(primaryFg).
			Padding(0, 1),
		HeaderTitle: lipgloss.NewStyle().
			Background(primaryBg).
			Foreground(primaryFg).
			Bold(true).
			Padding(0, 1),
		HeaderText: lipgloss.NewStyle().
			Background(primaryBg).
			Foreground(primaryFg),
		StatusBar: lipgloss.NewStyle().
			Background(secondaryBg).
			Foreground(primaryFg).
			Padding(0, 1),
		StatusMode: lipgloss.NewStyle().
			Background(primaryBg).
			Foreground(primaryFg).
			Bold(true).
			Padding(0, 1),
		StatusKey: lipgloss.NewStyle().
			Background(secondaryBg).
			Foreground(popFg).
			Bold(true),
		StatusDesc: lipgloss.NewStyle().
			Background(secondaryBg).
			Foreground(dimFg),
		StatusSep: lipgloss.NewStyle().
			Background(secondaryBg),
		PRNumber: lipgloss.NewStyle().
			Foreground(popFg).
			Bold(true),
		PRTitle: lipgloss.NewStyle().
			Foreground(primaryFg),
		PRMeta: lipgloss.NewStyle().
			Foreground(dimFg),
		PreviewBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dimFg).
			BorderBottom(false),
		PreviewContent: lipgloss.NewStyle().
			Foreground(primaryFg),
		StatusOpen: lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true),
		StatusClosed: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true),
		StatusMerged: lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Bold(true),
		ReviewApproved: lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true),
		ReviewChanges: lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true),
		ReviewRequired: lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true),
		Notification: lipgloss.NewStyle().
			Background(notificationBg).
			Foreground(notificationFg).
			Padding(0, 1).
			Bold(true),
		NotificationInfo: lipgloss.NewStyle().
			Background(infoBg).
			Foreground(notificationFg).
			Padding(0, 1).
			Bold(true),
		NotificationWarn: lipgloss.NewStyle().
			Background(warnBg).
			Foreground(notificationFg).
			Padding(0, 1).
			Bold(true),
		DiffAdd: lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")),
		DiffDelete: lipgloss.NewStyle().
			Foreground(lipgloss.Color("160")),
		DiffHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color("36")).
			Bold(true),
		DiffFileHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true),
		SecondaryBg: secondaryBg,
		ShadowFg:    shadowFg,
	}
}
