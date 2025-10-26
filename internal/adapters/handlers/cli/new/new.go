package newcmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func getStringFlag(cmd *cobra.Command, name string) (string, bool, error) {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", false, err
	}

	if cmd.Flags().Changed(name) {
		return value, true, nil
	}

	return "", false, nil
}

func getStringFlagWithLegacy(cmd *cobra.Command, primary, legacy, message string) (string, bool, error) {
	value, changed, err := getStringFlag(cmd, primary)
	if err != nil {
		return "", false, err
	}

	if changed {
		return value, true, nil
	}

	if legacy == "" {
		return "", false, nil
	}

	legacyValue, legacyChanged, err := getStringFlag(cmd, legacy)
	if err != nil {
		return "", false, err
	}

	if legacyChanged {
		emitDeprecatedWarning(cmd, message)
		return legacyValue, true, nil
	}

	return "", false, nil
}

func getStringArrayFlagWithLegacy(cmd *cobra.Command, primary, legacy, message string) ([]string, bool, error) {
	values, err := cmd.Flags().GetStringArray(primary)
	if err != nil {
		return nil, false, err
	}

	changed := cmd.Flags().Changed(primary)

	result := make([]string, 0, len(values))
	if changed {
		result = append(result, values...)
	}

	if legacy != "" {
		legacyValues, err := cmd.Flags().GetStringArray(legacy)
		if err != nil {
			return nil, false, err
		}

		if cmd.Flags().Changed(legacy) {
			emitDeprecatedWarning(cmd, message)

			result = append(result, legacyValues...)
			changed = true
		}
	}

	return result, changed, nil
}

func emitDeprecatedWarning(cmd *cobra.Command, message string) {
	fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s\n", message)
}

// Command constructs a fresh instance of the `pm new` Cobra command.

func run(cmd *cobra.Command, args []string) error {
	container, err := bootstrap.New(bootstrap.Options{Logger: bootstrap.WrapSlog(slog.Default())})
	if err != nil {
		return fmt.Errorf("bootstrap services: %w", err)
	}

	if container.ProjectInput == nil {
		return errors.New("project input service not available")
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
		return outputSkeleton(cmd, container, args, jsonSkeletonPath, services.SkeletonFormatJSON)
	}

	if yamlSkeletonSet {
		return outputSkeleton(cmd, container, args, yamlSkeletonPath, services.SkeletonFormatYAML)
	}

	cliInputPath, err := cmd.Flags().GetString(flagCliInput)
	if err != nil {
		return err
	}

	projectDescription, projectDescriptionSet, err := getStringFlag(cmd, flagProjectDescription)
	if err != nil {
		return err
	}

	projectShell, projectShellSet, err := getStringFlag(cmd, flagProjectShell)
	if err != nil {
		return err
	}

	projectEnvFile, projectEnvFileSet, err := getStringFlag(cmd, flagProjectEnvFile)
	if err != nil {
		return err
	}

	envFileValue, envFileSet, err := getStringFlagWithLegacy(
		cmd,
		flagEnvironmentEnvFile,
		flagEnvFileLegacy,
		"--env-file is deprecated; use --environment-env-file instead",
	)
	if err != nil {
		return err
	}

	envModeValue, envModeSet, err := getStringFlagWithLegacy(
		cmd,
		flagEnvironmentMode,
		flagEnvModeLegacy,
		"--env-mode is deprecated; use --environment-mode instead",
	)
	if err != nil {
		return err
	}

	envColorValue, envColorSet, err := getStringFlagWithLegacy(
		cmd,
		flagEnvironmentColor,
		flagEnvColorLegacy,
		"--env-color is deprecated; use --environment-color instead",
	)
	if err != nil {
		return err
	}

	envVarPairs, envVarFlagsSet, err := getStringArrayFlagWithLegacy(
		cmd,
		flagEnvironmentEnvVar,
		flagEnvVarLegacy,
		"--env-var is deprecated; use --environment-env-var instead",
	)
	if err != nil {
		return err
	}

	var envVarMap map[string]string
	if envVarFlagsSet && len(envVarPairs) > 0 {
		envVarMap, err = container.ProjectInput.ParseEnvVarFlags(
			envVarPairs,
			fmt.Sprintf("--%s", flagEnvironmentEnvVar),
		)
		if err != nil {
			return err
		}
	}

	var cfg *domain.ConfigInput
	if strings.TrimSpace(cliInputPath) != "" {
		cfg, err = container.ProjectInput.LoadConfig(cliInputPath)
		if err != nil {
			return err
		}

		if cfg.HasLegacyEnvironments() && !allowUnknown {
			return fmt.Errorf(
				"CLI input uses deprecated \"environments\" map; remove it or pass --allow-unknown to ignore legacy fields",
			)
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

	envFlags := domain.EnvironmentConfigFlags{}
	if envFileSet {
		envFlags.EnvVarsFile = strings.TrimSpace(envFileValue)
		envFlags.EnvFileSet = true
	}

	if envModeSet {
		envFlags.EnvVarsMode = strings.TrimSpace(envModeValue)
		envFlags.ModeSet = true
	}

	if envColorSet {
		envFlags.Color = strings.TrimSpace(envColorValue)
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
		normalizedEnvName, normErr := domain.NormalizeEnvironmentName(envName)
		if normErr != nil {
			return normErr
		}

		envFlags.Name = normalizedEnvName
		envFlags.NameSet = true
	}

	here, err := cmd.Flags().GetBool(flagHere)
	if err != nil {
		return err
	}

	if here && (pathProvided || envName != "") {
		return fmt.Errorf("--here cannot be used with additional positional arguments")
	}

	hereSet := cmd.Flags().Changed(flagHere)

	projectFlags := domain.ProjectConfigFlags{
		Name:        nameArg,
		Path:        pathArg,
		Here:        here,
		PathSet:     pathProvided,
		HereSet:     hereSet,
		Environment: envFlags,
	}

	if projectDescriptionSet {
		projectFlags.Description = strings.TrimSpace(projectDescription)
		projectFlags.DescriptionSet = true
	}

	if projectShellSet {
		projectFlags.Shell = strings.TrimSpace(projectShell)
		projectFlags.ShellSet = true
	}

	if projectEnvFileSet {
		projectFlags.EnvFile = strings.TrimSpace(projectEnvFile)
		projectFlags.EnvFileSet = true
	}

	merged, envVars, err := container.ProjectInput.MergeInputs(cfg, projectFlags)
	if err != nil {
		return err
	}

	if existingProject {
		merged.Path = strings.TrimSpace(probe.ProjectPath)
		merged.Here = false
	}

	envRequested := container.ProjectInput.EnvironmentProvided(merged.Environment)
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

	if existingProject &&
		(merged.Environment == nil || merged.Environment.Name == nil || strings.TrimSpace(*merged.Environment.Name) == "") {
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

	yesFlag, err := cmd.Flags().GetBool(flagYes)
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool(flagForce)
	if err != nil {
		return err
	}

	dryRun, err := cmd.Flags().GetBool(flagDryRun)
	if err != nil {
		return err
	}

	if dryRun && force {
		force = false
	}

	outputFormat, err := cmd.Flags().GetString(flagOutput)
	if err != nil {
		return err
	}

	if dryRun {
		return renderDryRun(cmd, merged, envVars, outputFormat)
	}

	if existingProject {
		if container.EnvironmentManager == nil {
			return errors.New("environment manager not available")
		}

		if merged.Environment == nil {
			return errors.New("environment input is required when adding environments")
		}

		if !yesFlag {
			envName := "environment"

			if merged.Environment.Name != nil {
				trimmed := strings.TrimSpace(*merged.Environment.Name)
				if trimmed != "" {
					envName = trimmed
				}
			}

			confirmed, err := promptYesNo(cmd, fmt.Sprintf("Add environment %q to project %q?", envName, merged.Name))
			if err != nil {
				return err
			}

			if !confirmed {
				fmt.Fprintln(cmd.OutOrStdout(), "Operation cancelled.")
				return nil
			}
		}

		if err := container.EnvironmentManager.Apply(merged.Name, merged.Environment, domain.EnvironmentApplyOptions{Force: force}); err != nil {
			return err
		}

		return nil
	}

	if container.ProjectCreator == nil {
		return errors.New("project creator not available")
	}

	if !yesFlag {
		message := fmt.Sprintf("Create project %q at %s?", merged.Name, merged.Path)

		confirmed, err := promptYesNo(cmd, message)
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Fprintln(cmd.OutOrStdout(), "Operation cancelled.")
			return nil
		}
	}

	_, err = container.ProjectCreator.Create(merged, envVars, domain.ProjectCreateOptions{Force: force})

	return err
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

			if len(envDetails.EnvVars) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    env vars: %d entries\n", len(envDetails.EnvVars))
			}
		}
	}

	return nil
}

func promptYesNo(cmd *cobra.Command, prompt string) (bool, error) {
	out := cmd.OutOrStdout()
	in := cmd.InOrStdin()

	if _, err := fmt.Fprintf(out, "%s (y/N): ", prompt); err != nil {
		return false, err
	}

	reader := bufio.NewReader(in)

	response, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}

		return false, err
	}

	answer := strings.TrimSpace(strings.ToLower(response))

	return answer == "y" || answer == "yes", nil
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

func resolveEnvVarsForDryRun(primary, fallback map[string]string) map[string]string {
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

type skeletonContext struct {
	ProjectName string
	ProjectPath string
}

func resolveSkeletonContext(container *bootstrap.Container, args []string) (skeletonContext, error) {
	ctx := skeletonContext{}

	if len(args) == 0 {
		return ctx, nil
	}

	name, err := normalizeSkeletonProjectName(args[0])
	if err != nil {
		return ctx, err
	}

	ctx.ProjectName = name

	existingProject := false
	projectPath := ""

	if name != "" && container != nil && container.ProjectService != nil {
		probe, probeErr := container.ProjectService.Probe(name)
		if probeErr != nil {
			return ctx, fmt.Errorf("probe project %q: %w", name, probeErr)
		}

		if probe.RegistryHit && probe.ProjectFileExists {
			existingProject = true
			projectPath = strings.TrimSpace(probe.ProjectPath)
		}
	}

	if existingProject {
		ctx.ProjectPath = projectPath
	}

	if len(args) > 1 {
		second := strings.TrimSpace(args[1])
		if existingProject {
			if _, envErr := normalizeSkeletonEnvironmentName(second); envErr != nil {
				return ctx, envErr
			}
		} else {
			ctx.ProjectPath = second
		}
	}

	return ctx, nil
}

func normalizeSkeletonProjectName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	def := &domain.ProjectDefinition{Name: trimmed, Here: true}
	if err := domain.ValidateProjectDefinition(def); err != nil {
		return "", fmt.Errorf("invalid project name %q: %s", raw, err)
	}

	return def.Name, nil
}

func normalizeSkeletonEnvironmentName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	normalized, err := domain.NormalizeEnvironmentName(trimmed)
	if err != nil {
		return "", err
	}

	return normalized, nil
}

func outputSkeleton(
	cmd *cobra.Command,
	container *bootstrap.Container,
	args []string,
	rawPath string,
	format services.SkeletonFormat,
) error {
	if len(args) > 2 {
		return fmt.Errorf("at most two positional arguments are allowed when generating CLI input skeletons")
	}

	ctx, err := resolveSkeletonContext(container, args)
	if err != nil {
		return err
	}

	opts := services.SkeletonOptions{
		ProjectName: ctx.ProjectName,
		ProjectPath: ctx.ProjectPath,
	}

	if projectDescription, changed, err := getStringFlag(cmd, flagProjectDescription); err != nil {
		return err
	} else if changed {
		opts.Description = strings.TrimSpace(projectDescription)
	}

	if projectShell, changed, err := getStringFlag(cmd, flagProjectShell); err != nil {
		return err
	} else if changed {
		opts.Shell = strings.TrimSpace(projectShell)
	}

	if projectEnvFile, changed, err := getStringFlag(cmd, flagProjectEnvFile); err != nil {
		return err
	} else if changed {
		opts.EnvFile = strings.TrimSpace(projectEnvFile)
	}

	envVarPairs, envVarsFlagSet, err := getStringArrayFlagWithLegacy(
		cmd,
		flagEnvironmentEnvVar,
		flagEnvVarLegacy,
		"--env-var is deprecated; use --environment-env-var instead",
	)
	if err != nil {
		return err
	}

	if envVarsFlagSet && len(envVarPairs) > 0 {
		envVars, err := container.ProjectInput.ParseEnvVarFlags(envVarPairs, fmt.Sprintf("--%s", flagEnvironmentEnvVar))
		if err != nil {
			return err
		}

		if len(envVars) > 0 {
			opts.EnvVars = envVars
		}
	}

	if len(envVarPairs) > 0 {
		envVarMap, parseErr := container.ProjectInput.ParseEnvVarFlags(
			envVarPairs,
			fmt.Sprintf("--%s", flagEnvironmentEnvVar),
		)
		if parseErr != nil {
			return parseErr
		}

		if len(envVarMap) > 0 {
			opts.EnvVars = envVarMap
		}
	}

	path := strings.TrimSpace(rawPath)
	if path == "" {
		path = "-"
	}

	if path == "-" {
		data, err := services.RenderProjectSkeletonWithOptions(format, opts)
		if err != nil {
			return err
		}

		if _, err := cmd.OutOrStdout().Write(data); err != nil {
			return fmt.Errorf("write skeleton to stdout: %w", err)
		}

		return nil
	}

	if err := services.GenerateProjectSkeletonWithOptions(path, format, opts); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "CLI input skeleton written to %s\n", path)

	return nil
}
