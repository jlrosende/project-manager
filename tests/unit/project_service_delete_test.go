//go:build unit
// +build unit

package unit_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
)

func TestProjectServiceDeleteProject_DryRunReturnsPlanWithoutSideEffects(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projectRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)
	fs := mocks.NewMockFilesystem(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	svc := services.NewProjectService(projectRepo, envRepo, gitRepo, fs, logger)

	options := domain.ProjectDeleteOptions{
		Target: domain.ProjectIdentifier{
			Name:       "sample-app",
			Path:       "/workspace/sample-app",
			RegistryID: "sample-app-id",
			Status: domain.ProjectStatus{
				Locked:   false,
				ReadOnly: false,
			},
		},
		Scope:  domain.DeleteScopeAll,
		DryRun: true,
	}

	plan := &domain.ProjectDeletePlan{
		Scope:     options.Scope,
		Artifacts: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry, Path: "registry://sample-app"}},
	}

	projectRepo.EXPECT().ResolveIdentifier(gomock.Any(), options.Target).Return(options.Target, nil)
	fs.EXPECT().PlanDeletion(gomock.Any(), options.Target, options.Scope, options.Backup).Return(plan, nil)
	envRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)
	gitRepo.EXPECT().RemoveHooks(gomock.Any(), gomock.Any()).Times(0)
	fs.EXPECT().CreateBackup(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	fs.EXPECT().ExecuteDeletion(gomock.Any(), gomock.Any()).Times(0)
	projectRepo.EXPECT().FinalizeDeletion(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	result, err := svc.DeleteProject(context.Background(), options)
	if err != nil {
		t.Fatalf("DeleteProject dry-run unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("DeleteProject dry-run returned nil result")
	}

	if result.Plan == nil {
		t.Fatal("DeleteProject dry-run should include plan details")
	}

	if result.Plan.Scope != domain.DeleteScopeAll {
		t.Fatalf("unexpected plan scope: %v", result.Plan.Scope)
	}

	if result.BackupPath != "" {
		t.Fatalf("dry-run should not create backup, got %q", result.BackupPath)
	}

	if len(result.ArtifactsRemoved) != 0 {
		t.Fatalf("dry-run should not remove artifacts, got %d", len(result.ArtifactsRemoved))
	}
}

func TestProjectServiceDeleteProject_AllScopeInvokesPorts(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projectRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)
	fs := mocks.NewMockFilesystem(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	svc := services.NewProjectService(projectRepo, envRepo, gitRepo, fs, logger)

	options := domain.ProjectDeleteOptions{
		Target: domain.ProjectIdentifier{
			Name:       "sample-app",
			Path:       "/workspace/sample-app",
			RegistryID: "sample-app-id",
			Status:     domain.ProjectStatus{},
		},
		Scope:  domain.DeleteScopeAll,
		Backup: &domain.BackupRequest{Destination: "/backups/sample-app.zip", IncludeWorkspace: true},
	}

	plan := &domain.ProjectDeletePlan{
		Scope: options.Scope,
		Artifacts: []domain.DeletionArtifact{
			{Type: domain.ArtifactBackup, Path: options.Target.Path},
			{Type: domain.ArtifactRegistry, Path: "registry://sample-app"},
			{Type: domain.ArtifactEnvVars, Path: options.Target.Path + "/.env"},
			{Type: domain.ArtifactWorkspace, Path: options.Target.Path},
		},
	}

	projectRepo.EXPECT().ResolveIdentifier(gomock.Any(), options.Target).Return(options.Target, nil)
	fs.EXPECT().PlanDeletion(gomock.Any(), options.Target, options.Scope, options.Backup).Return(plan, nil)
	fs.EXPECT().CreateBackup(gomock.Any(), options.Target, options.Backup).Return(&domain.BackupArtifact{FinalPath: options.Backup.Destination, Created: true}, nil)
	envRepo.EXPECT().Delete(gomock.Any(), plan.Artifacts[2].Path).Return(nil)
	gitRepo.EXPECT().RemoveHooks(gomock.Any(), options.Target).Return(nil)
	fs.EXPECT().ExecuteDeletion(gomock.Any(), plan).Return([]domain.DeletionArtifact{plan.Artifacts[1], plan.Artifacts[3]}, nil)
	projectRepo.EXPECT().FinalizeDeletion(gomock.Any(), options.Target, plan.Scope).Return(nil)
	logger.EXPECT().Info("delete artifact", gomock.Any(), gomock.Any()).Times(len(plan.Artifacts))

	result, err := svc.DeleteProject(context.Background(), options)
	if err != nil {
		t.Fatalf("DeleteProject all scope error: %v", err)
	}

	if result == nil {
		t.Fatal("DeleteProject result nil")
	}

	if result.BackupPath != options.Backup.Destination {
		t.Fatalf("expected backup path %q, got %q", options.Backup.Destination, result.BackupPath)
	}

	if len(result.ArtifactsRemoved) != 3 {
		t.Fatalf("expected 3 removed artifacts, got %d", len(result.ArtifactsRemoved))
	}
}
