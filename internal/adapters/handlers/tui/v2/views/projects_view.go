package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type projectsLoadedMsg struct{ Items []components.ListItem }

type ProjectsView struct {
	List        components.List
	TitleStyle  lipgloss.Style
	MarkerColor lipgloss.Color
	fetch       func() ([]components.ListItem, error)
}

func NewProjectsView(
	fetch func() ([]components.ListItem, error),
	title, item, sel lipgloss.Style,
	marker lipgloss.Color,
) ProjectsView {
	l := components.NewList(nil, item, sel).
		WithMarkerColor(func(_ components.ListItem) lipgloss.Color { return marker })

	return ProjectsView{List: l, TitleStyle: title, MarkerColor: marker, fetch: fetch}
}

func (v ProjectsView) Init() tea.Cmd {
	return func() tea.Msg {
		if v.fetch == nil {
			return projectsLoadedMsg{Items: nil}
		}

		it, _ := v.fetch()

		return projectsLoadedMsg{Items: it}
	}
}

func (v ProjectsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case projectsLoadedMsg:
		v.List.Items = m.Items
		return v, nil
	case tea.KeyMsg:
		if m.Type == tea.KeyEnter {
			if len(v.List.Items) == 0 {
				return v, nil
			}

			idx := v.List.Cursor
			if idx >= 0 && idx < len(v.List.Items) {
				it := v.List.Items[idx]
				if it.ID == "__new__" {
					return v, func() tea.Msg { return state.NewProjectMsg{} }
				}

				return v, func() tea.Msg { return state.SelectProjectMsg{ID: it.ID} }
			}
		}

		if m.Type == tea.KeyRunes {
			r := string(m.Runes)
			if r == "e" {
				if len(v.List.Items) == 0 {
					return v, nil
				}

				idx := v.List.Cursor
				if idx >= 0 && idx < len(v.List.Items) {
					it := v.List.Items[idx]
					if it.ID != "__new__" {
						return v, func() tea.Msg { return state.EditProjectMsg{ID: it.ID} }
					}
				}
			}
		}
	}

	l, cmd := v.List.Update(msg)
	v.List = l

	return v, cmd
}

func (v ProjectsView) View() string {
	b := strings.Builder{}
	b.WriteString(v.TitleStyle.Render(" Projects"))
	b.WriteString("\n")
	b.WriteString(v.List.View())

	return b.String()
}
