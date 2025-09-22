package components

import (
	bta "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextArea struct {
	txt      bta.Model
	style    lipgloss.Style
	valStyle lipgloss.Style
}

type TextAreaSubmitMsg struct{ Value string }

func NewTextArea(placeholder, value string, s, vs lipgloss.Style) TextArea {
	ta := bta.New()
	ta.Placeholder = placeholder
	ta.SetValue(value)

	return TextArea{txt: ta, style: s, valStyle: vs}
}

func (t TextArea) Init() tea.Cmd { return nil }

func (t TextArea) View() string { return t.txt.View() }

func (t TextArea) Value() string { return t.txt.Value() }

func (t *TextArea) SetValue(v string) { t.txt.SetValue(v) }

func (t *TextArea) SetWidth(w int) { t.txt.SetWidth(w) }

func (t *TextArea) SetHeight(h int) { t.txt.SetHeight(h) }

func (t *TextArea) Focus() { t.txt.Focus() }

func (t *TextArea) Blur() { t.txt.Blur() }

func (t TextArea) Update(msg tea.Msg) (TextArea, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		if k.Type == tea.KeyTab || k.Type == tea.KeyShiftTab {
			return t, nil
		}

		if k.Type == tea.KeyEnter && k.Alt {
			return t, func() tea.Msg { return TextAreaSubmitMsg{Value: t.txt.Value()} }
		}
	}

	ta, cmd := t.txt.Update(msg)
	t.txt = ta

	return t, cmd
}
