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

func writeProjectHCL(t *testing.T, dir, name string) {
	t.Helper()
	content := []byte("name = \"" + name + "\"\n" +
		"description = \"test project\"\n" +
		"shell = \"/bin/sh\"\n" +
		"env_vars_file = \".env\"\n")
	if err := os.WriteFile(filepath.Join(dir, ".project.hcl"), content, 0o644); err != nil {
		t.Fatalf("write hcl: %v", err)
	}
}

func TestProjectRepository_ListGetUpdateProject(t *testing.T) {
	home := t.TempDir()
	_ = os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	p1 := filepath.Join(home, "p1")
	p2 := filepath.Join(home, "p2")
	_ = os.MkdirAll(p1, 0o755)
	_ = os.MkdirAll(p2, 0o755)
	writeProjectHCL(t, p1, "P1")
	writeProjectHCL(t, p2, "P2")

	gc := config.NewConfig()
	gc.Raw.Section("includeIf").Subsection("gitdir/i:"+p1+"/").SetOption("path", filepath.Join(p1, ".p1.gitconfig"))
	gc.Raw.Section("includeIf").Subsection("gitdir/i:"+p2+"/").SetOption("path", filepath.Join(p2, ".p2.gitconfig"))
	b, _ := gc.Marshal()
	if err := os.WriteFile(filepath.Join(home, ".gitconfig"), b, 0o644); err != nil {
		t.Fatalf("write global git: %v", err)
	}

	repo, err := repositories.NewProjectRepository(nil)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	list, err := repo.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(list))
	}

	got, err := repo.Get("P1")
	if err != nil || got == nil || got.Path != p1 {
		t.Fatalf("get P1 failed: %v %#v", err, got)
	}

	got.Name = "P1-NEW"
	if err := repo.UpdateProject(&domain.Project{Path: p1, Name: "P1-NEW", Shell: "/bin/sh", EnvVarsFile: ".env"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if b, err := os.ReadFile(filepath.Join(p1, ".project.hcl")); err != nil || !contains(string(b), "P1-NEW") {
		t.Fatalf("update not persisted: %v %s", err, string(b))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && (len(s) > 0 && (s[0:len(sub)] == sub || contains(s[1:], sub)))))
}
