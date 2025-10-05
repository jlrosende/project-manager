//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
)

func TestCLIDelete_LogsArtifacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	projectRepo := mocks.NewMockProjectRepository(ctrl)
	envRepo := mocks.NewMockEnvVarsRepository(ctrl)
	gitRepo := mocks.NewMockGitRepository(ctrl)
	fs := mocks.NewMockFilesystem(ctrl)

	var logBuffer bytes.Buffer
	handler := slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := repositories.NewLogger(slog.New(handler))

	svc := services.NewProjectService(projectRepo, envRepo, gitRepo, fs, logger)

	identifier := domain.ProjectIdentifier{
		Name: "demo-app",
		Path: "/workspace/demo-app",
	}

	plan := &domain.ProjectDeletePlan{
		Scope: domain.DeleteScopeAll,
		Artifacts: []domain.DeletionArtifact{
			{Type: domain.ArtifactEnvVars, Path: "/workspace/demo-app/.env"},
			{Type: domain.ArtifactWorkspace, Path: "/workspace/demo-app"},
		},
	}

	options := domain.ProjectDeleteOptions{
		Target: domain.ProjectIdentifier{Name: "demo-app"},
		Scope:  domain.DeleteScopeAll,
	}

	projectRepo.EXPECT().ResolveIdentifier(gomock.Any(), options.Target).Return(identifier, nil)
	fs.EXPECT().PlanDeletion(gomock.Any(), identifier, options.Scope, options.Backup).Return(plan, nil)
	gitRepo.EXPECT().RemoveHooks(gomock.Any(), identifier).Return(nil)
	envRepo.EXPECT().Delete(gomock.Any(), plan.Artifacts[0].Path).Return(nil)
	fs.EXPECT().ExecuteDeletion(gomock.Any(), plan).Return([]domain.DeletionArtifact{plan.Artifacts[1]}, nil)
	projectRepo.EXPECT().FinalizeDeletion(gomock.Any(), identifier, plan.Scope).Return(nil)

	result, err := svc.DeleteProject(context.Background(), options)
	if err != nil {
		t.Fatalf("DeleteProject error: %v", err)
	}

	if result == nil {
		t.Fatal("DeleteProject returned nil result")
	}

	logs := strings.TrimSpace(logBuffer.String())
	if logs == "" {
		t.Fatal("expected logs to be captured, but buffer is empty")
	}

	var artifactLines []string
	for _, line := range strings.Split(logs, "\n") {
		if strings.Contains(line, "msg=\"delete artifact\"") {
			artifactLines = append(artifactLines, line)
		}
	}

	if len(artifactLines) != len(plan.Artifacts) {
		t.Fatalf("expected %d artifact log entries, got %d: %s", len(plan.Artifacts), len(artifactLines), logs)
	}

	for _, artifact := range plan.Artifacts {
		matched := false
		for _, line := range artifactLines {
			if strings.Contains(line, "type="+string(artifact.Type)) && strings.Contains(line, artifact.Path) {
				matched = true
				break
			}
		}
		if !matched {
			t.Fatalf("missing log entry for artifact %s", artifact.Path)
		}
	}
}
