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

// CreateOptions configures filesystem creation behaviour.
type CreateOptions struct {
	Force bool
}

// FileService coordinates filesystem updates required for project creation.
type FileService struct {
	projects ports.ProjectService
	fs       ports.Filesystem
}

// NewFileService constructs a FileService with the required collaborators.
func NewFileService(projects ports.ProjectService, fs ports.Filesystem) *FileService {
	return &FileService{projects: projects, fs: fs}
}

// Create writes project artifacts (.project.hcl, .env) using the provided definition and environment variables.
func (s *FileService) Create(
	def domain.ProjectDefinition,
	envVars domain.EnvVars,
	opts CreateOptions,
) (*domain.Project, error) {
	if s.projects == nil {
		return nil, errors.New("project service is not configured")
	}

	if s.fs == nil {
		return nil, errors.New("filesystem is not configured")
	}

	if opts.Force {
		if err := s.cleanupExisting(def); err != nil {
			return nil, err
		}
	}

	project, err := s.projects.Create(def.Name, def.Path, "", "", defaultEnvFile, envVars, nil)
	if err != nil {
		return nil, err
	}

	if err := s.ensureGitignore(def.Path); err != nil {
		return project, err
	}

	return project, nil
}

func (s *FileService) cleanupExisting(def domain.ProjectDefinition) error {
	paths := []string{
		s.fs.Join(def.Path, ".project.hcl"),
		s.fs.Join(def.Path, defaultEnvFile),
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

func (s *FileService) ensureGitignore(dir string) error {
	path := s.fs.Join(dir, ".gitignore")

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s.fs.WriteFile(path, []byte(defaultEnvFile+"\n"), 0o644)
		}

		return fmt.Errorf("read .gitignore: %w", err)
	}

	contents := string(data)
	if strings.Contains(contents, defaultEnvFile) {
		return nil
	}

	if !strings.HasSuffix(contents, "\n") && len(contents) > 0 {
		contents += "\n"
	}

	contents += defaultEnvFile + "\n"

	return s.fs.WriteFile(path, []byte(contents), 0o644)
}
