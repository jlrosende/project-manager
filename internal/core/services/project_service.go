package services

import (
	"log/slog"
	"path/filepath"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectService struct {
	git     ports.GitRepository
	envVars ports.EnvVarsRepository
	project ports.ProjectRepository
}

var _ ports.ProjectService = (*ProjectService)(nil)

func NewProjectService(project ports.ProjectRepository, envVars ports.EnvVarsRepository, git ports.GitRepository) *ProjectService {
	return &ProjectService{
		project: project,
		envVars: envVars,
		git:     git,
	}
}

func (svc *ProjectService) Get(name string) (*domain.Project, error) {

	// Get Project
	project, err := svc.project.Get(name)

	if err != nil {
		return nil, err
	}

	var envVarsPath string

	if filepath.IsAbs(project.EnvVarsFile) {
		envVarsPath = project.EnvVarsFile
	} else {
		envVarsPath = filepath.Join(project.Path, project.EnvVarsFile)
	}

	// load env vars
	project.EnvVars, err = svc.envVars.Load(envVarsPath)

	if err != nil {
		return nil, err
	}

	// Load environments EnvVars
	for _, env := range project.Environments {
		var envVarsPath string

		if filepath.IsAbs(project.EnvVarsFile) {
			envVarsPath = env.EnvVarsFile
		} else {
			envVarsPath = filepath.Join(project.Path, env.EnvVarsFile)
		}

		env.EnvVars, err = svc.envVars.Load(envVarsPath)

		if err != nil {
			slog.Warn("EnvVars file not found", slog.String("env", env.Name), slog.String("path", envVarsPath))
		}
	}

	slog.Debug("project service get", slog.Any("project", project))

	return project, nil
}

func (svc *ProjectService) List() ([]*domain.Project, error) {
	return svc.project.List()
}

func (svc *ProjectService) Create(name, path, subproject string, envVars domain.EnvVars, gitConfig *domain.GitConfig) (*domain.Project, error) {
	return svc.project.Create(name, path, subproject, envVars, gitConfig)
}

func (svc *ProjectService) Delete(name string) error {
	return svc.project.Delete(name)
}
