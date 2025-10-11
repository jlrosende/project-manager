package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/configs"
	"github.com/jlrosende/project-manager/internal"
	cmdDelete "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	cmdEdit "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/edit"
	cmdInit "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/init"
	cmdList "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/list"
	cmdNew "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/new"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	repositories "github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/adapters/repositories/shells"
	"github.com/jlrosende/project-manager/internal/bootstrap"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pm [project] [env|[path]] ",
		Short:        "pm is a tool to create and organize projects in your computer",
		Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
		Version:      internal.GetVersion(),
		Args:         cobra.MaximumNArgs(3),
		SilenceUsage: true,
		RunE:         root,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			root := cmd.Root()

			logLevel, err := root.PersistentFlags().GetString("log-level")
			if err != nil {
				return err
			}

			logFile, err := root.PersistentFlags().GetString("log-file")
			if err != nil {
				return err
			}

			_, err = repositories.SetupLogger(logLevel, logFile)
			if err != nil {
				return err
			}

			return nil
		},
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) >= 1 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return projectsList(), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().BoolP("list", "l", false, "List all the projects.")

	cmd.PersistentFlags().String("log-level", "info", "Change the log level (debug, info, warn, error)")
	cmd.PersistentFlags().String("log-file", "", "Path to log file (default: $XDG_CACHE_HOME/pm.log)")
	cmd.PersistentFlags().String("theme", "", "Theme for this run (nord, catppuccin, dracula, ayu)")
	cmd.PersistentFlags().String("config", "", "Path to pm config file")

	cmd.AddCommand(cmdInit.InitCmd)
	cmd.AddCommand(cmdNew.Command())
	cmd.AddCommand(cmdList.ListCmd)
	cmd.AddCommand(cmdEdit.Command())
	cmd.AddCommand(cmdDelete.Command())

	return cmd
}

func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		code := exitCodeFromError(err)
		slog.Error("something wrong happened", slog.Any("err", err), slog.Int("exit_code", code))
		os.Exit(code)
	}
}

// Root exposes the root command for tools like doc generators.
func Root() *cobra.Command { return newRootCommand() }

type exitCoder interface {
	ExitCode() int
}

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}

	var ec exitCoder
	if errors.As(err, &ec) {
		return ec.ExitCode()
	}

	return 1
}

func root(cmd *cobra.Command, args []string) error {
	container, err := bootstrap.New(bootstrap.Options{
		Logger: bootstrap.WrapSlog(slog.Default()),
	})
	if err != nil {
		return err
	}

	svc := container.ProjectService

	// NOTE: signal the waiting pm process to kill session shell
	if project, ok := os.LookupEnv("PM_ACTIVE_PROJECT"); ok {
		fmt.Fprintf(cmd.OutOrStderr(), "Already running pm, active project: %s\n", project)

		return nil
	}

	name := ""

	if len(args) > 0 {
		name = args[0]
	}

	env := ""

	if len(args) > 1 {
		env = args[1]
	}

	path := ""

	if len(args) > 2 {
		path = args[2]
	}

	// Check if env is a path
	if info, _ := os.Stat(env); info != nil && info.IsDir() {
		path = env
	}

	var project *domain.Project

	// Launch TUI if no args or project not exist
	if len(args) == 0 || name == "" {
		themeFlag, _ := cmd.PersistentFlags().GetString("theme")
		configFlag, _ := cmd.PersistentFlags().GetString("config")

		cfgPath := strings.TrimSpace(configFlag)
		if cfgPath == "" {
			if v, ok := os.LookupEnv("PM_CONFIG"); ok && strings.TrimSpace(v) != "" {
				cfgPath = v
			}
		}

		cfg, cfgErr := configs.GetConfig(cfgPath)
		if cfgPath != "" && cfgErr != nil {
			return cfgErr
		}

		theme := strings.TrimSpace(themeFlag)
		if theme == "" {
			if v, ok := os.LookupEnv("PM_THEME"); ok && strings.TrimSpace(v) != "" {
				theme = v
			} else if cfg != nil && strings.TrimSpace(cfg.Theme) != "" {
				theme = cfg.Theme
			}
		}

		ov := map[string]string{}

		if cfg != nil {
			for _, ct := range cfg.CustomThemes {
				if strings.EqualFold(ct.Name, theme) {
					if ct.Title != "" {
						ov["title"] = ct.Title
					}

					if ct.Section != "" {
						ov["section"] = ct.Section
					}

					if ct.Subtext != "" {
						ov["subtext"] = ct.Subtext
					}

					if ct.Text != "" {
						ov["text"] = ct.Text
					}

					if ct.Placeholder != "" {
						ov["placeholder"] = ct.Placeholder
					}

					if ct.Border != "" {
						ov["border"] = ct.Border
					}

					if ct.Error != "" {
						ov["error"] = ct.Error
					}

					if ct.ButtonDefFg != "" {
						ov["buttonDefFg"] = ct.ButtonDefFg
					}

					if ct.ButtonDefBg != "" {
						ov["buttonDefBg"] = ct.ButtonDefBg
					}

					if ct.ButtonSelFg != "" {
						ov["buttonSelFg"] = ct.ButtonSelFg
					}

					if ct.ButtonSelBg != "" {
						ov["buttonSelBg"] = ct.ButtonSelBg
					}

					if ct.SelectedFg != "" {
						ov["selectedFg"] = ct.SelectedFg
					}

					if ct.SelectedBg != "" {
						ov["selectedBg"] = ct.SelectedBg
					}

					if ct.Help != "" {
						ov["help"] = ct.Help
					}
				}
			}
		}

		window, err := tui.NewWindow(svc, tui.Options{Theme: theme, Overrides: ov})
		if err != nil {
			return err
		}

		p := tea.NewProgram(window, tea.WithAltScreen())

		m, err := p.Run()
		if err != nil {
			return err
		}

		var (
			selected *domain.Project
			selEnv   string
		)

		if w, ok := m.(*tui.Window); ok {
			selected = w.SelectedProject()
			selEnv = w.SelectedEnvironment()
		}

		if selected == nil {
			return nil
		}

		project = selected
		if selEnv != "" {
			env = selEnv
		} else if selected.DefaultEnv != "" {
			env = selected.DefaultEnv
		}
	} else {
		project, err = svc.Load(name)
		if err != nil {
			return err
		}

		if env == "" && project.DefaultEnv != "" {
			env = project.DefaultEnv
		}
	}

	shellRepo, err := shells.NewPseudoShellRepository(project, env, path)
	if err != nil {
		return err
	}

	shellSvc := services.NewShellService(shellRepo)

	process, err := shellSvc.Start()
	if err != nil {
		return err
	}

	slog.Debug("new project shell started", slog.Int("pid", process.Pid))

	exitCode, err := shellSvc.Wait()
	if err != nil {
		slog.Error("cmd.Wait", slog.Any("err", err))
	}

	slog.Debug(
		"End project session",
		slog.Int("exit_code", exitCode),
		slog.String("shell", project.Name),
		slog.Int("PID", process.Pid),
	)

	return nil
}

func projectsList() []string {
	container, err := bootstrap.New(bootstrap.Options{
		Logger: bootstrap.WrapSlog(slog.Default()),
	})
	if err != nil {
		return nil
	}

	projects, err := container.ProjectService.List()
	if err != nil {
		return nil
	}

	list := make([]string, 0, len(projects))
	for _, project := range projects {
		list = append(list, project.Name)
	}

	return list
}
