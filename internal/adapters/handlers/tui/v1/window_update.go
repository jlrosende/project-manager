package v1

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
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

				commitSign := strings.ToLower(strings.TrimSpace(f.commitGPGSign.Value())) != strFalse

				tagSign := strings.ToLower(strings.TrimSpace(f.tagGPGSign.Value())) != "false"
				if f.isEdit {
					proj, _ := m.projectSvc.Load(f.originalName)
					proj.Name = name
					proj.Shell = strings.TrimSpace(f.shell.Value())
					_ = m.projectSvc.UpdateProject(proj)
					// Validate but preserve raw content
					raw := strings.TrimSpace(f.envVars.Value())
					if raw != "" {
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
							projEnvPath := proj.EnvVarsFile
							if !filepath.IsAbs(projEnvPath) {
								projEnvPath = filepath.Join(proj.Path, projEnvPath)
							}

							_ = os.WriteFile(projEnvPath, []byte(f.envVars.Value()), 0o600)
						}
					}

					gitRepo, err := repositories.NewGitRepository()
					if err == nil {
						gitSvc := services.NewGitService(gitRepo)
						gitPath := filepath.Join(proj.Path, fmt.Sprintf(".%s.gitconfig", proj.Name))
						gc := &domain.GitConfig{
							User: domain.User{
								Name:       strings.TrimSpace(f.userName.Value()),
								Email:      strings.TrimSpace(f.userEmail.Value()),
								SigningKey: strings.TrimSpace(f.userSigningKey.Value()),
							},
							Commit: domain.Commit{GPGSign: commitSign},
							Tag:    domain.Tag{GPGSign: tagSign},
						}
						_ = gitSvc.Save(gitPath, gc)
					}

					projects, err := m.projectSvc.List()
					if err == nil {
						m.projects = projects
						m.total = len(projects) + 1
					}

					m.mode = 0
					m.form = nil

					return m, tea.ClearScreen
				}

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
				if proj != nil {
					proj.Shell = strings.TrimSpace(f.shell.Value())
					_ = m.projectSvc.UpdateProject(proj)
				}

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

				if f.isEdit {
					_ = m.projectSvc.UpdateEnvironment(m.envProjectName, f.originalName, env)
					// Persist env vars to env file
					project, _ := m.projectSvc.Load(m.envProjectName)
					// find updated env (by new name)
					var envPath string

					for _, e := range project.Environments {
						if e.Name == env.Name {
							envPath = e.EnvVarsFile
							break
						}
					}

					if envPath != "" {
						if !filepath.IsAbs(envPath) {
							envPath = filepath.Join(project.Path, envPath)
						}
						// validate but keep raw
						raw := strings.TrimSpace(f.envVars.Value())
						if raw != "" {
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
								_ = os.WriteFile(envPath, []byte(f.envVars.Value()), 0o600)
							}
						}
					}
				} else {
					_ = m.projectSvc.AddEnvironment(m.envProjectName, env, envs)
				}

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
		m.vp.Width = msg.Width

		hh := lipgloss.Height(m.help.View(m.keys))

		vh := msg.Height - hh
		if vh < 1 {
			vh = 1
		}

		m.vp.Height = vh
	case tea.KeyMsg:
		if m.mode == 0 {
			var vcmd tea.Cmd

			m.vp, vcmd = m.vp.Update(msg)
			cmd = tea.Batch(cmd, vcmd)
		}

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
			case "e":
				idx := mod(m.cursor, m.total)
				if m.focus == 0 && idx < len(m.projects) {
					project, _ := m.projectSvc.Load(m.projects[idx].Name)
					m.form = NewProjectEditFormModel(project)

					projEnvPath := project.EnvVarsFile
					if !filepath.IsAbs(projEnvPath) {
						projEnvPath = filepath.Join(project.Path, projEnvPath)
					}

					if b, err := os.ReadFile(projEnvPath); err == nil {
						m.form.envVars.SetValue(string(b))
					} else {
						pairs := project.EnvVars.ToSlice()
						sort.Strings(pairs)
						m.form.envVars.SetValue(strings.Join(pairs, "\n"))
					}

					gitRepo, err := repositories.NewGitRepository()
					if err == nil {
						gitSvc := services.NewGitService(gitRepo)

						gitPath := filepath.Join(project.Path, fmt.Sprintf(".%s.gitconfig", project.Name))
						if _, err := os.Stat(gitPath); os.IsNotExist(err) {
							if matches, _ := filepath.Glob(filepath.Join(project.Path, ".*.gitconfig")); len(matches) > 0 {
								gitPath = matches[0]
							}
						}

						if gc, e := gitSvc.Load(gitPath); e == nil {
							m.form.userName.SetValue(gc.User.Name)
							m.form.userEmail.SetValue(gc.User.Email)
							m.form.userSigningKey.SetValue(gc.User.SigningKey)

							if gc.Commit.GPGSign {
								m.form.commitGPGSign.SetValue(strTrue)
							} else {
								m.form.commitGPGSign.SetValue(strFalse)
							}

							if gc.Tag.GPGSign {
								m.form.tagGPGSign.SetValue(strTrue)
							} else {
								m.form.tagGPGSign.SetValue(strFalse)
							}
						}
					}

					m.mode = 1

					return m, tea.ClearScreen
				}

				if m.focus == 1 && idx < len(m.projects) {
					project, _ := m.projectSvc.Load(m.projects[idx].Name)
					if m.cursorEnv < len(project.Environments) {
						env := project.Environments[m.cursorEnv]
						m.envForm = NewEnvironmentEditFormModel(env)

						envPath := env.EnvVarsFile
						if !filepath.IsAbs(envPath) {
							envPath = filepath.Join(project.Path, envPath)
						}

						if b, err := os.ReadFile(envPath); err == nil {
							m.envForm.envVars.SetValue(string(b))
						} else if len(env.EnvVars) > 0 {
							pairs := env.EnvVars.ToSlice()
							sort.Strings(pairs)
							m.envForm.envVars.SetValue(strings.Join(pairs, "\n"))
						}

						m.envProjectName = project.Name
						m.mode = 3

						return m, tea.ClearScreen
					}
				}
			case "h":
				m.focus = 0
			case "l":
				mod := mod(m.cursor, m.total)
				if mod < len(m.projects) && len(m.projects) > 0 {
					m.focus = 1
					if len(m.projects[mod].Environments) == 0 {
						m.cursorEnv = 0
					}
				}
			case "k":
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
			case "j":
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
			case "q":
				m.selectedProject = nil
				return m, tea.Quit
			}
		}

		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			m.selectedProject = nil
			return m, tea.Quit
		case tea.KeyEnter:
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
		case tea.KeyLeft:
			m.focus = 0
		case tea.KeyRight:
			mod := mod(m.cursor, m.total)
			if mod < len(m.projects) && len(m.projects) > 0 {
				m.focus = 1
				if len(m.projects[mod].Environments) == 0 {
					m.cursorEnv = 0
				}
			}
		case tea.KeyUp:
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
		case tea.KeyDown:
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

	return m, cmd
}
