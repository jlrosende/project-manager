package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ButtonGroup struct {
	Buttons []Button
	Cursor  int
	Active  bool
}

type ButtonChosenMsg struct{ Label string }

type ButtonMovedMsg struct{ Cursor int }

func NewButtonGroup(btns []Button) ButtonGroup {
	return ButtonGroup{Buttons: btns}
}

func (g ButtonGroup) Init() tea.Cmd { return nil }

func (g ButtonGroup) View() string {
	out := "   "

	for i, b := range g.Buttons {
		if i > 0 {
			out += "         "
		}

		out += b.Render(i == g.Cursor && g.Active)
	}

	return out
}

func (g ButtonGroup) Update(msg tea.Msg) (ButtonGroup, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		s := k.String()

		switch s {
		case "left", "shift+tab":
			if g.Cursor > 0 {
				g.Cursor--
				return g, func() tea.Msg { return ButtonMovedMsg{Cursor: g.Cursor} }
			}
		case "right", "tab":
			if g.Cursor < len(g.Buttons)-1 {
				g.Cursor++
				return g, func() tea.Msg { return ButtonMovedMsg{Cursor: g.Cursor} }
			}
		case "enter":
			if len(g.Buttons) == 0 {
				return g, nil
			}

			b := g.Buttons[g.Cursor]
			if b.Disabled {
				return g, nil
			}

			return g, func() tea.Msg { return ButtonChosenMsg{Label: b.Label} }
		}
	}

	return g, nil
}
