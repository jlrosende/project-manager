package styles

import (
	"github.com/charmbracelet/lipgloss"
)

type Style string

const (
	Default    Style = "default"
	Dracula    Style = "dracula"
	Nord       Style = "nord"
	Catppuccin Style = "catppuccin"
)

func GetStyle(style Style) lipgloss.Style {
	switch style {
	case Catppuccin:
		return CatppuccinStyle
	case Dracula:
		return DraculaStyle
	case Nord:
		return NordStyle
	default:
		return DefaultStyle
	}
}

var CatppuccinStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

var DefaultStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

var DraculaStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

var NordStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)
