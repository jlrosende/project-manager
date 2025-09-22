//go:build unit
// +build unit

package unit_test

import (
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
	"go.uber.org/mock/gomock"
)

func TestProjectService_AddEnvironment_DefaultsAndIdempotency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)

	svc := services.NewProjectService(projRepo, envRepo, gitRepo)

	proj := &domain.Project{Name: "Foo", Path: t.TempDir(), EnvVarsFile: ".env"}
	projRepo.EXPECT().List().Return([]*domain.Project{proj}, nil).AnyTimes()
	projRepo.EXPECT().AddEnvironment("Foo", gomock.Any(), gomock.Any()).Return(nil)

	envRepo.EXPECT().Load(gomock.Any()).Return(domain.EnvVars{}, nil).AnyTimes()
	envRepo.EXPECT().Save(filepath.Join(proj.Path, ".dev.env"), gomock.Any()).Return(nil)

	env := &domain.Environment{Name: "Dev"}
	err := svc.AddEnvironment("Foo", env, domain.EnvVars{"K": "V"})
	if err != nil {
		t.Fatalf("AddEnvironment error: %v", err)
	}

	if env.EnvVarsFile == "" {
		t.Fatalf("expected default env file to be set")
	}
}
