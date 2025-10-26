package newcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/cli/internal/commandutil"
)

const (
	flagHere                 = "here"
	flagCliInput             = "cli-input"
	flagGenerateSkeletonJSON = "generate-cli-skeleton-json"
	flagGenerateSkeletonYAML = "generate-cli-skeleton-yaml"
	flagDryRun               = "dry-run"
	flagForce                = "force"
	flagAllowUnknown         = "allow-unknown"
	flagOutput               = "output"
	flagYes                  = "yes"

	flagProjectDescription = "project-description"
	flagProjectShell       = "project-shell"
	flagProjectEnvFile     = "project-env-file"

	flagEnvironmentEnvFile = "environment-env-file"
	flagEnvironmentMode    = "environment-mode"
	flagEnvironmentColor   = "environment-color"
	flagEnvironmentEnvVar  = "environment-env-var"

	flagEnvFileLegacy  = "env-file"
	flagEnvModeLegacy  = "env-mode"
	flagEnvColorLegacy = "env-color"
	flagEnvVarLegacy   = "env-var"
)

// Command constructs the `pm new` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <name> [path]",
		Short: "Create a new project from arguments or configuration",
		Long:  "Create or initialize a project directory, validating inputs from positional arguments, flags, and optional CLI input files.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := commandutil.EnforceExclusiveFlags(cmd, flagGenerateSkeletonJSON, flagGenerateSkeletonYAML); err != nil {
				return err
			}

			if cmd.Flags().Changed(flagGenerateSkeletonJSON) || cmd.Flags().Changed(flagGenerateSkeletonYAML) {
				if len(args) > 1 {
					return fmt.Errorf("only one skeleton output path may be provided")
				}

				return nil
			}

			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().Bool(flagHere, false, "Initialize the current directory instead of creating a new one")
	cmd.Flags().String(flagCliInput, "", "Path to JSON or YAML CLI input file")
	cmd.Flags().String(flagGenerateSkeletonJSON, "", "Write JSON CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().String(flagGenerateSkeletonYAML, "", "Write YAML CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().Bool(flagDryRun, false, "Preview actions without writing files")

	commandutil.SetDashDefault(cmd, flagGenerateSkeletonJSON, flagGenerateSkeletonYAML)
	cmd.Flags().Bool(flagForce, false, "Overwrite existing project files when rerun")
	cmd.Flags().Bool(flagAllowUnknown, false, "Ignore unknown fields in config files")
	cmd.Flags().String(flagOutput, "text", "Output format for dry runs (text or json)")

	cmd.Flags().String(flagProjectDescription, "", "Set project description metadata")
	cmd.Flags().String(flagProjectShell, "", "Set default shell for the project")
	cmd.Flags().String(flagProjectEnvFile, "", "Set default project-level environment vars file")

	cmd.Flags().String(flagEnvironmentEnvFile, "", "Set environment vars file when adding environments")
	cmd.Flags().String(flagEnvironmentMode, "", "Set environment vars merge mode (merge or replace)")
	cmd.Flags().String(flagEnvironmentColor, "", "Set environment color metadata when adding environments")
	cmd.Flags().StringArray(flagEnvironmentEnvVar, nil, "Environment variable in KEY=VALUE format (repeatable)")

	cmd.Flags().String(flagEnvFileLegacy, "", "[deprecated] Use --environment-env-file instead")
	cmd.Flags().String(flagEnvModeLegacy, "", "[deprecated] Use --environment-mode instead")
	cmd.Flags().String(flagEnvColorLegacy, "", "[deprecated] Use --environment-color instead")
	cmd.Flags().StringArray(flagEnvVarLegacy, nil, "[deprecated] Use --environment-env-var instead")

	_ = cmd.Flags().MarkDeprecated(flagEnvFileLegacy, "use --environment-env-file instead")
	_ = cmd.Flags().MarkDeprecated(flagEnvModeLegacy, "use --environment-mode instead")
	_ = cmd.Flags().MarkDeprecated(flagEnvColorLegacy, "use --environment-color instead")
	_ = cmd.Flags().MarkDeprecated(flagEnvVarLegacy, "use --environment-env-var instead")

	cmd.Flags().Bool(flagYes, false, "Automatically confirm project and environment creation prompts")

	return cmd
}
