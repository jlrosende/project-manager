package v1

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
	titleStyle := m.styles.Title
	defaultStyle := m.styles.Default

	left.WriteString(titleStyle.Render(title))
	left.WriteString("\n")

	mod := mod(m.cursor, m.total)
	for i, project := range m.projects {
		selected := mod == i

		marker := "  "
		text := defaultStyle.Render(project.Name)

		if selected {
			marker = lipgloss.NewStyle().Foreground(c("selectedBg")).Render("▌ ")
			text = lipgloss.NewStyle().Foreground(c("selectedFg")).Bold(true).Render(project.Name)
		}

		left.WriteString(marker)
		left.WriteString(text)
		left.WriteString("\n")
	}

	if mod == len(m.projects) {
		left.WriteString(lipgloss.NewStyle().Foreground(c("selectedBg")).Render("▌ "))
		left.WriteString(lipgloss.NewStyle().Foreground(c("selectedFg")).Bold(true).Render("+ New project"))
		left.WriteString("\n")
	} else {
		left.WriteString("  ")
		left.WriteString(lipgloss.NewStyle().Foreground(c("section")).Bold(true).Render("+ New project"))
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
				marker = lipgloss.NewStyle().Foreground(c("selectedBg")).Render("▌ ")
			}

			right.WriteString(marker)
			right.WriteString(lipgloss.NewStyle().Foreground(c("section")).Bold(true).Render("+ New environment"))
			right.WriteString("\n")
		} else {
			for i, env := range p.Environments {
				marker := "  "

				selected := m.focus == 1 && m.cursorEnv == i

				if selected {
					marker = lipgloss.NewStyle().Foreground(lipgloss.Color(env.Color)).Render("▌ ")
				}

				var content string

				if selected {
					bullet := lipgloss.NewStyle().Foreground(lipgloss.Color(env.Color)).Render("● ")
					name := lipgloss.NewStyle().Foreground(c("selectedFg")).Bold(true).Render(env.Name)
					content = bullet + name
				} else {
					content = defaultStyle.Foreground(lipgloss.Color(env.Color)).Render("● " + env.Name)
				}

				if p.DefaultEnv != "" && env.Name == p.DefaultEnv {
					content = content + " " + lipgloss.NewStyle().Foreground(c("subtext")).Render("(default)")
				}

				right.WriteString(marker)
				right.WriteString(content)
				right.WriteString("\n")
			}
		}

		if mod < len(m.projects) && len(m.projects) > 0 && len(m.projects[mod].Environments) > 0 {
			marker := "  "
			if m.focus == 1 && m.cursorEnv == len(m.projects[mod].Environments) {
				marker = lipgloss.NewStyle().Foreground(c("selectedBg")).Render("▌ ")
			}

			right.WriteString(marker)
			right.WriteString(lipgloss.NewStyle().Foreground(c("section")).Bold(true).Render("+ New environment"))
			right.WriteString("\n")
		}
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, left.String(), "  ", right.String())

	sepLen := m.width
	if sepLen < 1 {
		sepLen = 80
	}

	sep := lipgloss.NewStyle().Foreground(c("border")).Render(strings.Repeat("─", sepLen))
	content = content + "\n" + sep + "\n" + m.help.View(m.keys)

	v := m.vp
	v.SetContent(content)

	return v.View()
}
