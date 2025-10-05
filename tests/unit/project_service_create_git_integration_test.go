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

func TestProjectService_Create_WritesGitIncludeAndPerProjectConfig(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "proj")
	gitRepo, _ := repositories.NewGitRepository()
	envRepo, _ := repositories.NewEnvVarsRepository()
	projRepo, _ := repositories.NewProjectRepository()

	svc := services.NewProjectService(projRepo, envRepo, gitRepo, repositories.NewFilesystem(), nil)

	gitCfg := domain.New(domain.WithName("U"), domain.WithEmail("u@e"))
	proj, err := svc.Create("p", projDir, "", "/bin/sh", ".env", domain.EnvVars{"K": "V"}, gitCfg)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	g := filepath.Join(home, ".gitconfig")
	if _, err := os.Stat(g); err != nil {
		t.Fatalf("~/.gitconfig not written: %v", err)
	}

	per := filepath.Join(proj.Path, ".p.gitconfig")
	if _, err := os.Stat(per); err != nil {
		t.Fatalf("per-project gitconfig not written: %v", err)
	}
}
