package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type EnvVarsViewStyles struct {
	Title    lipgloss.Style
	Item     lipgloss.Style
	Selected lipgloss.Style
	Marker   lipgloss.Color
}

type EnvVarsViewStyleOption func(*EnvVarsViewStyles)

func NewEnvVarsViewStyles(opts ...EnvVarsViewStyleOption) EnvVarsViewStyles {
	s := EnvVarsViewStyles{}
	for _, o := range opts {
		o(&s)
	}

	return s
}

func WithEnvVarsTitle(s lipgloss.Style) EnvVarsViewStyleOption {
	return func(es *EnvVarsViewStyles) { es.Title = s }
}

func WithEnvVarsItem(s lipgloss.Style) EnvVarsViewStyleOption {
	return func(es *EnvVarsViewStyles) { es.Item = s }
}

func WithEnvVarsSelected(s lipgloss.Style) EnvVarsViewStyleOption {
	return func(es *EnvVarsViewStyles) { es.Selected = s }
}

func WithEnvVarsMarker(c lipgloss.Color) EnvVarsViewStyleOption {
	return func(es *EnvVarsViewStyles) { es.Marker = c }
}

type EnvVarsView struct {
	List   components.List
	Title  lipgloss.Style
	Item   lipgloss.Style
	Sel    lipgloss.Style
	Marker lipgloss.Color
}

func NewEnvVarsView(items []components.ListItem, styles EnvVarsViewStyles) EnvVarsView {
	l := components.NewList(items, styles.Item, styles.Selected).
		WithMarkerColor(func(it components.ListItem) lipgloss.Color {
			if it.ID == "__add__" {
				return styles.Marker
			}

			return lipgloss.Color(it.Color)
		})
	l = l.WithBullet(func(it components.ListItem) lipgloss.Color { return lipgloss.Color(it.Color) })

	return EnvVarsView{List: l, Title: styles.Title, Item: styles.Item, Sel: styles.Selected, Marker: styles.Marker}
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
