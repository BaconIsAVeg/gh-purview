package appstyles

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	Add            lipgloss.Style
	Delete         lipgloss.Style
	Header         lipgloss.Style
	FileHeader     lipgloss.Style
	PRNumber       lipgloss.Style
	PRTitle        lipgloss.Style
	PRMeta         lipgloss.Style
	ReviewApproved lipgloss.Style
	ReviewChanges  lipgloss.Style
	ReviewRequired lipgloss.Style
}

type colors struct {
	diffAdd        lipgloss.Color
	diffDelete     lipgloss.Color
	diffHeader     lipgloss.Color
	diffFileHeader lipgloss.Color
	popFg          lipgloss.Color
	primaryFg      lipgloss.Color
	dimFg          lipgloss.Color
	reviewApproved lipgloss.Color
	reviewChanges  lipgloss.Color
	reviewRequired lipgloss.Color
}

func NewPalette(darkBackground bool) *Palette {
	if darkBackground {
		return newPalette(darkColors())
	}
	return newPalette(lightColors())
}

func darkColors() colors {
	return colors{
		diffAdd:        lipgloss.Color("34"),
		diffDelete:     lipgloss.Color("160"),
		diffHeader:     lipgloss.Color("36"),
		diffFileHeader: lipgloss.Color("15"),
		popFg:          lipgloss.Color("178"),
		primaryFg:      lipgloss.Color("15"),
		dimFg:          lipgloss.Color("244"),
		reviewApproved: lipgloss.Color("2"),
		reviewChanges:  lipgloss.Color("1"),
		reviewRequired: lipgloss.Color("33"),
	}
}

func lightColors() colors {
	return colors{
		diffAdd:        lipgloss.Color("28"),
		diffDelete:     lipgloss.Color("124"),
		diffHeader:     lipgloss.Color("24"),
		diffFileHeader: lipgloss.Color("0"),
		popFg:          lipgloss.Color("172"),
		primaryFg:      lipgloss.Color("0"),
		dimFg:          lipgloss.Color("242"),
		reviewApproved: lipgloss.Color("2"),
		reviewChanges:  lipgloss.Color("1"),
		reviewRequired: lipgloss.Color("33"),
	}
}

func newPalette(c colors) *Palette {
	return &Palette{
		Add: lipgloss.NewStyle().
			Foreground(c.diffAdd),
		Delete: lipgloss.NewStyle().
			Foreground(c.diffDelete),
		Header: lipgloss.NewStyle().
			Foreground(c.diffHeader).
			Bold(true),
		FileHeader: lipgloss.NewStyle().
			Foreground(c.diffFileHeader).
			Bold(true),
		PRNumber: lipgloss.NewStyle().
			Foreground(c.popFg).
			Bold(true),
		PRTitle: lipgloss.NewStyle().
			Foreground(c.primaryFg),
		PRMeta: lipgloss.NewStyle().
			Foreground(c.dimFg),
		ReviewApproved: lipgloss.NewStyle().
			Foreground(c.reviewApproved).
			Bold(true),
		ReviewChanges: lipgloss.NewStyle().
			Foreground(c.reviewChanges).
			Bold(true),
		ReviewRequired: lipgloss.NewStyle().
			Foreground(c.reviewRequired).
			Bold(true),
	}
}
