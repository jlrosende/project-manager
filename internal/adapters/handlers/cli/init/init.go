package init

import (
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

func initCommand(_ *cobra.Command, _ []string) error {
	// NOTE: init pm config file and check requirements
	cfg := viper.GetViper().ConfigFileUsed()
	if cfg != "" {
		slog.Debug("config file", slog.String("path", cfg), slog.String("dir", filepath.Dir(cfg)))
	}

	keys := viper.AllKeys()
	slog.Debug("config keys loaded", slog.Int("count", len(keys)))

	// if err := viper.SafeWriteConfigAs(viper.GetViper().ConfigFileUsed()); err != nil {
	// 	return err
	// }

	return nil
}
