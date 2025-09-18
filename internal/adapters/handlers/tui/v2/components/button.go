package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ButtonVariant int

const (
	Primary ButtonVariant = iota
	Secondary
	Danger
)

type Button struct {
	Label         string
	Variant       ButtonVariant
	Disabled      bool
	style         lipgloss.Style
	styleDisabled lipgloss.Style
}

type ButtonPressedMsg struct{ Label string }

func NewButton(label string, v ButtonVariant, s, sd lipgloss.Style) Button {
	return Button{Label: label, Variant: v, style: s, styleDisabled: sd}
}

func (b Button) Init() tea.Cmd { return nil }

func (b Button) View() string {
	if b.Disabled {
		return b.styleDisabled.Render(b.Label)
	}

	return b.style.Render(b.Label)
}

func (b Button) Update(msg tea.Msg) (Button, tea.Cmd) {
	if b.Disabled {
		return b, nil
	}

	if m, ok := msg.(tea.KeyMsg); ok && m.Type == tea.KeyEnter {
		return b, func() tea.Msg { return ButtonPressedMsg{Label: b.Label} }
	}

	return b, nil
}
