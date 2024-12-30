package services

import (
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type EnvVarsService struct {
	repo ports.EnvVarsRepository
}

var _ ports.EnvVarsService = (*EnvVarsService)(nil)

func NewEnvVarsServiceService(repo ports.EnvVarsRepository) *EnvVarsService {
	return &EnvVarsService{
		repo: repo,
	}
}

func (e *EnvVarsService) Load(path string) (domain.EnvVars, error) {
	return e.repo.Load(path)
}

func (e *EnvVarsService) Save(path string, envVars map[string]string) error {
	return e.repo.Save(path, envVars)
}
