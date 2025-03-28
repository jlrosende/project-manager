package cli

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	cmdInit "github.com/jlrosende/project-manager/cmd/cli/init"
	cmdNew "github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/internal"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/adapters/repositories/shells"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "pm [project] [env|[path]] ",
		Short:        "pm is a tool to create and organize projects in your computer",
		Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
		Version:      internal.GetVersion(),
		Args:         cobra.MaximumNArgs(3),
		SilenceUsage: true,
		RunE:         root,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	logLevel, err := cmd.PersistentFlags().GetString("log-level")
		// 	if err != nil {
		// 		return err
		// 	}

		// 	level := slog.LevelVar{}
		// 	err = level.UnmarshalText([]byte(logLevel))

		// 	if err != nil {
		// 		return err
		// 	}

		// 	cache, err := os.UserCacheDir()

		// 	if err != nil {
		// 		return err
		// 	}

		// 	fp, err := os.OpenFile(filepath.Join(cache, "pm.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		// 	if err != nil {
		// 		return err
		// 	}

		// 	logger := slog.New(slog.NewTextHandler(fp, &slog.HandlerOptions{
		// 		Level: level.Level(),
		// 	}))

		// 	// logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		// 	// 	Level: level.Level(),
		// 	// }))

		// 	logger = logger.With(
		// 		slog.Group("ps",
		// 			slog.Int("pid", os.Getpid()),
		// 			slog.Int("ppid", os.Getppid()),
		// 			slog.String("project", os.Getenv("PM_ACTIVE_PROJECT")),
		// 		),
		// 	)

		// 	slog.SetDefault(logger)

		// 	slog.Info("----------------------------------------------------------------------")

		// 	return nil
		// },
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) >= 1 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return projectsList(), cobra.ShellCompDirectiveNoFileComp
		},
	}
)

func init() {

	rootCmd.Flags().BoolP("list", "l", false, "List all the projects.")

	// rootCmd.PersistentFlags().String("log-level", "info", "Change the log level (debug, info, warn, error)")

	rootCmd.AddCommand(cmdInit.InitCmd)
	rootCmd.AddCommand(cmdNew.NewCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("somethin wrong happend", slog.Any("err", err))
		os.Exit(1)
	}
}

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

	if list, err := cmd.Flags().GetBool("list"); err != nil {
		return err
	} else if list {
		fmt.Fprintln(cmd.OutOrStderr(), "Print list and return")
		projects, err := svc.List()
		if err != nil {
			return err
		}
		for _, p := range projects {

			fmt.Fprintln(cmd.OutOrStderr(), "---")
			fmt.Fprintf(cmd.OutOrStderr(), "Name: %s\n", p.Name)
			fmt.Fprintf(cmd.OutOrStderr(), "Description: %s\n", p.Description)
			fmt.Fprintln(cmd.OutOrStderr(), "---")
		}
		return nil
	}

	// TODO Signal the wating pm process to kill session shell
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
		// window, err := tui.NewWindow(svc)
		window, err := tui.NewSimpleWindow(svc)
		if err != nil {
			return err
		}
		p := tea.NewProgram(window, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}

		selected := window.SelectedProject()

		if selected == nil {
			return nil
		}
		// Get selected project
		project = selected
	} else {
		project, err = svc.Load(name)

		if err != nil {
			return err
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
