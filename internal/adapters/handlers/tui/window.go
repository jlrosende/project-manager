package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/pkg/ui/textinput"
)

type Window struct {
	projectSvc *services.ProjectService

	projects        []*domain.Project
	selectedProject *domain.Project

	cursor    int
	cursorEnv int
	focus     int // 0=projects, 1=environments
	total     int
	width     int
	height    int
}

func NewWindow(projectSvc *services.ProjectService) (*Window, error) {

	projects, err := projectSvc.List()

	if err != nil {
		return nil, err
	}

	total := len(projects) + 1
	if total == 0 {
		total = 1
	}

	return &Window{
		projectSvc: projectSvc,
		projects:   projects,
		total:      total,
		cursor:     0,
	}, nil
}

func (m Window) Init() tea.Cmd {
	return tea.SetWindowTitle("Project Manager")
}

func (m *Window) newProjectFlow() {
	name := ""
	if m2, err := tea.NewProgram(textinput.NewTextInput("Project name:", "my-project")).Run(); err == nil {
		if ti, ok := m2.(*textinput.TextInput); ok {
			name = ti.Value()
		}
	}
	if name == "" {
		return
	}
	path := ""
	if m3, err := tea.NewProgram(textinput.NewTextInput("Path:", "~/code/my-project")).Run(); err == nil {
		if ti, ok := m3.(*textinput.TextInput); ok {
			path = ti.Value()
		}
	}
	if path == "" {
		return
	}
	_, _ = m.projectSvc.Create(name, path, "", domain.EnvVars{}, domain.New())
	projects, err := m.projectSvc.List()
	if err == nil {
		m.projects = projects
		m.total = len(projects) + 1
	}
}

func (m *Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focus == 0 {
				idx := mod(m.cursor, m.total)
				if idx == len(m.projects) { // '+ New project'
					m.newProjectFlow()
					return m, nil
				}
				project, _ := m.projectSvc.Load(m.projects[idx].Name)
				m.selectedProject = project
				return m, tea.Quit
			}
			if m.focus == 1 {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 && len(m.projects[mod].Environments) > 0 {
					project, _ := m.projectSvc.Load(m.projects[mod].Name)
					m.selectedProject = project
					return m, tea.Quit
				}
			}

			// The "up" and "k" keys move the cursor up
		case "left", "h":
			m.focus = 0
		case "right", "l":
			mod := mod(m.cursor, m.total)
			if mod < len(m.projects) && len(m.projects) > 0 && len(m.projects[mod].Environments) > 0 {
				m.focus = 1
			}
		case "up", "k":
			if m.focus == 0 {
				m.cursor--
				m.cursorEnv = 0
			} else {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					envCount := len(m.projects[mod].Environments)
					if envCount > 0 {
						m.cursorEnv--
						if m.cursorEnv < 0 {
							m.cursorEnv = envCount - 1
						}
					}
				}
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.focus == 0 {
				m.cursor++
				m.cursorEnv = 0
			} else {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					envCount := len(m.projects[mod].Environments)
					if envCount > 0 {
						m.cursorEnv++
						if m.cursorEnv >= envCount {
							m.cursorEnv = 0
						}
					}
				}
			}
		}
	}

	return m, tea.Batch(
		cmd,
		tea.Printf("Let's go to %d!", m.cursor),
	)
}

func (m Window) View() string {
	left := strings.Builder{}
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
	left.WriteString(titleStyle.Render(title))
	left.WriteString("\n")
	mod := mod(m.cursor, m.total)
	for i, project := range m.projects {
		if mod == i {
			left.WriteString("->")
			left.WriteString(selectedStyle.Render(project.Name))
			left.WriteString("\n")
		} else {
			left.WriteString(" ")
			left.WriteString(defaultStyle.Render(project.Name))
			left.WriteString("\n")
		}
	}
	if mod == len(m.projects) {
		left.WriteString("->")
		left.WriteString(selectedStyle.Render("+ New project"))
		left.WriteString("\n")
	} else {
		left.WriteString(" ")
		left.WriteString(defaultStyle.Render("+ New project"))
		left.WriteString("\n")
	}
	right := strings.Builder{}
	rightTitle := "󱄑 Environments"
	right.WriteString(titleStyle.Render(rightTitle))
	right.WriteString("\n")
	if mod < len(m.projects) && len(m.projects) > 0 {
		p := m.projects[mod]
		if len(p.Environments) == 0 {
			right.WriteString(defaultStyle.Render("No environments"))
			right.WriteString("\n")
		} else {
			for i, env := range p.Environments {
				marker := "  "
				selected := m.focus == 1 && m.cursorEnv%max(1, len(p.Environments)) == i
				if selected {
					marker = "->"
				}
				rowStyle := lipgloss.NewStyle()
				if selected {
					rowStyle = rowStyle.Background(lipgloss.Color(env.Color)).Padding(0, 0, 0, 1)
				}
				content := defaultStyle.Foreground(lipgloss.Color(env.Color)).Render("● " + env.Name)
				if p.DefaultEnv != "" && env.Name == p.DefaultEnv {
					content = content + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(default)")
				}
				right.WriteString(marker)
				right.WriteString(rowStyle.Render(content))
				right.WriteString("\n")
			}
		}
	}
	content := lipgloss.JoinHorizontal(lipgloss.Left, left.String(), "  ", right.String())
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Keys: ↑/k ↓/j navigate  ←/h →/l focus  Enter select  Esc/Ctrl+C exit")
	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

func (m *Window) SelectedProject() *domain.Project {
	return m.selectedProject
}

func (m *Window) SelectedEnvironment() string {
	if len(m.projects) == 0 {
		return ""
	}
	idx := mod(m.cursor, m.total)
	p := m.projects[idx]
	if len(p.Environments) == 0 {
		return ""
	}
	e := m.cursorEnv % len(p.Environments)
	if e < 0 {
		e = 0
	}
	return p.Environments[e].Name
}

func mod(a, b int) int { return (a%b + b) % b }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
