package ports

import "github.com/jlrosende/project-manager/internal/core/domain"

type EnvVarsService interface {
	Load(path string) (domain.EnvVars, error)
	Save(path string, envVars map[string]string) error
}

type EnvVarsRepository interface {
	Load(path string) (domain.EnvVars, error)
	Save(path string, envVars map[string]string) error
}
