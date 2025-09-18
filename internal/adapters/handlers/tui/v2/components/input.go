package components

import (
	bti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	label       string
	placeholder string
	value       string
	txt         bti.Model
	style       lipgloss.Style
	valStyle    lipgloss.Style
}

type InputChangedMsg struct{ Value string }

type InputFocusMsg struct{}

type InputBlurMsg struct{}

type InputSubmitMsg struct{ Value string }

func NewInput(label, placeholder, value string, s, vs lipgloss.Style) Input {
	ti := bti.New()
	ti.Prompt = s.Render(label + ": ")
	ti.PromptStyle = lipgloss.NewStyle()
	ti.PlaceholderStyle = vs
	ti.TextStyle = vs
	ti.Placeholder = placeholder
	ti.SetValue(value)

	return Input{label: label, placeholder: placeholder, value: value, txt: ti, style: s, valStyle: vs}
}

func (i Input) Init() tea.Cmd { return bti.Blink }

func (i Input) View() string { return i.txt.View() }

func (i *Input) SetPrompt(p string) { i.txt.Prompt = p }

func (i *Input) SetPromptStyle(s lipgloss.Style) { i.txt.PromptStyle = s }

func (i *Input) SetPlaceholderStyle(s lipgloss.Style) { i.txt.PlaceholderStyle = s }

func (i *Input) SetTextStyle(s lipgloss.Style) { i.txt.TextStyle = s }

func (i Input) Value() string { return i.txt.Value() }

func (i *Input) SetValue(s string) { i.txt.SetValue(s) }

func (i *Input) SetWidth(w int) { i.txt.Width = w }

func (i *Input) Focus() { i.txt.Focus() }

func (i *Input) Blur() { i.txt.Blur() }

func (i Input) Focused() bool { return i.txt.Focused() }

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		if m.Type == tea.KeyEnter {
			return i, func() tea.Msg { return InputSubmitMsg{Value: i.txt.Value()} }
		}
	}

	t, cmd := i.txt.Update(msg)
	i.txt = t

	return i, cmd
}
