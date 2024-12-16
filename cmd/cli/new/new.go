package new

import (
	"fmt"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new projects and your workspace",
	Long:  `new`,
	Args:  cobra.NoArgs,
	RunE:  newCommand,
}

func init() {

}

func newCommand(cmd *cobra.Command, args []string) error {
	globalConfig, err := config.LoadConfig(config.GlobalScope)

	if err != nil {
		return err
	}

	newConfig := config.NewConfig()

	newConfig.User.Email = "test-3"
	newConfig.User.Name = "test"

	newConfig.Raw.AddOption("commit", "", "gpgsign", "true")

	newGitConf, err := newConfig.Marshal()

	if err != nil {
		return err
	}

	fmt.Printf("%+s\n", newGitConf)
	fmt.Println("-----------------------------")

	globalConfig.Raw.AddOption("includeIf", "gitdir:/workspaces/project-manager/temp/test-2/", "path", "/workspaces/project-manager/temp/test-2/.gitconfig-test-2")

	gitconf, err := globalConfig.Marshal()

	if err != nil {
		return err
	}

	fmt.Printf("%+s\n", gitconf)

	return nil
}
