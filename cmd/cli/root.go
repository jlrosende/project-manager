package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	initialize "github.com/jlrosende/project-manager/cmd/cli/init"
	"github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/configs"
	"github.com/jlrosende/project-manager/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:          "pm",
		Short:        "pm is a tool to create and organize projects in your computer",
		Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
		Version:      internal.GetVersion(),
		SilenceUsage: true,
		RunE:         root,
	}
)

func init() {

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pm/config)")

	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(new.NewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).

		viper.AddConfigPath(path.Join(home, ".config/pm/"))
		viper.SetConfigName("config")
		viper.SetConfigType("hcl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func root(cmd *cobra.Command, args []string) error {

	config, err := configs.GetConfig(cfgFile)

	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%+v\n", config))

	return nil
}
