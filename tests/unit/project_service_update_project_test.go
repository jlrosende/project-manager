//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestProjectService_UpdateProject_ChangesPersist(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "proj-upd")
	gitRepo, _ := repositories.NewGitRepository()
	envRepo, _ := repositories.NewEnvVarsRepository()
	projRepo, _ := repositories.NewProjectRepository(nil)
	svc := services.NewProjectService(projRepo, envRepo, gitRepo, repositories.NewFilesystem(), nil)

	proj, err := svc.Create("p", projDir, "", "/bin/sh", ".env", nil, domain.New())
	if err != nil || proj == nil {
		t.Fatalf("Create failed: %v", err)
	}

	proj.Name = "p2"
	if err := svc.UpdateProject(proj); err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(projDir, ".project.hcl"))
	if err != nil {
		t.Fatalf("read hcl: %v", err)
	}
	if !strings.Contains(string(b), "p2") {
		t.Fatalf("update not persisted: %s", string(b))
	}
}
