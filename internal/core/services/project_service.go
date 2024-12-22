package services

import (
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectService struct {
	repo ports.ProjectRepository
}

func NewProjectService(repo ports.ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (p *ProjectService) Get(name string) (domain.Project, error) {
	return p.repo.Get(name)
}

func (p *ProjectService) List() ([]domain.Project, error) {
	return p.repo.List()
}

func (p *ProjectService) Create() (domain.Project, error) {
	return p.repo.Create()
}

func (p *ProjectService) Edit(name string) (domain.Project, error) {
	return p.repo.Edit(name)
}

func (p *ProjectService) Delete(name string) error {
	return p.repo.Delete(name)
}
