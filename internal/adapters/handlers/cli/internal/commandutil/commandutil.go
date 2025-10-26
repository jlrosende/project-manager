package commandutil

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// EnforceExclusiveFlags returns an error when both flags are supplied.
func EnforceExclusiveFlags(cmd *cobra.Command, flagA, flagB string) error {
	flags := cmd.Flags()
	if flags.Changed(flagA) && flags.Changed(flagB) {
		return fmt.Errorf("cannot combine --%s with --%s", flagA, flagB)
	}

	return nil
}

// RequireNonEmptyArg trims the value and ensures it is not empty.
func RequireNonEmptyArg(value, label string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", label)
	}

	return trimmed, nil
}

// SetDashDefault configures NoOptDefVal to "-" for the provided string flags.
func SetDashDefault(cmd *cobra.Command, flags ...string) {
	for _, name := range flags {
		if flag := cmd.Flags().Lookup(name); flag != nil {
			flag.NoOptDefVal = "-"
		}
	}
}
