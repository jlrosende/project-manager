package exportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	project := args[0]

	// TODO: Resolve bootstrap container and acquire project service for export operations.
	// TODO: Build export options including destination, format, and environment inclusion preferences.
	// TODO: Fetch project definition and write it to the selected output target.

	fmt.Fprintf(cmd.ErrOrStderr(), "pm export for project %s is not implemented yet\n", project)

	return nil
}
