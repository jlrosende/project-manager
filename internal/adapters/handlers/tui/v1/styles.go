package v1

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title    lipgloss.Style
	Selected lipgloss.Style
	Default  lipgloss.Style
	Help     lipgloss.Style
}

func BuildStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(c("title")).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(c("border")).
			Padding(0, 1),
		Selected: lipgloss.NewStyle().Foreground(c("selectedFg")).Background(c("selectedBg")).Padding(0, 1),
		Default:  lipgloss.NewStyle().Foreground(c("subtext")).Padding(0, 1),
		Help:     lipgloss.NewStyle().Foreground(c("help")),
	}
}
