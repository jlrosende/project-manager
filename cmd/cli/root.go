package cli

import (
	"fmt"
	"os"

	initialize "github.com/jlrosende/project-manager/cmd/cli/init"
	"github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "pm",
	Short:        "pm is a tool to create and organize projects in your computer",
	Long:         `A tool to manage the configuration and estructure of multiple projects inside your computer`,
	Version:      internal.GetVersion(),
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(new.NewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
