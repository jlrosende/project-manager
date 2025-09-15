package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m *Window) newProjectFlow() {
	m.form = NewProjectFormModel()
	m.mode = 1
}

func (m *Window) newEnvironmentFlow(projectName string) {
	m.envForm = NewEnvironmentFormModel()
	m.envProjectName = projectName
	m.mode = 3
}

type postCreatePrompt struct {
	idx, choice int
	projectName string
}

func newPostCreatePrompt(projectName string) *postCreatePrompt {
	return &postCreatePrompt{projectName: projectName}
}

func (p *postCreatePrompt) Init() tea.Cmd { return nil }

func (p *postCreatePrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if p.idx > 0 {
				p.idx--
			}
		case "down", "j":
			if p.idx < 2 {
				p.idx++
			}
		case "enter":
			p.choice = p.idx
			return p, tea.Quit
		case "esc", "q", "ctrl+c":
			p.choice = 0
			return p, tea.Quit
		}
	case tea.WindowSizeMsg:
		return p, nil
	}
	return p, nil
}

func (p *postCreatePrompt) View() string {
	styleTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Align(lipgloss.Center).Border(lipgloss.NormalBorder(), false, false, true).Padding(0, 1)
	styleSel := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Padding(0, 1)
	styleDef := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1)
	b := strings.Builder{}
	b.WriteString(styleTitle.Render("Project created. What next?"))
	b.WriteString("\n")
	opts := []string{"Return to list", "Start this project", "Exit"}
	for i, o := range opts {
		if i == p.idx {
			b.WriteString("➜ ")
			b.WriteString(styleSel.Render(o))
		} else {
			b.WriteString(" ")
			b.WriteString(styleDef.Render(o))
		}
		b.WriteString("\n")
	}
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Keys: ↑/k ↓/j navigate  Enter select  Esc cancel")
	return lipgloss.JoinVertical(lipgloss.Left, b.String(), help)
}
