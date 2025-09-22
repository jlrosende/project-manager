//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/config"
	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestProjectRepository_UpdateEnvironment_RenameRelativeAbsolute(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	dir := filepath.Join(home, "proj")
	_ = os.MkdirAll(dir, 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".project.hcl"), []byte("name=\"P\"\ndescription=\"d\"\nshell=\"/bin/sh\"\nenv_vars_file=\".env\"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".old.env"), []byte("K=V\n"), 0o600)

	gc := config.NewConfig()
	gc.Raw.Section("includeIf").Subsection("gitdir/i:" + dir + "/").SetOption("path", filepath.Join(dir, ".p.gitconfig"))
	b, _ := gc.Marshal()
	_ = os.WriteFile(filepath.Join(home, ".gitconfig"), b, 0o644)

	repo, _ := repositories.NewProjectRepository()
	_ = repo.AddEnvironment("P", &domain.Environment{Name: "Dev", EnvVarsMode: domain.EnvVarsModeMerge, EnvVarsFile: ".old.env"}, nil)

	if err := repo.UpdateEnvironment("P", "Dev", &domain.Environment{Name: "Dev", EnvVarsMode: domain.EnvVarsModeMerge, EnvVarsFile: ".new.env"}); err != nil {
		t.Fatalf("update env: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".new.env")); err != nil {
		t.Fatalf("new env path missing: %v", err)
	}
}
