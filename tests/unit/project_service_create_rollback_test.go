//go:build unit
// +build unit

package unit_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

type failingGitRepo struct{ repositories.GitRepository }

func (f failingGitRepo) Load(string) (*domain.GitConfig, error)       { return &domain.GitConfig{}, nil }
func (f failingGitRepo) Save(string, *domain.GitConfig) error         { return nil }
func (f failingGitRepo) LoadGlobal() error                            { return nil }
func (f failingGitRepo) UpdateIncludeIf(string, string, string) error { return errors.New("boom") }
func (f failingGitRepo) SaveGlobal(string) error                      { return nil }
func (f failingGitRepo) RemoveHooks(context.Context, domain.ProjectIdentifier) error {
	return nil
}

func TestService_Create_RollbackOnGitFailure(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "p-rollback")
	gitRepo := failingGitRepo{}
	envRepo, _ := repositories.NewEnvVarsRepository()
	projRepo, _ := repositories.NewProjectRepository()
	svc := services.NewProjectService(projRepo, envRepo, gitRepo, repositories.NewFilesystem(), nil)

	_, err := svc.Create("p", projDir, "", "/bin/sh", ".env", domain.EnvVars{"K": "V"}, domain.New())
	if err == nil {
		t.Fatalf("expected failure, got nil")
	}
	if _, err := os.Stat(filepath.Join(projDir, ".env")); !os.IsNotExist(err) {
		t.Fatalf("env file should be removed on rollback")
	}
	if _, err := os.Stat(filepath.Join(projDir, ".project.hcl")); !os.IsNotExist(err) {
		t.Fatalf("project hcl should be removed on rollback")
	}
}
