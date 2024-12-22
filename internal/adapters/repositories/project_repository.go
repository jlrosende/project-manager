package repositories

import (
	"context"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectRepository struct {
	ctx context.Context
}

var _ ports.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

func (project *ProjectRepository) Get(name string) (domain.Project, error) {
	return domain.Project{}, nil
}

func (project *ProjectRepository) List() ([]domain.Project, error) {
	return []domain.Project{}, nil

}

func (project *ProjectRepository) Create() (domain.Project, error) {
	return domain.Project{}, nil

}

func (project *ProjectRepository) Edit(name string) (domain.Project, error) {
	return domain.Project{}, nil

}

func (project *ProjectRepository) Delete(name string) error {
	return nil

}
