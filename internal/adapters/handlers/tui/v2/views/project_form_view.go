package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

type ProjectFormView struct {
	OriginalName   string
	Name           components.Input
	Path           components.Input
	Subproject     components.Input
	Shell          components.Input
	UserName       components.Input
	UserEmail      components.Input
	UserSigningKey components.Input
	CommitGPGSign  components.Input
	TagGPGSign     components.Input
	EnvVars        components.TextArea
	SaveBtn        components.Button
	Cancel         components.Button
	Buttons        components.ButtonGroup
	Title          lipgloss.Style
	IsEdit         bool
	Focused        int
}

func NewProjectFormView(
	name, path string,
	title, input, inputVal, btnPri, btnSec, btnDis lipgloss.Style,
) ProjectFormView {
	v := ProjectFormView{
		Name:           components.NewInput("Name*", "my-awesome-app", name, input, inputVal),
		Path:           components.NewInput("Path*", "~/my-awesome-app", path, input, inputVal),
		Subproject:     components.NewInput("Subproject", "services/api", "", input, inputVal),
		Shell:          components.NewInput("Shell", "/bin/bash", "", input, inputVal),
		UserName:       components.NewInput("Git user.name", "Jane Doe", "", input, inputVal),
		UserEmail:      components.NewInput("Git user.email", "jane@example.com", "", input, inputVal),
		UserSigningKey: components.NewInput("Git user.signingkey", "0xDEADBEEF", "", input, inputVal),
		CommitGPGSign:  components.NewInput("commit.gpgsign (true/false)", "true/false", "true", input, inputVal),
		TagGPGSign:     components.NewInput("tag.gpgsign (true/false)", "true/false", "true", input, inputVal),
		EnvVars:        components.NewTextArea("# One per line (like .env)\nAPP_ENV=development\nDATABASE_URL=postgres://user:pass@localhost:5432/app\n# comments allowed", "", input, inputVal),
		SaveBtn:        components.NewButton("Save", components.Primary, btnPri, btnDis),
		Cancel:         components.NewButton("Cancel", components.Secondary, btnSec, btnDis),
		Title:          title,
	}

	v.Buttons = components.NewButtonGroup(
		[]components.ButtonSpec{{Label: "Save", Primary: true}, {Label: "Cancel"}},
		btnPri,
		btnSec,
		btnDis,
	)

	v.Name.SetWidth(40)
	v.Path.SetWidth(40)
	v.Subproject.SetWidth(40)
	v.Shell.SetWidth(40)
	v.UserName.SetWidth(40)
	v.UserEmail.SetWidth(40)
	v.UserSigningKey.SetWidth(40)
	v.CommitGPGSign.SetWidth(40)
	v.TagGPGSign.SetWidth(40)
	v.EnvVars.SetHeight(6)
	v.EnvVars.SetWidth(60)

	v.Name.Focus()
	v.Focused = 0

	if strings.TrimSpace(name) != "" && strings.TrimSpace(path) != "" {
		v.IsEdit = true
		v.OriginalName = name
	}
	return v
}

func (v ProjectFormView) Init() tea.Cmd { return nil }

func (v ProjectFormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.Focused >= 10 {
			switch msg.Type {
			case tea.KeyLeft:
				if v.Focused == 11 {
					v.Focused = 10
				}
				return v, nil
			case tea.KeyRight:
				if v.Focused == 10 {
					v.Focused = 11
				}
				return v, nil
			case tea.KeyEnter:
				if v.Focused == 10 {
					return v, func() tea.Msg {
						return state.SaveProjectMsg{
							Name:          v.Name.Value(),
							Path:          v.Path.Value(),
							Subproject:    v.Subproject.Value(),
							Shell:         v.Shell.Value(),
							GitUserName:   v.UserName.Value(),
							GitUserEmail:  v.UserEmail.Value(),
							GitSigningKey: v.UserSigningKey.Value(),
							CommitGPGSign: v.CommitGPGSign.Value(),
							TagGPGSign:    v.TagGPGSign.Value(),
							EnvVarsRaw:    v.EnvVars.Value(),
						}
					}
				}
				return v, func() tea.Msg { return state.CancelMsg{} }
			}
		}

		if msg.Type == tea.KeyCtrlS {
			return v, func() tea.Msg {
				return state.SaveProjectMsg{
					Name:          v.Name.Value(),
					Path:          v.Path.Value(),
					Subproject:    v.Subproject.Value(),
					Shell:         v.Shell.Value(),
					GitUserName:   v.UserName.Value(),
					GitUserEmail:  v.UserEmail.Value(),
					GitSigningKey: v.UserSigningKey.Value(),
					CommitGPGSign: v.CommitGPGSign.Value(),
					TagGPGSign:    v.TagGPGSign.Value(),
					EnvVarsRaw:    v.EnvVars.Value(),
				}
			}
		}

		if msg.Type == tea.KeyTab || msg.Type == tea.KeyDown || msg.Type == tea.KeyEnter || msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
			f := v.Focused
			if msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
				f--
			} else {
				if !(msg.Type == tea.KeyEnter && f >= 9) {
					f++
				}
			}
			if v.IsEdit && f == 1 {
				if msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
					f--
				} else {
					f++
				}
			}
			if f < 0 {
				f = 11
			}
			if f > 11 {
				f = 0
			}
			if (msg.Type == tea.KeyUp || msg.Type == tea.KeyDown) && v.Focused == 9 {
				f = v.Focused
			}
			v.blurAll()
			v.Focused = f
			v.focusCurrent()
		}

	case components.InputSubmitMsg:
		val := strings.TrimSpace(msg.Value)
		if val == "" {
			return v, nil
		}

		return v, func() tea.Msg {
			return state.SaveProjectMsg{
				Name:          v.Name.Value(),
				Path:          v.Path.Value(),
				Subproject:    v.Subproject.Value(),
				Shell:         v.Shell.Value(),
				GitUserName:   v.UserName.Value(),
				GitUserEmail:  v.UserEmail.Value(),
				GitSigningKey: v.UserSigningKey.Value(),
				CommitGPGSign: v.CommitGPGSign.Value(),
				TagGPGSign:    v.TagGPGSign.Value(),
				EnvVarsRaw:    v.EnvVars.Value(),
			}
		}
	case components.ButtonPressedMsg:
		if msg.Label == "Save" {
			return v, func() tea.Msg {
				return state.SaveProjectMsg{
					Name:          v.Name.Value(),
					Path:          v.Path.Value(),
					Subproject:    v.Subproject.Value(),
					Shell:         v.Shell.Value(),
					GitUserName:   v.UserName.Value(),
					GitUserEmail:  v.UserEmail.Value(),
					GitSigningKey: v.UserSigningKey.Value(),
					CommitGPGSign: v.CommitGPGSign.Value(),
					TagGPGSign:    v.TagGPGSign.Value(),
					EnvVarsRaw:    v.EnvVars.Value(),
				}
			}
		}

		if msg.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	case components.ButtonChosenMsg:
		if msg.Label == "Save" {
			return v, func() tea.Msg {
				return state.SaveProjectMsg{
					Name:          v.Name.Value(),
					Path:          v.Path.Value(),
					Subproject:    v.Subproject.Value(),
					Shell:         v.Shell.Value(),
					GitUserName:   v.UserName.Value(),
					GitUserEmail:  v.UserEmail.Value(),
					GitSigningKey: v.UserSigningKey.Value(),
					CommitGPGSign: v.CommitGPGSign.Value(),
					TagGPGSign:    v.TagGPGSign.Value(),
					EnvVarsRaw:    v.EnvVars.Value(),
				}
			}
		}

		if msg.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	}

	oldName := strings.TrimSpace(v.Name.Value())
	n, ncmd := v.Name.Update(msg)
	v.Name = n
	if !v.IsEdit && v.Name.Focused() {
		newName := strings.TrimSpace(v.Name.Value())
		if newName != oldName {
			slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(newName, " ", "-")))
			cur := strings.TrimSpace(v.Path.Value())
			if slug == "" {
				if cur == "" || strings.HasPrefix(cur, "~/") {
					v.Path.SetValue("~/")
				}
			} else {
				if cur == "" || strings.HasPrefix(cur, "~/") {
					v.Path.SetValue("~/" + slug)
				}
			}
		}
	}

	p, pcmd := v.Path.Update(msg)
	v.Path = p
	s, scmd := v.Subproject.Update(msg)
	v.Subproject = s
	sh, shcmd := v.Shell.Update(msg)
	v.Shell = sh
	u, ucmd := v.UserName.Update(msg)
	v.UserName = u
	uem, uemcmd := v.UserEmail.Update(msg)
	v.UserEmail = uem
	uk, ukcmd := v.UserSigningKey.Update(msg)
	v.UserSigningKey = uk
	cg, cgcmd := v.CommitGPGSign.Update(msg)
	v.CommitGPGSign = cg
	tg, tgcmd := v.TagGPGSign.Update(msg)
	v.TagGPGSign = tg
	e, ecmd := v.EnvVars.Update(msg)
	v.EnvVars = e

	var bgcmd tea.Cmd
	if v.Focused >= 10 {
		bg, bcmd := v.Buttons.Update(msg)
		v.Buttons = bg
		bgcmd = bcmd
	}

	return v, tea.Batch(ncmd, pcmd, scmd, shcmd, ucmd, uemcmd, ukcmd, cgcmd, tgcmd, ecmd, bgcmd)
}

func (v *ProjectFormView) blurAll() {
	v.Name.Blur()
	v.Path.Blur()
	v.Subproject.Blur()
	v.Shell.Blur()
	v.UserName.Blur()
	v.UserEmail.Blur()
	v.UserSigningKey.Blur()
	v.CommitGPGSign.Blur()
	v.TagGPGSign.Blur()
	v.EnvVars.Blur()
}

func (v *ProjectFormView) focusCurrent() {
	switch v.Focused {
	case 0:
		v.Name.Focus()
	case 1:
		if !v.IsEdit {
			v.Path.Focus()
		}
	case 2:
		v.Subproject.Focus()
	case 3:
		v.Shell.Focus()
	case 4:
		v.UserName.Focus()
	case 5:
		v.UserEmail.Focus()
	case 6:
		v.UserSigningKey.Focus()
	case 7:
		v.CommitGPGSign.Focus()
	case 8:
		v.TagGPGSign.Focus()
	case 9:
		v.EnvVars.Focus()
	case 10:
		v.Buttons.Cursor = 0
	case 11:
		v.Buttons.Cursor = 1
	}
}

func (v ProjectFormView) View() string {
	b := strings.Builder{}
	title := "New project"
	if v.IsEdit {
		title = "Edit project"
	}
	b.WriteString(v.Title.Render(title))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Project details"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")).Render("Basic information"))
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle()
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.Name.SetPrompt(p.Render("Name") + lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A")).Render("*") + p.Render(": "))
		v.Name.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Name.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Name.SetTextStyle(val)
	}
	b.WriteString(v.Name.View())
	b.WriteString("\n")
	if v.IsEdit {
		p := strings.TrimSpace(v.Path.Value())
		if p == "" { p = "~/my-awesome-app" }
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(p + " (read-only)"))
	} else {
		{
			p := lipgloss.NewStyle()
			val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
			v.Path.SetPrompt(p.Render("Path") + lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A")).Render("*") + p.Render(": "))
			v.Path.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
			v.Path.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
			v.Path.SetTextStyle(val)
		}
		b.WriteString(v.Path.View())
	}
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Project options"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")).Render("Optional settings"))
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.Subproject.SetPrompt(p.Render("Subproject: "))
		v.Subproject.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Subproject.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Subproject.SetTextStyle(val)
	}
	b.WriteString(v.Subproject.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.Shell.SetPrompt(p.Render("Shell: "))
		v.Shell.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.Shell.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.Shell.SetTextStyle(val)
	}
	b.WriteString(v.Shell.View())
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Git settings"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")).Render("Applied to this project only"))
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.UserName.SetPrompt(p.Render("Git user.name: "))
		v.UserName.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.UserName.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.UserName.SetTextStyle(val)
	}
	b.WriteString(v.UserName.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.UserEmail.SetPrompt(p.Render("Git user.email: "))
		v.UserEmail.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.UserEmail.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.UserEmail.SetTextStyle(val)
	}
	b.WriteString(v.UserEmail.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.UserSigningKey.SetPrompt(p.Render("Git user.signingkey: "))
		v.UserSigningKey.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.UserSigningKey.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.UserSigningKey.SetTextStyle(val)
	}
	b.WriteString(v.UserSigningKey.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.CommitGPGSign.SetPrompt(p.Render("commit.gpgsign (true/false): "))
		v.CommitGPGSign.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.CommitGPGSign.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.CommitGPGSign.SetTextStyle(val)
	}
	b.WriteString(v.CommitGPGSign.View())
	b.WriteString("\n")
	{
		p := lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))
		v.TagGPGSign.SetPrompt(p.Render("tag.gpgsign (true/false): "))
		v.TagGPGSign.SetPromptStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")))
		v.TagGPGSign.SetPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")))
		v.TagGPGSign.SetTextStyle(val)
	}
	b.WriteString(v.TagGPGSign.View())
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true).Render("Environment variables"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C818C")).Render("One KEY=VALUE per line; '#' comments allowed"))
	b.WriteString("\n")
	b.WriteString(v.EnvVars.View())
	b.WriteString("\n\n")
	saveLbl := "Save"
	cancelLbl := "Cancel"
	if v.Focused == 10 {
		saveLbl = lipgloss.NewStyle().Bold(true).Render("Save")
	}
	if v.Focused == 11 {
		cancelLbl = lipgloss.NewStyle().Bold(true).Render("Cancel")
	}
	b.WriteString("   " + saveLbl + "         " + cancelLbl)
	b.WriteString("\n\n")
	return b.String()
}
