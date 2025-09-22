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
)

func TestGitRepository_SaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".proj.gitconfig")

	repo, err := repositories.NewGitRepository()
	if err != nil {
		t.Fatalf("new git repo: %v", err)
	}

	cfg := domain.New(
		domain.WithName("Test User"),
		domain.WithEmail("test@example.com"),
		domain.WithSigningKey("ABCDEF"),
		domain.WithCommitSign(true),
		domain.WithTagSign(true),
	)

	if err := repo.Save(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := repo.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if loaded.User.Name != "Test User" || loaded.User.Email != "test@example.com" || loaded.User.SigningKey != "ABCDEF" {
		t.Fatalf("user mismatch: %#v", loaded.User)
	}
	if !loaded.Commit.GPGSign || !loaded.Tag.GPGSign {
		t.Fatalf("gpg flags mismatch: %+v", loaded)
	}
}

func TestGitRepository_GlobalIncludeIf(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	repo, err := repositories.NewGitRepository()
	if err != nil {
		t.Fatalf("new git repo: %v", err)
	}

	if err := repo.LoadGlobal(); err != nil {
		t.Fatalf("load global: %v", err)
	}

	projDir := filepath.Join(home, "work", "p1")
	perProj := filepath.Join(projDir, ".p1.gitconfig")
	gitdir := "gitdir/i:" + projDir + "/"

	if err := repo.UpdateIncludeIf(gitdir, perProj, ""); err != nil {
		t.Fatalf("update includeIf: %v", err)
	}

	if err := repo.SaveGlobal(home); err != nil {
		t.Fatalf("save global: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(home, ".gitconfig"))
	if err != nil {
		t.Fatalf("read ~/.gitconfig: %v", err)
	}

	s := string(b)
	if !strings.Contains(s, gitdir) || !strings.Contains(s, perProj) {
		t.Fatalf("includeIf not written correctly: %s", s)
	}
}
