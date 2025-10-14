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
	flagEnvFile              = "env-file"
	flagEnvMode              = "env-mode"
	flagEnvColor             = "env-color"
	flagEnvVar               = "env-var"
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
	cmd.Flags().String(flagEnvFile, "", "Set environment vars file when adding environments")
	cmd.Flags().String(flagEnvMode, "", "Set environment vars merge mode (merge or replace)")
	cmd.Flags().String(flagEnvColor, "", "Set environment color metadata when adding environments")
	cmd.Flags().StringArray(flagEnvVar, nil, "Environment variable in KEY=VALUE format (repeatable)")

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

	envFileFlag, err := cmd.Flags().GetString(flagEnvFile)
	if err != nil {
		return err
	}

	envModeFlag, err := cmd.Flags().GetString(flagEnvMode)
	if err != nil {
		return err
	}

	envColorFlag, err := cmd.Flags().GetString(flagEnvColor)
	if err != nil {
		return err
	}

	envVarPairs, err := cmd.Flags().GetStringArray(flagEnvVar)
	if err != nil {
		return err
	}

	envVarMap, err := parseEnvVarFlags(envVarPairs)
	if err != nil {
		return err
	}

	var cfg *domain.ConfigInput
	if strings.TrimSpace(cliInputPath) != "" {
		cfg, err = services.LoadProjectConfig(cliInputPath)
		if err != nil {
			return err
		}

		if cfg.HasLegacyEnvironments() && !allowUnknown {
			return fmt.Errorf("CLI input uses deprecated \"environments\" map; remove it or pass --allow-unknown to ignore legacy fields")
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

	envFlags := services.EnvironmentConfigFlags{}
	if cmd.Flags().Changed(flagEnvFile) {
		envFlags.EnvVarsFile = envFileFlag
		envFlags.EnvFileSet = true
	}

	if cmd.Flags().Changed(flagEnvMode) {
		envFlags.EnvVarsMode = envModeFlag
		envFlags.ModeSet = true
	}

	if cmd.Flags().Changed(flagEnvColor) {
		envFlags.Color = envColorFlag
		envFlags.ColorSet = true
	}

	if len(envVarMap) > 0 {
		envFlags.EnvVars = envVarMap
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

	if envName != "" {
		envFlags.Name = envName
		envFlags.NameSet = true
	}

	here, err := cmd.Flags().GetBool(flagHere)
	if err != nil {
		return err
	}

	if here && (pathProvided || envName != "") {
		return fmt.Errorf("--here cannot be used with additional positional arguments")
	}

	merged, envVars, err := services.MergeProjectInputs(cfg, services.ProjectConfigFlags{
		Name:        nameArg,
		Path:        pathArg,
		Here:        here,
		PathSet:     pathProvided,
		HereSet:     cmd.Flags().Changed(flagHere),
		Environment: envFlags,
	})
	if err != nil {
		return err
	}

	if existingProject {
		merged.Path = strings.TrimSpace(probe.ProjectPath)
		merged.Here = false
	}

	envRequested := environmentProvided(merged.Environment)
	if !existingProject && envRequested {
		envLabel := "environment input"
		if merged.Environment != nil && merged.Environment.Name != nil {
			if trimmed := strings.TrimSpace(*merged.Environment.Name); trimmed != "" {
				envLabel = fmt.Sprintf("environment %q", trimmed)
			}
		}

		return fmt.Errorf("%s requires an existing project; create the project before adding environments", envLabel)
	}

	if existingProject && !envRequested {
		fmt.Fprintf(cmd.OutOrStdout(), "project %s already exists; nothing to do\n", merged.Name)
		return nil
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

	if existingProject {
		envService := services.NewEnvironmentService(container.ProjectService, container.EnvVarsRepo, container.Filesystem)
		if merged.Environment == nil {
			return errors.New("environment input is required when adding environments")
		}

		if err := envService.Apply(merged.Name, merged.Environment, services.EnvironmentApplyOptions{Force: force}); err != nil {
			return err
		}

		return nil
	}

	creator := services.NewFileService(container.ProjectService, container.Filesystem)
	if creator == nil {
		return errors.New("project file service not available")
	}

	_, err = creator.Create(merged, envVars, services.CreateOptions{Force: force})

	return err
}

func parseEnvVarFlags(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}

	vars := make(map[string]string, len(pairs))
	for _, raw := range pairs {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		idx := strings.Index(trimmed, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("invalid --%s value %q; expected KEY=VALUE", flagEnvVar, raw)
		}

		key := strings.TrimSpace(trimmed[:idx])
		value := trimmed[idx+1:]
		if key == "" {
			return nil, fmt.Errorf("--%s requires a non-empty key: %q", flagEnvVar, raw)
		}

		vars[key] = value
	}

	if len(vars) == 0 {
		return nil, nil
	}

	return vars, nil
}

func environmentProvided(input *domain.EnvironmentInput) bool {
	if input == nil {
		return false
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		return true
	}

	if input.EnvVarsFile != nil && strings.TrimSpace(*input.EnvVarsFile) != "" {
		return true
	}

	if input.EnvVarsMode != nil && strings.TrimSpace(*input.EnvVarsMode) != "" {
		return true
	}

	if input.Color != nil && strings.TrimSpace(*input.Color) != "" {
		return true
	}

	return len(input.EnvVars) > 0
}

func renderDryRun(cmd *cobra.Command, def domain.ProjectDefinition, envVars domain.EnvVars, format string) error {
	envDetails := summarizeEnvironment(def.Environment, envVars)
	projectEnvVars := envVars
	if envDetails != nil {
		projectEnvVars = nil
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		payload := map[string]any{
			"name": def.Name,
			"path": def.Path,
			"here": def.Here,
		}

		if len(projectEnvVars) > 0 {
			payload["env_vars_count"] = len(projectEnvVars)
			payload["env_vars"] = copyEnvVars(projectEnvVars)
		}

		if envDetails != nil {
			payload["environment"] = envDetails.asMap()
		}

		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("format dry-run output: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(data))
	default:
		fmt.Fprintf(cmd.OutOrStdout(), "pm new (dry-run)\n  name: %s\n  path: %s\n", def.Name, def.Path)

		if def.Here {
			fmt.Fprintln(cmd.OutOrStdout(), "  here: true")
		}

		if len(projectEnvVars) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  env vars: %d entries\n", len(projectEnvVars))
		}

		if envDetails != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "  environment:")
			fmt.Fprintf(cmd.OutOrStdout(), "    name: %s\n", envDetails.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "    env vars file: %s\n", envDetails.EnvVarsFile)
			fmt.Fprintf(cmd.OutOrStdout(), "    env vars mode: %s\n", envDetails.EnvVarsMode)

			if envDetails.Color != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "    color: %s\n", envDetails.Color)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "    env vars: %d entries\n", len(envDetails.EnvVars))
		}
	}

	return nil
}

type environmentDryRunDetails struct {
	Name        string
	EnvVarsFile string
	EnvVarsMode string
	Color       string
	EnvVars     map[string]string
}

func (d environmentDryRunDetails) asMap() map[string]any {
	payload := map[string]any{
		"name":           d.Name,
		"env_vars_file":  d.EnvVarsFile,
		"env_vars_mode":  d.EnvVarsMode,
		"env_vars_count": len(d.EnvVars),
	}

	if d.Color != "" {
		payload["color"] = d.Color
	}

	payload["env_vars"] = copyEnvVars(d.EnvVars)

	return payload
}

func summarizeEnvironment(input *domain.EnvironmentInput, merged domain.EnvVars) *environmentDryRunDetails {
	if input == nil {
		return nil
	}

	summary := &environmentDryRunDetails{
		Name: trimPtr(input.Name),
	}

	file := trimPtr(input.EnvVarsFile)
	if file == "" && summary.Name != "" {
		file = dryRunDefaultEnvironmentFile(summary.Name)
	}
	summary.EnvVarsFile = file

	mode := strings.ToLower(trimPtr(input.EnvVarsMode))
	if mode == "" {
		mode = domain.EnvVarsModeMerge
	}
	summary.EnvVarsMode = mode
	summary.Color = trimPtr(input.Color)

	summary.EnvVars = resolveEnvVarsForDryRun(merged, input.EnvVars)

	return summary
}

func copyEnvVars(vars map[string]string) map[string]string {
	if len(vars) == 0 {
		return map[string]string{}
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}

	return out
}

func resolveEnvVarsForDryRun(primary map[string]string, fallback map[string]string) map[string]string {
	if len(primary) > 0 {
		return copyEnvVars(primary)
	}

	if len(fallback) > 0 {
		return copyEnvVars(fallback)
	}

	return map[string]string{}
}

func trimPtr(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func dryRunDefaultEnvironmentFile(name string) string {
	slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(name, " ", "-")))
	if slug == "" {
		return ".env"
	}

	return "." + slug + ".env"
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
