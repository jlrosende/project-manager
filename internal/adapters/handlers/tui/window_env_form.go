package tui

import (
	"strings"

	bta "github.com/charmbracelet/bubbles/textarea"
	bti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

const (
	errNameRequired  = "name is required"
	errModeMustBeVal = "mode must be merge or replace"
)

type NewEnvironmentForm struct {
	name         bti.Model
	color        bti.Model
	mode         bti.Model
	envVars      bta.Model
	focused      int
	submitted    bool
	canceled     bool
	err          string
	isEdit       bool
	originalName string
}

func NewEnvironmentFormModel() *NewEnvironmentForm {
	name := bti.New()
	name.Prompt = "Env name*: "
	name.Placeholder = "staging"
	name.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	name.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	name.Width = 60
	name.Focus()

	color := bti.New()
	color.Prompt = "Color (name/#hex/0-255): "
	color.Placeholder = "teal"
	color.SetValue("grey")
	color.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	color.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	color.TextStyle = lipgloss.NewStyle().Foreground(c("subtext"))
	color.Width = 60
	mode := bti.New()
	mode.Prompt = "Env vars mode* (merge/replace): "
	mode.Placeholder = domain.EnvVarsModeMerge
	mode.SetValue(domain.EnvVarsModeMerge)
	mode.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	mode.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	mode.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	mode.Width = 60
	env := bta.New()
	env.Placeholder = "# One per line (like .env)\nAPI_URL=https://api.example.com\nLOG_LEVEL=info\n# comments allowed"
	env.SetHeight(6)
	env.SetWidth(70)

	return &NewEnvironmentForm{name: name, color: color, mode: mode, envVars: env}
}

func NewEnvironmentEditFormModel(e *domain.Environment) *NewEnvironmentForm {
	f := NewEnvironmentFormModel()
	f.isEdit = true
	f.originalName = e.Name
	f.name.SetValue(e.Name)
	f.color.SetValue(e.Color)
	f.mode.SetValue(e.EnvVarsMode)

	return f
}

func (f *NewEnvironmentForm) Init() tea.Cmd { return bti.Blink }

func (f *NewEnvironmentForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		w := wm.Width - 10
		if w < 30 {
			w = 30
		}

		f.name.Width = w
		f.color.Width = w
		f.mode.Width = w
		f.envVars.SetWidth(w + 10)
	}

	if m, ok := msg.(tea.KeyMsg); ok {
		switch m.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			f.canceled = true
			f.submitted = false

			return f, nil
		case tea.KeyTab, tea.KeyShiftTab:
			if m.Type == tea.KeyTab {
				f.focused = (f.focused + 1) % 6
			} else {
				f.focused = (f.focused + 5) % 6
			}

			f.blurAll()
			f.focusCurrent()

			return f, nil
		case tea.KeyEnter:
			switch {
			case f.focused < 3:
				f.focused = (f.focused + 1) % 6
				f.blurAll()
				f.focusCurrent()

				return f, nil
			case f.focused == 4:
				f.err = ""
				if strings.TrimSpace(f.name.Value()) == "" {
					f.err = errNameRequired
					return f, nil
				}

				if strings.TrimSpace(f.mode.Value()) != domain.EnvVarsModeMerge &&
					strings.TrimSpace(f.mode.Value()) != domain.EnvVarsModeReplace {
					f.err = errModeMustBeVal
					return f, nil
				}

				f.submitted = true
				f.canceled = false

				return f, nil
			case f.focused == 5:
				f.canceled = true
				f.submitted = false

				return f, nil
			}
		case tea.KeyUp:
			if f.focused != 3 {
				f.focused--
				if f.focused < 0 {
					f.focused = 5
				}

				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyDown:
			if f.focused != 3 {
				f.focused++
				if f.focused > 5 {
					f.focused = 0
				}

				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyLeft:
			if f.focused == 5 {
				f.focused = 4
				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyRight:
			if f.focused == 4 {
				f.focused = 5
				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyCtrlS:
			f.err = ""
			if strings.TrimSpace(f.name.Value()) == "" {
				f.err = errNameRequired
				return f, nil
			}

			if strings.TrimSpace(f.mode.Value()) != domain.EnvVarsModeMerge &&
				strings.TrimSpace(f.mode.Value()) != domain.EnvVarsModeReplace {
				f.err = "mode must be merge or replace"
				return f, nil
			}

			f.submitted = true
			f.canceled = false

			return f, nil
		}
	}

	var cmd tea.Cmd

	switch f.focused {
	case 0:
		f.name, cmd = f.name.Update(msg)
	case 1:
		f.color, cmd = f.color.Update(msg)
		cv := strings.TrimSpace(f.color.Value())

		col := "240"

		if isValidColorInput(cv) {
			col = normalizeColorInput(cv)
		}

		f.color.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(col))
	case 2:
		f.mode, cmd = f.mode.Update(msg)
	case 3:
		f.envVars, cmd = f.envVars.Update(msg)
	case 4:
		// save button
	case 5:
		// cancel button
	}

	return f, cmd
}

func (f *NewEnvironmentForm) blurAll() {
	f.name.Blur()
	f.color.Blur()
	f.mode.Blur()
	f.envVars.Blur()
}

func (f *NewEnvironmentForm) focusCurrent() {
	switch f.focused {
	case 0:
		f.name.Focus()
	case 1:
		f.color.Focus()
	case 2:
		f.mode.Focus()
	case 3:
		f.envVars.Focus()
	}
}

func (f *NewEnvironmentForm) View() string {
	box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(c("border")).Padding(1, 2)
	help := lipgloss.NewStyle().
		Foreground(c("help")).
		Render(`Move: Tab/↑/↓  Buttons: ←/→
Next/Newline: Enter  Save: Ctrl+S  Cancel: Esc`)
	styleTitle := lipgloss.NewStyle().
		Foreground(c("title")).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true).
		BorderForeground(c("border")).
		Padding(0, 1)
	styleSection := lipgloss.NewStyle().Foreground(c("section")).Bold(true)
	styleSub := lipgloss.NewStyle().Foreground(c("subtext"))
	btnDef := lipgloss.NewStyle().Foreground(c("buttonDefFg")).Background(c("buttonDefBg")).Padding(0, 2)
	btnSel := lipgloss.NewStyle().Foreground(c("buttonSelFg")).Background(c("buttonSelBg")).Padding(0, 2)
	b := strings.Builder{}

	title := "New environment"
	if f.isEdit {
		title = "Edit environment"
	}

	b.WriteString(styleTitle.Render(title))
	b.WriteString("\n")
	b.WriteString(styleSection.Render("Environment details"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("Display and behavior"))
	b.WriteString("\n")
	b.WriteString(f.name.View())
	b.WriteString("\n")
	b.WriteString(f.color.View())
	b.WriteString("\n\n")
	b.WriteString(styleSection.Render("Variables"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("One KEY=VALUE per line; '#' comments allowed"))
	b.WriteString("\n")
	b.WriteString(f.mode.View())
	b.WriteString("\n")
	b.WriteString(f.envVars.View())
	b.WriteString("\n\n")

	btnSave := btnDef.Render(" Save ")
	if f.focused == 4 {
		btnSave = btnSel.Render(" Save ")
	}

	btnCancel := btnDef.Render(" Cancel ")
	if f.focused == 5 {
		btnCancel = btnSel.Render(" Cancel ")
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, btnSave, "   ", btnCancel))
	b.WriteString("\n\n")

	if f.err != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(c("error")).Render(f.err))
		b.WriteString("\n")
	}

	b.WriteString(help)

	return box.Render(b.String())
}
