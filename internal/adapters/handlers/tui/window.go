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
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 2)

	width := docStyle.GetWidth()
	// height := docStyle.GetWidth()

	// Help Block
	helpPanel := m.helpPanelRender()

	historyA := "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
	historyB := "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."

	listBlockStyle := lipgloss.NewStyle().
		Width(width/3).
		//Height(height).
		Align(lipgloss.Left).
		Padding(0, 2)

	viewBlockStyle := lipgloss.NewStyle().
		Width((width/3)*2).
		//Height(height).
		Align(lipgloss.Left).
		Padding(0, 2).
		Foreground(lipgloss.Color("#fff"))

	// TODO Debug
	if false {
		docStyle = docStyle.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("grey")).
			Background(lipgloss.Color("#f0f"))
		listBlockStyle = listBlockStyle.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#f00"))

		viewBlockStyle = viewBlockStyle.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#0ff"))
	}

	bodyBlockStyle := lipgloss.JoinHorizontal(lipgloss.Top,
		listBlockStyle.Render(historyA),
		viewBlockStyle.Render(historyB),
	)

	main := lipgloss.JoinVertical(
		lipgloss.Left,
		titleBlockStyle.Render("Projects"),
		bodyBlockStyle,
		helpPanel,
	)

	return main

	// return lipgloss.Place(
	// 	m.width, m.height,
	// 	lipgloss.Left, lipgloss.Top,

	// return docStyle.
	// 	Width(m.width - 2).
	// 	Height(m.height - 2).
	// 	Render("a")

	// return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, cards...))
}

func (w Window) helpPanelRender() string {
	helpBlockStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777777"))

	return helpBlockStyle.Render("help")
}

func (w Window) sidebarPanelRender() string {
	// Title block
	titleBlockStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#f0f")).
		Padding(1, 2).
		Render("Projects")

	return helpBlockStyle.Render("help")
}
