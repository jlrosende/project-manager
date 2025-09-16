package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const markerArrow = "➜"

func (m Window) View() string {
	if m.mode == 1 && m.form != nil {
		return m.form.View()
	}

	if m.mode == 2 && m.prompt != nil {
		return m.prompt.View()
	}

	if m.mode == 3 && m.envForm != nil {
		return m.envForm.View()
	}

	left := strings.Builder{}
	title := " Projects"
	titleStyle := lipgloss.NewStyle().
		Foreground(c("title")).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true).
		BorderForeground(c("border")).
		Padding(0, 1)
	selectedStyle := lipgloss.NewStyle().
		Foreground(c("selectedFg")).
		Background(c("selectedBg")).Padding(0, 1)
	defaultStyle := lipgloss.NewStyle().
		Foreground(c("subtext")).Padding(0, 1)

	left.WriteString(titleStyle.Render(title))
	left.WriteString("\n")

	mod := mod(m.cursor, m.total)
	for i, project := range m.projects {
		if mod == i {
			left.WriteString(markerArrow + " ")
			left.WriteString(selectedStyle.Render(project.Name))
			left.WriteString("\n")
		} else {
			left.WriteString(" ")
			left.WriteString(defaultStyle.Render(project.Name))
			left.WriteString("\n")
		}
	}

	if mod == len(m.projects) {
		left.WriteString(markerArrow + " ")
		left.WriteString(selectedStyle.Render("+ New project"))
		left.WriteString("\n")
	} else {
		left.WriteString(" ")
		left.WriteString(defaultStyle.Render("+ New project"))
		left.WriteString("\n")
	}

	right := strings.Builder{}
	rightTitle := "󱄑 Environments"
	right.WriteString(titleStyle.Render(rightTitle))
	right.WriteString("\n")

	if mod < len(m.projects) && len(m.projects) > 0 {
		p := m.projects[mod]
		if len(p.Environments) == 0 {
			marker := "  "
			if m.focus == 1 && m.cursorEnv == 0 {
				marker = "➜"
			}

			right.WriteString(marker)
			right.WriteString(defaultStyle.Render("+ New environment"))
			right.WriteString("\n")
		} else {
			for i, env := range p.Environments {
				marker := "  "

				selected := m.focus == 1 && m.cursorEnv == i

				if selected {
					marker = markerArrow
				}

				rowStyle := lipgloss.NewStyle()
				if selected {
					rowStyle = rowStyle.Background(lipgloss.Color(env.Color)).Padding(0, 0, 0, 1)
				}

				content := defaultStyle.Foreground(lipgloss.Color(env.Color)).Render("● " + env.Name)
				if p.DefaultEnv != "" && env.Name == p.DefaultEnv {
					content = content + " " + lipgloss.NewStyle().Foreground(c("subtext")).Render("(default)")
				}

				right.WriteString(marker)
				right.WriteString(rowStyle.Render(content))
				right.WriteString("\n")
			}
		}
		// Add new environment entry when there are environments
		if mod < len(m.projects) && len(m.projects) > 0 && len(m.projects[mod].Environments) > 0 {
			marker := "  "
			if m.focus == 1 && m.cursorEnv == len(m.projects[mod].Environments) {
				marker = markerArrow
			}

			right.WriteString(marker)
			right.WriteString(defaultStyle.Render("+ New environment"))
			right.WriteString("\n")
		}
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, left.String(), "  ", right.String())
	help := lipgloss.NewStyle().
		Foreground(c("help")).
		Render(`Navigate: ↑/k ↓/j  Focus: ←/h →/l
Select: Enter  Edit: e  Exit: Esc/Ctrl+C/q`)

	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}
