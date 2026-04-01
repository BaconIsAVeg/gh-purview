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

type colors struct {
	primaryBg      lipgloss.Color
	secondaryBg    lipgloss.Color
	primaryFg      lipgloss.Color
	dimFg          lipgloss.Color
	popFg          lipgloss.Color
	notificationBg lipgloss.Color
	notificationFg lipgloss.Color
	infoBg         lipgloss.Color
	warnBg         lipgloss.Color
	shadowFg       lipgloss.Color
	statusOpen     lipgloss.Color
	statusClosed   lipgloss.Color
	statusMerged   lipgloss.Color
	reviewApproved lipgloss.Color
	reviewChanges  lipgloss.Color
	reviewRequired lipgloss.Color
	diffAdd        lipgloss.Color
	diffDelete     lipgloss.Color
	diffHeader     lipgloss.Color
	diffFileHeader lipgloss.Color
}

func NewPalette(darkBackground bool) *Palette {
	if darkBackground {
		return newPalette(darkColors())
	}
	return newPalette(lightColors())
}

func darkColors() colors {
	return colors{
		primaryBg:      lipgloss.Color("97"),
		secondaryBg:    lipgloss.Color("234"),
		primaryFg:      lipgloss.Color("15"),
		dimFg:          lipgloss.Color("244"),
		popFg:          lipgloss.Color("178"),
		notificationBg: lipgloss.Color("28"),
		notificationFg: lipgloss.Color("15"),
		infoBg:         lipgloss.Color("25"),
		warnBg:         lipgloss.Color("124"),
		shadowFg:       lipgloss.Color("#333333"),
		statusOpen:     lipgloss.Color("2"),
		statusClosed:   lipgloss.Color("5"),
		statusMerged:   lipgloss.Color("4"),
		reviewApproved: lipgloss.Color("2"),
		reviewChanges:  lipgloss.Color("1"),
		reviewRequired: lipgloss.Color("33"),
		diffAdd:        lipgloss.Color("34"),
		diffDelete:     lipgloss.Color("160"),
		diffHeader:     lipgloss.Color("36"),
		diffFileHeader: lipgloss.Color("15"),
	}
}

func lightColors() colors {
	return colors{
		primaryBg:      lipgloss.Color("97"),
		secondaryBg:    lipgloss.Color("254"),
		primaryFg:      lipgloss.Color("0"),
		dimFg:          lipgloss.Color("242"),
		popFg:          lipgloss.Color("172"),
		notificationBg: lipgloss.Color("28"),
		notificationFg: lipgloss.Color("15"),
		infoBg:         lipgloss.Color("25"),
		warnBg:         lipgloss.Color("124"),
		shadowFg:       lipgloss.Color("#cccccc"),
		statusOpen:     lipgloss.Color("2"),
		statusClosed:   lipgloss.Color("1"),
		statusMerged:   lipgloss.Color("4"),
		reviewApproved: lipgloss.Color("2"),
		reviewChanges:  lipgloss.Color("1"),
		reviewRequired: lipgloss.Color("33"),
		diffAdd:        lipgloss.Color("28"),
		diffDelete:     lipgloss.Color("124"),
		diffHeader:     lipgloss.Color("24"),
		diffFileHeader: lipgloss.Color("0"),
	}
}

func newPalette(c colors) *Palette {
	return &Palette{
		Header: lipgloss.NewStyle().
			Background(c.primaryBg).
			Foreground(c.primaryFg).
			Padding(0, 1),
		HeaderTitle: lipgloss.NewStyle().
			Background(c.primaryBg).
			Foreground(c.primaryFg).
			Bold(true).
			Padding(0, 1),
		HeaderText: lipgloss.NewStyle().
			Background(c.primaryBg).
			Foreground(c.primaryFg),
		StatusBar: lipgloss.NewStyle().
			Background(c.secondaryBg).
			Foreground(c.primaryFg).
			Padding(0, 1),
		StatusMode: lipgloss.NewStyle().
			Background(c.primaryBg).
			Foreground(c.primaryFg).
			Bold(true).
			Padding(0, 1),
		StatusKey: lipgloss.NewStyle().
			Background(c.secondaryBg).
			Foreground(c.popFg).
			Bold(true),
		StatusDesc: lipgloss.NewStyle().
			Background(c.secondaryBg).
			Foreground(c.dimFg),
		StatusSep: lipgloss.NewStyle().
			Background(c.secondaryBg),
		PRNumber: lipgloss.NewStyle().
			Foreground(c.popFg).
			Bold(true),
		PRTitle: lipgloss.NewStyle().
			Foreground(c.primaryFg),
		PRMeta: lipgloss.NewStyle().
			Foreground(c.dimFg),
		PreviewBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.dimFg).
			BorderBottom(false),
		PreviewContent: lipgloss.NewStyle().
			Foreground(c.primaryFg),
		StatusOpen: lipgloss.NewStyle().
			Foreground(c.statusOpen).
			Bold(true),
		StatusClosed: lipgloss.NewStyle().
			Foreground(c.statusClosed).
			Bold(true),
		StatusMerged: lipgloss.NewStyle().
			Foreground(c.statusMerged).
			Bold(true),
		ReviewApproved: lipgloss.NewStyle().
			Foreground(c.reviewApproved).
			Bold(true),
		ReviewChanges: lipgloss.NewStyle().
			Foreground(c.reviewChanges).
			Bold(true),
		ReviewRequired: lipgloss.NewStyle().
			Foreground(c.reviewRequired).
			Bold(true),
		Notification: lipgloss.NewStyle().
			Background(c.notificationBg).
			Foreground(c.notificationFg).
			Padding(0, 1).
			Bold(true),
		NotificationInfo: lipgloss.NewStyle().
			Background(c.infoBg).
			Foreground(c.notificationFg).
			Padding(0, 1).
			Bold(true),
		NotificationWarn: lipgloss.NewStyle().
			Background(c.warnBg).
			Foreground(c.notificationFg).
			Padding(0, 1).
			Bold(true),
		DiffAdd: lipgloss.NewStyle().
			Foreground(c.diffAdd),
		DiffDelete: lipgloss.NewStyle().
			Foreground(c.diffDelete),
		DiffHeader: lipgloss.NewStyle().
			Foreground(c.diffHeader).
			Bold(true),
		DiffFileHeader: lipgloss.NewStyle().
			Foreground(c.diffFileHeader).
			Bold(true),
		SecondaryBg: c.secondaryBg,
		ShadowFg:    c.shadowFg,
	}
}
