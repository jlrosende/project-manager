package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ButtonSpec struct {
	Label    string
	Primary  bool
	Disabled bool
}

type ButtonGroup struct {
	Buttons      []ButtonSpec
	Cursor       int
	primaryStyle lipgloss.Style
	secondStyle  lipgloss.Style
	disStyle     lipgloss.Style
}

type ButtonChosenMsg struct{ Label string }

type ButtonMovedMsg struct{ Cursor int }

func NewButtonGroup(btns []ButtonSpec, primary, secondary, disabled lipgloss.Style) ButtonGroup {
	return ButtonGroup{Buttons: btns, primaryStyle: primary, secondStyle: secondary, disStyle: disabled}
}

func (g ButtonGroup) Init() tea.Cmd { return nil }

func (g ButtonGroup) View() string {
	out := ""

	for i, b := range g.Buttons {
		style := g.secondStyle

		if b.Primary {
			style = g.primaryStyle
		}

		if b.Disabled {
			style = g.disStyle
		}

		label := style.Render(" " + b.Label + " ")

		if i == g.Cursor {
			label = style.Bold(true).Render(" " + b.Label + " ")
		}

		if i > 0 {
			out += "   "
		}

		out += label
	}

	return out
}

func (g ButtonGroup) Update(msg tea.Msg) (ButtonGroup, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		s := k.String()

		switch s {
		case "left":
			if g.Cursor > 0 {
				g.Cursor--

				return g, func() tea.Msg { return ButtonMovedMsg{Cursor: g.Cursor} }
			}
		case "right":
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
