package cli

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	cmdEdit "github.com/jlrosende/project-manager/cmd/cli/edit"
	cmdInit "github.com/jlrosende/project-manager/cmd/cli/init"
	cmdList "github.com/jlrosende/project-manager/cmd/cli/list"
	cmdNew "github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/internal"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/adapters/repositories/shells"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/internal/logger"
)

var rootCmd = &cobra.Command{
	Use:          "pm [project] [env|[path]] ",
	Short:        "pm is a tool to create and organize projects in your computer",
	Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
	Version:      internal.GetVersion(),
	Args:         cobra.MaximumNArgs(3),
	SilenceUsage: true,
	RunE:         root,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		logLevel, err := cmd.PersistentFlags().GetString("log-level")
		if err != nil {
			return err
		}
		logFile, err := cmd.PersistentFlags().GetString("log-file")
		if err != nil {
			return err
		}

		_, err = logger.Setup(logLevel, logFile)
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

func init() {
	rootCmd.Flags().BoolP("list", "l", false, "List all the projects.")

	rootCmd.PersistentFlags().String("log-level", "info", "Change the log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-file", "", "Path to log file (default: $XDG_CACHE_HOME/pm.log)")

	rootCmd.AddCommand(cmdInit.InitCmd)
	rootCmd.AddCommand(cmdNew.NewCmd)
	rootCmd.AddCommand(cmdList.ListCmd)
	rootCmd.AddCommand(cmdEdit.EditCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("somethin wrong happened", slog.Any("err", err))
		os.Exit(1)
	}
}

// Root exposes the root command for tools like doc generators.
func Root() *cobra.Command { return rootCmd }

func root(cmd *cobra.Command, args []string) error {
	repoProject, err := repositories.NewProjectRepository()
	if err != nil {
		return err
	}

	repoEnvVars, err := repositories.NewEnvVarsRepository()
	if err != nil {
		return err
	}

	repoGitConfig, err := repositories.NewGitRepository()
	if err != nil {
		return err
	}

	svc := services.NewProjectService(repoProject, repoEnvVars, repoGitConfig)

	// NOTE: signal the waiting pm process to kill session shell
	if projet, ok := os.LookupEnv("PM_ACTIVE_PROJECT"); ok {
		fmt.Fprintf(cmd.OutOrStderr(), "Already runnign pm, active project: %s\n", projet)

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

	// Launch TUI if no args or project not exsit
	if len(args) == 0 || name == "" {
		window, err := tui.NewWindow(svc)
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
		log.Printf("cmd.Wait: %v", err)
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
	repoProject, _ := repositories.NewProjectRepository()
	repoEnvVars, _ := repositories.NewEnvVarsRepository()
	repoGitConfig, _ := repositories.NewGitRepository()
	svc := services.NewProjectService(repoProject, repoEnvVars, repoGitConfig)

	projects, _ := svc.List()

	list := make([]string, len(projects))
	for _, project := range projects {
		list = append(list, project.Name)
	}

	return list
}
