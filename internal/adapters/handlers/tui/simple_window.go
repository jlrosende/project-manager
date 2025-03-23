package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type SimpleWindow struct {
	projectSvc *services.ProjectService

	projects        []*domain.Project
	selectedProject *domain.Project

	cursor    int
	cursorEnv int
	total     int
	width     int
	height    int
}

func NewSimpleWindow(projectSvc *services.ProjectService) (*SimpleWindow, error) {

	projects, err := projectSvc.List()

	if err != nil {
		return nil, err
	}

	total := 0

	if len(projects) > 0 {
		total = len(projects)
	} else {
		total = 1
	}

	return &SimpleWindow{
		projectSvc: projectSvc,
		projects:   projects,
		total:      total,
		cursor:     0,
	}, nil
}

func (m SimpleWindow) Init() tea.Cmd {
	return tea.SetWindowTitle("Project Manager")
}

func (m *SimpleWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.selectedProject = nil
			return m, tea.Quit
		case "enter":
			project, _ := m.projectSvc.Load(m.projects[m.cursor%m.total].Name)
			m.selectedProject = project
			// Todo check project if have environments and create a new view with environment selection
			return m, tea.Quit

			// The "up" and "k" keys move the cursor up
		case "up", "k":
			m.cursor--
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			m.cursor++
		}
	}

	return m, tea.Batch(
		cmd,
		tea.Printf("Let's go to %d!", m.cursor),
	)
}

func (m SimpleWindow) View() string {
	s := strings.Builder{}

	title := " Projects"
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).Padding(0, 1)

	defaultStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).Padding(0, 1)

	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n")

	mod := mod(m.cursor, m.total)

	for i, project := range m.projects {
		if mod == i {
			s.WriteString("->")
			s.WriteString(selectedStyle.Render(project.Name))
			s.WriteString("\n")
			for _, env := range project.Environments {
				s.WriteString("    \u2022")
				s.WriteString(defaultStyle.Foreground(lipgloss.Color(env.Color)).Render(env.Name))
				s.WriteString("\n")
			}
			s.WriteString("\n")
		} else {
			s.WriteString("-")
			s.WriteString(defaultStyle.Render(project.Name))
			s.WriteString("\n")
		}

	}

	// s.WriteString(fmt.Sprintf("Cursor: %d, mod: %d", m.cursor, mod))
	s.WriteString("\n")

	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, s.String())
}

func (m *SimpleWindow) SelectedProject() *domain.Project {
	return m.selectedProject
}

func mod(a, b int) int {
	return (a%b + b) % b
}
