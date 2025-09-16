package tui

import (
	"strings"

	helpcomp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

var currentPalette = map[string]string{
	"title":       "#88C0D0",
	"section":     "#81A1C1",
	"subtext":     "#7C818C",
	"text":        "#D8DEE9",
	"placeholder": "#7C818C",
	"border":      "#4C566A",
	"error":       "#BF616A",
	"buttonDefFg": "#2E3440",
	"buttonDefBg": "#4C566A",
	"buttonSelFg": "#2E3440",
	"buttonSelBg": "#A3BE8C",
	"selectedFg":  "#ECEFF4",
	"selectedBg":  "#5E81AC",
	"help":        "#7C818C",
}

var presetPalettes = map[string]map[string]string{
	"nord": {
		"title":       "#88C0D0",
		"section":     "#81A1C1",
		"subtext":     "#7C818C",
		"text":        "#D8DEE9",
		"placeholder": "#7C818C",
		"border":      "#4C566A",
		"error":       "#BF616A",
		"buttonDefFg": "#2E3440",
		"buttonDefBg": "#4C566A",
		"buttonSelFg": "#2E3440",
		"buttonSelBg": "#A3BE8C",
		"selectedFg":  "#ECEFF4",
		"selectedBg":  "#5E81AC",
		"help":        "#7C818C",
	},
	"catppuccin": {
		"title":       "#89B4FA",
		"section":     "#B4BEFE",
		"subtext":     "#A6ADC8",
		"text":        "#CDD6F4",
		"placeholder": "#6C7086",
		"border":      "#585B70",
		"error":       "#F38BA8",
		"buttonDefFg": "#1E1E2E",
		"buttonDefBg": "#585B70",
		"buttonSelFg": "#1E1E2E",
		"buttonSelBg": "#A6E3A1",
		"selectedFg":  "#CDD6F4",
		"selectedBg":  "#89B4FA",
		"help":        "#A6ADC8",
	},
	"dracula": {
		"title":       "#BD93F9",
		"section":     "#8BE9FD",
		"subtext":     "#6272A4",
		"text":        "#F8F8F2",
		"placeholder": "#6272A4",
		"border":      "#44475A",
		"error":       "#FF5555",
		"buttonDefFg": "#282A36",
		"buttonDefBg": "#44475A",
		"buttonSelFg": "#282A36",
		"buttonSelBg": "#50FA7B",
		"selectedFg":  "#F8F8F2",
		"selectedBg":  "#6272A4",
		"help":        "#6272A4",
	},
	"ayu": {
		"title":       "#59C2FF",
		"section":     "#D4BFFF",
		"subtext":     "#5C6773",
		"text":        "#CBCCC6",
		"placeholder": "#5C6773",
		"border":      "#3D4754",
		"error":       "#D95757",
		"buttonDefFg": "#1F2430",
		"buttonDefBg": "#3D4754",
		"buttonSelFg": "#1F2430",
		"buttonSelBg": "#AAD94C",
		"selectedFg":  "#CBCCC6",
		"selectedBg":  "#59C2FF",
		"help":        "#5C6773",
	},
}

func init() {
	presetPalettes["nord-light"] = map[string]string{
		"title":       "#5E81AC",
		"section":     "#81A1C1",
		"subtext":     "#4C566A",
		"text":        "#2E3440",
		"placeholder": "#7C818C",
		"border":      "#D8DEE9",
		"error":       "#BF616A",
		"buttonDefFg": "#2E3440",
		"buttonDefBg": "#D8DEE9",
		"buttonSelFg": "#2E3440",
		"buttonSelBg": "#A3BE8C",
		"selectedFg":  "#2E3440",
		"selectedBg":  "#88C0D0",
		"help":        "#6C6F7D",
	}

	presetPalettes["catppuccin-light"] = map[string]string{
		"title":       "#1E66F5",
		"section":     "#7287FD",
		"subtext":     "#6C6F85",
		"text":        "#4C4F69",
		"placeholder": "#9CA0B0",
		"border":      "#BCC0CC",
		"error":       "#D20F39",
		"buttonDefFg": "#4C4F69",
		"buttonDefBg": "#BCC0CC",
		"buttonSelFg": "#4C4F69",
		"buttonSelBg": "#40A02B",
		"selectedFg":  "#4C4F69",
		"selectedBg":  "#8CAAEE",
		"help":        "#6C6F85",
	}

	presetPalettes["dracula-light"] = map[string]string{
		"title":       "#6272A4",
		"section":     "#8BE9FD",
		"subtext":     "#6D7086",
		"text":        "#282A36",
		"placeholder": "#A0A0A0",
		"border":      "#E5E5E5",
		"error":       "#FF5555",
		"buttonDefFg": "#282A36",
		"buttonDefBg": "#E5E5E5",
		"buttonSelFg": "#282A36",
		"buttonSelBg": "#50FA7B",
		"selectedFg":  "#282A36",
		"selectedBg":  "#8BE9FD",
		"help":        "#6D7086",
	}

	presetPalettes["ayu-light"] = map[string]string{
		"title":       "#55B4D4",
		"section":     "#D4BFFF",
		"subtext":     "#8A9199",
		"text":        "#5C6773",
		"placeholder": "#9AA5B1",
		"border":      "#E6E9EF",
		"error":       "#F07178",
		"buttonDefFg": "#5C6773",
		"buttonDefBg": "#E6E9EF",
		"buttonSelFg": "#1F2430",
		"buttonSelBg": "#AAD94C",
		"selectedFg":  "#5C6773",
		"selectedBg":  "#FFCC66",
		"help":        "#8A9199",
	}
}

func setPaletteByName(name string) {
	if p, ok := presetPalettes[strings.ToLower(strings.TrimSpace(name))]; ok {
		for k, v := range p {
			currentPalette[k] = v
		}
	}
}

type Options struct {
	Theme     string
	Overrides map[string]string
}

type keyMap struct {
	Quit  key.Binding
	Edit  key.Binding
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Help  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{k.Quit, k.Edit, k.Enter, k.Help} }
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Left, k.Right, k.Up, k.Down}, {k.Edit, k.Enter, k.Quit, k.Help}}
}

type Window struct {
	projectSvc *services.ProjectService

	projects        []*domain.Project
	selectedProject *domain.Project

	cursor         int
	cursorEnv      int
	focus          int
	total          int
	width          int
	height         int
	shouldQuit     bool
	mode           int
	form           *NewProjectForm
	prompt         *postCreatePrompt
	envForm        *NewEnvironmentForm
	envProjectName string
	styles         Styles
	vp             viewport.Model
	help           helpcomp.Model
	keys           keyMap
}

func NewWindow(projectSvc *services.ProjectService, opts Options) (*Window, error) {
	if strings.TrimSpace(opts.Theme) != "" {
		setPaletteByName(opts.Theme)
	}

	for k, v := range opts.Overrides {
		if strings.TrimSpace(v) != "" {
			currentPalette[k] = v
		}
	}

	s := BuildStyles()

	projects, err := projectSvc.List()
	if err != nil {
		return nil, err
	}

	total := len(projects) + 1
	if total == 0 {
		total = 1
	}

	vp := viewport.New(0, 0)
	h := helpcomp.New()

	w := &Window{
		projectSvc: projectSvc,
		projects:   projects,
		total:      total,
		cursor:     0,
		styles:     s,
		vp:         vp,
		help:       h,
	}

	w.keys.Quit = key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit"))
	w.keys.Edit = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
	w.keys.Left = key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "focus projects"))
	w.keys.Right = key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "focus envs"))
	w.keys.Up = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up"))
	w.keys.Down = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down"))
	w.keys.Enter = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))
	w.keys.Help = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))

	return w, nil
}

func c(k string) lipgloss.Color { return lipgloss.Color(currentPalette[k]) }

func (m Window) Init() tea.Cmd { return tea.SetWindowTitle("Project Manager") }

func (m *Window) SelectedProject() *domain.Project { return m.selectedProject }

func (m *Window) SelectedEnvironment() string {
	if len(m.projects) == 0 {
		return ""
	}

	idx := mod(m.cursor, m.total)
	if idx < 0 || idx >= len(m.projects) {
		return ""
	}

	p := m.projects[idx]
	if len(p.Environments) == 0 {
		return ""
	}

	e := m.cursorEnv % len(p.Environments)
	if e < 0 {
		e = 0
	}

	return p.Environments[e].Name
}

func mod(a, b int) int { return (a%b + b) % b }
