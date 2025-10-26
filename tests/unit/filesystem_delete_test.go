//go:build unit
// +build unit

package unit_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestFilesystemBackupPlannerDryRunSkipsArchiveCreation(t *testing.T) {
	repo := repositories.NewFilesystem()

	targetDir := t.TempDir()

	req := &domain.BackupRequest{
		Destination: filepath.Join(t.TempDir(), "sample-app.zip"),
	}

	identifier := domain.ProjectIdentifier{
		Name: "sample-app",
		Path: targetDir,
	}

	artifact, err := repo.PlanBackup(context.Background(), identifier, req, true)
	if err != nil {
		t.Fatalf("PlanBackup dry-run error: %v", err)
	}

	if artifact == nil {
		t.Fatal("PlanBackup returned nil artifact in dry-run")
	}

	if artifact.Created {
		t.Fatalf("dry-run should not mark backup as created: %+v", artifact)
	}

	if _, err := os.Stat(req.Destination); err == nil || !os.IsNotExist(err) {
		t.Fatalf("dry-run should not create archive at %q", req.Destination)
	}
}
