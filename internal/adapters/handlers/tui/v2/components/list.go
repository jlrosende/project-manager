package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct{ ID, Label, Color string }

type List struct {
	Items       []ListItem
	Cursor      int
	style       lipgloss.Style
	styleSel    lipgloss.Style
	markerColor func(ListItem) lipgloss.Color
	bulletColor func(ListItem) lipgloss.Color
}

type ListMovedMsg struct{ Cursor int }

type ListSelectedMsg struct{ Item ListItem }

func NewList(items []ListItem, s, ss lipgloss.Style) List {
	return List{Items: items, style: s, styleSel: ss}
}

func (l List) WithMarkerColor(fn func(ListItem) lipgloss.Color) List {
	l.markerColor = fn
	return l
}

func (l List) WithBullet(fn func(ListItem) lipgloss.Color) List {
	l.bulletColor = fn
	return l
}

func (l List) Init() tea.Cmd { return nil }

func (l List) View() string {
	var out string

	for i, it := range l.Items {
		prefix := "  "
		if i == l.Cursor && l.markerColor != nil {
			prefix = lipgloss.NewStyle().Foreground(l.markerColor(it)).Render("▌ ")
		}

		bullet := ""

		if l.bulletColor != nil {
			c := l.bulletColor(it)
			if c != lipgloss.Color("") {
				bullet = lipgloss.NewStyle().Foreground(c).Render("● ")
			}
		}

		if i == l.Cursor {
			rendered := l.styleSel.Render(bullet + it.Label)
			rendered = strings.TrimPrefix(rendered, " ")

			out += prefix + rendered + "\n"
		} else {
			out += prefix + l.style.Render(bullet+it.Label) + "\n"
		}
	}

	return out
}

func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if len(l.Items) == 0 {
				return l, nil
			}

			if l.Cursor > 0 {
				l.Cursor--
			} else {
				l.Cursor = len(l.Items) - 1
			}

			return l, func() tea.Msg { return ListMovedMsg{Cursor: l.Cursor} }
		case "down", "j":
			if len(l.Items) == 0 {
				return l, nil
			}

			if l.Cursor < len(l.Items)-1 {
				l.Cursor++
			} else {
				l.Cursor = 0
			}

			return l, func() tea.Msg { return ListMovedMsg{Cursor: l.Cursor} }
		case "enter":
			if len(l.Items) == 0 {
				return l, nil
			}

			return l, func() tea.Msg { return ListSelectedMsg{Item: l.Items[l.Cursor]} }
		}
	}

	return l, nil
}
