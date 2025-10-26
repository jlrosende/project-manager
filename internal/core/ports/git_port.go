package ports

//go:generate go tool mockgen -source=git_port.go -destination=../../../mocks/mock_git_port.go -package=mocks

import (
	"context"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

type GitService interface {
	Load(path string) (*domain.GitConfig, error)
	Save(path string, gitConfig *domain.GitConfig) error
	LoadGlobal() error
	UpdateIncludeIf(gitdir, perProjectPath, subproject string) error
	SaveGlobal(home string) error
	RemoveHooks(ctx context.Context, project domain.ProjectIdentifier) error
}

type GitRepository interface {
	Load(path string) (*domain.GitConfig, error)
	Save(path string, gitConfig *domain.GitConfig) error
	LoadGlobal() error
	UpdateIncludeIf(gitdir, perProjectPath, subproject string) error
	SaveGlobal(home string) error
	RemoveHooks(ctx context.Context, project domain.ProjectIdentifier) error
}
