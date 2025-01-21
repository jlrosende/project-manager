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

func (m Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Card) View() string {
	var view string

	view = lipgloss.JoinVertical(
		lipgloss.Left,
		m.Styles.Title.Render(m.Title),
		m.Styles.Subtitle.Render(m.SubTitle),
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
