package tui

import (
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func (m *Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.shouldQuit {
		return m, tea.Quit
	}
	var cmd tea.Cmd
	if m.mode == 1 && m.form != nil {
		fm, fcmd := m.form.Update(msg)
		if f, ok := fm.(*NewProjectForm); ok {
			m.form = f
			if f.canceled {
				m.mode = 0
				m.form = nil
				return m, tea.ClearScreen
			}
			if f.submitted {
				name := strings.TrimSpace(f.name.Value())
				path := strings.TrimSpace(f.path.Value())
				if name == "" || path == "" {
					return m, nil
				}
				envs := domain.EnvVars{}
				if strings.TrimSpace(f.envVars.Value()) != "" {
					envs = parseEnvLines(f.envVars.Value())
				}
				commitSign := strings.ToLower(strings.TrimSpace(f.commitGPGSign.Value())) != "false"
				tagSign := strings.ToLower(strings.TrimSpace(f.tagGPGSign.Value())) != "false"
				proj, _ := m.projectSvc.Create(
					name,
					path,
					f.subproject.Value(),
					envs,
					domain.New(
						domain.WithName(strings.TrimSpace(f.userName.Value())),
						domain.WithEmail(strings.TrimSpace(f.userEmail.Value())),
						domain.WithSigningKey(strings.TrimSpace(f.userSigningKey.Value())),
						domain.WithCommitSign(commitSign),
						domain.WithTagSign(tagSign),
					),
				)
				projects, err := m.projectSvc.List()
				if err == nil {
					m.projects = projects
					m.total = len(projects) + 1
				}
				m.prompt = newPostCreatePrompt(name)
				m.mode = 2
				_ = proj
				return m, tea.ClearScreen
			}
			return m, fcmd
		}
	}
	if m.mode == 2 && m.prompt != nil {
		pm, pcmd := m.prompt.Update(msg)
		if p, ok := pm.(*postCreatePrompt); ok {
			m.prompt = p
			if p.choice == 0 { // return
				m.mode = 0
				m.prompt = nil
				return m, tea.ClearScreen
			}
			if p.choice == 1 { // start project
				proj, _ := m.projectSvc.Load(p.projectName)
				m.selectedProject = proj
				return m, tea.Quit
			}
			if p.choice == 2 { // exit
				m.selectedProject = nil
				return m, tea.Quit
			}
			return m, pcmd
		}
	}
	if m.mode == 3 && m.envForm != nil {
		efm, ecmd := m.envForm.Update(msg)
		if f, ok := efm.(*NewEnvironmentForm); ok {
			m.envForm = f
			if f.canceled {
				m.mode = 0
				m.envForm = nil
				return m, tea.ClearScreen
			}
			if f.submitted {
				name := strings.TrimSpace(f.name.Value())
				if name == "" || m.envProjectName == "" {
					return m, nil
				}
				env := &domain.Environment{
					Name:        name,
					Color:       normalizeColorInput(strings.TrimSpace(f.color.Value())),
					EnvVarsMode: strings.TrimSpace(f.mode.Value()),
				}
				envs := domain.EnvVars{}
				if strings.TrimSpace(f.envVars.Value()) != "" {
					envs = parseEnvLines(f.envVars.Value())
				}
				_ = m.projectSvc.AddEnvironment(m.envProjectName, env, envs)
				projects, err := m.projectSvc.List()
				if err == nil {
					m.projects = projects
					m.total = len(projects) + 1
				}
				m.mode = 0
				m.envForm = nil
				m.envProjectName = ""
				return m, tea.ClearScreen
			}
			return m, ecmd
		}
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.selectedProject = nil
			return m, tea.Quit
		case "enter":
			if m.focus == 0 {
				idx := mod(m.cursor, m.total)
				if idx == len(m.projects) { // '+ New project'
					m.newProjectFlow()
					return m, tea.ClearScreen
				}
				project, _ := m.projectSvc.Load(m.projects[idx].Name)
				m.selectedProject = project
				return m, tea.Quit
			}
			if m.focus == 1 {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					envCount := len(m.projects[mod].Environments)
					if m.cursorEnv == envCount { // '+ New environment'
						m.newEnvironmentFlow(m.projects[mod].Name)
						return m, tea.ClearScreen
					}
					if envCount > 0 {
						project, _ := m.projectSvc.Load(m.projects[mod].Name)
						m.selectedProject = project
						return m, tea.Quit
					}
				}
			}

			// The "up" and "k" keys move the cursor up
		case "left", "h":
			m.focus = 0
		case "right", "l":
			mod := mod(m.cursor, m.total)
			if mod < len(m.projects) && len(m.projects) > 0 {
				m.focus = 1
				if len(m.projects[mod].Environments) == 0 {
					m.cursorEnv = 0
				}
			}
		case "up", "k":
			if m.focus == 0 {
				m.cursor--
				m.cursorEnv = 0
			} else {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					envCount := len(m.projects[mod].Environments)
					envCountPlus := envCount + 1
					if envCountPlus > 0 {
						m.cursorEnv--
						if m.cursorEnv < 0 {
							m.cursorEnv = envCountPlus - 1
						}
					}
				}
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.focus == 0 {
				m.cursor++
				m.cursorEnv = 0
			} else {
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					envCount := len(m.projects[mod].Environments)
					envCountPlus := envCount + 1
					if envCountPlus > 0 {
						m.cursorEnv++
						if m.cursorEnv >= envCountPlus {
							m.cursorEnv = 0
						}
					}
				}
			}
		}
	}

	return m, tea.Batch(
		cmd,
		tea.Printf("Let's go to %d!", m.cursor),
	)
}
