package edit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/cli/internal/commandutil"
	"github.com/jlrosende/project-manager/internal/bootstrap"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

type projectSkeletonService interface {
	RenderProjectEditSkeleton(context.Context, string, string, services.SkeletonFormat) ([]byte, error)
	GenerateProjectEditSkeleton(context.Context, string, string, string, services.SkeletonFormat) error
}

type projectEditService interface {
	projectSkeletonService
	ProjectEdit(context.Context, services.ProjectEditOptions) (*services.ProjectEditResult, error)
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	projectName, err := commandutil.RequireNonEmptyArg(args[0], "project name")
	if err != nil {
		return err
	}

	environmentName := ""
	if len(args) > 1 {
		environmentName, err = commandutil.RequireNonEmptyArg(args[1], "environment name")
		if err != nil {
			return err
		}
	}

	container, err := bootstrap.New(bootstrap.Options{Logger: bootstrap.WrapSlog(slog.Default())})
	if err != nil {
		return fmt.Errorf("bootstrap services: %w", err)
	}

	svc := container.ProjectService

	editSvc, ok := svc.(projectEditService)
	if !ok {
		return errors.New("project service does not support edit operations")
	}

	jsonSkeleton := cmd.Flags().Changed(flagGenerateSkeletonJSON)
	yamlSkeleton := cmd.Flags().Changed(flagGenerateSkeletonYAML)

	if jsonSkeleton && yamlSkeleton {
		return fmt.Errorf("cannot combine --%s with --%s", flagGenerateSkeletonJSON, flagGenerateSkeletonYAML)
	}

	if jsonSkeleton {
		destination, err := cmd.Flags().GetString(flagGenerateSkeletonJSON)
		if err != nil {
			return err
		}

		return outputProjectSkeleton(
			ctx,
			cmd,
			editSvc,
			projectName,
			environmentName,
			services.SkeletonFormatJSON,
			destination,
		)
	}

	if yamlSkeleton {
		destination, err := cmd.Flags().GetString(flagGenerateSkeletonYAML)
		if err != nil {
			return err
		}

		return outputProjectSkeleton(
			ctx,
			cmd,
			editSvc,
			projectName,
			environmentName,
			services.SkeletonFormatYAML,
			destination,
		)
	}

	envFlagUsed := cmd.Flags().Changed(flagEnvColor) || cmd.Flags().Changed(flagEnvEnvVarsMode) ||
		cmd.Flags().Changed(flagEnvEnvVarsFile)
	if envFlagUsed && environmentName == "" {
		return fmt.Errorf("environment flags require an environment argument")
	}

	cliInput, err := cmd.Flags().GetString(flagCLIInput)
	if err != nil {
		return err
	}

	allowUnknown, err := cmd.Flags().GetBool(flagAllowUnknown)
	if err != nil {
		return err
	}

	dryRun, err := cmd.Flags().GetBool(flagDryRun)
	if err != nil {
		return err
	}

	outputFormat, err := cmd.Flags().GetString(flagOutput)
	if err != nil {
		return err
	}

	outputFormat = strings.ToLower(strings.TrimSpace(outputFormat))
	if outputFormat == "" {
		outputFormat = "text"
	}

	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	description, err := cmd.Flags().GetString(flagProjectDescription)
	if err != nil {
		return err
	}

	shell, err := cmd.Flags().GetString(flagProjectShell)
	if err != nil {
		return err
	}

	projectEnvFile, err := cmd.Flags().GetString(flagProjectEnvVarsFile)
	if err != nil {
		return err
	}

	defaultEnv, err := cmd.Flags().GetString(flagProjectDefaultEnv)
	if err != nil {
		return err
	}

	envColor, err := cmd.Flags().GetString(flagEnvColor)
	if err != nil {
		return err
	}

	envMode, err := cmd.Flags().GetString(flagEnvEnvVarsMode)
	if err != nil {
		return err
	}

	envFile, err := cmd.Flags().GetString(flagEnvEnvVarsFile)
	if err != nil {
		return err
	}

	options := services.ProjectEditOptions{
		Name:            projectName,
		EnvironmentName: environmentName,
		AllowUnknown:    allowUnknown,
		CLIInputPath:    cliInput,
		DryRun:          dryRun,
		Flags: services.ProjectEditFlags{
			Description:    description,
			DescriptionSet: cmd.Flags().Changed(flagProjectDescription),
			Shell:          shell,
			ShellSet:       cmd.Flags().Changed(flagProjectShell),
			EnvVarsFile:    projectEnvFile,
			EnvVarsFileSet: cmd.Flags().Changed(flagProjectEnvVarsFile),
			DefaultEnv:     defaultEnv,
			DefaultEnvSet:  cmd.Flags().Changed(flagProjectDefaultEnv),
		},
		EnvFlags: services.EnvironmentEditFlags{
			Color:       envColor,
			ColorSet:    cmd.Flags().Changed(flagEnvColor),
			EnvVarsMode: envMode,
			ModeSet:     cmd.Flags().Changed(flagEnvEnvVarsMode),
			EnvVarsFile: envFile,
			EnvFileSet:  cmd.Flags().Changed(flagEnvEnvVarsFile),
		},
	}

	result, err := editSvc.ProjectEdit(ctx, options)
	if err != nil {
		if errors.Is(err, services.ErrProjectEditNoChanges) {
			fmt.Fprintln(cmd.OutOrStdout(), "No changes applied; nothing to update.")
			return nil
		}

		var validationErr *services.ProjectEditValidationError
		if errors.As(err, &validationErr) {
			return newExitError(2, formatValidationFailure(validationErr.Result()), err)
		}

		if errors.Is(err, services.ErrEnvironmentNotFound) {
			return newExitError(3, err.Error(), err)
		}

		if errors.Is(err, domain.ErrProjectLocked) {
			return newExitError(4, err.Error(), err)
		}

		return newExitError(1, err.Error(), err)
	}

	return renderProjectEditResult(cmd, result, outputFormat)
}

func outputProjectSkeleton(
	ctx context.Context,
	cmd *cobra.Command,
	svc projectSkeletonService,
	project string,
	environment string,
	format services.SkeletonFormat,
	destination string,
) error {
	destination = strings.TrimSpace(destination)
	if destination == "" || destination == "-" {
		data, err := svc.RenderProjectEditSkeleton(ctx, project, environment, format)
		if err != nil {
			return err
		}

		if _, err := cmd.OutOrStdout().Write(data); err != nil {
			return fmt.Errorf("write skeleton to stdout: %w", err)
		}

		return nil
	}

	if err := svc.GenerateProjectEditSkeleton(ctx, project, environment, destination, format); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "CLI input skeleton written to %s\n", destination)

	return nil
}

func renderProjectEditResult(cmd *cobra.Command, result *services.ProjectEditResult, format string) error {
	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("format json output: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(data))

		return nil
	default:
		out := cmd.OutOrStdout()
		changeCount := len(result.Changes)
		changeWord := "change"
		if changeCount != 1 {
			changeWord = "changes"
		}

		if result.DryRun {
			if result.EnvironmentName != "" {
				fmt.Fprintf(out, "Previewing environment edit for %q in project %q (dry-run)\n", result.EnvironmentName, result.ProjectName)
			} else {
				fmt.Fprintf(out, "Previewing project edit for %q (dry-run)\n", result.ProjectName)
			}

			fmt.Fprintln(out, "No changes have been written; showing preview only.")

			if changeCount == 0 {
				fmt.Fprintln(out, "  No changes detected; nothing would be modified.")

				if result.EnvironmentName != "" {
					fmt.Fprintln(out, "  Other environments unchanged; project metadata untouched.")
				}

				return nil
			}

			fmt.Fprintf(out, "Planned changes (%d %s):\n", changeCount, changeWord)
		} else {
			if result.EnvironmentName != "" {
				fmt.Fprintf(out, "Updated environment %q in project %q (%d %s):\n", result.EnvironmentName, result.ProjectName, changeCount, changeWord)
			} else {
				fmt.Fprintf(out, "Updated project %q (%d %s):\n", result.ProjectName, changeCount, changeWord)
			}

			if changeCount == 0 {
				fmt.Fprintln(out, "  No changes detected; nothing to update.")

				if result.EnvironmentName != "" {
					fmt.Fprintln(out, "  Other environments unchanged; project metadata untouched.")
				}

				return nil
			}
		}

		for _, change := range result.Changes {
			fmt.Fprintf(
				out,
				"  %s: %s -> %s\n",
				change.Label,
				quoteValue(change.Old),
				quoteValue(change.New),
			)
		}

		if result.EnvironmentName != "" {
			fmt.Fprintln(out, "  Other environments unchanged; project metadata untouched.")
		}

		return nil
	}
}

func quoteValue(value string) string {
	return fmt.Sprintf("%q", value)
}

type exitError struct {
	code    int
	message string
	cause   error
}

func newExitError(code int, message string, cause error) error {
	return &exitError{code: code, message: message, cause: cause}
}

func (e *exitError) Error() string {
	if e == nil {
		return ""
	}

	if strings.TrimSpace(e.message) != "" {
		return e.message
	}

	if e.cause != nil {
		return e.cause.Error()
	}

	return ""
}

func (e *exitError) ExitCode() int {
	if e == nil || e.code == 0 {
		return 1
	}

	return e.code
}

func (e *exitError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

func formatValidationFailure(result *services.ProjectEditResult) string {
	if result == nil {
		return "Validation failed; no changes applied."
	}

	scope := fmt.Sprintf("project %q", result.ProjectName)
	if strings.TrimSpace(result.EnvironmentName) != "" {
		scope = fmt.Sprintf("environment %q in project %q", result.EnvironmentName, result.ProjectName)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Validation failed; no changes applied to %s.", scope))

	if len(result.Errors) == 0 {
		return builder.String()
	}

	builder.WriteString("\n")

	for _, item := range result.Errors {
		label := strings.TrimSpace(item.Label)
		if label == "" {
			label = strings.TrimSpace(item.Field)
		}

		builder.WriteString("  - ")

		if label != "" {
			builder.WriteString(label)
		} else {
			builder.WriteString("field")
		}

		message := strings.TrimSpace(item.New)
		if message != "" {
			builder.WriteString(": ")
			builder.WriteString(message)
		}

		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}
