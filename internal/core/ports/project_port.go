package ports

import "github.com/jlrosende/project-manager/internal/core/domain"

type ProjectService interface {
	Get(name string) (*domain.Project, error)
	List() ([]domain.Project, error)
	Create(name, path, subproject string, envVars domain.EnvVars, git domain.GitConfig) (*domain.Project, error)
	Delete(name string) error
}

type ProjectRepository interface {
	Get(name string) (*domain.Project, error)
	List() ([]domain.Project, error)
	Create(name, path, subproject string, envVars domain.EnvVars, git domain.GitConfig) (*domain.Project, error)
	Delete(name string) error
}
