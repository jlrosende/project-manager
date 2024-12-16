package init

import (
	"fmt"

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
	// TODO init pm config file and check requirements
	fmt.Println("TODO")
	return nil
}
