package importcmd

import (
	"github.com/spf13/cobra"
)

const (
	flagSourcePath  = "source"
	flagFormat      = "format"
	flagProjectName = "project-name"
)

// Command constructs the `pm import` Cobra command following the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "import",
		Short:        "Import a project definition into pm",
		Long:         "Import a project definition from a file or stdin and register it with pm for local use.",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().StringP(flagSourcePath, "s", "", "Path to the project definition file (default: stdin)")
	cmd.Flags().String(flagFormat, "", "Explicit input format (json, yaml); auto-detected when omitted")
	cmd.Flags().String(flagProjectName, "", "Override the project name when importing")

	return cmd
}
