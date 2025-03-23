package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/pkg/ui/card"
)

type Window struct {
	projects []card.Card

	width  int
	height int
}

var _ tea.Model = (*Window)(nil)

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

func (m Window) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("project manager"),
	)
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

	// Help Block
	helpPanel := m.helpPanelRender()

	// Sidebar
	sidebarPanel := lipgloss.Place(
		m.width/5,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		m.sidebarPanelRender(),
	)

	// View
	viewPanel := lipgloss.Place(
		(m.width/5)*4,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		m.viewPanelRender(),
	)

	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left,
		sidebarPanel,
		viewPanel,
	)

	// return mainPanel
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		mainPanel,
		helpPanel,
	)

	return view

}

func (m Window) helpPanelRender() string {
	helpBlockStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777777"))

	return helpBlockStyle.Render("? help")
}

func (m Window) sidebarPanelRender() string {

	var sidebarRender string

	// Title block
	titleBlockStyle := lipgloss.NewStyle().
		Width(m.width/5).
		Foreground(lipgloss.Color("14")).
		Padding(1, 2, 0).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Align(lipgloss.Center).
		Render(" Projects")

	projectsBlockStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Render("[ ] - project 1")

	sidebarRender = lipgloss.JoinVertical(
		lipgloss.Left,
		titleBlockStyle,
		projectsBlockStyle,
	)

	return sidebarRender
}

func (m Window) viewPanelRender() string {

	viewPanelStyle := lipgloss.NewStyle().
		Padding(1, 2, 0).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder(),
			false, false, true, true,
		)

	descriptionTitle := viewPanelStyle.
		Foreground(lipgloss.Color("10")).
		Render("󰦪 Description")

	historyB := "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac.Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac.	Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."

	descriptionText := viewPanelStyle.
		Width((m.width / 5) * 4).
		BorderBottom(false).
		Render(historyB)

	envsTitle := viewPanelStyle.
		Foreground(lipgloss.Color("11")).
		Render("󱄑 Environments")

	// envsList := viewPanelStyle.BorderBottom(false).PaddingTop(0)

	viewBlockStyle := lipgloss.JoinVertical(
		lipgloss.Left,
		descriptionTitle,
		descriptionText,
		envsTitle,
		list.New(
			"dev",
			"pre",
			"pro",
		).Enumerator(list.Roman).Value(),
	)

	return viewBlockStyle
}
