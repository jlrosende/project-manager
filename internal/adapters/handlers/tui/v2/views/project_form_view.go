package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

const (
	idxName = iota
	idxPath
	idxSubproject
	idxShell
	idxUserName
	idxUserEmail
	idxUserSigningKey
	idxCommitGPGSign
	idxTagGPGSign
	idxEnvFile
	idxEnvVars
	idxBtnSave
	idxBtnCancel
)

type ProjectFormViewStyles struct {
	Title               lipgloss.Style
	Section             lipgloss.Style
	Subtext             lipgloss.Style
	Accent              lipgloss.Style
	Input               lipgloss.Style
	InputVal            lipgloss.Style
	BtnPrimary          lipgloss.Style
	BtnSecondary        lipgloss.Style
	BtnPrimaryFocused   lipgloss.Style
	BtnSecondaryFocused lipgloss.Style
	BtnDisabled         lipgloss.Style
}

type ProjectFormViewStyleOption func(*ProjectFormViewStyles)

func NewProjectFormViewStyles(opts ...ProjectFormViewStyleOption) ProjectFormViewStyles {
	s := ProjectFormViewStyles{}
	for _, o := range opts {
		o(&s)
	}

	return s
}

func WithProjectFormTitle(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.Title = s }
}

func WithProjectFormInput(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.Input = s }
}

func WithProjectFormInputVal(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.InputVal = s }
}

func WithProjectFormBtnPrimary(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.BtnPrimary = s }
}

func WithProjectFormBtnSecondary(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.BtnSecondary = s }
}

func WithProjectFormBtnPrimaryFocused(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.BtnPrimaryFocused = s }
}

func WithProjectFormBtnSecondaryFocused(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.BtnSecondaryFocused = s }
}

func WithProjectFormBtnDisabled(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.BtnDisabled = s }
}

func WithProjectFormSection(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.Section = s }
}

func WithProjectFormSubtext(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.Subtext = s }
}

func WithProjectFormAccent(s lipgloss.Style) ProjectFormViewStyleOption {
	return func(ps *ProjectFormViewStyles) { ps.Accent = s }
}

type ProjectFormInit struct {
	OriginalName   string
	Name           string
	Path           string
	Subproject     string
	Shell          string
	GitUserName    string
	GitUserEmail   string
	GitSigningKey  string
	CommitGPGSign  string
	TagGPGSign     string
	EnvVarsFile    string
	EnvVarsRaw     string
}

func NewProjectFormViewWith(init ProjectFormInit, styles ProjectFormViewStyles) ProjectFormView {
	v := NewProjectFormView(init.Name, init.Path, styles)

	if strings.TrimSpace(init.Subproject) != "" {
		v.Subproject.SetValue(init.Subproject)
	}
	if strings.TrimSpace(init.Shell) != "" {
		v.Shell.SetValue(init.Shell)
	}
	if strings.TrimSpace(init.GitUserName) != "" {
		v.UserName.SetValue(init.GitUserName)
	}
	if strings.TrimSpace(init.GitUserEmail) != "" {
		v.UserEmail.SetValue(init.GitUserEmail)
	}
	if strings.TrimSpace(init.GitSigningKey) != "" {
		v.UserSigningKey.SetValue(init.GitSigningKey)
	}
	if strings.TrimSpace(init.CommitGPGSign) != "" {
		v.CommitGPGSign.SetValue(init.CommitGPGSign)
	}
	if strings.TrimSpace(init.TagGPGSign) != "" {
		v.TagGPGSign.SetValue(init.TagGPGSign)
	}
	if strings.TrimSpace(init.EnvVarsFile) != "" {
		v.EnvFile.SetValue(init.EnvVarsFile)
	}
	if strings.TrimSpace(init.EnvVarsRaw) != "" {
		v.EnvVars.SetValue(init.EnvVarsRaw)
	}
	if strings.TrimSpace(init.OriginalName) != "" {
		v.IsEdit = true
		v.OriginalName = init.OriginalName
	}

	return v
}

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
	EnvFile        components.Input
	EnvVars        components.TextArea
	SaveBtn        components.Button
	Cancel         components.Button
	Buttons        components.ButtonGroup
	Title          lipgloss.Style
	Section        lipgloss.Style
	Subtext        lipgloss.Style
	Accent         lipgloss.Style
	InputVal       lipgloss.Style
	BtnPri         lipgloss.Style
	BtnSec         lipgloss.Style
	IsEdit         bool
	Focused        int
}

func NewProjectFormView(
	name, path string,
	styles ProjectFormViewStyles,
) ProjectFormView {
	v := ProjectFormView{
		Name:           components.NewInput("Name*", "my-awesome-app", name, styles.Input, styles.InputVal),
		Path:           components.NewInput("Path*", "~/my-awesome-app", path, styles.Input, styles.InputVal),
		Subproject:     components.NewInput("Subproject", "services/api", "", styles.Input, styles.InputVal),
		Shell:          components.NewInput("Shell", "/bin/bash", "", styles.Input, styles.InputVal),
		UserName:       components.NewInput("Git user.name", "Jane Doe", "", styles.Input, styles.InputVal),
		UserEmail:      components.NewInput("Git user.email", "jane@example.com", "", styles.Input, styles.InputVal),
		UserSigningKey: components.NewInput("Git user.signingkey", "0xDEADBEEF", "", styles.Input, styles.InputVal),
		CommitGPGSign: components.NewInput(
			"commit.gpgsign (true/false)",
			"true/false",
			"true",
			styles.Input,
			styles.InputVal,
		),
		TagGPGSign: components.NewInput(
			"tag.gpgsign (true/false)",
			"true/false",
			"true",
			styles.Input,
			styles.InputVal,
		),
		EnvFile: components.NewInput("Env vars file", ".env", ".env", styles.Input, styles.InputVal),
		EnvVars: components.NewTextArea(
			"# One per line (like .env)\nAPP_ENV=development\nDATABASE_URL=postgres://user:pass@localhost:5432/app\n# comments allowed",
			"",
			styles.Input,
			styles.InputVal,
		),
		SaveBtn: components.NewButton(
			"Save",
			components.Primary,
			styles.BtnPrimary,
			styles.BtnPrimaryFocused,
			styles.BtnDisabled,
		),
		Cancel: components.NewButton(
			"Cancel",
			components.Secondary,
			styles.BtnSecondary,
			styles.BtnSecondaryFocused,
			styles.BtnDisabled,
		),
		Title:    styles.Title,
		Section:  styles.Section,
		Subtext:  styles.Subtext,
		Accent:   styles.Accent,
		InputVal: styles.InputVal,
		BtnPri:   styles.BtnPrimary,
		BtnSec:   styles.BtnSecondary,
	}

	v.Buttons = components.NewButtonGroup([]components.Button{v.SaveBtn, v.Cancel})
	v.Buttons.Active = false

	v.Name.SetWidth(40)
	v.Path.SetWidth(40)
	v.Subproject.SetWidth(40)
	v.Shell.SetWidth(40)
	v.UserName.SetWidth(40)
	v.UserEmail.SetWidth(40)
	v.UserSigningKey.SetWidth(40)
	v.CommitGPGSign.SetWidth(40)
	v.TagGPGSign.SetWidth(40)
	v.EnvFile.SetWidth(40)
	v.EnvVars.SetHeight(6)
	v.EnvVars.SetWidth(60)

	v.Name.Focus()
	v.Focused = idxName

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
		if cmd, ok := handleButtonGroupNav(idxBtnSave, idxBtnCancel, &v.Focused, &v.Buttons, msg); ok {
			v.blurAll()
			v.Buttons.Active = v.Focused >= idxBtnSave
			v.focusCurrent()

			return v, cmd
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
					EnvVarsFile:   v.EnvFile.Value(),
					EnvVarsRaw:    v.EnvVars.Value(),
				}
			}
		}

		if msg.Type == tea.KeyTab || msg.Type == tea.KeyDown || msg.Type == tea.KeyEnter || msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
			allowed := v.Focused != idxEnvVars || (msg.Type == tea.KeyTab || msg.Type == tea.KeyShiftTab)
			if allowed {
				f := v.Focused

				if msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
					f--
				} else {
					if msg.Type != tea.KeyEnter || f < idxEnvVars {
						f++
					}
				}

				if v.IsEdit && f == idxPath {
					if msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyUp {
						f--
					} else {
						f++
					}
				}

				if f < 0 {
					f = idxBtnCancel
				}

				if f > idxBtnCancel {
					f = idxName
				}

				if msg.Type == tea.KeyUp && v.Focused == idxEnvVars {
					f = v.Focused
				}

				v.blurAll()
				v.Focused = f
				v.Buttons.Active = v.Focused >= idxBtnSave
				v.focusCurrent()

				return v, nil
			}
		}

	case components.InputSubmitMsg:
		f := v.Focused

		f++

		if v.IsEdit && f == idxPath {
			f++
		}

		if f > idxBtnCancel {
			f = idxName
		}

		v.blurAll()
		v.Focused = f
		v.Buttons.Active = v.Focused >= idxBtnSave
		v.focusCurrent()

		return v, nil
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
					EnvVarsFile:   v.EnvFile.Value(),
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
					EnvVarsFile:   v.EnvFile.Value(),
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
	envf, envfcmd := v.EnvFile.Update(msg)
	v.EnvFile = envf
	e, ecmd := v.EnvVars.Update(msg)
	v.EnvVars = e

	return v, tea.Batch(ncmd, pcmd, scmd, shcmd, envfcmd, ucmd, uemcmd, ukcmd, cgcmd, tgcmd, ecmd)
}

func handleButtonGroupNav(base, last int, focused *int, g *components.ButtonGroup, k tea.KeyMsg) (tea.Cmd, bool) {
	if *focused < base {
		return nil, false
	}

	switch k.Type {
	case tea.KeyLeft, tea.KeyRight, tea.KeyEnter:
		bg, cmd := g.Update(k)
		*g = bg
		*focused = base + g.Cursor
		g.Active = true

		return cmd, true
	case tea.KeyShiftTab, tea.KeyTab:
		f := *focused
		if k.Type == tea.KeyShiftTab {
			f--
		} else {
			f++
		}

		if f < 0 {
			f = last
		}

		if f > last {
			f = idxName
		}

		*focused = f

		g.Active = *focused >= base
		if *focused >= base && *focused <= last {
			g.Cursor = *focused - base
		}

		return nil, true
	}

	return nil, false
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
	v.EnvFile.Blur()
	v.EnvVars.Blur()
}

func (v *ProjectFormView) focusCurrent() {
	switch v.Focused {
	case idxName:
		v.Name.Focus()
	case idxPath:
		if !v.IsEdit {
			v.Path.Focus()
		}
	case idxSubproject:
		v.Subproject.Focus()
	case idxShell:
		v.Shell.Focus()
	case idxUserName:
		v.UserName.Focus()
	case idxUserEmail:
		v.UserEmail.Focus()
	case idxUserSigningKey:
		v.UserSigningKey.Focus()
	case idxCommitGPGSign:
		v.CommitGPGSign.Focus()
	case idxTagGPGSign:
		v.TagGPGSign.Focus()
	case idxEnvFile:
		v.EnvFile.Focus()
	case idxEnvVars:
		v.EnvVars.Focus()
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
	b.WriteString(v.Section.Bold(true).Render("Project details"))
	b.WriteString("\n")
	b.WriteString(v.Subtext.Render("Basic information"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal
		v.Name.SetPrompt(
			p.Render("Name") + v.Accent.Render("*") + p.Render(": "),
		)
		v.Name.SetPromptStyle(v.Section)
		v.Name.SetPlaceholderStyle(v.Subtext)
		v.Name.SetTextStyle(val)
	}

	b.WriteString(v.Name.View())
	b.WriteString("\n")

	if v.IsEdit {
		p := strings.TrimSpace(v.Path.Value())
		if p == "" {
			p = "~/my-awesome-app"
		}

		b.WriteString(v.Subtext.Render(p + " (read-only)"))
	} else {
		{
			p := v.Section
			val := v.InputVal
			v.Path.SetPrompt(p.Render("Path") + v.Accent.Render("*") + p.Render(": "))
			v.Path.SetPromptStyle(v.Section)
			v.Path.SetPlaceholderStyle(v.Subtext)
			v.Path.SetTextStyle(val)
		}

		b.WriteString(v.Path.View())
	}

	b.WriteString("\n\n")
	b.WriteString(v.Section.Bold(true).Render("Project options"))
	b.WriteString("\n")
	b.WriteString(v.Subtext.Render("Optional settings"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.Subproject.SetPrompt(p.Render("Subproject: "))
		v.Subproject.SetPromptStyle(v.Section)
		v.Subproject.SetPlaceholderStyle(v.Subtext)
		v.Subproject.SetTextStyle(val)
	}

	b.WriteString(v.Subproject.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.Shell.SetPrompt(p.Render("Shell: "))
		v.Shell.SetPromptStyle(v.Section)
		v.Shell.SetPlaceholderStyle(v.Subtext)
		v.Shell.SetTextStyle(val)
	}

	b.WriteString(v.Shell.View())
	b.WriteString("\n\n")
	b.WriteString(v.Section.Bold(true).Render("Git settings"))
	b.WriteString("\n")
	b.WriteString(v.Subtext.Render("Applied to this project only"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.UserName.SetPrompt(p.Render("Git user.name: "))
		v.UserName.SetPromptStyle(v.Section)
		v.UserName.SetPlaceholderStyle(v.Subtext)
		v.UserName.SetTextStyle(val)
	}

	b.WriteString(v.UserName.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.UserEmail.SetPrompt(p.Render("Git user.email: "))
		v.UserEmail.SetPromptStyle(v.Section)
		v.UserEmail.SetPlaceholderStyle(v.Subtext)
		v.UserEmail.SetTextStyle(val)
	}

	b.WriteString(v.UserEmail.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.UserSigningKey.SetPrompt(p.Render("Git user.signingkey: "))
		v.UserSigningKey.SetPromptStyle(v.Section)
		v.UserSigningKey.SetPlaceholderStyle(v.Subtext)
		v.UserSigningKey.SetTextStyle(val)
	}

	b.WriteString(v.UserSigningKey.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.CommitGPGSign.SetPrompt(p.Render("commit.gpgsign (true/false): "))
		v.CommitGPGSign.SetPromptStyle(v.Section)
		v.CommitGPGSign.SetPlaceholderStyle(v.Subtext)
		v.CommitGPGSign.SetTextStyle(val)
	}

	b.WriteString(v.CommitGPGSign.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal

		v.TagGPGSign.SetPrompt(p.Render("tag.gpgsign (true/false): "))
		v.TagGPGSign.SetPromptStyle(v.Section)
		v.TagGPGSign.SetPlaceholderStyle(v.Subtext)
		v.TagGPGSign.SetTextStyle(val)
	}

	b.WriteString(v.TagGPGSign.View())
	b.WriteString("\n\n")
	b.WriteString(v.Section.Bold(true).Render("Environment variables"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.InputVal
		v.EnvFile.SetPrompt(p.Render("Env vars file: "))
		v.EnvFile.SetPromptStyle(v.Section)
		v.EnvFile.SetPlaceholderStyle(v.Subtext)
		v.EnvFile.SetTextStyle(val)
	}
	b.WriteString(v.EnvFile.View())
	b.WriteString("\n")
	b.WriteString(v.Subtext.Render("One KEY=VALUE per line; '#' comments allowed"))
	b.WriteString("\n")
	b.WriteString(v.EnvVars.View())
	b.WriteString("\n\n")
	b.WriteString(v.Buttons.View())
	b.WriteString("\n\n")

	return b.String()
}
