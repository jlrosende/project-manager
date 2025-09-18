package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type EnvVarsView struct {
	List   components.List
	Title  lipgloss.Style
	Item   lipgloss.Style
	Sel    lipgloss.Style
	Marker lipgloss.Color
}

func NewEnvVarsView(items []components.ListItem, title, item, sel lipgloss.Style, marker lipgloss.Color) EnvVarsView {
	l := components.NewList(items, item, sel).
		WithMarkerColor(func(it components.ListItem) lipgloss.Color {
			if it.ID == "__add__" {
				return marker
			}

			return lipgloss.Color(it.Color)
		})
	l = l.WithBullet(func(it components.ListItem) lipgloss.Color { return lipgloss.Color(it.Color) })

	return EnvVarsView{List: l, Title: title, Item: item, Sel: sel, Marker: marker}
}

func (v EnvVarsView) Init() tea.Cmd { return nil }

func (v EnvVarsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		s := m.String()
		if s == "enter" {
			if len(v.List.Items) == 0 {
				return v, nil
			}
			it := v.List.Items[v.List.Cursor]
			if it.ID == "__add__" {
				return v, func() tea.Msg { return state.AddEnvVarMsg{} }
			}
			return v, func() tea.Msg { return state.ExitMsg{} }
		}
		if s == "e" {
			if len(v.List.Items) == 0 {
				return v, nil
			}
			it := v.List.Items[v.List.Cursor]
			if it.ID != "__add__" {
				return v, func() tea.Msg { return state.EditEnvVarMsg{Name: it.ID} }
			}
			return v, nil
		}
		if s == "d" {
			if len(v.List.Items) == 0 {
				return v, nil
			}
			it := v.List.Items[v.List.Cursor]
			if it.ID != "__add__" {
				return v, func() tea.Msg { return state.RemoveEnvVarMsg{Name: it.ID} }
			}
		}
	}

	l, cmd := v.List.Update(msg)
	v.List = l

	return v, cmd
}

func (v EnvVarsView) View() string { return v.List.View() }
