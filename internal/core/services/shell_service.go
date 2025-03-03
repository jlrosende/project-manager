package services

import (
	"os"

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

func (s *ShellService) Start() (*os.Process, error) {
	return s.repo.Start()
}

func (s *ShellService) Wait() (int, error) {
	return s.repo.Wait()
}

func (s *ShellService) Kill() error {
	return s.repo.Kill()
}
