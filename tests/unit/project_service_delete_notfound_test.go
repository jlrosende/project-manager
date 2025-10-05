//go:build unit
// +build unit

package unit_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
)

func TestProjectServiceDeleteProjectMissingTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projectRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)
	fs := mocks.NewMockFilesystem(ctrl)

	svc := services.NewProjectService(projectRepo, envRepo, gitRepo, fs, nil)

	opts := domain.ProjectDeleteOptions{
		Target: domain.ProjectIdentifier{Name: "ghost"},
	}

	notFoundErr := domain.ErrProjectNotFound

	projectRepo.EXPECT().ResolveIdentifier(gomock.Any(), opts.Target).Return(domain.ProjectIdentifier{}, notFoundErr)
	envRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)
	fs.EXPECT().PlanDeletion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	gitRepo.EXPECT().RemoveHooks(gomock.Any(), gomock.Any()).Times(0)

	result, err := svc.DeleteProject(context.Background(), opts)
	if !errors.Is(err, notFoundErr) {
		t.Fatalf("expected not found error, got %v", err)
	}

	if result != nil {
		t.Fatalf("expected nil result on not-found, got %+v", result)
	}
}
