package ports

import "github.com/jlrosende/project-manager/internal/core/domain"

type GitService interface {
	Load(path string) (*domain.GitConfig, error)
	Save(path string, gitConfig *domain.GitConfig) error
}

type GitRepository interface {
	Load(path string) (*domain.GitConfig, error)
	Save(path string, gitConfig *domain.GitConfig) error
}
