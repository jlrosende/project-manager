package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Modal struct {
	Title   string
	Body    string
	Footer  string
	Visible bool
	wrap    lipgloss.Style
	title   lipgloss.Style
	body    lipgloss.Style
	footer  lipgloss.Style
}

type ModalClosedMsg struct{}

type ModalActionMsg struct{ Action string }

func NewModal(title, body, footer string, wrap, t, b, f lipgloss.Style) Modal {
	return Modal{Title: title, Body: body, Footer: footer, Visible: false, wrap: wrap, title: t, body: b, footer: f}
}

func (m Modal) Init() tea.Cmd { return nil }

func (m Modal) View() string {
	if !m.Visible {
		return ""
	}

	return m.wrap.Render(m.title.Render(m.Title) + "\n" + m.body.Render(m.Body) + "\n" + m.footer.Render(m.Footer))
}

func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	if !m.Visible {
		return m, nil
	}

	if k, ok := msg.(tea.KeyMsg); ok {
		s := k.String()
		if s == "esc" {
			m.Visible = false
			return m, func() tea.Msg { return ModalClosedMsg{} }
		}

		if s == "enter" {
			m.Visible = false
			return m, func() tea.Msg { return ModalActionMsg{Action: "enter"} }
		}
	}

	return m, nil
}
