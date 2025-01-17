package services

import (
	"fmt"
	"log/slog"

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

	// load env vars
	// project.EnvVars, err = svc.envVars.Load(project.Path)

	// if err != nil {
	// 	return nil, err
	// }

	// // load git
	// project.GitConfig, err = svc.git.Load(project.Path)

	// if err != nil {
	// 	return nil, err
	// }

	slog.Info(fmt.Sprintf("%+v", project))

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
