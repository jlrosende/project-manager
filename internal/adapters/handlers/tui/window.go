package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/pkg/ui/card"
)

type Window struct {
	projects []card.Card

	width  int
	height int
}

func (m Window) Init() tea.Cmd {
	return nil
}

func (m Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m Window) View() string {

	docStyle := lipgloss.NewStyle().
		Padding(0, 2) /*.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("grey"))*/

	var cards []string

	for _, project := range m.projects {
		cards = append(cards, project.View())
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Top,
		docStyle.
			Width(m.width-2).
			Height(m.height-2).
			Render(lipgloss.JoinVertical(lipgloss.Left, cards...)),
	)
	// return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, cards...))
}

func NewWindow(projectSvc *services.ProjectService) (*Window, error) {
	projects, err := projectSvc.List()

	if err != nil {
		return nil, err
	}

	projectCards := []card.Card{}

	for _, project := range projects {
		projectCards = append(projectCards, card.NewCard(project.Name, project.Description))
	}

	return &Window{
		projects: projectCards,
	}, nil
}
