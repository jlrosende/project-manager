package list

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/bootstrap"
)

var ListCmd = &cobra.Command{
	Use:   "list <project>",
	Short: "list projects",
	Long:  `List all projects`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	RunE:  list,
}

func init() {
	ListCmd.Flags().String("format", "", "output format")

	ListCmd.Flags().String("user.name", "", "git user.name (default git --global)")
	ListCmd.Flags().String("user.email", "", "git user.email (default git --global)")
	ListCmd.Flags().String("user.signingkey", "", "git user.signingkey (default git --global)")

	ListCmd.Flags().Bool("commit.gpgsign", true, "git commit.gpgsign (default git --global)")
	ListCmd.Flags().Bool("tag.gpgsign", true, "git tag.gpgsign (default git --global)")

	ListCmd.Flags().StringToString("env-vars", nil, "List of ENV_VARS to add to the environment")

	ListCmd.Flags().
		String("shell", "", "Shell of the project, need be installed in the system) (default to $SHELL env var))")

	//	if err := ListCmd.MarkFlagRequired("name"); err != nil {
	//		log.Fatal(err)
	//	}
}

func list(cmd *cobra.Command, _ []string) error {
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
