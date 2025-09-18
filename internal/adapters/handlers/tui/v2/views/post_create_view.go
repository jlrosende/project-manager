package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type PostCreateView struct {
	idx         int
	projectName string
	Title       lipgloss.Style
	Item        lipgloss.Style
	Sel         lipgloss.Style
	Help        lipgloss.Style
}

func NewPostCreateView(projectName string, title, item, sel, help lipgloss.Style) PostCreateView {
	return PostCreateView{projectName: projectName, Title: title, Item: item, Sel: sel, Help: help}
}

func (v PostCreateView) Init() tea.Cmd { return nil }

func (v PostCreateView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyUp:
			if v.idx > 0 {
				v.idx--
			}
		case tea.KeyDown:
			if v.idx < 2 {
				v.idx++
			}
		case tea.KeyEnter:
			choice := v.idx
			return v, func() tea.Msg { return state.PostCreateChoiceMsg{Choice: choice, ProjectName: v.projectName} }
		case tea.KeyEsc, tea.KeyCtrlC:
			return v, func() tea.Msg { return state.PostCreateChoiceMsg{Choice: 0, ProjectName: v.projectName} }
		}

		if m.Type == tea.KeyRunes {
			switch string(m.Runes) {
			case "k":
				if v.idx > 0 {
					v.idx--
				}
			case "j":
				if v.idx < 2 {
					v.idx++
				}
			case "q":
				return v, func() tea.Msg { return state.PostCreateChoiceMsg{Choice: 0, ProjectName: v.projectName} }
			}
		}
	}

	return v, nil
}

func (v PostCreateView) View() string {
	b := strings.Builder{}
	b.WriteString(v.Title.Render("Project created. What next?"))
	b.WriteString("\n")

	opts := []string{"Return to list", "Start this project", "Exit"}
	for i, o := range opts {
		if i == v.idx {
			b.WriteString("➜ ")
			b.WriteString(v.Sel.Render(o))
		} else {
			b.WriteString(" ")
			b.WriteString(v.Item.Render(o))
		}

		b.WriteString("\n")
	}

	help := v.Help.Render(`Navigate: ↑/k ↓/j
Select: Enter  Cancel: Esc/q/Ctrl+C`)

	return lipgloss.JoinVertical(lipgloss.Left, b.String(), help)
}
