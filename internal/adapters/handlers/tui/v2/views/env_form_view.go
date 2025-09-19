package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

type EnvFormView struct {
	OriginalName string
	Name         components.Input
	Color        components.Input
	Mode         components.Input
	EnvVars      components.TextArea
	SaveBtn      components.Button
	Cancel       components.Button
	Buttons      components.ButtonGroup
	Title        lipgloss.Style
	Item         lipgloss.Style
	BtnPri       lipgloss.Style
	BtnSec       lipgloss.Style
	Focused      int
	Err          string
}

func NewEnvFormView(orig, name, color, mode, envRaw string, title, item, btnPri, btnSec lipgloss.Style) *EnvFormView {
	v := &EnvFormView{
		OriginalName: orig,
		Name:         components.NewInput("Env name*", "staging", name, item, item),
		Color:        components.NewInput("Color (name/#hex/0-255)", "teal", color, item, item),
		Mode:         components.NewInput("Env vars mode* (merge/replace)", "merge", mode, item, item),
		EnvVars: components.NewTextArea(
			"# One per line (like .env)\nAPI_URL=https://api.example.com\nLOG_LEVEL=info\n# comments allowed",
			envRaw,
			item,
			item,
		),
		SaveBtn: components.NewButton("Save", components.Primary, btnPri, btnPri.Bold(true).Underline(true), item),
		Cancel:  components.NewButton("Cancel", components.Secondary, btnSec, btnSec.Bold(true).Underline(true), item),
		Title:   title,
		Item:    item,
		BtnPri:  btnPri,
		BtnSec:  btnSec,
		Focused: 0,
	}

	v.Name.SetWidth(60)
	v.Color.SetWidth(60)
	v.Mode.SetWidth(60)
	v.EnvVars.SetHeight(6)
	v.EnvVars.SetWidth(70)
	v.Buttons = components.NewButtonGroup([]components.Button{v.SaveBtn, v.Cancel})
	v.Buttons.Active = false
	v.Name.Focus()

	if strings.TrimSpace(orig) == "" {
		if strings.TrimSpace(v.Color.Value()) == "" {
			v.Color.SetValue("grey")
		}

		if strings.TrimSpace(v.Mode.Value()) == "" {
			v.Mode.SetValue("merge")
		}
	}

	return v
}

func (v *EnvFormView) Init() tea.Cmd { return nil }

func (v *EnvFormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case components.ButtonPressedMsg:
		if m.Label == "Save" {
			v.Err = ""
			name := strings.TrimSpace(v.Name.Value())

			mode := strings.TrimSpace(v.Mode.Value())
			if name == "" {
				v.Err = "name is required"
				return v, nil
			}

			if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
				v.Err = "mode must be merge or replace"
				return v, nil
			}

			return v, func() tea.Msg {
				return state.SaveEnvFormMsg{
					OriginalName: v.OriginalName,
					Name:         name,
					Color:        strings.TrimSpace(v.Color.Value()),
					Mode:         mode,
					EnvVarsRaw:   v.EnvVars.Value(),
				}
			}
		}

		if m.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	case components.ButtonChosenMsg:
		if m.Label == "Save" {
			v.Err = ""
			name := strings.TrimSpace(v.Name.Value())

			mode := strings.TrimSpace(v.Mode.Value())
			if name == "" {
				v.Err = "name is required"
				return v, nil
			}

			if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
				v.Err = "mode must be merge or replace"
				return v, nil
			}

			return v, func() tea.Msg {
				return state.SaveEnvFormMsg{
					OriginalName: v.OriginalName,
					Name:         name,
					Color:        strings.TrimSpace(v.Color.Value()),
					Mode:         mode,
					EnvVarsRaw:   v.EnvVars.Value(),
				}
			}
		}

		if m.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	case tea.KeyMsg:
		if cmd, ok := handleButtonGroupNav(4, 5, &v.Focused, &v.Buttons, m); ok {
			return v, cmd
		}

		s := m.String()
		switch s {
		case "ctrl+s":
			v.Err = ""
			name := strings.TrimSpace(v.Name.Value())

			mode := strings.TrimSpace(v.Mode.Value())
			if name == "" {
				v.Err = "name is required"
				return v, nil
			}

			if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
				v.Err = "mode must be merge or replace"
				return v, nil
			}

			return v, func() tea.Msg {
				return state.SaveEnvFormMsg{
					OriginalName: v.OriginalName,
					Name:         name,
					Color:        strings.TrimSpace(v.Color.Value()),
					Mode:         mode,
					EnvVarsRaw:   v.EnvVars.Value(),
				}
			}
		case "tab":
			v.Focused = (v.Focused + 1) % 6
			v.blurAll()
			v.Buttons.Active = v.Focused >= 4
			v.focusCurrent()

			return v, nil
		case "shift+tab":
			v.Focused = (v.Focused + 5) % 6
			v.blurAll()
			v.Buttons.Active = v.Focused >= 4
			v.focusCurrent()

			return v, nil
		}
	}

	n, nameCmd := v.Name.Update(msg)
	v.Name = n
	c, colorCmd := v.Color.Update(msg)
	v.Color = c
	mo, modeCmd := v.Mode.Update(msg)
	v.Mode = mo
	e, envCmd := v.EnvVars.Update(msg)
	v.EnvVars = e

	return v, tea.Batch(nameCmd, colorCmd, modeCmd, envCmd)
}

func (v *EnvFormView) View() string {
	b := strings.Builder{}

	title := "New environment"
	if strings.TrimSpace(v.OriginalName) != "" {
		title = "Edit environment"
	}

	b.WriteString(v.Title.Render(title))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Environment details"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")).Render("Display and behavior"))
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle()
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))

		v.Name.SetPrompt(p.Render("Env name*") + p.Render(": "))
		v.Name.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Name.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Name.SetTextStyle(val)
	}

	b.WriteString(v.Name.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle()
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.Color.SetPrompt(p.Render("Color (name/#hex/0-255): "))
		v.Color.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Color.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Color.SetTextStyle(val)
	}

	b.WriteString(v.Color.View())
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Environment variables"))
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle()
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))

		v.Mode.SetPrompt(p.Render("Env vars mode* (merge/replace): "))
		v.Mode.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Mode.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Mode.SetTextStyle(val)
	}

	b.WriteString(v.Mode.View())
	b.WriteString("\n")
	b.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C818C")).
			Render("One KEY=VALUE per line; '#' comments allowed"),
	)
	b.WriteString("\n")
	b.WriteString(v.EnvVars.View())
	b.WriteString("\n\n")
	b.WriteString(v.Buttons.View())
	b.WriteString("\n\n")

	if strings.TrimSpace(v.Err) != "" {
		b.WriteString(v.Err)
	}

	return b.String()
}

func (v *EnvFormView) blurAll() {
	v.Name.Blur()
	v.Color.Blur()
	v.Mode.Blur()
	v.EnvVars.Blur()
}

func (v *EnvFormView) focusCurrent() {
	switch v.Focused {
	case 0:
		v.Name.Focus()
	case 1:
		v.Color.Focus()
	case 2:
		v.Mode.Focus()
	case 3:
		v.EnvVars.Focus()
	case 4:
		v.Buttons.Cursor = 0
	case 5:
		v.Buttons.Cursor = 1
	}
}
