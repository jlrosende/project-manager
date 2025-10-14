package newcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/internal/bootstrap"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
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
)

// Command constructs a fresh instance of the `pm new` Cobra command.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <name> [path]",
		Short: "Create a new project from arguments or configuration",
		Long:  "Create or initialize a project directory, validating inputs from positional arguments, flags, and optional CLI input files.",
		Args: func(cmd *cobra.Command, args []string) error {
			skeletonJSON := cmd.Flags().Changed(flagGenerateSkeletonJSON)
			skeletonYAML := cmd.Flags().Changed(flagGenerateSkeletonYAML)

			if skeletonJSON || skeletonYAML {
				if len(args) > 1 {
					return fmt.Errorf("only one skeleton output path may be provided")
				}

				return nil
			}

			return cobra.RangeArgs(1, 2)(cmd, args)
		},

		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	cmd.Flags().Bool(flagHere, false, "Initialize the current directory instead of creating a new one")
	cmd.Flags().String(flagCliInput, "", "Path to JSON or YAML CLI input file")
	cmd.Flags().String(flagGenerateSkeletonJSON, "", "Write JSON CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().String(flagGenerateSkeletonYAML, "", "Write YAML CLI input skeleton to path (stdout if omitted)")
	cmd.Flags().Bool(flagDryRun, false, "Preview actions without writing files")

	cmd.Flags().Lookup(flagGenerateSkeletonJSON).NoOptDefVal = "-"
	cmd.Flags().Lookup(flagGenerateSkeletonYAML).NoOptDefVal = "-"
	cmd.Flags().Bool(flagForce, false, "Overwrite existing project files when rerun")
	cmd.Flags().Bool(flagAllowUnknown, false, "Ignore unknown fields in config files")
	cmd.Flags().String(flagOutput, "text", "Output format for dry runs (text or json)")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	container, err := bootstrap.New(bootstrap.Options{Logger: bootstrap.WrapSlog(slog.Default())})
	if err != nil {
		return fmt.Errorf("bootstrap services: %w", err)
	}

	allowUnknown, err := cmd.Flags().GetBool(flagAllowUnknown)
	if err != nil {
		return err
	}

	jsonSkeletonPath, err := cmd.Flags().GetString(flagGenerateSkeletonJSON)
	if err != nil {
		return err
	}

	yamlSkeletonPath, err := cmd.Flags().GetString(flagGenerateSkeletonYAML)
	if err != nil {
		return err
	}

	jsonSkeletonSet := cmd.Flags().Changed(flagGenerateSkeletonJSON)
	yamlSkeletonSet := cmd.Flags().Changed(flagGenerateSkeletonYAML)

	if jsonSkeletonSet && yamlSkeletonSet {
		return fmt.Errorf("cannot combine --%s with --%s", flagGenerateSkeletonJSON, flagGenerateSkeletonYAML)
	}

	if jsonSkeletonSet {
		path := jsonSkeletonPath
		if len(args) > 0 {
			path = args[0]

			if len(args) > 1 {
				return fmt.Errorf("positional arguments are not allowed when generating CLI input skeletons")
			}
		}

		return outputSkeleton(cmd, path, services.SkeletonFormatJSON)
	}

	if yamlSkeletonSet {
		path := yamlSkeletonPath
		if len(args) > 0 {
			path = args[0]

			if len(args) > 1 {
				return fmt.Errorf("positional arguments are not allowed when generating CLI input skeletons")
			}
		}

		return outputSkeleton(cmd, path, services.SkeletonFormatYAML)
	}

	cliInputPath, err := cmd.Flags().GetString(flagCliInput)
	if err != nil {
		return err
	}

	var cfg *domain.ConfigInput
	if strings.TrimSpace(cliInputPath) != "" {
		cfg, err = services.LoadProjectConfig(cliInputPath)
		if err != nil {
			return err
		}

		if unknown := cfg.UnknownFields(); len(unknown) > 0 && !allowUnknown {
			keys := make([]string, 0, len(unknown))
			for k := range unknown {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			return fmt.Errorf(
				"CLI input contains unknown fields: %s (use --allow-unknown to ignore)",
				strings.Join(keys, ", "),
			)
		}
	}

	nameArg := strings.TrimSpace(args[0])

	probe, err := container.ProjectService.Probe(nameArg)
	if err != nil {
		return fmt.Errorf("probe project existence: %w", err)
	}

	existingProject := probe.RegistryHit && probe.ProjectFileExists
	pathArg := ""
	pathProvided := false
	envName := ""

	if len(args) > 1 {
		second := strings.TrimSpace(args[1])
		if existingProject {
			envName = second
		} else {
			pathArg = second
			pathProvided = true
		}
	}

	here, err := cmd.Flags().GetBool(flagHere)
	if err != nil {
		return err
	}

	if here && pathProvided {
		return fmt.Errorf("--here cannot be used when a path argument is provided")
	}

	merged, envVars, err := services.MergeProjectInputs(cfg, services.ProjectConfigFlags{
		Name:    nameArg,
		Path:    pathArg,
		Here:    here,
		PathSet: pathProvided,
		HereSet: cmd.Flags().Changed(flagHere),
	})
	if err != nil {
		return err
	}

	if envName != "" {
		nameCopy := envName
		if merged.Environment == nil {
			merged.Environment = &domain.EnvironmentInput{}
		}
		merged.Environment.Name = &nameCopy
	}

	if existingProject {
		merged.Path = strings.TrimSpace(probe.ProjectPath)
		merged.Here = false
	}

	if strings.TrimSpace(merged.Path) == "" && !merged.Here {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return fmt.Errorf("resolve current directory: %w", cwdErr)
		}

		merged.Path = filepath.Join(cwd, merged.Name)
	}

	if err := domain.ValidateProjectDefinition(&merged); err != nil {
		return err
	}

	if merged.Here {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return fmt.Errorf("resolve current directory: %w", cwdErr)
		}

		merged.Path = cwd
	}

	if existingProject && (merged.Environment == nil || merged.Environment.Name == nil || strings.TrimSpace(*merged.Environment.Name) == "") {
		fmt.Fprintf(cmd.OutOrStdout(), "project %s already exists; nothing to do\n", merged.Name)
		return nil
	}

	if !existingProject {
		existingProjects, listErr := container.ProjectService.List()
		if listErr != nil {
			return fmt.Errorf("list projects: %w", listErr)
		}

		if err := services.EnsureProjectUniqueness(existingProjects, merged); err != nil {
			return err
		}
	}

	force, err := cmd.Flags().GetBool(flagForce)
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

	if dryRun {
		return renderDryRun(cmd, merged, envVars, outputFormat)
	}

	creator := services.NewFileService(container.ProjectService, container.Filesystem)
	if creator == nil {
		return errors.New("project file service not available")
	}

	_, err = creator.Create(merged, envVars, services.CreateOptions{Force: force})

	return err
}

func renderDryRun(cmd *cobra.Command, def domain.ProjectDefinition, envVars domain.EnvVars, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		payload := map[string]any{
			"name": def.Name,
			"path": def.Path,
			"here": def.Here,
		}

		if def.Environment != nil {
			payload["environment"] = def.Environment
		}

		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("format dry-run output: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(data))
	default:
		fmt.Fprintf(cmd.OutOrStdout(), "pm new (dry-run)\n  name: %s\n  path: %s\n", def.Name, def.Path)

		if len(envVars) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  env vars: %d entries\n", len(envVars))
		}
	}

	return nil
}

func outputSkeleton(cmd *cobra.Command, rawPath string, format services.SkeletonFormat) error {
	path := strings.TrimSpace(rawPath)

	data, err := services.RenderProjectSkeleton(format)
	if err != nil {
		return err
	}

	if path == "" || path == "-" {
		if _, err := cmd.OutOrStdout().Write(data); err != nil {
			return fmt.Errorf("write skeleton to stdout: %w", err)
		}

		return nil
	}

	if err := services.GenerateProjectSkeleton(path, format); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "CLI input skeleton written to %s\n", path)

	return nil
}
