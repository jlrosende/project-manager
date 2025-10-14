package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

var allowedProjectConfigKeys = map[string]struct{}{
	"name":        {},
	"here":        {},
	"path":        {},
	"environment": {},
	"metadata":    {},
}

// SkeletonFormat represents the serialization format for generated CLI input files.
type SkeletonFormat string

const (
	SkeletonFormatJSON SkeletonFormat = "json"
	SkeletonFormatYAML SkeletonFormat = "yaml"
)

// ProjectConfigFlags captures CLI-provided overrides for project creation.
type ProjectConfigFlags struct {
	Name        string
	Path        string
	Here        bool
	PathSet     bool
	HereSet     bool
	Environment EnvironmentConfigFlags
}

// EnvironmentConfigFlags captures CLI-provided overrides for environment creation.
type EnvironmentConfigFlags struct {
	Name        string
	NameSet     bool
	EnvVarsFile string
	EnvFileSet  bool
	EnvVarsMode string
	ModeSet     bool
	Color       string
	ColorSet    bool
	EnvVars     map[string]string
}

// HasInput reports whether any environment-related flag values were provided.
func (f EnvironmentConfigFlags) HasInput() bool {
	return f.NameSet || f.EnvFileSet || f.ModeSet || f.ColorSet || len(f.EnvVars) > 0
}

// MergeProjectInputs combines configuration file inputs with CLI arguments,
// preferring explicit flag values. It returns the merged project definition and
// the environment variables destined for the .env file.
func MergeProjectInputs(
	cfg *domain.ConfigInput,
	flags ProjectConfigFlags,
) (domain.ProjectDefinition, domain.EnvVars, error) {
	var (
		def     domain.ProjectDefinition
		envVars = domain.EnvVars{}
	)

	var env *domain.EnvironmentInput
	if cfg != nil {
		if cfg.Name != nil {
			def.Name = strings.TrimSpace(*cfg.Name)
		}

		if cfg.Here != nil {
			def.Here = *cfg.Here
		}

		if cfg.Path != nil {
			def.Path = strings.TrimSpace(*cfg.Path)
		}

		if cfg.Environment != nil {
			env = cloneEnvironmentInput(cfg.Environment)
		}

		if len(cfg.Metadata) > 0 {
			def.Metadata = copyStringMap(cfg.Metadata)
		}
	}

	envFlags := flags.Environment
	if envFlags.HasInput() {
		if env == nil {
			env = &domain.EnvironmentInput{}
		}

		if envFlags.NameSet {
			env.Name = stringPtr(strings.TrimSpace(envFlags.Name))
		}

		if envFlags.EnvFileSet {
			env.EnvVarsFile = stringPtr(strings.TrimSpace(envFlags.EnvVarsFile))
		}

		if envFlags.ModeSet {
			env.EnvVarsMode = stringPtr(strings.TrimSpace(envFlags.EnvVarsMode))
		}

		if envFlags.ColorSet {
			env.Color = stringPtr(strings.TrimSpace(envFlags.Color))
		}

		if len(envFlags.EnvVars) > 0 {
			if env.EnvVars == nil {
				env.EnvVars = make(map[string]string, len(envFlags.EnvVars))
			}

			for k, v := range envFlags.EnvVars {
				trimmedKey := strings.TrimSpace(k)
				if trimmedKey == "" {
					continue
				}
				env.EnvVars[trimmedKey] = v
			}
		}
	}

	if env != nil {
		def.Environment = env
		if len(env.EnvVars) > 0 {
			envVars = domain.EnvVars(copyStringMap(env.EnvVars))
		}
	}

	if flags.PathSet {
		def.Path = strings.TrimSpace(flags.Path)
		if def.Here {
			def.Here = false
		}
	}

	if flags.HereSet {
		def.Here = flags.Here
		if flags.Here && !flags.PathSet {
			def.Path = ""
		}
	}

	if strings.TrimSpace(flags.Name) != "" {
		def.Name = strings.TrimSpace(flags.Name)
	}

	return def, envVars, nil
}

const (
	projectConflictExitCode = 3
	projectConflictCodeName = "NEW-CONFLICT-NAME"
	projectConflictCodePath = "NEW-CONFLICT-PATH"
)

type projectConflictKind int

const (
	projectConflictKindName projectConflictKind = iota + 1
	projectConflictKindPath
)

// ProjectConflictError describes a uniqueness violation during project creation.
type ProjectConflictError struct {
	kind          projectConflictKind
	projectName   string
	existingPath  string
	candidatePath string
}

// Error implements the error interface.
func (e *ProjectConflictError) Error() string {
	if e == nil {
		return ""
	}

	switch e.kind {
	case projectConflictKindName:
		path := strings.TrimSpace(e.existingPath)
		if path == "" {
			return fmt.Sprintf("[%s] project %q already exists", e.Code(), e.projectName)
		}

		return fmt.Sprintf(
			"[%s] project %q already exists at %s; rerun without specifying a path to add environments or choose a different name.",
			e.Code(),
			e.projectName,
			path,
		)
	case projectConflictKindPath:
		name := strings.TrimSpace(e.projectName)
		if name == "" {
			return fmt.Sprintf(
				"[%s] path %s is already registered to another project; choose a different destination.",
				e.Code(),
				e.candidatePath,
			)
		}

		return fmt.Sprintf(
			"[%s] path %s is already registered to project %q; choose a different destination.",
			e.Code(),
			e.candidatePath,
			name,
		)
	default:
		return "project conflict detected"
	}
}

// ExitCode exposes the CLI exit code associated with the conflict.
func (e *ProjectConflictError) ExitCode() int {
	if e == nil {
		return 0
	}

	return projectConflictExitCode
}

// Code returns the symbolic identifier for the conflict.
func (e *ProjectConflictError) Code() string {
	if e == nil {
		return ""
	}

	switch e.kind {
	case projectConflictKindName:
		return projectConflictCodeName
	case projectConflictKindPath:
		return projectConflictCodePath
	default:
		return ""
	}
}

// Field indicates which field triggered the conflict.
func (e *ProjectConflictError) Field() string {
	if e == nil {
		return ""
	}

	switch e.kind {
	case projectConflictKindName:
		return "name"
	case projectConflictKindPath:
		return "path"
	default:
		return ""
	}
}

// ProjectName returns the conflicting project name.
func (e *ProjectConflictError) ProjectName() string {
	if e == nil {
		return ""
	}

	return e.projectName
}

// ExistingPath returns the path associated with the conflicting project.
func (e *ProjectConflictError) ExistingPath() string {
	if e == nil {
		return ""
	}

	return e.existingPath
}

// CandidatePath returns the requested path that triggered the conflict.
func (e *ProjectConflictError) CandidatePath() string {
	if e == nil {
		return ""
	}

	return e.candidatePath
}

func newProjectConflictError(kind projectConflictKind, name, existingPath, candidatePath string) *ProjectConflictError {
	return &ProjectConflictError{
		kind:          kind,
		projectName:   strings.TrimSpace(name),
		existingPath:  strings.TrimSpace(existingPath),
		candidatePath: strings.TrimSpace(candidatePath),
	}
}

// EnsureProjectUniqueness verifies that the candidate project definition does not
// collide with existing registry entries. Entries lacking a .project.hcl file are
// ignored to allow recreation after manual cleanup.
func EnsureProjectUniqueness(projects []*domain.Project, candidate domain.ProjectDefinition) error {
	trimmedName := strings.TrimSpace(candidate.Name)
	trimmedPath := strings.TrimSpace(candidate.Path)

	for _, existing := range projects {
		if existing == nil {
			continue
		}

		existingName := strings.TrimSpace(existing.Name)
		existingPath := strings.TrimSpace(existing.Path)
		if existingName == "" && existingPath == "" {
			continue
		}

		hasMetadata, err := projectMetadataExists(existingPath)
		if err != nil {
			return fmt.Errorf("check project metadata for %s: %w", existingPath, err)
		}

		if !hasMetadata {
			continue
		}

		if trimmedName != "" && existingName == trimmedName {
			return newProjectConflictError(projectConflictKindName, existingName, existingPath, trimmedPath)
		}

		if trimmedPath != "" && pathsEqual(existingPath, trimmedPath) {
			return newProjectConflictError(projectConflictKindPath, existingName, existingPath, trimmedPath)
		}
	}

	return nil
}

func projectMetadataExists(root string) (bool, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return false, nil
	}

	meta := filepath.Join(root, ".project.hcl")
	info, err := os.Stat(meta)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	if info.IsDir() {
		return false, nil
	}

	return true, nil
}

func pathsEqual(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == "" || b == "" {
		return false
	}

	cleanA := filepath.Clean(a)
	cleanB := filepath.Clean(b)

	if runtime.GOOS == "windows" {
		return strings.EqualFold(cleanA, cleanB)
	}

	return cleanA == cleanB
}

// LoadProjectConfig reads a JSON or YAML configuration file from disk.
func LoadProjectConfig(path string) (*domain.ConfigInput, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file %s is empty", path)
	}

	ext := strings.ToLower(filepath.Ext(path))

	var (
		cfg domain.ConfigInput
		raw map[string]any
	)

	switch ext {
	case ".json":
		raw, err = decodeProjectJSON(data, &cfg)
	case ".yaml", ".yml":
		raw, err = decodeProjectYAML(data, &cfg)
	default:
		raw, err = decodeProjectJSON(data, &cfg)
		if err != nil {
			raw, err = decodeProjectYAML(data, &cfg)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}

	legacyEnvironments := false
	if raw != nil {
		if _, ok := raw["environments"]; ok {
			legacyEnvironments = true
			delete(raw, "environments")
		}
	}

	cfg.SetLegacyEnvironments(legacyEnvironments)
	cfg.SetUnknownFields(filterProjectUnknown(raw))

	return &cfg, nil
}

func decodeProjectJSON(data []byte, cfg *domain.ConfigInput) (map[string]any, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var parsed domain.ConfigInput
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	*cfg = parsed

	return raw, nil
}

func decodeProjectYAML(data []byte, cfg *domain.ConfigInput) (map[string]any, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var parsed domain.ConfigInput
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	*cfg = parsed

	return raw, nil
}

func filterProjectUnknown(raw map[string]any) map[string]any {
	if len(raw) == 0 {
		return nil
	}

	unknown := map[string]any{}

	for k, v := range raw {
		if _, ok := allowedProjectConfigKeys[k]; !ok {
			unknown[k] = v
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	return unknown
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}

func stringPtr(value string) *string {
	v := value
	return &v
}

func cloneEnvironmentInput(in *domain.EnvironmentInput) *domain.EnvironmentInput {
	if in == nil {
		return nil
	}

	out := &domain.EnvironmentInput{}

	if in.Name != nil {
		name := *in.Name
		out.Name = &name
	}

	if in.EnvVarsFile != nil {
		file := *in.EnvVarsFile
		out.EnvVarsFile = &file
	}

	if in.EnvVarsMode != nil {
		mode := *in.EnvVarsMode
		out.EnvVarsMode = &mode
	}

	if in.Color != nil {
		color := *in.Color
		out.Color = &color
	}

	if len(in.EnvVars) > 0 {
		out.EnvVars = copyStringMap(in.EnvVars)
	}

	return out
}

// GenerateProjectSkeleton writes a CLI input skeleton in the requested format to the

// provided path. The generated file can be edited and passed to --cli-input for
// future project creation runs.
func GenerateProjectSkeleton(path string, format SkeletonFormat) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("skeleton output path is empty")
	}

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}

		return fmt.Errorf("file %s already exists", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check skeleton destination: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("prepare skeleton destination: %w", err)
	}

	data, err := RenderProjectSkeleton(format)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write skeleton file: %w", err)
	}

	return nil
}

// RenderProjectSkeleton returns the serialized CLI input skeleton for the requested format.
func RenderProjectSkeleton(format SkeletonFormat) ([]byte, error) {
	template := defaultProjectConfigSkeleton()

	var (
		data []byte
		err  error
	)

	switch format {
	case SkeletonFormatJSON:
		data, err = json.MarshalIndent(template, "", "  ")
	case SkeletonFormatYAML:
		data, err = yaml.Marshal(template)
	default:
		return nil, fmt.Errorf("unsupported skeleton format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("render skeleton: %w", err)
	}

	if !bytes.HasSuffix(data, []byte("\n")) {
		data = append(data, '\n')
	}

	return data, nil
}

func defaultProjectConfigSkeleton() domain.ConfigInput {
	name := "your-project-name"
	path := "/absolute/path/to/your-project"
	here := false
	envName := "example"
	envFile := ".env.example"
	envMode := domain.EnvVarsModeMerge

	return domain.ConfigInput{
		Name: &name,
		Path: &path,
		Here: &here,
		Environment: &domain.EnvironmentInput{
			Name:        &envName,
			EnvVarsFile: &envFile,
			EnvVarsMode: &envMode,
			EnvVars: map[string]string{
				"EXAMPLE_ENV_VAR": "value",
			},
		},
		Metadata: map[string]string{
			"description": "Describe your project",
		},
	}
}

// DeleteCLIFlags captures the parsed CLI flags for project deletion.
type DeleteCLIFlags struct {
	All                    bool
	KeepFiles              bool
	OnlyEnv                bool
	DryRun                 bool
	Force                  bool
	Backup                 bool
	BackupDestination      string
	DefaultBackupDirectory string
}

var errConflictingDeleteScopes = errors.New("conflicting delete scope flags")

// ResolveDeleteScope determines which scope should be executed given the parsed
// flags.
func ResolveDeleteScope(flags DeleteCLIFlags) (domain.DeleteScope, error) {
	selected := 0
	if flags.All {
		selected++
	}

	if flags.KeepFiles {
		selected++
	}

	if flags.OnlyEnv {
		selected++
	}

	if selected > 1 {
		return domain.DeleteScopeMetadata, errConflictingDeleteScopes
	}

	switch {
	case flags.All:
		return domain.DeleteScopeAll, nil
	case flags.KeepFiles:
		return domain.DeleteScopeKeepFiles, nil
	case flags.OnlyEnv:
		return domain.DeleteScopeEnvOnly, nil
	default:
		return domain.DeleteScopeMetadata, nil
	}
}

// BuildDeleteOptions creates a domain-level ProjectDeleteOptions structure based
// on the resolved CLI flags.
func BuildDeleteOptions(target domain.ProjectIdentifier, flags DeleteCLIFlags) (domain.ProjectDeleteOptions, error) {
	scope, err := ResolveDeleteScope(flags)
	if err != nil {
		return domain.ProjectDeleteOptions{}, err
	}

	options := domain.ProjectDeleteOptions{
		Target: target,
		Scope:  scope,
		DryRun: flags.DryRun,
		Force:  flags.Force,
	}

	if flags.Backup {
		destination := strings.TrimSpace(flags.BackupDestination)

		baseDir := strings.TrimSpace(flags.DefaultBackupDirectory)
		if baseDir == "" {
			baseDir = filepath.Join("~", ".pm", "backups")
		}

		if destination == "" {
			destination = defaultBackupDestination(baseDir, target)
		}

		options.Backup = &domain.BackupRequest{
			Destination:      destination,
			IncludeWorkspace: scope.IncludesWorkspace(),
		}
	}

	if err := options.Validate(); err != nil {
		return domain.ProjectDeleteOptions{}, err
	}

	return options, nil
}

func defaultBackupDestination(baseDir string, target domain.ProjectIdentifier) string {
	label := strings.TrimSpace(target.Name)
	if label == "" {
		label = strings.TrimSpace(target.Path)
		if label != "" {
			label = filepath.Base(label)
		}
	}

	slug := sanitizeBackupLabel(label)
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.zip", timestamp, slug)

	if strings.TrimSpace(baseDir) == "" {
		baseDir = filepath.Join("~", ".pm", "backups")
	}

	return filepath.Join(baseDir, filename)
}

func sanitizeBackupLabel(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return "project"
	}

	var builder strings.Builder

	lastDash := false

	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)

			lastDash = false
		case r == '-' || r == '_' || r == '.':
			builder.WriteRune(r)

			lastDash = false
		case unicode.IsSpace(r):
			if !lastDash {
				builder.WriteByte('-')

				lastDash = true
			}
		default:
			if !lastDash {
				builder.WriteByte('-')

				lastDash = true
			}
		}
	}

	slug := strings.Trim(builder.String(), "-_")
	if slug == "" {
		return "project"
	}

	return slug
}
