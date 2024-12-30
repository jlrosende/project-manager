package services

import (
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type GitService struct {
	repo ports.GitRepository
}

var _ ports.GitService = (*GitService)(nil)

func NewGiService(repo ports.GitRepository) *GitService {
	return &GitService{
		repo: repo,
	}
}

func (g *GitService) Load(path string) (*domain.GitConfig, error) {
	return g.repo.Load(path)
}

func (g *GitService) Save(path string, gitConfig *domain.GitConfig) error {
	return g.repo.Save(path, gitConfig)
}
