package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UpdateFlowView struct {
	Title lipgloss.Style
	Body  lipgloss.Style
	Logs  []string
}

func NewUpdateFlowView(title, body lipgloss.Style) UpdateFlowView {
	return UpdateFlowView{Title: title, Body: body, Logs: nil}
}

func (v UpdateFlowView) Init() tea.Cmd { return nil }

func (v UpdateFlowView) Update(tea.Msg) (tea.Model, tea.Cmd) { return v, nil }

func (v UpdateFlowView) View() string {
	b := strings.Builder{}
	b.WriteString(v.Title.Render("Update Flow"))
	b.WriteString("\n")

	if len(v.Logs) == 0 {
		b.WriteString(v.Body.Render("Update in progress..."))
	} else {
		b.WriteString(v.Body.Render(strings.Join(v.Logs, "\n")))
	}

	return b.String()
}
