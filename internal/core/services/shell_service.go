package services

import (
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ShellService struct {
	repo ports.ShellRepository
}

var _ ports.ShellService = (*ShellService)(nil)

func NewShellService(repo ports.ShellRepository) *ShellService {
	return &ShellService{
		repo: repo,
	}
}

func (s *ShellService) Start(project domain.Project) error {
	return s.repo.Start(project)
}
