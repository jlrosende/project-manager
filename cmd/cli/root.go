package cli

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	initialize "github.com/jlrosende/project-manager/cmd/cli/init"
	"github.com/jlrosende/project-manager/cmd/cli/new"
	"github.com/jlrosende/project-manager/configs"
	"github.com/jlrosende/project-manager/internal"
	"github.com/jlrosende/project-manager/pkg/ui/styles/list"
	"github.com/spf13/cobra"
)

var (
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

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/pm/config.hcl)")

	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(new.NewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func root(cmd *cobra.Command, args []string) error {

	cfgFile, err := cmd.PersistentFlags().GetString("config")

	if err != nil {
		return err
	}

	config, err := configs.GetConfig(cfgFile)

	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%+v\n", config))

	items := []list.Item{}

	for project := range config.Projects {
		slog.Info(project)
		items = append(items, list.Item{
			Name: project,
			Desc: config.Projects[project].Path,
		})
	}

	m := list.NewList("Projects", items)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}
