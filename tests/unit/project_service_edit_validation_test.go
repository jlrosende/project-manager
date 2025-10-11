//go:build unit
// +build unit

package unit_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
)

func TestProjectEdit_ProjectValidationErrorPreventsPersistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projectRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)
	fs := mocks.NewMockFilesystem(ctrl)

	svc := services.NewProjectService(projectRepo, envRepo, gitRepo, fs, nil)

	identifier := domain.ProjectIdentifier{Name: "demo", Path: "/workspace/demo"}
	project := &domain.Project{
		Name:        "demo",
		Path:        identifier.Path,
		Description: "Sample",
		Shell:       "/bin/bash",
		EnvVarsFile: ".env",
		Environments: []*domain.Environment{
			{Name: "staging", EnvVarsMode: domain.EnvVarsModeMerge, EnvVarsFile: ".env.staging"},
		},
	}

	projectRepo.EXPECT().ResolveIdentifier(gomock.Any(), gomock.Any()).Return(identifier, nil)
	projectRepo.EXPECT().LoadProjectDefinition(gomock.Any(), identifier).Return(project, nil)
	projectRepo.EXPECT().AcquireEditLock(gomock.Any(), gomock.Any()).Times(0)
	projectRepo.EXPECT().ApplyEditChangeSet(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	result, err := svc.ProjectEdit(context.Background(), services.ProjectEditOptions{
		Name: identifier.Name,
		Flags: services.ProjectEditFlags{
			EnvVarsFileSet: true,
			EnvVarsFile:    "",
		},
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	var validationErr *services.ProjectEditValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ProjectEditValidationError, got %T", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result describing validation failure")
	}

	if len(result.Errors) == 0 {
		t.Fatal("expected validation errors to be reported")
	}

	var aggregated domain.ProjectValidationErrors
	if !errors.As(err, &aggregated) {
		t.Fatal("expected underlying ProjectValidationErrors")
	}

	found := false
	for _, e := range aggregated {
		if e.Field == domain.ProjectFieldEnvVarsFile {
			if !strings.Contains(e.Message, "cannot be empty") {
				t.Fatalf("unexpected message for env vars file: %s", e.Message)
			}
			found = true
		}
	}

	if !found {
		t.Fatal("expected env vars file validation error")
	}
}
