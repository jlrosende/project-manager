package cli

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

	log.Println(os.Getpid(), os.Getppid())

	name := ""

	if len(args) > 0 {
		name = args[0]
	}

	repo, err := repositories.NewProjectRepository()

	if err != nil {
		return err
	}

	svc := services.NewProjectService(repo)

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
		project.Path = args[1]
	}

	log.Printf("Start project %s shell in %s", project.Name, project.Path)

	shell := exec.Command(os.Getenv("SHELL"))
	shell.Dir = filepath.Dir(project.Path)
	shell.Env = append(os.Environ(), project.EnvVarsSlice()...)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	shell.Start()
	log.Println(shell.Process.Pid)

	if err := shell.Wait(); err != nil {
		log.Println(err)
	}

	log.Println(shell.Process.Pid)

	log.Printf("End project session %s shell", project.Name)

	return nil
}
