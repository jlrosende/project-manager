package card

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	Title    string
	SubTitle string

	Styles Styles
}

func (m Card) Init() tea.Cmd {
	return nil
}

func (m Card) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Card) View() string {
	var view string

	view = lipgloss.JoinVertical(
		lipgloss.Left,
		m.Styles.Title.Render(m.Title),
		m.Styles.Subtitle.Render(ellipsis(m.SubTitle, 15)),
	)

	view = m.Styles.Border.Render(view)

	return view
}

func NewCard(title, subtitle string) Card {
	styles := DefaultStyles()

	return Card{
		Title:    title,
		SubTitle: subtitle,
		Styles:   styles,
	}
}

func ellipsis(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	if maxLen < 3 {
		maxLen = 3
	}

	return string(runes[:maxLen-3]) + "..."
}
