package edit

import (
	"fmt"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	EditCmd = &cobra.Command{
		Use:   "edit <project>",
		Short: "edit project",
		Long:  `Edit all projects`,
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:  edit,
	}
)

func init() {

}

func edit(cmd *cobra.Command, args []string) error {

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

	if err != nil {
		return err
	}

	projects, err := svc.List()
	if err != nil {
		return err
	}
	for i, p := range projects {

		fmt.Fprintln(cmd.OutOrStderr(), "---")
		fmt.Fprintf(cmd.OutOrStderr(), "Name: %s\n", p.Name)
		fmt.Fprintf(cmd.OutOrStderr(), "Description: %s\n", p.Description)
		if len(projects)-1 == i {
			fmt.Fprintln(cmd.OutOrStderr(), "---")
		}
	}
	return nil
}
