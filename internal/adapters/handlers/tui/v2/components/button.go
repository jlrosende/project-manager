package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ButtonVariant int

const (
	Primary ButtonVariant = iota
	Secondary
	Success
	Info
	Warning
	Danger
)

type Button struct {
	Label         string
	Variant       ButtonVariant
	Disabled      bool
	style         lipgloss.Style
	styleFocused  lipgloss.Style
	styleDisabled lipgloss.Style
}

func (b Button) Render(focused bool) string {
	if b.Disabled {
		return b.styleDisabled.Render(b.Label)
	}

	if focused {
		return b.styleFocused.Render(b.Label)
	}

	return b.style.Render(b.Label)
}

type ButtonPressedMsg struct{ Label string }

func NewButton(label string, v ButtonVariant, s, sf, sd lipgloss.Style) Button {
	// Ensure minimum horizontal padding for better size/visibility
	s = s.Padding(0, 2)
	sf = sf.Padding(0, 2)
	sd = sd.Padding(0, 2)

	return Button{Label: label, Variant: v, style: s, styleFocused: sf, styleDisabled: sd}
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
