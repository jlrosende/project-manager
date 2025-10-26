package importcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	// TODO: Wire bootstrap container and fetch required services for import operations.
	// TODO: Read source input (file or stdin) and parse according to the detected format.
	// TODO: Delegate parsed definition to a dedicated domain service for persistence.

	fmt.Fprintln(cmd.ErrOrStderr(), "pm import is not implemented yet")

	return nil
}
