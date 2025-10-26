package services

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

const defaultEnvFile = ".env"

// FileService coordinates filesystem updates required for project creation.
type FileService struct {
	projects ports.ProjectService
	fs       ports.Filesystem
}

// NewFileService constructs a FileService with the required collaborators.
func NewFileService(projects ports.ProjectService, fs ports.Filesystem) *FileService {
	return &FileService{projects: projects, fs: fs}
}

var _ ports.ProjectCreator = (*FileService)(nil)

// Create writes project artifacts (.project.hcl, .env) using the provided definition and environment variables.
func (s *FileService) Create(
	def domain.ProjectDefinition,
	envVars domain.EnvVars,
	opts domain.ProjectCreateOptions,
) (*domain.Project, error) {
	if s.projects == nil {
		return nil, errors.New("project service is not configured")
	}

	if s.fs == nil {
		return nil, errors.New("filesystem is not configured")
	}

	existing, err := s.projects.List()
	if err != nil {
		return nil, err
	}

	if err := EnsureProjectUniqueness(existing, def); err != nil {
		return nil, err
	}

	envFile := strings.TrimSpace(def.EnvVarsFile)
	if envFile == "" {
		envFile = defaultEnvFile
	}

	shell := strings.TrimSpace(def.Shell)

	if opts.Force {
		if err := s.cleanupExisting(def, envFile); err != nil {
			return nil, err
		}
	}

	project, err := s.projects.Create(def.Name, def.Path, "", shell, envFile, envVars, nil)
	if err != nil {
		return nil, err
	}

	description := strings.TrimSpace(def.Description)
	if description != "" {
		project.Description = description
		if err := s.projects.UpdateProject(project); err != nil {
			return project, err
		}
	}

	if err := s.ensureGitignore(def.Path, envFile); err != nil {
		return project, err
	}

	return project, nil
}

func (s *FileService) cleanupExisting(def domain.ProjectDefinition, envFile string) error {
	trimmed := strings.TrimSpace(envFile)
	if trimmed == "" {
		trimmed = defaultEnvFile
	}

	paths := []string{
		s.fs.Join(def.Path, ".project.hcl"),
		s.fs.Join(def.Path, trimmed),
		s.fs.Join(def.Path, fmt.Sprintf(".%s.gitconfig", def.Name)),
	}

	for _, p := range paths {
		err := s.fs.Remove(p)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove stale artifact %s: %w", p, err)
		}
	}

	return nil
}

func (s *FileService) ensureGitignore(dir, envFile string) error {
	trimmed := strings.TrimSpace(envFile)
	if trimmed == "" {
		trimmed = defaultEnvFile
	}

	path := s.fs.Join(dir, ".gitignore")

	data, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s.fs.WriteFile(path, []byte(trimmed+"\n"), 0o644)
		}

		return fmt.Errorf("read .gitignore: %w", err)
	}

	contents := string(data)
	if strings.Contains(contents, trimmed) {
		return nil
	}

	if !strings.HasSuffix(contents, "\n") && len(contents) > 0 {
		contents += "\n"
	}

	contents += trimmed + "\n"

	return s.fs.WriteFile(path, []byte(contents), 0o644)
}
