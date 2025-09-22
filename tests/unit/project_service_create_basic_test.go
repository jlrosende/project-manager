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

func TestProjectService_Create_Basic(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "proj-basic")
	gitRepo, _ := repositories.NewGitRepository()
	envRepo, _ := repositories.NewEnvVarsRepository()
	projRepo, _ := repositories.NewProjectRepository()
	svc := services.NewProjectService(projRepo, envRepo, gitRepo)

	gitCfg := domain.New(domain.WithName("User"), domain.WithEmail("u@e"))
	proj, err := svc.Create("p", projDir, "", "/bin/sh", ".env", domain.EnvVars{"K": "V"}, gitCfg)
	if err != nil || proj == nil {
		t.Fatalf("Create failed: %v", err)
	}
	if proj.Name != "p" || proj.Path != projDir {
		t.Fatalf("unexpected project: %#v", proj)
	}

	if _, err := os.Stat(filepath.Join(projDir, ".project.hcl")); err != nil {
		t.Fatalf("missing .project.hcl: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, ".env")); err != nil {
		t.Fatalf("missing .env: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".gitconfig")); err != nil {
		t.Fatalf("missing ~/.gitconfig: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, ".p.gitconfig")); err != nil {
		t.Fatalf("missing per-project gitconfig: %v", err)
	}
}
