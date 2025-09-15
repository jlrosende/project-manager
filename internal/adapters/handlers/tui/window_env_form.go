package tui

import (
	"strings"

	bta "github.com/charmbracelet/bubbles/textarea"
	bti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

type NewEnvironmentForm struct {
	name       bti.Model
	color      bti.Model
	mode       bti.Model
	envVars    bta.Model
	focused    int
	submitted  bool
	canceled   bool
	err        string
}

func NewEnvironmentFormModel() *NewEnvironmentForm {
	name := bti.New()
	name.Prompt = "Env name: "
	name.Placeholder = "dev"
	name.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	name.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	name.Focus()
	color := bti.New()
	color.Prompt = "Color (0-255): "
	color.Placeholder = "39"
	color.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	color.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	color.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	mode := bti.New()
	mode.Prompt = "Env vars mode (merge/replace): "
	mode.Placeholder = domain.ENV_VARS_MODE_MERGE
	mode.SetValue(domain.ENV_VARS_MODE_MERGE)
	mode.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	mode.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	mode.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	env := bta.New()
	env.Placeholder = "# One per line (like .env)\nKEY=VALUE\nFOO=bar\n# comments allowed"
	return &NewEnvironmentForm{name: name, color: color, mode: mode, envVars: env}
}

func (f *NewEnvironmentForm) Init() tea.Cmd { return bti.Blink }

func (f *NewEnvironmentForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "esc", "ctrl+c":
			f.canceled = true
			f.submitted = false
			return f, nil
		case "tab", "shift+tab":
			if m.String() == "tab" {
				f.focused = (f.focused + 1) % 4
			} else {
				f.focused = (f.focused + 3) % 4
			}
			f.blurAll()
			f.focusCurrent()
			return f, nil
		case "enter":
			if f.focused < 3 {
				f.focused = (f.focused + 1) % 4
				f.blurAll()
				f.focusCurrent()
				return f, nil
			}
			f.err = ""
			if strings.TrimSpace(f.name.Value()) == "" {
				f.err = "name is required"
				return f, nil
			}
			if strings.TrimSpace(f.mode.Value()) != domain.ENV_VARS_MODE_MERGE && strings.TrimSpace(f.mode.Value()) != domain.ENV_VARS_MODE_REPLACE {
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
	case 2:
		f.mode, cmd = f.mode.Update(msg)
	case 3:
		f.envVars, cmd = f.envVars.Update(msg)
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
	box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Tab switch  Enter next/submit  Esc cancel")
	b := strings.Builder{}
	b.WriteString(f.name.View())
	b.WriteString("\n")
	b.WriteString(f.color.View())
	b.WriteString("\n")
	b.WriteString(f.mode.View())
	b.WriteString("\n")
	b.WriteString(f.envVars.View())
	b.WriteString("\n\n")
	if f.err != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(f.err))
		b.WriteString("\n")
	}
	b.WriteString(help)
	return box.Render(b.String())
}
