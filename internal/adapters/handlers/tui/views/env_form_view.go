package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/state"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

type EnvFormViewStyles struct {
	Title               lipgloss.Style
	Item                lipgloss.Style
	Section             lipgloss.Style
	Subtext             lipgloss.Style
	BtnPrimary          lipgloss.Style
	BtnSecondary        lipgloss.Style
	BtnPrimaryFocused   lipgloss.Style
	BtnSecondaryFocused lipgloss.Style
	BtnDisabled         lipgloss.Style
}

type EnvFormViewStyleOption func(*EnvFormViewStyles)

func NewEnvFormViewStyles(opts ...EnvFormViewStyleOption) EnvFormViewStyles {
	s := EnvFormViewStyles{}
	for _, o := range opts {
		o(&s)
	}

	return s
}

func WithEnvFormTitle(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.Title = s }
}

func WithEnvFormItem(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.Item = s }
}

func WithEnvFormSection(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.Section = s }
}

func WithEnvFormSubtext(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.Subtext = s }
}

func WithEnvFormBtnPrimary(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.BtnPrimary = s }
}

func WithEnvFormBtnSecondary(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.BtnSecondary = s }
}

func WithEnvFormBtnPrimaryFocused(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.BtnPrimaryFocused = s }
}

func WithEnvFormBtnSecondaryFocused(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.BtnSecondaryFocused = s }
}

func WithEnvFormBtnDisabled(s lipgloss.Style) EnvFormViewStyleOption {
	return func(es *EnvFormViewStyles) { es.BtnDisabled = s }
}

const (
	envIdxName = iota
	envIdxColor
	envIdxEnvFile
	envIdxMode
	envIdxEnvVars
	envIdxBtnSave
	envIdxBtnCancel
)

type EnvFormView struct {
	OriginalName string
	Name         components.Input
	Color        components.Input
	Mode         components.Input
	EnvFile      components.Input
	EnvVars      components.TextArea
	SaveBtn      components.Button
	Cancel       components.Button
	Buttons      components.ButtonGroup
	Title        lipgloss.Style
	Item         lipgloss.Style
	Section      lipgloss.Style
	Subtext      lipgloss.Style
	BtnPri       lipgloss.Style
	BtnSec       lipgloss.Style
	Focused      int
	Err          string
}

func NewEnvFormView(orig, name, color, mode, envRaw string, styles EnvFormViewStyles) *EnvFormView {
	v := &EnvFormView{
		OriginalName: orig,
		Name:         components.NewInput("Env name*", "staging", name, styles.Item, styles.Item),
		Color:        components.NewInput("Color (name/#hex/0-255)", "teal", color, styles.Item, styles.Item),
		Mode:         components.NewInput("Env vars mode* (merge/replace)", "merge", mode, styles.Item, styles.Item),
		EnvFile:      components.NewInput("Env vars file", ".<env>.env", "", styles.Item, styles.Item),
		EnvVars: components.NewTextArea(
			"# One per line (like .env)\nAPI_URL=https://api.example.com\nLOG_LEVEL=info\n# comments allowed",
			envRaw,
			styles.Item,
			styles.Item,
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
		Title:   styles.Title,
		Item:    styles.Item,
		Section: styles.Section,
		Subtext: styles.Subtext,
		BtnPri:  styles.BtnPrimary,
		BtnSec:  styles.BtnSecondary,
		Focused: 0,
	}

	v.Name.SetWidth(60)
	v.Color.SetWidth(60)
	v.Mode.SetWidth(60)
	v.EnvFile.SetWidth(60)
	v.EnvVars.SetHeight(6)
	v.EnvVars.SetWidth(70)
	v.Buttons = components.NewButtonGroup([]components.Button{v.SaveBtn, v.Cancel})
	v.Buttons.Active = false
	v.Name.Focus()
	v.Focused = envIdxName

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
			if cmd, ok := v.makeSaveCmd(); ok {
				return v, cmd
			}

			return v, nil
		}

		if m.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	case components.InputSubmitMsg:
		f := v.Focused

		f++

		if f > envIdxBtnCancel {
			f = envIdxName
		}

		v.blurAll()
		v.Focused = f
		v.Buttons.Active = v.Focused >= envIdxBtnSave
		v.focusCurrent()

		return v, nil
	case components.ButtonChosenMsg:
		if m.Label == "Save" {
			if cmd, ok := v.makeSaveCmd(); ok {
				return v, cmd
			}

			return v, nil
		}

		if m.Label == "Cancel" {
			return v, func() tea.Msg { return state.CancelMsg{} }
		}
	case tea.KeyMsg:
		if cmd, ok := handleButtonGroupNav(envIdxBtnSave, envIdxBtnCancel, &v.Focused, &v.Buttons, m); ok {
			v.blurAll()
			v.Buttons.Active = v.Focused >= envIdxBtnSave
			v.focusCurrent()

			return v, cmd
		}

		if m.Type == tea.KeyTab || m.Type == tea.KeyDown || m.Type == tea.KeyEnter || m.Type == tea.KeyShiftTab || m.Type == tea.KeyUp {
			allowed := v.Focused != envIdxEnvVars || (m.Type == tea.KeyTab || m.Type == tea.KeyShiftTab)
			if allowed {
				f := v.Focused

				if m.Type == tea.KeyShiftTab || m.Type == tea.KeyUp {
					f--
				} else {
					if m.Type != tea.KeyEnter || f < envIdxEnvVars {
						f++
					}
				}

				if f < envIdxName {
					f = envIdxBtnCancel
				}

				if f > envIdxBtnCancel {
					f = envIdxName
				}

				if m.Type == tea.KeyUp && v.Focused == envIdxEnvVars {
					f = v.Focused
				}

				v.blurAll()
				v.Focused = f
				v.Buttons.Active = v.Focused >= envIdxBtnSave
				v.focusCurrent()

				return v, nil
			}
		}

		s := m.String()
		switch s {
		case "ctrl+s":
			if cmd, ok := v.makeSaveCmd(); ok {
				return v, cmd
			}

			return v, nil
		}
	}

	oldName := strings.TrimSpace(v.Name.Value())
	n, nameCmd := v.Name.Update(msg)
	v.Name = n

	newName := strings.TrimSpace(v.Name.Value())
	if newName != oldName {
		slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(newName, " ", "-")))
		cur := strings.TrimSpace(v.EnvFile.Value())

		oldSlug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(oldName, " ", "-")))
		if slug == "" {
			if cur == "" || cur == "."+oldSlug+".env" {
				v.EnvFile.SetValue("")
			}
		} else {
			if cur == "" || cur == "."+oldSlug+".env" {
				v.EnvFile.SetValue("." + slug + ".env")
			}
		}
	}

	c, colorCmd := v.Color.Update(msg)
	v.Color = c
	mo, modeCmd := v.Mode.Update(msg)
	v.Mode = mo
	ef, efCmd := v.EnvFile.Update(msg)
	v.EnvFile = ef
	e, envCmd := v.EnvVars.Update(msg)
	v.EnvVars = e

	return v, tea.Batch(nameCmd, colorCmd, modeCmd, efCmd, envCmd)
}

func (v *EnvFormView) makeSaveCmd() (tea.Cmd, bool) {
	v.Err = ""
	name := strings.TrimSpace(v.Name.Value())

	mode := strings.TrimSpace(v.Mode.Value())
	if name == "" {
		v.Err = "name is required"
		return nil, false
	}

	if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
		v.Err = "mode must be merge or replace"
		return nil, false
	}

	return func() tea.Msg {
		return state.SaveEnvFormMsg{
			OriginalName: v.OriginalName,
			Name:         name,
			Color:        strings.TrimSpace(v.Color.Value()),
			Mode:         mode,
			EnvVarsFile:  strings.TrimSpace(v.EnvFile.Value()),
			EnvVarsRaw:   v.EnvVars.Value(),
		}
	}, true
}

func (v *EnvFormView) View() string {
	b := strings.Builder{}

	title := "New environment"
	if strings.TrimSpace(v.OriginalName) != "" {
		title = "Edit environment"
	}

	b.WriteString(v.Title.Render(title))
	b.WriteString("\n")
	b.WriteString(v.Section.Bold(true).Render("Environment details"))
	b.WriteString("\n")
	b.WriteString(v.Subtext.Render("Display and behavior"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.Item

		v.Name.SetPrompt(p.Render("Env name*") + p.Render(": "))
		v.Name.SetPromptStyle(v.Section)
		v.Name.SetPlaceholderStyle(v.Subtext)
		v.Name.SetTextStyle(val)
	}

	b.WriteString(v.Name.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.Item

		cc := normalizeEnvColorInput(strings.TrimSpace(v.Color.Value()))
		if cc != "" {
			val = val.Foreground(lipgloss.Color(cc))
		}

		v.Color.SetPrompt(p.Render("Color (name/#hex/0-255): "))
		v.Color.SetPromptStyle(v.Section)
		v.Color.SetPlaceholderStyle(v.Subtext)
		v.Color.SetTextStyle(val)
	}

	b.WriteString(v.Color.View())
	b.WriteString("\n\n")
	b.WriteString(v.Section.Bold(true).Render("Environment variables"))
	b.WriteString("\n")
	{
		p := v.Section
		val := v.Item
		v.EnvFile.SetPrompt(p.Render("Env vars file: "))
		v.EnvFile.SetPromptStyle(v.Section)
		v.EnvFile.SetPlaceholderStyle(v.Subtext)
		v.EnvFile.SetTextStyle(val)
	}

	b.WriteString(v.EnvFile.View())
	b.WriteString("\n")
	{
		p := v.Section
		val := v.Item

		v.Mode.SetPrompt(p.Render("Env vars mode* (merge/replace): "))
		v.Mode.SetPromptStyle(v.Section)
		v.Mode.SetPlaceholderStyle(v.Subtext)
		v.Mode.SetTextStyle(val)
	}

	b.WriteString(v.Mode.View())
	b.WriteString("\n")
	b.WriteString(
		v.Subtext.Render("One KEY=VALUE per line; '#' comments allowed"),
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
	v.EnvFile.Blur()
	v.EnvVars.Blur()
}

func (v *EnvFormView) focusCurrent() {
	switch v.Focused {
	case envIdxName:
		v.Name.Focus()
	case envIdxColor:
		v.Color.Focus()
	case envIdxMode:
		v.Mode.Focus()
	case envIdxEnvFile:
		v.EnvFile.Focus()
	case envIdxEnvVars:
		v.EnvVars.Focus()
	case envIdxBtnSave:
		v.Buttons.Cursor = 0
	case envIdxBtnCancel:
		v.Buttons.Cursor = 1
	}
}

var envColorNameMap = map[string]string{
	"black":   "0",
	"white":   "15",
	"red":     "196",
	"green":   "46",
	"blue":    "21",
	"yellow":  "226",
	"magenta": "201",
	"purple":  "93",
	"cyan":    "51",
	"teal":    "30",
	"orange":  "208",
	"pink":    "205",
	"grey":    "240",
	"gray":    "240",
}

func normalizeEnvColorInput(s string) string {
	ss := strings.ToLower(strings.TrimSpace(s))
	if ss == "" {
		return ""
	}

	if strings.HasPrefix(ss, "#") {
		return ss
	}

	if v, ok := envColorNameMap[ss]; ok {
		return v
	}

	isNum := true

	for _, r := range ss {
		if r < '0' || r > '9' {
			isNum = false
			break
		}
	}

	if isNum {
		return ss
	}

	return ""
}
