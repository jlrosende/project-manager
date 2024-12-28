package textinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	textInput textinput.Model
	err       error
}

func NewTextInput(prompt, sugestion string) *TextInput {
	ti := textinput.New()
	ti.Placeholder = sugestion
	ti.Prompt = prompt
	ti.Focus()
	return &TextInput{
		textInput: ti,
		err:       nil,
	}
}

func (m *TextInput) Value() string {
	return m.textInput.Value()
}

func (m *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *TextInput) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
