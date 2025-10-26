package list

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/bootstrap"
)

const (
	flagFormat        = "format"
	flagGitUserName   = "user.name"
	flagGitUserEmail  = "user.email"
	flagGitSigningKey = "user.signingkey"
	flagCommitGPGSign = "commit.gpgsign"
	flagTagGPGSign    = "tag.gpgsign"
	flagEnvVars       = "env-vars"
	flagShell         = "shell"
)

// Command constructs the `pm list` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list <project>",
		Short:        "list projects",
		Long:         "List all projects",
		Args:         cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().String(flagFormat, "", "output format")
	cmd.Flags().String(flagGitUserName, "", "git user.name (default git --global)")
	cmd.Flags().String(flagGitUserEmail, "", "git user.email (default git --global)")
	cmd.Flags().String(flagGitSigningKey, "", "git user.signingkey (default git --global)")
	cmd.Flags().Bool(flagCommitGPGSign, true, "git commit.gpgsign (default git --global)")
	cmd.Flags().Bool(flagTagGPGSign, true, "git tag.gpgsign (default git --global)")
	cmd.Flags().StringToString(flagEnvVars, nil, "List of ENV_VARS to add to the environment")
	cmd.Flags().
		String(flagShell, "", "Shell of the project, need be installed in the system) (default to $SHELL env var))")

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
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
