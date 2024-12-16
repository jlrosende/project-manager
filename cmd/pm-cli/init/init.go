package init

import (
	"fmt"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a your workspace",
	Long:  `Init`,
	Args:  cobra.NoArgs,
	RunE:  initCommand,
}

func init() {

}

func initCommand(cmd *cobra.Command, args []string) error {
	globalConfig, err := config.LoadConfig(config.GlobalScope)

	if err != nil {
		return err
	}

	newConfig := config.NewConfig()

	newConfig.User.Email = "test-3"
	newConfig.User.Name = "test"

	gitcong, err := newConfig.Marshal()

	if err != nil {
		return err
	}

	fmt.Printf("%+s\n", gitcong)

	globalConfig.Raw.AddOption("includeIf", "gitdir:/workspaces/project-manager/temp/test-3/", "path", "/workspaces/project-manager/temp/test-1/.gitconfig-test-3")

	section := globalConfig.Raw.Section("includeIf")

	for _, sub := range section.Subsections {
		fmt.Printf("%+v\n", sub.Name)
		fmt.Printf("%+v\n", sub.Option("path"))
	}

	gitcong, err = globalConfig.Marshal()

	if err != nil {
		return err
	}

	fmt.Printf("%+s\n", gitcong)

	return nil
}
