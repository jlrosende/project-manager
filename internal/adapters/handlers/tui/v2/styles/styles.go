package styles

import "github.com/charmbracelet/lipgloss"

type ComponentStyles struct {
	InputPrompt     lipgloss.Style
	InputValue      lipgloss.Style
	ListItem        lipgloss.Style
	ListSelected    lipgloss.Style
	ModalWrap       lipgloss.Style
	ModalTitle      lipgloss.Style
	ModalBody       lipgloss.Style
	ModalFooter     lipgloss.Style
	ButtonPrimary   lipgloss.Style
	ButtonSecondary lipgloss.Style
	ButtonDisabled  lipgloss.Style
}

func BuildComponentStyles(t Theme) ComponentStyles {
	return ComponentStyles{
		InputPrompt:     t.Default,
		InputValue:      t.Default,
		ListItem:        t.Default,
		ListSelected:    t.Selected,
		ModalWrap:       t.Default,
		ModalTitle:      t.Title,
		ModalBody:       t.Default,
		ModalFooter:     t.Help,
		ButtonPrimary:   t.ButtonPrimary,
		ButtonSecondary: t.ButtonSecondary,
		ButtonDisabled:  t.Default,
	}
}
