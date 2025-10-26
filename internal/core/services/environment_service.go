package services

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

// ErrEnvironmentExists indicates the requested environment already exists and cannot be recreated without force.
var ErrEnvironmentExists = errors.New("environment already exists")

// EnvironmentService orchestrates environment creation and updates for existing projects.
type EnvironmentService struct {
	projects ports.ProjectService
	envVars  ports.EnvVarsRepository
	fs       ports.Filesystem
}

// NewEnvironmentService constructs an EnvironmentService with the required collaborators.
func NewEnvironmentService(
	projects ports.ProjectService,
	envVars ports.EnvVarsRepository,
	fs ports.Filesystem,
) *EnvironmentService {
	if projects == nil || envVars == nil || fs == nil {
		panic("environment service dependencies must not be nil")
	}

	return &EnvironmentService{projects: projects, envVars: envVars, fs: fs}
}

var _ ports.EnvironmentManager = (*EnvironmentService)(nil)

// Apply adds or updates an environment for the specified project based on the provided input.
func (s *EnvironmentService) Apply(
	projectName string,
	input *domain.EnvironmentInput,
	opts domain.EnvironmentApplyOptions,
) error {
	if input == nil {
		return errors.New("environment input is nil")
	}

	project, err := s.projects.Load(projectName)
	if err != nil {
		return fmt.Errorf("load project %s: %w", projectName, err)
	}

	if project == nil {
		return fmt.Errorf("project %s not found", projectName)
	}

	root, err := s.projectRoot(project.Path)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}

	envName := ""
	if input.Name != nil {
		normalized, err := domain.NormalizeEnvironmentName(*input.Name)
		if err != nil {
			return err
		}

		envName = normalized
	}

	if envName == "" {
		return errors.New("environment name is required")
	}

	relFile, absFile, err := s.resolveEnvFile(root, envName, input.EnvVarsFile)
	if err != nil {
		return err
	}

	envMode := domain.EnvVarsModeMerge
	if input.EnvVarsMode != nil && strings.TrimSpace(*input.EnvVarsMode) != "" {
		mode := strings.ToLower(strings.TrimSpace(*input.EnvVarsMode))
		if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
			return fmt.Errorf(
				"invalid env vars mode %q; expected %q or %q",
				mode,
				domain.EnvVarsModeMerge,
				domain.EnvVarsModeReplace,
			)
		}

		envMode = mode
	}

	color := ""
	if input.Color != nil {
		color = strings.TrimSpace(*input.Color)
	}

	envVars := cloneEnvVars(input.EnvVars)

	existing := findEnvironment(project.Environments, envName)
	oldEnvPath := ""
	if existing != nil {
		oldEnvPath = s.resolveExistingEnvFile(project.Path, root, existing.EnvVarsFile)
	}

	envDefinition := &domain.Environment{
		Name:        envName,
		Color:       color,
		EnvVarsMode: envMode,
		EnvVarsFile: relFile,
	}

	if existing != nil {
		if !opts.Force {
			return fmt.Errorf("%w: %s", ErrEnvironmentExists, envName)
		}

		if err := s.projects.UpdateEnvironment(projectName, existing.Name, envDefinition); err != nil {
			return fmt.Errorf("update environment %s: %w", envName, err)
		}

		if envVars == nil {
			envVars = domain.EnvVars{}
		}

		if err := s.envVars.Save(absFile, envVars); err != nil {
			return fmt.Errorf("write environment vars file: %w", err)
		}

		if oldEnvPath != "" && !samePath(oldEnvPath, absFile) {
			if err := s.fs.Remove(oldEnvPath); err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("remove previous env vars file: %w", err)
			}
		}

		return nil
	}

	if envVars == nil {
		envVars = domain.EnvVars{}
	}

	if err := s.projects.AddEnvironment(projectName, envDefinition, envVars); err != nil {
		return fmt.Errorf("add environment %s: %w", envName, err)
	}

	return nil
}

func (s *EnvironmentService) projectRoot(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", errors.New("project path is empty")
	}

	expanded := s.fs.ExpandHome(trimmed)
	var abs string
	if s.fs.IsAbs(expanded) {
		abs = expanded
	} else {
		resolved, err := s.fs.Abs(expanded)
		if err != nil {
			return "", err
		}

		abs = resolved
	}

	return filepath.Clean(abs), nil
}

func (s *EnvironmentService) resolveEnvFile(
	projectRoot string,
	envName string,
	filePtr *string,
) (relative string, absolute string, err error) {
	candidate := ""
	if filePtr != nil {
		candidate = strings.TrimSpace(*filePtr)
	}

	if candidate == "" {
		candidate = defaultEnvironmentFile(envName)
	}

	expanded := s.fs.ExpandHome(candidate)
	var absPath string
	if s.fs.IsAbs(expanded) {
		absPath = expanded
	} else {
		absPath = s.fs.Join(projectRoot, expanded)
	}

	absPath = filepath.Clean(absPath)
	root := filepath.Clean(projectRoot)

	relPath, relErr := filepath.Rel(root, absPath)
	if relErr != nil {
		return "", "", fmt.Errorf("resolve env vars file: %w", relErr)
	}

	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("env vars file %q must reside within project directory %s", candidate, projectRoot)
	}

	if relPath == "." {
		relPath = defaultEnvironmentFile(envName)
		absPath = filepath.Join(root, relPath)
	}

	return relPath, absPath, nil
}

func defaultEnvironmentFile(name string) string {
	slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(name, " ", "-")))
	if slug == "" {
		return ".env"
	}

	return "." + slug + ".env"
}

func cloneEnvVars(in map[string]string) domain.EnvVars {
	if len(in) == 0 {
		return nil
	}

	out := make(domain.EnvVars, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}

func findEnvironment(environments []*domain.Environment, target string) *domain.Environment {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil
	}

	for _, env := range environments {
		if env != nil && env.Name == target {
			return env
		}
	}

	return nil
}

func (s *EnvironmentService) resolveExistingEnvFile(projectPath, projectRoot, candidate string) string {
	trimmed := strings.TrimSpace(candidate)
	if trimmed == "" {
		return ""
	}

	expanded := s.fs.ExpandHome(trimmed)
	if s.fs.IsAbs(expanded) {
		return filepath.Clean(expanded)
	}

	base := strings.TrimSpace(projectRoot)
	if base == "" {
		base = strings.TrimSpace(projectPath)
		if base == "" {
			return ""
		}

		if !s.fs.IsAbs(base) {
			abs, err := s.fs.Abs(base)
			if err != nil {
				return ""
			}

			base = abs
		}
	}

	return filepath.Clean(filepath.Join(base, expanded))
}

func samePath(a, b string) bool {
	if strings.TrimSpace(a) == "" || strings.TrimSpace(b) == "" {
		return false
	}

	cleanA := filepath.Clean(a)
	cleanB := filepath.Clean(b)

	if runtime.GOOS == "windows" {
		return strings.EqualFold(cleanA, cleanB)
	}

	return cleanA == cleanB
}
