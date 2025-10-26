package init

import (
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command constructs the `pm init` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "init",
		Short:        "Initialize your workspace",
		Long:         "Init",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         run,
	}

	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	cfg := viper.GetViper().ConfigFileUsed()
	if cfg != "" {
		slog.Debug("config file", slog.String("path", cfg), slog.String("dir", filepath.Dir(cfg)))
	}

	keys := viper.AllKeys()
	slog.Debug("config keys loaded", slog.Int("count", len(keys)))

	return nil
}
