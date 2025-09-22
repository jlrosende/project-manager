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
	ButtonSuccess   lipgloss.Style
	ButtonInfo      lipgloss.Style
	ButtonWarning   lipgloss.Style
	ButtonDanger    lipgloss.Style
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
		ButtonSuccess:   t.ButtonSuccess,
		ButtonInfo:      t.ButtonInfo,
		ButtonWarning:   t.ButtonWarning,
		ButtonDanger:    t.ButtonDanger,
		ButtonDisabled:  t.Default,
	}
}
