package repositories

import (
	"path/filepath"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
	"github.com/joho/godotenv"
)

type EnvVarsRepository struct {
}

var _ ports.EnvVarsRepository = (*EnvVarsRepository)(nil)

func NewEnvVarsRepository() (*EnvVarsRepository, error) {

	return &EnvVarsRepository{}, nil
}

func (e *EnvVarsRepository) Load(path string) (domain.EnvVars, error) {
	return godotenv.Read(path)
}

func (e *EnvVarsRepository) Save(path string, envVars map[string]string) error {
	return godotenv.Write(envVars, filepath.Join(path, ".env"))
}
