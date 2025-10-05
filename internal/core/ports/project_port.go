package ports

//go:generate go tool mockgen -source=project_port.go -destination=../../../mocks/mock_project_port.go -package=mocks

import (
	"context"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

type ProjectService interface {
	Load(name string) (*domain.Project, error)
	List() ([]*domain.Project, error)
	Create(
		name, path, subproject, shell, envFile string,
		envVars domain.EnvVars,
		git *domain.GitConfig,
	) (*domain.Project, error)
	AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error
	UpdateProject(project *domain.Project) error
	UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error
	Delete(name string) error
	DeleteProject(ctx context.Context, options domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error)
}

type ProjectRepository interface {
	List() ([]*domain.Project, error)
	Create(
		name, path, subproject, shell, envFile string,
		envVars domain.EnvVars,
		git *domain.GitConfig,
	) (*domain.Project, error)
	AddEnvironment(projectName string, env *domain.Environment, envVars domain.EnvVars) error
	UpdateProject(project *domain.Project) error
	UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error
	Delete(name string) error
	ResolveIdentifier(ctx context.Context, lookup domain.ProjectIdentifier) (domain.ProjectIdentifier, error)
	FinalizeDeletion(ctx context.Context, identifier domain.ProjectIdentifier, scope domain.DeleteScope) error
}
