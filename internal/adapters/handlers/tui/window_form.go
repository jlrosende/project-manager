package tui

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	bti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NewProjectForm struct {
	name          bti.Model
	path          bti.Model
	subproject    bti.Model
	userName      bti.Model
	userEmail     bti.Model
	userSigningKey bti.Model
	commitGPGSign bti.Model
	tagGPGSign    bti.Model
	envVars       bti.Model
	focused       int
	width         int
	height        int
	submitted     bool
	canceled      bool
	err           string
}

func NewProjectFormModel() *NewProjectForm {
	ni := bti.New()
	ni.Prompt = "Name: "
	ni.Placeholder = "my-project"
	ni.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	ni.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ni.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	ni.Focus()
	pi := bti.New()
	pi.Prompt = "Path: "
	pi.Placeholder = "~/code/my-project"
	pi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	pi.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	pi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	spi := bti.New()
	spi.Prompt = "Subproject: "
	spi.Placeholder = ""
	spi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	spi.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	spi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	un := bti.New()
	un.Prompt = "Git user.name: "
	un.Placeholder = ""
	un.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	un.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	un.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	uem := bti.New()
	uem.Prompt = "Git user.email: "
	uem.Placeholder = ""
	uem.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	uem.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	uem.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	usk := bti.New()
	usk.Prompt = "Git user.signingkey: "
	usk.Placeholder = ""
	usk.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	usk.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	usk.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	csg := bti.New()
	csg.Prompt = "commit.gpgsign (true/false): "
	csg.Placeholder = "true"
	csg.SetValue("true")
	csg.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	csg.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	csg.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	tsg := bti.New()
	tsg.Prompt = "tag.gpgsign (true/false): "
	tsg.Placeholder = "true"
	tsg.SetValue("true")
	tsg.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	tsg.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	tsg.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	ev := bti.New()
	ev.Prompt = "Env vars (K=V,K2=V2): "
	ev.Placeholder = ""
	ev.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	ev.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ev.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	return &NewProjectForm{name: ni, path: pi, subproject: spi, userName: un, userEmail: uem, userSigningKey: usk, commitGPGSign: csg, tagGPGSign: tsg, envVars: ev, focused: 0}
}

func (f *NewProjectForm) Init() tea.Cmd { return bti.Blink }

func (f *NewProjectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			f.canceled = true
			f.submitted = false
			return f, nil
		case "tab", "shift+tab":
			if msg.String() == "tab" {
				f.focused++
				if f.focused > 8 {
					f.focused = 0
				}
			} else {
				f.focused--
				if f.focused < 0 {
					f.focused = 8
				}
			}
			f.blurAll()
			f.focusCurrent()
			return f, nil
		case "enter":
			if f.focused < 8 {
				f.focused++
				f.blurAll()
				f.focusCurrent()
				return f, nil
			}
			f.err = ""
			if strings.TrimSpace(f.name.Value()) == "" {
				f.err = "name is required"
				return f, nil
			}
			if strings.TrimSpace(f.path.Value()) == "" {
				f.err = "path is required"
				return f, nil
			}
			if ok, reason := canCreatePath(f.path.Value()); !ok {
				f.err = reason
				return f, nil
			}
			cg := strings.ToLower(strings.TrimSpace(f.commitGPGSign.Value()))
			if cg != "" && cg != "true" && cg != "false" {
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
	var cmd tea.Cmd
	switch f.focused {
	case 0:
		f.name, cmd = f.name.Update(msg)
	case 1:
		f.path, cmd = f.path.Update(msg)
	case 2:
		f.subproject, cmd = f.subproject.Update(msg)
	case 3:
		f.userName, cmd = f.userName.Update(msg)
	case 4:
		f.userEmail, cmd = f.userEmail.Update(msg)
	case 5:
		f.userSigningKey, cmd = f.userSigningKey.Update(msg)
	case 6:
		f.commitGPGSign, cmd = f.commitGPGSign.Update(msg)
	case 7:
		f.tagGPGSign, cmd = f.tagGPGSign.Update(msg)
	case 8:
		f.envVars, cmd = f.envVars.Update(msg)
	}
	return f, cmd
}

func (f *NewProjectForm) blurAll() {
	f.name.Blur()
	f.path.Blur()
	f.subproject.Blur()
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
		f.path.Focus()
	case 2:
		f.subproject.Focus()
	case 3:
		f.userName.Focus()
	case 4:
		f.userEmail.Focus()
	case 5:
		f.userSigningKey.Focus()
	case 6:
		f.commitGPGSign.Focus()
	case 7:
		f.tagGPGSign.Focus()
	case 8:
		f.envVars.Focus()
	}
}

func (f *NewProjectForm) View() string {
	box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Tab switch  Enter next/submit  Esc cancel")
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	b := strings.Builder{}
	b.WriteString(f.name.View())
	b.WriteString("\n")
	b.WriteString(f.path.View())
	b.WriteString("\n")
	b.WriteString(f.subproject.View())
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
	b.WriteString("\n")
	b.WriteString(f.envVars.View())
	b.WriteString("\n\n")
	if f.err != "" {
		b.WriteString(errStyle.Render(f.err))
		b.WriteString("\n")
	}
	b.WriteString(help)
	return box.Render(b.String())
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
