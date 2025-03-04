package repositories

import (
	"os"
	"path/filepath"
	"strings"

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

	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
	}

	return godotenv.Read(path)
}

func (e *EnvVarsRepository) Save(path string, envVars map[string]string) error {

	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
	}

	return godotenv.Write(envVars, filepath.Join(path, ".env"))
}
