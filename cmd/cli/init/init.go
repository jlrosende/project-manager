package init

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	slog.Debug(viper.GetViper().ConfigFileUsed())
	slog.Debug(filepath.Dir(viper.GetViper().ConfigFileUsed()))
	slog.Debug(fmt.Sprintf("%+v", viper.AllSettings()))

	// if err := viper.SafeWriteConfigAs(viper.GetViper().ConfigFileUsed()); err != nil {
	// 	return err
	// }

	return nil
}
