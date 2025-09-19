package v2

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	helpcomp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/router"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
	stylespkg "github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/styles"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/views"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

var currentPalette = map[string]string{
	"title":           "#88C0D0",
	"section":         "#81A1C1",
	"subtext":         "#7C818C",
	"text":            "#D8DEE9",
	"placeholder":     "#7C818C",
	"border":          "#4C566A",
	"error":           "#BF616A",
	"buttonDefFg":     "#2E3440",
	"buttonDefBg":     "#4C566A",
	"buttonSelFg":     "#2E3440",
	"buttonSelBg":     "#A3BE8C",
	"buttonSuccessFg": "#2E3440",
	"buttonSuccessBg": "#A3BE8C",
	"buttonInfoFg":    "#2E3440",
	"buttonInfoBg":    "#81A1C1",
	"buttonWarningFg": "#2E3440",
	"buttonWarningBg": "#EBCB8B",
	"buttonDangerFg":  "#2E3440",
	"buttonDangerBg":  "#BF616A",
	"selectedFg":      "#ECEFF4",
	"selectedBg":      "#5E81AC",
	"help":            "#7C818C",
}

var presetPalettes = map[string]map[string]string{
	"nord": {
		"title":           "#88C0D0",
		"section":         "#81A1C1",
		"subtext":         "#7C818C",
		"text":            "#D8DEE9",
		"placeholder":     "#7C818C",
		"border":          "#4C566A",
		"error":           "#BF616A",
		"buttonDefFg":     "#2E3440",
		"buttonDefBg":     "#4C566A",
		"buttonSelFg":     "#2E3440",
		"buttonSelBg":     "#A3BE8C",
		"buttonSuccessFg": "#2E3440",
		"buttonSuccessBg": "#A3BE8C",
		"buttonInfoFg":    "#2E3440",
		"buttonInfoBg":    "#81A1C1",
		"buttonWarningFg": "#2E3440",
		"buttonWarningBg": "#EBCB8B",
		"buttonDangerFg":  "#2E3440",
		"buttonDangerBg":  "#BF616A",
		"selectedFg":      "#ECEFF4",
		"selectedBg":      "#5E81AC",
		"help":            "#7C818C",
	},
	"nord-light": {
		"title":           "#5E81AC",
		"section":         "#81A1C1",
		"subtext":         "#4C566A",
		"text":            "#2E3440",
		"placeholder":     "#7C818C",
		"border":          "#D8DEE9",
		"error":           "#BF616A",
		"buttonDefFg":     "#2E3440",
		"buttonDefBg":     "#D8DEE9",
		"buttonSelFg":     "#2E3440",
		"buttonSelBg":     "#A3BE8C",
		"buttonSuccessFg": "#2E3440",
		"buttonSuccessBg": "#A3BE8C",
		"buttonInfoFg":    "#2E3440",
		"buttonInfoBg":    "#5E81AC",
		"buttonWarningFg": "#2E3440",
		"buttonWarningBg": "#EBCB8B",
		"buttonDangerFg":  "#2E3440",
		"buttonDangerBg":  "#BF616A",
		"selectedFg":      "#2E3440",
		"selectedBg":      "#88C0D0",
		"help":            "#6C6F7D",
	},
	"catppuccin": {
		"title":           "#89B4FA",
		"section":         "#B4BEFE",
		"subtext":         "#A6ADC8",
		"text":            "#CDD6F4",
		"placeholder":     "#6C7086",
		"border":          "#585B70",
		"error":           "#F38BA8",
		"buttonDefFg":     "#1E1E2E",
		"buttonDefBg":     "#585B70",
		"buttonSelFg":     "#1E1E2E",
		"buttonSelBg":     "#A6E3A1",
		"buttonSuccessFg": "#1E1E2E",
		"buttonSuccessBg": "#A6E3A1",
		"buttonInfoFg":    "#1E1E2E",
		"buttonInfoBg":    "#89B4FA",
		"buttonWarningFg": "#1E1E2E",
		"buttonWarningBg": "#FAB387",
		"buttonDangerFg":  "#1E1E2E",
		"buttonDangerBg":  "#F38BA8",
		"selectedFg":      "#CDD6F4",
		"selectedBg":      "#89B4FA",
		"help":            "#A6ADC8",
	},
	"catppuccin-light": {
		"title":           "#1E66F5",
		"section":         "#7287FD",
		"subtext":         "#6C6F85",
		"text":            "#4C4F69",
		"placeholder":     "#9CA0B0",
		"border":          "#BCC0CC",
		"error":           "#D20F39",
		"buttonDefFg":     "#4C4F69",
		"buttonDefBg":     "#BCC0CC",
		"buttonSelFg":     "#4C4F69",
		"buttonSelBg":     "#40A02B",
		"buttonSuccessFg": "#4C4F69",
		"buttonSuccessBg": "#40A02B",
		"buttonInfoFg":    "#4C4F69",
		"buttonInfoBg":    "#8CAAEE",
		"buttonWarningFg": "#4C4F69",
		"buttonWarningBg": "#DF8E1D",
		"buttonDangerFg":  "#4C4F69",
		"buttonDangerBg":  "#D20F39",
		"selectedFg":      "#4C4F69",
		"selectedBg":      "#8CAAEE",
		"help":            "#6C6F85",
	},
	"dracula": {
		"title":           "#BD93F9",
		"section":         "#8BE9FD",
		"subtext":         "#6272A4",
		"text":            "#F8F8F2",
		"placeholder":     "#6272A4",
		"border":          "#44475A",
		"error":           "#FF5555",
		"buttonDefFg":     "#282A36",
		"buttonDefBg":     "#44475A",
		"buttonSelFg":     "#282A36",
		"buttonSelBg":     "#50FA7B",
		"buttonSuccessFg": "#282A36",
		"buttonSuccessBg": "#50FA7B",
		"buttonInfoFg":    "#282A36",
		"buttonInfoBg":    "#8BE9FD",
		"buttonWarningFg": "#282A36",
		"buttonWarningBg": "#F1FA8C",
		"buttonDangerFg":  "#282A36",
		"buttonDangerBg":  "#FF5555",
		"selectedFg":      "#F8F8F2",
		"selectedBg":      "#6272A4",
		"help":            "#6272A4",
	},
	"dracula-light": {
		"title":           "#6272A4",
		"section":         "#8BE9FD",
		"subtext":         "#6D7086",
		"text":            "#282A36",
		"placeholder":     "#A0A0A0",
		"border":          "#E5E5E5",
		"error":           "#FF5555",
		"buttonDefFg":     "#282A36",
		"buttonDefBg":     "#E5E5E5",
		"buttonSelFg":     "#282A36",
		"buttonSelBg":     "#50FA7B",
		"buttonSuccessFg": "#282A36",
		"buttonSuccessBg": "#50FA7B",
		"buttonInfoFg":    "#282A36",
		"buttonInfoBg":    "#8BE9FD",
		"buttonWarningFg": "#282A36",
		"buttonWarningBg": "#F1FA8C",
		"buttonDangerFg":  "#282A36",
		"buttonDangerBg":  "#FF5555",
		"selectedFg":      "#282A36",
		"selectedBg":      "#8BE9FD",
		"help":            "#6D7086",
	},
	"ayu": {
		"title":           "#59C2FF",
		"section":         "#D4BFFF",
		"subtext":         "#5C6773",
		"text":            "#CBCCC6",
		"placeholder":     "#5C6773",
		"border":          "#3D4754",
		"error":           "#D95757",
		"buttonDefFg":     "#1F2430",
		"buttonDefBg":     "#3D4754",
		"buttonSelFg":     "#1F2430",
		"buttonSelBg":     "#AAD94C",
		"buttonSuccessFg": "#1F2430",
		"buttonSuccessBg": "#AAD94C",
		"buttonInfoFg":    "#1F2430",
		"buttonInfoBg":    "#59C2FF",
		"buttonWarningFg": "#1F2430",
		"buttonWarningBg": "#FFCC66",
		"buttonDangerFg":  "#1F2430",
		"buttonDangerBg":  "#D95757",
		"selectedFg":      "#CBCCC6",
		"selectedBg":      "#59C2FF",
		"help":            "#5C6773",
	},
	"ayu-light": {
		"title":           "#55B4D4",
		"section":         "#D4BFFF",
		"subtext":         "#8A9199",
		"text":            "#5C6773",
		"placeholder":     "#9AA5B1",
		"border":          "#E6E9EF",
		"error":           "#F07178",
		"buttonDefFg":     "#5C6773",
		"buttonDefBg":     "#E6E9EF",
		"buttonSelFg":     "#1F2430",
		"buttonSelBg":     "#AAD94C",
		"buttonSuccessFg": "#1F2430",
		"buttonSuccessBg": "#AAD94C",
		"buttonInfoFg":    "#1F2430",
		"buttonInfoBg":    "#8CAAEE",
		"buttonWarningFg": "#1F2430",
		"buttonWarningBg": "#FFCC66",
		"buttonDangerFg":  "#1F2430",
		"buttonDangerBg":  "#F07178",
		"selectedFg":      "#5C6773",
		"selectedBg":      "#FFCC66",
		"help":            "#8A9199",
	},
}

func init() {
}

func setPaletteByName(name string) {
	if p, ok := presetPalettes[strings.ToLower(strings.TrimSpace(name))]; ok {
		for k, v := range p {
			currentPalette[k] = v
		}
	}
}

func (m *Window) buildProjectFormStyles() views.ProjectFormViewStyles {
	return views.NewProjectFormViewStyles(
		views.WithProjectFormTitle(m.theme.Title),
		views.WithProjectFormSection(lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["section"]))),
		views.WithProjectFormSubtext(lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["subtext"]))),
		views.WithProjectFormAccent(lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["error"]))),
		views.WithProjectFormInput(m.cs.ListItem),
		views.WithProjectFormInputVal(m.cs.ListItem),
		views.WithProjectFormBtnPrimary(m.cs.ButtonPrimary),
		views.WithProjectFormBtnSecondary(m.cs.ButtonSecondary),
		views.WithProjectFormBtnPrimaryFocused(
			m.cs.ButtonPrimary.Background(lipgloss.Color(currentPalette["selectedBg"])),
		),
		views.WithProjectFormBtnSecondaryFocused(
			m.cs.ButtonSecondary.Background(lipgloss.Color(currentPalette["selectedBg"])),
		),
		views.WithProjectFormBtnDisabled(m.cs.ButtonDisabled),
	)
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

type formKeyMap struct {
	Save    key.Binding
	Cancel  key.Binding
	Next    key.Binding
	Prev    key.Binding
	Buttons key.Binding
	Help    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{k.Quit, k.Edit, k.Enter, k.Help} }
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Left, k.Right, k.Up, k.Down}, {k.Edit, k.Enter, k.Quit, k.Help}}
}

func (k formKeyMap) ShortHelp() []key.Binding { return []key.Binding{k.Next, k.Save, k.Help} }
func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Next, k.Prev, k.Buttons}, {k.Save, k.Cancel, k.Help}}
}

type Window struct {
	projectSvc *services.ProjectService
	r          *router.Router

	projects        []*domain.Project
	selectedProject *domain.Project
	selectedEnv     string

	cursor int
	total  int
	width  int

	height         int
	shouldQuit     bool
	mode           int
	focus          int
	form           tea.Model
	prompt         tea.Model
	envForm        tea.Model
	projectsView   tea.Model
	envVarsView    tea.Model
	envProjectName string
	theme          stylespkg.Theme
	cs             stylespkg.ComponentStyles
	vp             viewport.Model
	help           helpcomp.Model
	keys           keyMap
	fkeys          formKeyMap
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

	tokens := stylespkg.NewTokensFromHex(currentPalette)
	theme := stylespkg.BuildTheme(tokens)
	cs := stylespkg.BuildComponentStyles(theme)

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
		theme:      theme,
		cs:         cs,
		vp:         vp,
		help:       h,
		r:          router.New(router.RouteProjects),
	}

	w.keys.Quit = key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit"))
	w.keys.Edit = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
	w.keys.Left = key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("\u2190/h", "focus projects"))
	w.keys.Right = key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("\u2192/l", "focus envs"))
	w.keys.Up = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("\u2191/k", "up"))
	w.keys.Down = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("\u2193/j", "down"))
	w.keys.Enter = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))
	w.keys.Help = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))

	w.fkeys.Save = key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
	w.fkeys.Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	w.fkeys.Next = key.NewBinding(key.WithKeys("tab", "down", "enter"), key.WithHelp("tab/↓/enter", "next"))
	w.fkeys.Prev = key.NewBinding(key.WithKeys("shift+tab", "up"), key.WithHelp("shift+tab/↑", "prev"))
	w.fkeys.Buttons = key.NewBinding(key.WithKeys("left", "right"), key.WithHelp("←/→", "buttons"))
	w.fkeys.Help = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))

	fetch := func() ([]components.ListItem, error) {
		items := make([]components.ListItem, 0, len(w.projects)+1)
		for _, p := range w.projects {
			items = append(items, components.ListItem{ID: p.Name, Label: p.Name})
		}

		items = append(items, components.ListItem{ID: "__new__", Label: "+ New project"})

		return items, nil
	}

	projStyles := views.NewProjectsViewStyles(
		views.WithProjectsTitle(w.theme.Title),
		views.WithProjectsItem(w.cs.ListItem),
		views.WithProjectsSelected(w.cs.ListSelected),
		views.WithProjectsMarker(lipgloss.Color(currentPalette["selectedBg"])),
	)

	w.projectsView = views.NewProjectsView(
		fetch,
		projStyles,
	)

	if pv, ok := w.projectsView.(views.ProjectsView); ok {
		if items, err := fetch(); err == nil {
			pv.List.Items = items
			w.projectsView = pv
		}
	}

	return w, nil
}

func (m Window) Init() tea.Cmd {
	cmds := []tea.Cmd{tea.SetWindowTitle("Project Manager")}

	if m.projectsView != nil {
		cmds = append(cmds, m.projectsView.Init())
	}

	return tea.Batch(cmds...)
}

func (m *Window) SelectedProject() *domain.Project { return m.selectedProject }

func (m *Window) SelectedEnvironment() string {
	if m.selectedProject == nil {
		return ""
	}

	if strings.TrimSpace(m.selectedEnv) != "" {
		return m.selectedEnv
	}

	return strings.TrimSpace(m.selectedProject.DefaultEnv)
}

func parseEnvLines(s string) domain.EnvVars {
	res := domain.EnvVars{}

	for _, line := range strings.Split(s, "\n") {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		kv := strings.SplitN(l, "=", 2)
		if len(kv) == 2 {
			res[strings.TrimSpace(kv[0])] = kv[1]
		}
	}

	return res
}

var colorNameMap = map[string]string{
	"black":   "0",
	"white":   "15",
	"red":     "196",
	"green":   "46",
	"blue":    "21",
	"yellow":  "226",
	"magenta": "201",
	"purple":  "93",
	"cyan":    "51",
	"teal":    "30",
	"orange":  "208",
	"pink":    "205",
	"grey":    "240",
	"gray":    "240",
}

func normalizeColorInput(s string) string {
	ss := strings.ToLower(strings.TrimSpace(s))
	if ss == "" || ss == "grey" || ss == "gray" {
		return "240"
	}

	if strings.HasPrefix(ss, "#") {
		return ss
	}

	if v, ok := colorNameMap[ss]; ok {
		return v
	}

	return ss
}

func (m *Window) newProjectFlow() {
	projFormStyles := m.buildProjectFormStyles()

	m.form = views.NewProjectFormView(
		"",
		"",
		projFormStyles,
	)

	m.mode = 1
	if m.r != nil {
		m.r.NavigateTo(router.RouteProjectForm, nil)
	}
}

func (m *Window) newEnvironmentFlow(projectName string) {
	m.envForm = views.NewEnvFormView(
		"",
		"",
		"",
		"merge",
		"",
		m.theme.Title,
		m.cs.ListItem,
		m.cs.ButtonPrimary,
		m.cs.ButtonSecondary,
	)
	m.envProjectName = projectName

	m.mode = 3
	if m.r != nil {
		m.r.NavigateTo(router.RouteEnvVars, router.Params{"project": projectName})
	}
}

func (m *Window) currentView() tea.Model {
	switch m.r.Current() {
	case router.RouteProjectForm:
		if m.form != nil {
			return m.form
		}
	case router.RouteEnvVars:
		if m.envForm != nil {
			return m.envForm
		}

		projName := ""
		if m.r != nil && m.r.Params() != nil {
			projName = m.r.Params()["project"]
		}

		if projName == "" {
			projName = m.envProjectName
		}

		if projName != "" {
			p, _ := m.projectSvc.Load(projName)
			if p != nil {
				items := make([]components.ListItem, 0, len(p.Environments)+1)
				for _, e := range p.Environments {
					items = append(items, components.ListItem{ID: e.Name, Label: e.Name, Color: e.Color})
				}

				items = append(items, components.ListItem{ID: "__add__", Label: "+ Add environment"})

				m.envVarsView = views.NewEnvVarsView(
					items,
					views.NewEnvVarsViewStyles(
						views.WithEnvVarsTitle(m.theme.Title),
						views.WithEnvVarsItem(m.cs.ListItem),
						views.WithEnvVarsSelected(m.cs.ListSelected),
						views.WithEnvVarsMarker(lipgloss.Color(currentPalette["selectedBg"])),
					),
				)

				return m.envVarsView
			}
		}
	case router.RouteUpdateFlow:
		if m.prompt != nil {
			return m.prompt
		}
	}

	if m.projectsView != nil {
		return m.projectsView
	}

	return nil
}

func (m Window) View() string {
	if m.mode == 0 && m.r != nil && m.r.Current() == router.RouteProjects {
		var (
			leftList  string
			rightList string
		)

		if pv, ok := m.projectsView.(views.ProjectsView); ok {
			idx := pv.List.Cursor

			{
				b := strings.Builder{}

				for i, p := range m.projects {
					selected := idx == i
					marker := "  "

					text := m.cs.ListItem.Render(p.Name)
					if selected {
						marker = lipgloss.NewStyle().
							Foreground(lipgloss.Color(currentPalette["selectedBg"])).
							Render("▌ ")
						text = lipgloss.NewStyle().
							Foreground(lipgloss.Color(currentPalette["selectedFg"])).
							Bold(true).
							Render(p.Name)
					}

					b.WriteString(marker)
					b.WriteString(text)
					b.WriteString("\n")
				}

				if idx == len(m.projects) {
					b.WriteString(
						lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["selectedBg"])).Render("▌ "),
					)
					b.WriteString(
						lipgloss.NewStyle().
							Foreground(lipgloss.Color(currentPalette["selectedFg"])).
							Bold(true).
							Render("+ New project"),
					)
					b.WriteString("\n")
				} else {
					b.WriteString("  ")
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["section"])).Bold(true).Render("+ New project"))
					b.WriteString("\n")
				}

				leftList = b.String()
			}

			if idx >= 0 && idx < len(pv.List.Items) {
				it := pv.List.Items[idx]
				if it.ID != "__new__" {
					p, _ := m.projectSvc.Load(it.ID)
					if p != nil {
						cursor := 0
						if ev, ok := m.envVarsView.(views.EnvVarsView); ok {
							cursor = ev.List.Cursor
						}

						b := strings.Builder{}

						if len(p.Environments) == 0 {
							marker := "  "
							if m.focus == 1 && cursor == 0 {
								marker = lipgloss.NewStyle().
									Foreground(lipgloss.Color(currentPalette["selectedBg"])).
									Render("▌ ")
							}

							b.WriteString(marker)
							b.WriteString(
								lipgloss.NewStyle().
									Foreground(lipgloss.Color(currentPalette["section"])).
									Bold(true).
									Render("+ New environment"),
							)
							b.WriteString("\n")
						} else {
							for i, e := range p.Environments {
								selected := m.focus == 1 && cursor == i

								marker := "  "
								if selected {
									marker = lipgloss.NewStyle().Foreground(lipgloss.Color(e.Color)).Render("▌ ")
								}

								var content string

								if selected {
									bullet := lipgloss.NewStyle().Foreground(lipgloss.Color(e.Color)).Render("● ")
									name := lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["selectedFg"])).Bold(true).Render(e.Name)
									content = bullet + name
								} else {
									content = m.cs.ListItem.Foreground(lipgloss.Color(e.Color)).Render("● " + e.Name)
								}

								if p.DefaultEnv != "" && e.Name == p.DefaultEnv {
									content = content + " " + lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["subtext"])).Render("(default)")
								}

								b.WriteString(marker)
								b.WriteString(content)
								b.WriteString("\n")
							}

							marker := "  "
							if m.focus == 1 && cursor == len(p.Environments) {
								marker = lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["selectedBg"])).Render("▌ ")
							}

							b.WriteString(marker)
							b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["section"])).Bold(true).Render("+ New environment"))
							b.WriteString("\n")
						}

						rightList = b.String()
					}
				}
			}
		}

		left := m.theme.Title.Render(" Projects") + "\n" + leftList
		right := m.theme.Title.Render("󱄑 Environments") + "\n" + rightList

		ll := strings.Split(strings.TrimRight(left, "\n"), "\n")
		rl := strings.Split(strings.TrimRight(right, "\n"), "\n")

		lw := 0

		for _, s := range ll {
			w := lipgloss.Width(strings.TrimRight(s, " "))
			if w > lw {
				lw = w
			}
		}

		var rows []string

		maxRows := len(ll)
		if len(rl) > maxRows {
			maxRows = len(rl)
		}

		for i := 0; i < maxRows; i++ {
			ls := ""
			if i < len(ll) {
				ls = strings.TrimRight(ll[i], " ")
			}

			rs := ""
			if i < len(rl) {
				rs = rl[i]
			}

			pad := lw - lipgloss.Width(ls)
			if pad < 0 {
				pad = 0
			}

			rows = append(rows, ls+strings.Repeat(" ", pad)+"  "+rs)
		}

		content := strings.Join(rows, "\n")

		sepLen := m.width
		if sepLen < 1 {
			sepLen = 80
		}

		sep := lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["border"]))
		content = content + "\n\n" + sep.Render(strings.Repeat("─", sepLen)) + "\n" + m.help.View(m.keys)

		v := m.vp
		v.SetContent(content)

		return v.View()
	}

	v := m.currentView()
	if v != nil {
		content := v.View()

		sepLen := 80
		sep := lipgloss.NewStyle().Foreground(lipgloss.Color(currentPalette["border"]))

		pad := "\n\n"
		if m.mode == 1 || m.mode == 3 {
			pad = "\n\n\n"

			if m.mode == 3 {
				if f, ok := m.envForm.(*views.EnvFormView); ok {
					if strings.TrimSpace(f.Err) != "" {
						pad = "\n\n"
					}
				}
			}
		}

		content = strings.TrimRight(
			content,
			"\n",
		) + pad + sep.Render(
			strings.Repeat("─", sepLen),
		) + "\n" + m.help.View(
			m.fkeys,
		)

		return content
	}

	return ""
}

func (m *Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.shouldQuit {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	m.reduce(msg)

	if m.shouldQuit {
		return m, tea.Quit
	}

	if km, ok := msg.(tea.KeyMsg); ok {
		if km.Type == tea.KeyRunes {
			r := string(km.Runes)
			if r == "?" {
				m.help.ShowAll = !m.help.ShowAll
				hh := lipgloss.Height(m.help.View(m.keys))

				vh := m.height - hh
				if vh < 1 {
					vh = 1
				}

				m.vp.Height = vh

				return m, nil
			}
		}
	}

	if m.mode == 1 && m.form != nil {
		fm, fcmd := m.form.Update(msg)
		m.form = fm
		cmd = tea.Batch(cmd, fcmd)
	}

	if m.mode == 3 && m.envForm != nil {
		efm, ecmd := m.envForm.Update(msg)
		m.envForm = efm
		cmd = tea.Batch(cmd, ecmd)
	}

	if m.mode == 0 && m.projectsView != nil && m.r != nil && m.r.Current() == router.RouteProjects {
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.Type {
			case tea.KeyRight:
				if pv, ok := m.projectsView.(views.ProjectsView); ok {
					idx := pv.List.Cursor
					if idx >= 0 && idx < len(pv.List.Items) {
						it := pv.List.Items[idx]
						if it.ID != "__new__" {
							m.focus = 1
						}
					}
				}
			case tea.KeyLeft:
				m.focus = 0
			}

			if km.Type == tea.KeyRunes {
				r := string(km.Runes)
				if r == "l" {
					if pv, ok := m.projectsView.(views.ProjectsView); ok {
						idx := pv.List.Cursor
						if idx >= 0 && idx < len(pv.List.Items) {
							it := pv.List.Items[idx]
							if it.ID != "__new__" {
								m.focus = 1
							}
						}
					}
				}

				if r == "h" {
					m.focus = 0
				}
			}
		}

		if m.focus == 1 && m.envVarsView == nil {
			if pv, ok := m.projectsView.(views.ProjectsView); ok {
				idx := pv.List.Cursor
				if idx >= 0 && idx < len(pv.List.Items) {
					it := pv.List.Items[idx]
					if it.ID != "__new__" {
						m.envProjectName = it.ID

						p, _ := m.projectSvc.Load(it.ID)
						if p != nil {
							items := make([]components.ListItem, 0, len(p.Environments)+1)
							for _, e := range p.Environments {
								label := e.Name
								if p.DefaultEnv != "" && e.Name == p.DefaultEnv {
									label = label + "  (default)"
								}

								items = append(items, components.ListItem{ID: e.Name, Label: label, Color: e.Color})
							}

							items = append(items, components.ListItem{ID: "__add__", Label: "+ New environment"})
							m.envVarsView = views.NewEnvVarsView(
								items,
								views.NewEnvVarsViewStyles(
									views.WithEnvVarsTitle(m.theme.Title),
									views.WithEnvVarsItem(m.cs.ListItem),
									views.WithEnvVarsSelected(m.cs.ListSelected),
									views.WithEnvVarsMarker(lipgloss.Color(currentPalette["selectedBg"])),
								),
							)
						}
					}
				}
			}
		}

		var vcmd tea.Cmd

		if m.focus == 0 {
			pv, pvcmd := m.projectsView.Update(msg)
			m.projectsView = pv
			vcmd = pvcmd
		} else if m.envVarsView != nil {
			ev, evcmd := m.envVarsView.Update(msg)
			m.envVarsView = ev
			vcmd = evcmd
		}

		cmd = tea.Batch(cmd, vcmd)
	}

	if m.r != nil && m.r.Current() == router.RouteEnvVars && m.envForm == nil && m.envVarsView != nil {
		ev, evcmd := m.envVarsView.Update(msg)
		m.envVarsView = ev
		cmd = tea.Batch(cmd, evcmd)
	}

	switch msg := msg.(type) {
	case components.ListMovedMsg:
		if m.mode == 0 && m.r != nil && m.r.Current() == router.RouteProjects && m.focus == 0 {
			m.envVarsView = nil
			m.envProjectName = ""
		}
	case state.NewProjectMsg:
		m.newProjectFlow()
		return m, tea.ClearScreen
	case state.SelectProjectMsg:
		project, _ := m.projectSvc.Load(msg.ID)
		m.selectedProject = project

		return m, tea.Quit
	case state.SaveProjectMsg:
		name := strings.TrimSpace(msg.Name)
		path := strings.TrimSpace(msg.Path)

		if name == "" || path == "" {
			return m, nil
		}

		sub := strings.TrimSpace(msg.Subproject)

		envs := domain.EnvVars{}
		if strings.TrimSpace(msg.EnvVarsRaw) != "" {
			envs = parseEnvLines(msg.EnvVarsRaw)
		}

		commitOn := strings.ToLower(strings.TrimSpace(msg.CommitGPGSign))
		tagOn := strings.ToLower(strings.TrimSpace(msg.TagGPGSign))
		commitEnable := commitOn == "true" || commitOn == "1" || commitOn == "yes" || commitOn == "y"
		tagEnable := tagOn == "true" || tagOn == "1" || tagOn == "yes" || tagOn == "y"

		gitCfg := domain.New(
			domain.WithName(strings.TrimSpace(msg.GitUserName)),
			domain.WithEmail(strings.TrimSpace(msg.GitUserEmail)),
			domain.WithSigningKey(strings.TrimSpace(msg.GitSigningKey)),
			domain.WithCommitSign(commitEnable),
			domain.WithTagSign(tagEnable),
		)

		if f, ok := m.form.(*views.ProjectFormView); ok && f.IsEdit && strings.TrimSpace(f.OriginalName) != "" {
			proj, _ := m.projectSvc.Load(f.OriginalName)
			if proj != nil {
				proj.Name = name
				proj.Shell = strings.TrimSpace(msg.Shell)
				_ = m.projectSvc.UpdateProject(proj)
			}

			m.mode = 0
			if m.r != nil {
				m.r.NavigateTo(router.RouteProjects, nil)
			}

			if m.projectsView != nil {
				return m, tea.Batch(cmd, m.projectsView.Init())
			}

			return m, tea.ClearScreen
		}

		_, _ = m.projectSvc.Create(name, path, sub, envs, gitCfg)
		m.prompt = views.NewPostCreateView(name, m.theme.Title, m.cs.ListItem, m.cs.ListSelected, m.theme.Help)

		m.mode = 2
		if m.r != nil {
			m.r.NavigateTo(router.RouteUpdateFlow, nil)
		}

		return m, tea.ClearScreen
	case state.PostCreateChoiceMsg:
		if msg.Choice == 0 {
			m.mode = 0
			if m.r != nil {
				m.r.NavigateTo(router.RouteProjects, nil)
			}

			m.prompt = nil

			if m.projectsView != nil {
				return m, tea.Batch(cmd, m.projectsView.Init())
			}

			return m, tea.ClearScreen
		}

		if msg.Choice == 1 {
			proj, _ := m.projectSvc.Load(msg.ProjectName)
			m.selectedProject = proj

			return m, tea.Quit
		}

		if msg.Choice == 2 {
			m.selectedProject = nil
			return m, tea.Quit
		}

		return m, nil
	case state.CancelMsg:
		m.mode = 0
		if m.r != nil {
			m.r.NavigateTo(router.RouteProjects, nil)
		}

		if m.projectsView != nil {
			return m, tea.Batch(cmd, m.projectsView.Init())
		}

		return m, tea.ClearScreen
	case state.EditProjectMsg:
		project, _ := m.projectSvc.Load(msg.ID)
		if project == nil {
			return m, nil
		}

		projFormStyles := m.buildProjectFormStyles()
		m.form = views.NewProjectFormView(project.Name, project.Path, projFormStyles)
		m.mode = 1

		if m.r != nil {
			m.r.NavigateTo(router.RouteProjectForm, nil)
		}

		return m, tea.ClearScreen
	case state.AddEnvVarMsg:
		pname := m.envProjectName
		if pname == "" && m.r != nil && m.r.Params() != nil {
			pname = m.r.Params()["project"]
		}

		if pname != "" {
			m.newEnvironmentFlow(pname)
			return m, tea.ClearScreen
		}

		return m, nil
	case state.EditEnvVarMsg:
		pname := m.envProjectName
		if pname == "" && m.r != nil && m.r.Params() != nil {
			pname = m.r.Params()["project"]
		}

		if pname == "" {
			if pv, ok := m.projectsView.(views.ProjectsView); ok {
				idx := pv.List.Cursor
				if idx >= 0 && idx < len(pv.List.Items) {
					it := pv.List.Items[idx]
					if it.ID != "__new__" {
						pname = it.ID
					}
				}
			}
		}

		if pname != "" {
			proj, _ := m.projectSvc.Load(pname)
			if proj != nil {
				for _, e := range proj.Environments {
					if e.Name == msg.Name {
						m.envForm = views.NewEnvFormView(e.Name, e.Name, e.Color, e.EnvVarsMode, "", m.theme.Title, m.cs.ListItem, m.cs.ButtonPrimary, m.cs.ButtonSecondary)
						m.envProjectName = pname

						envPath := e.EnvVarsFile
						if !filepath.IsAbs(envPath) {
							envPath = filepath.Join(proj.Path, envPath)
						}

						if b, err := os.ReadFile(envPath); err == nil {
							if f, ok := m.envForm.(*views.EnvFormView); ok {
								f.EnvVars.SetValue(string(b))
							}
						} else if len(e.EnvVars) > 0 {
							pairs := e.EnvVars.ToSlice()
							sort.Strings(pairs)

							if f, ok := m.envForm.(*views.EnvFormView); ok {
								f.EnvVars.SetValue(strings.Join(pairs, "\n"))
							}
						}

						m.mode = 3
						if m.r != nil {
							m.r.NavigateTo(router.RouteEnvVars, router.Params{"project": pname})
						}

						return m, tea.ClearScreen
					}
				}
			}
		}

		return m, nil
	case state.RemoveEnvVarMsg:
		return m, nil
	case state.SaveEnvFormMsg:
		pname := m.envProjectName
		if pname == "" && m.r != nil && m.r.Params() != nil {
			pname = m.r.Params()["project"]
		}

		if pname == "" {
			return m, nil
		}

		env := &domain.Environment{
			Name:        strings.TrimSpace(msg.Name),
			Color:       normalizeColorInput(strings.TrimSpace(msg.Color)),
			EnvVarsMode: strings.TrimSpace(msg.Mode),
		}

		envs := domain.EnvVars{}
		if strings.TrimSpace(msg.EnvVarsRaw) != "" {
			envs = parseEnvLines(msg.EnvVarsRaw)
		}

		if strings.TrimSpace(msg.OriginalName) != "" {
			_ = m.projectSvc.UpdateEnvironment(pname, msg.OriginalName, env)
		} else {
			_ = m.projectSvc.AddEnvironment(pname, env, envs)
		}

		raw := strings.TrimSpace(msg.EnvVarsRaw)
		if raw != "" {
			proj, _ := m.projectSvc.Load(pname)
			if proj != nil {
				var envPath string

				for _, e := range proj.Environments {
					if e.Name == env.Name {
						envPath = e.EnvVarsFile
						break
					}
				}

				if envPath != "" {
					if !filepath.IsAbs(envPath) {
						envPath = filepath.Join(proj.Path, envPath)
					}

					lines := strings.Split(raw, "\n")
					ok := true

					for _, line := range lines {
						l := strings.TrimSpace(line)
						if l == "" || strings.HasPrefix(l, "#") {
							continue
						}

						kv := strings.SplitN(l, "=", 2)
						if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" {
							ok = false
							break
						}
					}

					if ok {
						_ = os.WriteFile(envPath, []byte(raw), 0o600)
					}
				}
			}
		}

		if projects, err := m.projectSvc.List(); err == nil {
			m.projects = projects
			m.total = len(projects) + 1
		}

		m.mode = 0
		if m.r != nil {
			m.r.NavigateTo(router.RouteProjects, nil)
		}

		m.envForm = nil
		m.envProjectName = ""

		if m.projectsView != nil {
			return m, tea.Batch(cmd, m.projectsView.Init())
		}

		return m, tea.ClearScreen
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.vp.Width = msg.Width

		hh := lipgloss.Height(m.help.View(m.keys))

		vh := msg.Height - hh
		if vh < 1 {
			vh = 1
		}

		m.vp.Height = vh
	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes {
			r := string(msg.Runes)
			switch r {
			case "?":
				m.help.ShowAll = !m.help.ShowAll

				hh := lipgloss.Height(m.help.View(m.keys))

				vh := m.height - hh
				if vh < 1 {
					vh = 1
				}

				m.vp.Height = vh

				return m, nil
			case "q":
				if m.r != nil && m.r.Current() == router.RouteProjects {
					m.selectedProject = nil
					return m, tea.Quit
				}

				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			if m.r != nil && m.r.Current() != router.RouteProjects {
				m.mode = 0
				m.form = nil
				m.envForm = nil
				m.r.NavigateTo(router.RouteProjects, nil)

				if m.projectsView != nil {
					return m, tea.Batch(cmd, m.projectsView.Init())
				}

				return m, tea.ClearScreen
			}

			m.selectedProject = nil

			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *Window) reduce(msg tea.Msg) {
	switch v := msg.(type) {
	case state.NavigateToMsg:
		if m.r != nil {
			m.r.NavigateTo(v.Route, v.Params)
		}
	case state.BackMsg:
		if m.r != nil {
			m.r.Back()
		}
	case state.ExitMsg:
		if m.r != nil {
			if m.r.Current() == router.RouteEnvVars || (m.r.Current() == router.RouteProjects && m.focus == 1) {
				pname := m.envProjectName
				if pname == "" {
					if pv, ok := m.projectsView.(views.ProjectsView); ok {
						idx := pv.List.Cursor
						if idx >= 0 && idx < len(pv.List.Items) {
							it := pv.List.Items[idx]
							if it.ID != "__new__" {
								pname = it.ID
							}
						}
					}
				}

				if pname != "" {
					proj, _ := m.projectSvc.Load(pname)
					m.selectedProject = proj
				}

				if ev, ok := m.envVarsView.(views.EnvVarsView); ok {
					if len(ev.List.Items) > 0 {
						it := ev.List.Items[ev.List.Cursor]
						if it.ID != "__add__" {
							m.selectedEnv = it.ID
						}
					}
				}
			}
		}

		m.shouldQuit = true
	}
}
