package new

import (
	"log"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
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

	// if name == "" {
	// 	// Ask for name
	// 	nameInput := textinput.NewTextInput("Whats the name of the new project?", "Project name")
	// 	p := tea.NewProgram(nameInput)
	// 	if _, err := p.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	name = nameInput.Value()
	// }

	repo, err := repositories.NewProjectRepository()

	if err != nil {
		return err
	}

	svc := services.NewProjectService(repo)

	projects, _ := svc.List()

	if err != nil {
		return err
	}

	log.Println(projects)

	_, _ = svc.Create(name, path, subproject, envVars)

	if err != nil {
		return err
	}

	return nil
}
