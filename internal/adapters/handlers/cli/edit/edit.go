package edit

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/bootstrap"
)

var EditCmd = &cobra.Command{
	Use:   "edit <project>",
	Short: "edit project",
	Long:  `Edit all projects`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:  edit,
}

func init() {
}

func edit(cmd *cobra.Command, _ []string) error {
	container, err := bootstrap.New(bootstrap.Options{
		Logger: bootstrap.WrapSlog(slog.Default()),
	})
	if err != nil {
		return err
	}

	projects, err := container.ProjectService.List()
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
