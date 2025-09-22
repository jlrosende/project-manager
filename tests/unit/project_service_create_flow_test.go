//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestService_Create_SuccessPath(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "p-success")
	gitRepo, _ := repositories.NewGitRepository()
	envRepo, _ := repositories.NewEnvVarsRepository()
	projRepo, _ := repositories.NewProjectRepository()
	svc := services.NewProjectService(projRepo, envRepo, gitRepo)

	gitCfg := domain.New(domain.WithName("U"), domain.WithEmail("u@e"))
	proj, err := svc.Create("p", projDir, "", "/bin/sh", ".env", domain.EnvVars{"K": "V"}, gitCfg)
	if err != nil || proj == nil {
		t.Fatalf("Create failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, ".gitconfig")); err != nil {
		t.Fatalf("global git missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, ".env")); err != nil {
		t.Fatalf("env file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, ".p.gitconfig")); err != nil {
		t.Fatalf("per-project git missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, ".project.hcl")); err != nil {
		t.Fatalf("project hcl missing: %v", err)
	}
}
