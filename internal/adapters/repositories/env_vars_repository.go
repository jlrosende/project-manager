package repositories

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type EnvVarsRepository struct{}

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

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	keys := make([]string, 0, len(envVars))
	for key := range envVars {
		if strings.TrimSpace(key) == "" {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	var builder strings.Builder

	for _, key := range keys {
		value := envVars[key]
		builder.WriteString(key)
		builder.WriteByte('=')

		if needsQuote(value) {
			builder.WriteString(strconv.Quote(value))
		} else {
			builder.WriteString(value)
		}

		builder.WriteByte('\n')
	}

	return os.WriteFile(path, []byte(builder.String()), 0o600)
}

func needsQuote(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if r <= 32 || r == '#' || r == '=' || r == '"' || r == '\'' {
			return true
		}
	}

	return false
}
