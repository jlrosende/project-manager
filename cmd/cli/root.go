package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cmdInit "github.com/jlrosende/project-manager/cmd/cli/init"
	cmdNew "github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/internal"
	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "pm [project] [path]",
		Short:        "pm is a tool to create and organize projects in your computer",
		Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
		Version:      internal.GetVersion(),
		Args:         cobra.MaximumNArgs(2),
		SilenceUsage: true,
		RunE:         root,
	}
)

func init() {

	rootCmd.Flags().BoolP("list", "l", false, "List all the projects.")

	rootCmd.AddCommand(cmdInit.InitCmd)
	rootCmd.AddCommand(cmdNew.NewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
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
			fmt.Fprintf(cmd.OutOrStderr(), "Description: %s\n", p.Name)
			fmt.Fprintln(cmd.OutOrStderr(), "---")
		}
		return nil
	}

	log.Printf("PID: %d PPID: %d PM_ACTIVE_PROJECT: %s", os.Getpid(), os.Getppid(), os.Getenv("PM_ACTIVE_PROJECT"))

	if os.Getenv("PM_ACTIVE_PROJECT") != "" {

		log.Println("Interrupt parent process")

		p, err := os.FindProcess(os.Getppid())
		if err != nil {
			return err
		}

		if err := p.Kill(); err != nil {
			return err
		}
		return nil
	}

	name := ""

	if len(args) > 0 {
		name = args[0]
	}

	project, err := svc.Get(name)

	if err != nil {
		return err
	}

	// Launch TUI if project is empty
	if project == nil {
		log.Println("No project defined")
		return nil
	}

	if len(args) > 1 {
		// project.Path = args[1]
	}

	// log.Printf("Start project %s shell in %s", project.Name, project.Path)

	shell := exec.Command(os.Getenv("SHELL"))

	// shell.Dir = filepath.Dir(project.Path)

	shell.Env = append(
		os.Environ(),
		// project.EnvVars.ToSlice()...,
	)

	shell.Env = append(
		shell.Env,
		fmt.Sprintf("PM_ACTIVE_PROJECT=%s", project.Name),
	)

	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	if err := shell.Start(); err != nil {
		return err
	}

	log.Printf("New SHELL PID: %d\n", shell.Process.Pid)

	if err := shell.Wait(); err != nil {
		log.Println(err)
	}

	log.Printf("End project session %s shell, (PID: %d)", project.Name, shell.Process.Pid)

	return nil
}
