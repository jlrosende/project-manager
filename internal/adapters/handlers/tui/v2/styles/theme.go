package styles

import "github.com/charmbracelet/lipgloss"

type Tokens struct {
	Title           lipgloss.Color
	Section         lipgloss.Color
	Subtext         lipgloss.Color
	Text            lipgloss.Color
	Placeholder     lipgloss.Color
	Border          lipgloss.Color
	Error           lipgloss.Color
	ButtonDefFg     lipgloss.Color
	ButtonDefBg     lipgloss.Color
	ButtonSelFg     lipgloss.Color
	ButtonSelBg     lipgloss.Color
	SelectedFg      lipgloss.Color
	SelectedBg      lipgloss.Color
	Help            lipgloss.Color
	ButtonSuccessFg lipgloss.Color
	ButtonSuccessBg lipgloss.Color
	ButtonInfoFg    lipgloss.Color
	ButtonInfoBg    lipgloss.Color
	ButtonWarningFg lipgloss.Color
	ButtonWarningBg lipgloss.Color
	ButtonDangerFg  lipgloss.Color
	ButtonDangerBg  lipgloss.Color
}

type Theme struct {
	Title           lipgloss.Style
	Selected        lipgloss.Style
	Default         lipgloss.Style
	Help            lipgloss.Style
	Error           lipgloss.Style
	ButtonPrimary   lipgloss.Style
	ButtonSecondary lipgloss.Style
	ButtonSuccess   lipgloss.Style
	ButtonInfo      lipgloss.Style
	ButtonWarning   lipgloss.Style
	ButtonDanger    lipgloss.Style
}

func NewTokensFromHex(p map[string]string) Tokens {
	return Tokens{
		Title:           lipgloss.Color(p["title"]),
		Section:         lipgloss.Color(p["section"]),
		Subtext:         lipgloss.Color(p["subtext"]),
		Text:            lipgloss.Color(p["text"]),
		Placeholder:     lipgloss.Color(p["placeholder"]),
		Border:          lipgloss.Color(p["border"]),
		Error:           lipgloss.Color(p["error"]),
		ButtonDefFg:     lipgloss.Color(p["buttonDefFg"]),
		ButtonDefBg:     lipgloss.Color(p["buttonDefBg"]),
		ButtonSelFg:     lipgloss.Color(p["buttonSelFg"]),
		ButtonSelBg:     lipgloss.Color(p["buttonSelBg"]),
		SelectedFg:      lipgloss.Color(p["selectedFg"]),
		SelectedBg:      lipgloss.Color(p["selectedBg"]),
		Help:            lipgloss.Color(p["help"]),
		ButtonSuccessFg: lipgloss.Color(p["buttonSuccessFg"]),
		ButtonSuccessBg: lipgloss.Color(p["buttonSuccessBg"]),
		ButtonInfoFg:    lipgloss.Color(p["buttonInfoFg"]),
		ButtonInfoBg:    lipgloss.Color(p["buttonInfoBg"]),
		ButtonWarningFg: lipgloss.Color(p["buttonWarningFg"]),
		ButtonWarningBg: lipgloss.Color(p["buttonWarningBg"]),
		ButtonDangerFg:  lipgloss.Color(p["buttonDangerFg"]),
		ButtonDangerBg:  lipgloss.Color(p["buttonDangerBg"]),
	}
}

func BuildTheme(t Tokens) Theme {
	return Theme{
		Title: lipgloss.NewStyle().
			Foreground(t.Title).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(t.Border).
			Padding(0, 1),
		Selected:        lipgloss.NewStyle().Foreground(t.SelectedFg).Background(t.SelectedBg),
		Default:         lipgloss.NewStyle().Foreground(t.Subtext).Padding(0, 1),
		Help:            lipgloss.NewStyle().Foreground(t.Help),
		Error:           lipgloss.NewStyle().Foreground(t.Error),
		ButtonPrimary:   lipgloss.NewStyle().Foreground(t.ButtonSelFg).Background(t.ButtonSelBg),
		ButtonSecondary: lipgloss.NewStyle().Foreground(t.ButtonDefFg).Background(t.Error),
		ButtonSuccess:   lipgloss.NewStyle().Foreground(t.ButtonSuccessFg).Background(t.ButtonSuccessBg),
		ButtonInfo:      lipgloss.NewStyle().Foreground(t.ButtonInfoFg).Background(t.ButtonInfoBg),
		ButtonWarning:   lipgloss.NewStyle().Foreground(t.ButtonWarningFg).Background(t.ButtonWarningBg),
		ButtonDanger:    lipgloss.NewStyle().Foreground(t.ButtonDangerFg).Background(t.ButtonDangerBg),
	}
}
