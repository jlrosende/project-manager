package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type projectsLoadedMsg struct{ Items []components.ListItem }

type ProjectsViewStyles struct {
	Title    lipgloss.Style
	Item     lipgloss.Style
	Selected lipgloss.Style
	Marker   lipgloss.Color
}

type ProjectsViewStyleOption func(*ProjectsViewStyles)

func NewProjectsViewStyles(opts ...ProjectsViewStyleOption) ProjectsViewStyles {
	s := ProjectsViewStyles{}
	for _, o := range opts {
		o(&s)
	}

	return s
}

func WithProjectsTitle(s lipgloss.Style) ProjectsViewStyleOption {
	return func(ps *ProjectsViewStyles) { ps.Title = s }
}

func WithProjectsItem(s lipgloss.Style) ProjectsViewStyleOption {
	return func(ps *ProjectsViewStyles) { ps.Item = s }
}

func WithProjectsSelected(s lipgloss.Style) ProjectsViewStyleOption {
	return func(ps *ProjectsViewStyles) { ps.Selected = s }
}

func WithProjectsMarker(c lipgloss.Color) ProjectsViewStyleOption {
	return func(ps *ProjectsViewStyles) { ps.Marker = c }
}

type ProjectsView struct {
	List        components.List
	TitleStyle  lipgloss.Style
	MarkerColor lipgloss.Color
	fetch       func() ([]components.ListItem, error)
}

func NewProjectsView(
	fetch func() ([]components.ListItem, error),
	styles ProjectsViewStyles,
) ProjectsView {
	l := components.NewList(nil, styles.Item, styles.Selected).
		WithMarkerColor(func(_ components.ListItem) lipgloss.Color { return styles.Marker })

	return ProjectsView{List: l, TitleStyle: styles.Title, MarkerColor: styles.Marker, fetch: fetch}
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
