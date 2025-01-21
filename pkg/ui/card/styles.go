package card

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style

	Border lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1),
		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Padding(0, 2),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("lightgrey")),
	}
}
