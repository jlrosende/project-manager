package exportcmd

import (
	"github.com/spf13/cobra"
)

const (
	flagDestinationPath = "destination"
	flagFormat          = "format"
	flagProjectName     = "project-name"
	flagIncludeEnv      = "include-env"
)

// Command constructs the `pm export` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "export <project>",
		Short:        "Export a project definition from pm",
		Long:         "Export a managed project definition to a file or stdout so it can be shared or backed up.",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().StringP(flagDestinationPath, "d", "", "Path to write the exported definition (default: stdout)")
	cmd.Flags().String(flagFormat, "", "Explicit output format (json, yaml); auto-detected from destination when omitted")
	cmd.Flags().String(flagProjectName, "", "Override the exported project name")
	cmd.Flags().Bool(flagIncludeEnv, false, "Include environment variable values in the export (default false)")

	return cmd
}
