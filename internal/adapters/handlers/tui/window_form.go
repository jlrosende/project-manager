package tui

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	helpcomp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	bta "github.com/charmbracelet/bubbles/textarea"
	bti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

const (
	strTrue         = "true"
	strFalse        = "false"
	errPathRequired = "path is required"
)

type formKeyMap struct {
	Save    key.Binding
	Cancel  key.Binding
	Next    key.Binding
	Prev    key.Binding
	Buttons key.Binding
	Help    key.Binding
}

func (k formKeyMap) ShortHelp() []key.Binding { return []key.Binding{k.Next, k.Save, k.Help} }
func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Next, k.Prev, k.Buttons}, {k.Save, k.Cancel, k.Help}}
}

type NewProjectForm struct {
	name           bti.Model
	path           bti.Model
	subproject     bti.Model
	shell          bti.Model
	userName       bti.Model
	userEmail      bti.Model
	userSigningKey bti.Model
	commitGPGSign  bti.Model
	tagGPGSign     bti.Model
	envVars        bta.Model
	focused        int
	width          int
	height         int
	pathDirty      bool
	submitted      bool
	canceled       bool
	err            string
	isEdit         bool
	originalName   string
	help           helpcomp.Model
	keys           formKeyMap
}

func NewProjectFormModel() *NewProjectForm {
	ni := bti.New()
	sc := lipgloss.NewStyle().Foreground(c("title"))
	sr := lipgloss.NewStyle().Foreground(c("error"))
	ni.Prompt = sc.Render("Name") + sr.Render("*") + sc.Render(": ")
	ni.Placeholder = "my-awesome-app"
	ni.PromptStyle = lipgloss.NewStyle()
	ni.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	ni.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	ni.Width = 40
	ni.Focus()

	pi := bti.New()
	pi.Prompt = sc.Render("Path") + sr.Render("*") + sc.Render(": ")
	pi.Placeholder = "~/my-awesome-app"
	pi.PromptStyle = lipgloss.NewStyle()
	pi.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	pi.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	pi.Width = 40
	spi := bti.New()
	spi.Prompt = "Subproject: "
	spi.Placeholder = "services/api"
	spi.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	spi.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	spi.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	spi.Width = 40
	un := bti.New()
	un.Prompt = "Git user.name: "
	un.Placeholder = "Jane Doe"
	un.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	un.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	un.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	un.Width = 40
	uem := bti.New()
	uem.Prompt = "Git user.email: "
	uem.Placeholder = "jane@example.com"
	uem.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	uem.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	uem.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	uem.Width = 40
	usk := bti.New()
	usk.Prompt = "Git user.signingkey: "
	usk.Placeholder = "0xDEADBEEF"
	usk.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	usk.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	usk.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	usk.Width = 40
	csg := bti.New()
	csg.Prompt = "commit.gpgsign (true/false): "
	csg.Placeholder = strTrue
	csg.SetValue(strTrue)
	csg.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	csg.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	csg.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	csg.Width = 40
	tsg := bti.New()
	tsg.Prompt = "tag.gpgsign (true/false): "
	tsg.Placeholder = strTrue
	tsg.SetValue(strTrue)
	tsg.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	tsg.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	tsg.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	tsg.Width = 40
	sh := bti.New()
	sh.Prompt = "Shell: "
	sh.Placeholder = "/bin/bash"
	sh.PromptStyle = lipgloss.NewStyle().Foreground(c("title"))
	sh.PlaceholderStyle = lipgloss.NewStyle().Foreground(c("placeholder"))
	sh.TextStyle = lipgloss.NewStyle().Foreground(c("text"))
	sh.Width = 40
	ev := bta.New()
	ev.Placeholder = "# One per line (like .env)\nAPP_ENV=development\nDATABASE_URL=postgres://user:pass@localhost:5432/app\n# comments allowed"
	ev.SetHeight(6)
	ev.SetWidth(60)

	h := helpcomp.New()
	km := formKeyMap{}
	km.Save = key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save"))
	km.Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	km.Next = key.NewBinding(key.WithKeys("tab", "down", "enter"), key.WithHelp("tab/↓/enter", "next"))
	km.Prev = key.NewBinding(key.WithKeys("shift+tab", "up"), key.WithHelp("shift+tab/↑", "prev"))
	km.Buttons = key.NewBinding(key.WithKeys("left", "right"), key.WithHelp("←/→", "buttons"))
	km.Help = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))

	return &NewProjectForm{
		name:           ni,
		path:           pi,
		subproject:     spi,
		shell:          sh,
		userName:       un,
		userEmail:      uem,
		userSigningKey: usk,
		commitGPGSign:  csg,
		tagGPGSign:     tsg,
		envVars:        ev,
		focused:        0,
		help:           h,
		keys:           km,
	}
}

func (f *NewProjectForm) Init() tea.Cmd { return bti.Blink }

func (f *NewProjectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

		w := msg.Width - 10
		if w < 20 {
			w = 20
		}

		f.name.Width = w
		f.path.Width = w
		f.subproject.Width = w
		f.shell.Width = w
		f.userName.Width = w
		f.userEmail.Width = w
		f.userSigningKey.Width = w
		f.commitGPGSign.Width = w
		f.tagGPGSign.Width = w
		f.envVars.SetWidth(w + 10)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			f.canceled = true
			f.submitted = false

			return f, nil
		case tea.KeyTab, tea.KeyShiftTab:
			if msg.Type == tea.KeyTab {
				f.focused++
				if f.focused > 11 {
					f.focused = 0
				}
			} else {
				f.focused--
				if f.focused < 0 {
					f.focused = 11
				}
			}

			if f.isEdit && f.focused == 1 {
				if msg.Type == tea.KeyTab {
					f.focused++
				} else {
					f.focused--
				}
			}

			f.blurAll()
			f.focusCurrent()

			return f, nil
		case tea.KeyEnter:
			switch {
			case f.focused < 9:
				if f.isEdit && f.focused == 1 {
					f.focused++
				}

				f.focused++
				f.blurAll()
				f.focusCurrent()

				return f, nil
			case f.focused == 10:
				f.err = ""
				if strings.TrimSpace(f.name.Value()) == "" {
					f.err = "name is required"
					return f, nil
				}

				if !f.isEdit {
					if strings.TrimSpace(f.path.Value()) == "" {
						f.err = errPathRequired
						return f, nil
					}

					if ok, reason := canCreatePath(f.path.Value()); !ok {
						f.err = reason
						return f, nil
					}
				}

				cg := strings.ToLower(strings.TrimSpace(f.commitGPGSign.Value()))
				if cg != "" && cg != strTrue && cg != strFalse {
					f.err = "commit.gpgsign must be true or false"
					return f, nil
				}

				tg := strings.ToLower(strings.TrimSpace(f.tagGPGSign.Value()))
				if tg != "" && tg != strTrue && tg != strFalse {
					f.err = "tag.gpgsign must be true or false"
					return f, nil
				}

				f.submitted = true
				f.canceled = false

				return f, nil
			case f.focused == 11:
				f.canceled = true
				f.submitted = false

				return f, nil
			}
		case tea.KeyUp:
			if f.focused != 9 {
				f.focused--
				if f.isEdit && f.focused == 1 {
					f.focused--
				}

				if f.focused < 0 {
					f.focused = 11
				}

				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyDown:
			if f.focused != 9 {
				f.focused++
				if f.isEdit && f.focused == 1 {
					f.focused++
				}

				if f.focused > 11 {
					f.focused = 0
				}

				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyLeft:
			if f.focused == 11 {
				f.focused = 10
				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyRight:
			if f.focused == 10 {
				f.focused = 11
				f.blurAll()
				f.focusCurrent()

				return f, nil
			}
		case tea.KeyCtrlS:
			f.err = ""
			if strings.TrimSpace(f.name.Value()) == "" {
				f.err = "name is required"

				return f, nil
			}

			if !f.isEdit {
				if strings.TrimSpace(f.path.Value()) == "" {
					f.err = errPathRequired

					return f, nil
				}

				if ok, reason := canCreatePath(f.path.Value()); !ok {
					f.err = reason

					return f, nil
				}
			}

			cg := strings.ToLower(strings.TrimSpace(f.commitGPGSign.Value()))
			if cg != "" && cg != strTrue && cg != strFalse {
				f.err = "commit.gpgsign must be true or false"

				return f, nil
			}

			tg := strings.ToLower(strings.TrimSpace(f.tagGPGSign.Value()))
			if tg != "" && tg != "true" && tg != "false" {
				f.err = "tag.gpgsign must be true or false"

				return f, nil
			}

			f.submitted = true
			f.canceled = false

			return f, nil
		}
	}

	if m, ok := msg.(tea.KeyMsg); ok {
		if m.Type == tea.KeyRunes && len(m.Runes) > 0 && m.Runes[0] == '?' {
			f.help.ShowAll = !f.help.ShowAll
			return f, nil
		}
	}

	var cmd tea.Cmd

	switch f.focused {
	case 0:
		oldName := f.name.Value()

		f.name, cmd = f.name.Update(msg)

		if f.name.Value() != oldName && !f.pathDirty && !f.isEdit {
			slug := slugify(strings.TrimSpace(f.name.Value()))
			if slug != "" {
				cur := strings.TrimSpace(f.path.Value())

				var newPath string

				switch {
				case cur == "":
					newPath = filepath.Join("~/", slug)
				case strings.HasSuffix(cur, string(os.PathSeparator)):
					newPath = cur + slug
				default:
					dir := filepath.Dir(cur)
					newPath = filepath.Join(dir, slug)
				}

				f.path.SetValue(newPath)
			}
		}
	case 1:
		if f.isEdit {
			break
		}

		prev := f.path.Value()

		f.path, cmd = f.path.Update(msg)

		if f.path.Value() != prev {
			f.pathDirty = true
		}
	case 2:
		f.subproject, cmd = f.subproject.Update(msg)
	case 3:
		f.shell, cmd = f.shell.Update(msg)
	case 4:
		f.userName, cmd = f.userName.Update(msg)
	case 5:
		f.userEmail, cmd = f.userEmail.Update(msg)
	case 6:
		f.userSigningKey, cmd = f.userSigningKey.Update(msg)
	case 7:
		f.commitGPGSign, cmd = f.commitGPGSign.Update(msg)
	case 8:
		f.tagGPGSign, cmd = f.tagGPGSign.Update(msg)
	case 9:
		f.envVars, cmd = f.envVars.Update(msg)
	case 10:
		// save button: nothing to update
	case 11:
		// cancel button: nothing to update
	}

	return f, cmd
}

func (f *NewProjectForm) blurAll() {
	f.name.Blur()
	f.path.Blur()
	f.subproject.Blur()
	f.shell.Blur()
	f.userName.Blur()
	f.userEmail.Blur()
	f.userSigningKey.Blur()
	f.commitGPGSign.Blur()
	f.tagGPGSign.Blur()
	f.envVars.Blur()
}

func (f *NewProjectForm) focusCurrent() {
	switch f.focused {
	case 0:
		f.name.Focus()
	case 1:
		if !f.isEdit {
			f.path.Focus()
		}
	case 2:
		f.subproject.Focus()
	case 3:
		f.shell.Focus()
	case 4:
		f.userName.Focus()
	case 5:
		f.userEmail.Focus()
	case 6:
		f.userSigningKey.Focus()
	case 7:
		f.commitGPGSign.Focus()
	case 8:
		f.tagGPGSign.Focus()
	case 9:
		f.envVars.Focus()
	}
}

func (f *NewProjectForm) View() string {
	errStyle := lipgloss.NewStyle().Foreground(c("error"))
	styleTitle := lipgloss.NewStyle().
		Foreground(c("title")).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true).
		BorderForeground(c("border")).
		Padding(0, 1)
	styleSection := lipgloss.NewStyle().Foreground(c("section")).Bold(true)
	styleSub := lipgloss.NewStyle().Foreground(c("subtext"))
	btnDef := lipgloss.NewStyle().Foreground(c("buttonDefFg")).Background(c("buttonDefBg")).Padding(0, 2)
	btnSel := lipgloss.NewStyle().Foreground(c("buttonSelFg")).Background(c("buttonSelBg")).Padding(0, 2)
	b := strings.Builder{}

	title := "New project"
	if f.isEdit {
		title = "Edit project"
	}

	b.WriteString(styleTitle.Render(title))
	b.WriteString("\n")
	b.WriteString(styleSection.Render("Project details"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("Basic information"))
	b.WriteString("\n")
	b.WriteString(f.name.View())
	b.WriteString("\n")

	pathView := f.path.View()
	if f.isEdit {
		pathView = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(f.path.Value() + " (read-only)")
	}

	b.WriteString(pathView)
	b.WriteString("\n\n")
	b.WriteString(styleSection.Render("Project options"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("Optional settings"))
	b.WriteString("\n")
	b.WriteString(f.subproject.View())
	b.WriteString("\n")
	b.WriteString(f.shell.View())
	b.WriteString("\n\n")
	b.WriteString(styleSection.Render("Git settings"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("Applied to this project only"))
	b.WriteString("\n")
	b.WriteString(f.userName.View())
	b.WriteString("\n")
	b.WriteString(f.userEmail.View())
	b.WriteString("\n")
	b.WriteString(f.userSigningKey.View())
	b.WriteString("\n")
	b.WriteString(f.commitGPGSign.View())
	b.WriteString("\n")
	b.WriteString(f.tagGPGSign.View())
	b.WriteString("\n\n")
	b.WriteString(styleSection.Render("Environment variables"))
	b.WriteString("\n")
	b.WriteString(styleSub.Render("One KEY=VALUE per line; '#' comments allowed"))
	b.WriteString("\n")
	b.WriteString(f.envVars.View())
	b.WriteString("\n\n")

	btnSave := btnDef.Render(" Save ")
	if f.focused == 10 {
		btnSave = btnSel.Render(" Save ")
	}

	btnCancel := btnDef.Render(" Cancel ")
	if f.focused == 11 {
		btnCancel = btnSel.Render(" Cancel ")
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, btnSave, "   ", btnCancel))
	b.WriteString("\n\n")

	if f.err != "" {
		b.WriteString(errStyle.Render(f.err))
		b.WriteString("\n")
	}

	content := b.String()

	sepLen := f.width
	if sepLen < 1 {
		sepLen = 80
	}

	sep := lipgloss.NewStyle().Foreground(c("border")).Render(strings.Repeat("─", sepLen))

	return lipgloss.JoinVertical(lipgloss.Left, content, sep, f.help.View(f.keys))
}

func NewProjectEditFormModel(p *domain.Project) *NewProjectForm {
	f := NewProjectFormModel()
	f.isEdit = true
	f.originalName = p.Name
	f.name.SetValue(p.Name)
	f.path.SetValue(p.Path)
	f.shell.SetValue(p.Shell)

	return f
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return ""
	}

	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9-_]+`)
	s = re.ReplaceAllString(s, "-")
	re2 := regexp.MustCompile(`-+`)
	s = re2.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-_")

	return s
}

func canCreatePath(p string) (bool, string) {
	pp := strings.TrimSpace(p)
	if pp == "" {
		return false, "path is required"
	}

	if strings.HasPrefix(pp, "~/") {
		h, _ := os.UserHomeDir()
		pp = filepath.Join(h, pp[2:])
	}

	info, err := os.Stat(pp)
	if err == nil {
		if !info.IsDir() {
			return false, "path exists and is not a directory"
		}

		empty, e := isDirEmpty(pp)
		if e != nil {
			return false, e.Error()
		}

		if !empty {
			return false, "directory is not empty"
		}

		return true, ""
	}

	parent := filepath.Dir(pp)

	pi, perr := os.Stat(parent)
	if perr != nil || !pi.IsDir() {
		return false, "parent directory does not exist"
	}

	return true, ""
}
