package edit

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/cli/internal/commandutil"
)

const (
	flagCLIInput             = "cli-input"
	flagGenerateSkeletonJSON = "generate-cli-skeleton-json"
	flagGenerateSkeletonYAML = "generate-cli-skeleton-yaml"
	flagAllowUnknown         = "allow-unknown"
	flagDryRun               = "dry-run"
	flagOutput               = "output"
	flagProjectDescription   = "project-description"
	flagProjectShell         = "project-shell"
	flagProjectEnvVarsFile   = "project-env-vars-file"
	flagProjectDefaultEnv    = "project-default-env"
	flagEnvColor             = "env-color"
	flagEnvEnvVarsMode       = "env-env-vars-mode"
	flagEnvEnvVarsFile       = "env-env-vars-file"
)

// Command constructs the `pm edit` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "edit <project> [env]",
		Short:         "Edit project configuration",
		SilenceUsage:  true,
		SilenceErrors: false,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := commandutil.EnforceExclusiveFlags(cmd, flagGenerateSkeletonJSON, flagGenerateSkeletonYAML); err != nil {
				return err
			}

			if len(args) == 0 {
				return fmt.Errorf("project name is required")
			}

			if _, err := commandutil.RequireNonEmptyArg(args[0], "project name"); err != nil {
				return err
			}

			if len(args) > 2 {
				return cobra.RangeArgs(1, 2)(cmd, args)
			}

			if len(args) == 2 {
				if _, err := commandutil.RequireNonEmptyArg(args[1], "environment name"); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: run,
	}

	cmd.Flags().String(flagCLIInput, "", "Path to JSON or YAML CLI input file")
	cmd.Flags().String(flagGenerateSkeletonJSON, "", "Write JSON CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().String(flagGenerateSkeletonYAML, "", "Write YAML CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().Bool(flagAllowUnknown, false, "Ignore unknown fields in CLI input files")
	cmd.Flags().Bool(flagDryRun, false, "Preview changes without persisting them")
	cmd.Flags().String(flagOutput, "text", "Output format for results (text or json)")

	commandutil.SetDashDefault(cmd, flagGenerateSkeletonJSON, flagGenerateSkeletonYAML)

	cmd.Flags().String(flagProjectDescription, "", "Set the project description")
	cmd.Flags().String(flagProjectShell, "", "Set the default project shell")
	cmd.Flags().String(flagProjectEnvVarsFile, "", "Set the project-level env vars file path")
	cmd.Flags().String(flagProjectDefaultEnv, "", "Set the default environment")

	cmd.Flags().String(flagEnvColor, "", "Set the environment color (requires <env>)")
	cmd.Flags().String(flagEnvEnvVarsMode, "", "Set the environment env vars mode (merge or replace, requires <env>)")
	cmd.Flags().String(flagEnvEnvVarsFile, "", "Set the environment env vars file path (requires <env>)")

	return cmd
}
