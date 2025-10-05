package services

import (
	"context"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type GitService struct {
	repo ports.GitRepository
}

var _ ports.GitService = (*GitService)(nil)

func NewGitService(repo ports.GitRepository) *GitService {
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

func (g *GitService) LoadGlobal() error { return g.repo.LoadGlobal() }
func (g *GitService) UpdateIncludeIf(gitdir, perProjectPath, subproject string) error {
	return g.repo.UpdateIncludeIf(gitdir, perProjectPath, subproject)
}
func (g *GitService) SaveGlobal(home string) error { return g.repo.SaveGlobal(home) }
func (g *GitService) RemoveHooks(ctx context.Context, project domain.ProjectIdentifier) error {
	return g.repo.RemoveHooks(ctx, project)
}
