package ports

//go:generate go tool mockgen -source=env_vars_port.go -destination=../../../mocks/mock_env_vars_port.go -package=mocks

import "github.com/jlrosende/project-manager/internal/core/domain"

type EnvVarsService interface {
	Load(path string) (domain.EnvVars, error)
	Save(path string, envVars map[string]string) error
}

type EnvVarsRepository interface {
	Load(path string) (domain.EnvVars, error)
	Save(path string, envVars map[string]string) error
}
