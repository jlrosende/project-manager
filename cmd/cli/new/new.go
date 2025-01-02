package new

import (
	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	NewCmd = &cobra.Command{
		Use:   "new <project> [path]",
		Short: "Create new projects",
		Long:  `Create a new project and all the basic configuration files`,
		Args:  cobra.MatchAll(cobra.RangeArgs(1, 2), cobra.OnlyValidArgs),
		RunE:  new,
	}
)

func init() {
	NewCmd.Flags().String("subproject", "", "Set this new project as subproject")

	NewCmd.Flags().String("user.name", "", "git user.name (default git --global)")
	NewCmd.Flags().String("user.email", "", "git user.email (default git --global)")
	NewCmd.Flags().String("user.signingkey", "", "git user.signingkey (default git --global)")

	NewCmd.Flags().Bool("commit.gpgsign", true, "git commit.gpgsign (default git --global)")
	NewCmd.Flags().Bool("tag.gpgsign", true, "git tag.gpgsign (default git --global)")

	NewCmd.Flags().StringToString("env-vars", nil, "List of ENV_VARS to add to the environment")

	// if err := NewCmd.MarkFlagRequired("name"); err != nil {
	// 	log.Fatal(err)
	// }

}

func new(cmd *cobra.Command, args []string) error {

	// ask for a name if not set
	name := args[0]
	path := ""

	subproject, err := cmd.Flags().GetString("subproject")

	if err != nil {
		return err
	}

	envVars, err := cmd.Flags().GetStringToString("env-vars")

	if err != nil {
		return err
	}

	if len(args) > 1 {
		path = args[1]
	}

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

	gitUserName, err := cmd.Flags().GetString("user.name")

	if err != nil {
		return err
	}

	gitUserEmail, err := cmd.Flags().GetString("user.email")

	if err != nil {
		return err
	}

	gitUserSigningKey, err := cmd.Flags().GetString("user.signingkey")

	if err != nil {
		return err
	}

	gitCommitGPGsign, err := cmd.Flags().GetBool("commit.gpgsign")

	if err != nil {
		return err
	}

	_, err = svc.Create(
		name,
		path,
		subproject,
		envVars,
		domain.New(
			domain.WithName(gitUserName),
			domain.WithEmail(gitUserEmail),
			domain.WithSigningKey(gitUserSigningKey),
			domain.WithCommitSign(gitCommitGPGsign),
			domain.WithTagSign(gitCommitGPGsign),
		),
	)

	if err != nil {
		return err
	}

	return nil
}
