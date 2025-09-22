//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5/config"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestProjectRepository_UpdateProject_RenamesPerProjectGitconfigAndUpdatesIncludeIf(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	projDir := filepath.Join(home, "proj")
	_ = os.MkdirAll(projDir, 0o755)

	// write initial .project.hcl with name "p"
	if err := os.WriteFile(
		filepath.Join(projDir, ".project.hcl"),
		[]byte("name = \"p\"\n"+
			"description = \"t\"\n"+
			"shell = \"/bin/sh\"\n"+
			"env_vars_file = \".env\"\n"),
		0o600,
	); err != nil {
		t.Fatalf("write hcl: %v", err)
	}

	// create existing per-project gitconfig for old name
	oldGit := filepath.Join(projDir, ".p.gitconfig")
	if err := os.WriteFile(oldGit, []byte("[user]\n\tname = test\n"), 0o600); err != nil {
		t.Fatalf("write old gitconfig: %v", err)
	}

	// prepare global includeIf pointing to old path
	gc := config.NewConfig()
	gc.Raw.Section("includeIf").Subsection("gitdir/i:" + projDir + "/").SetOption("path", oldGit)
	b, _ := gc.Marshal()
	if err := os.WriteFile(filepath.Join(home, ".gitconfig"), b, 0o600); err != nil {
		t.Fatalf("write global git: %v", err)
	}

	repo, err := repositories.NewProjectRepository()
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	// perform rename to "p2"
	if err := repo.UpdateProject(&domain.Project{Path: projDir, Name: "p2", Shell: "/bin/sh", EnvVarsFile: ".env"}); err != nil {
		t.Fatalf("update project: %v", err)
	}

	// .project.hcl updated
	if b, err := os.ReadFile(filepath.Join(projDir, ".project.hcl")); err != nil || !strings.Contains(string(b), "p2") {
		t.Fatalf("hcl not updated: %v %s", err, string(b))
	}

	// per-project gitconfig renamed
	newGit := filepath.Join(projDir, ".p2.gitconfig")
	if _, err := os.Stat(newGit); err != nil {
		t.Fatalf("new per-project gitconfig missing: %v", err)
	}
	if _, err := os.Stat(oldGit); err == nil {
		t.Fatalf("old per-project gitconfig still exists")
	}

	// includeIf updated to new path
	gc2, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		t.Fatalf("reload global git: %v", err)
	}
	ss := gc2.Raw.Section("includeIf").Subsection("gitdir/i:" + projDir + "/")
	if got := ss.Option("path"); got != newGit {
		t.Fatalf("includeIf path not updated: got %q want %q", got, newGit)
	}
}
