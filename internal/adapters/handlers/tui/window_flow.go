package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

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

func isValidColorInput(s string) bool {
	ss := strings.ToLower(strings.TrimSpace(s))
	if ss == "" {
		return false
	}

	if strings.HasPrefix(ss, "#") {
		h := ss[1:]
		if len(h) != 3 && len(h) != 6 {
			return false
		}

		if _, err := strconv.ParseUint(h, 16, 64); err != nil {
			return false
		}

		return true
	}

	if _, ok := colorNameMap[ss]; ok {
		return true
	}

	if n, err := strconv.Atoi(ss); err == nil && n >= 0 && n <= 255 {
		return true
	}

	return false
}

func (m *Window) newProjectFlow() {
	m.form = NewProjectFormModel()
	m.mode = 1
}

func (m *Window) newEnvironmentFlow(projectName string) {
	m.envForm = NewEnvironmentFormModel()
	m.envProjectName = projectName
	m.mode = 3
}

type postCreatePrompt struct {
	idx, choice int
	projectName string
}

func newPostCreatePrompt(projectName string) *postCreatePrompt {
	return &postCreatePrompt{projectName: projectName}
}

func (p *postCreatePrompt) Init() tea.Cmd { return nil }

func (p *postCreatePrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyUp:
			if p.idx > 0 {
				p.idx--
			}
		case tea.KeyDown:
			if p.idx < 2 {
				p.idx++
			}
		case tea.KeyEnter:
			p.choice = p.idx

			return p, tea.Quit
		case tea.KeyEsc, tea.KeyCtrlC:
			p.choice = 0

			return p, tea.Quit
		}

		if m.Type == tea.KeyRunes {
			switch string(m.Runes) {
			case "k":
				if p.idx > 0 {
					p.idx--
				}
			case "j":
				if p.idx < 2 {
					p.idx++
				}
			case "q":
				p.choice = 0
				return p, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		return p, nil
	}

	return p, nil
}

func (p *postCreatePrompt) View() string {
	styleTitle := lipgloss.NewStyle().
		Foreground(c("title")).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true).
		BorderForeground(c("border")).
		Padding(0, 1)
	styleSel := lipgloss.NewStyle().Foreground(c("selectedFg")).Background(c("selectedBg")).Padding(0, 1)
	styleDef := lipgloss.NewStyle().Foreground(c("subtext")).Padding(0, 1)
	b := strings.Builder{}
	b.WriteString(styleTitle.Render("Project created. What next?"))
	b.WriteString("\n")

	opts := []string{"Return to list", "Start this project", "Exit"}
	for i, o := range opts {
		if i == p.idx {
			b.WriteString("➜ ")
			b.WriteString(styleSel.Render(o))
		} else {
			b.WriteString(" ")
			b.WriteString(styleDef.Render(o))
		}

		b.WriteString("\n")
	}

	help := lipgloss.NewStyle().
		Foreground(c("help")).
		Render(`Navigate: ↑/k ↓/j
Select: Enter  Cancel: Esc/q/Ctrl+C`)

	return lipgloss.JoinVertical(lipgloss.Left, b.String(), help)
}
