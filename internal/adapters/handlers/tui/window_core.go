package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

type Window struct {
	projectSvc *services.ProjectService

	projects        []*domain.Project
	selectedProject *domain.Project

	cursor    int
	cursorEnv int
	focus     int
	total     int
	width     int
	height    int
	shouldQuit bool
	mode       int
	form       *NewProjectForm
	prompt     *postCreatePrompt
	envForm    *NewEnvironmentForm
	envProjectName string
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

func (m Window) Init() tea.Cmd { return tea.SetWindowTitle("Project Manager") }

func (m *Window) SelectedProject() *domain.Project { return m.selectedProject }

func (m *Window) SelectedEnvironment() string {
	if len(m.projects) == 0 {
		return ""
	}
	idx := mod(m.cursor, m.total)
	if idx < 0 || idx >= len(m.projects) {
		return ""
	}
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

