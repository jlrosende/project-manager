//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/repositories"
)

func TestEnvVarsRepository_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	_ = os.Setenv("HOME", dir)
	defer os.Unsetenv("HOME")

	repo, err := repositories.NewEnvVarsRepository()
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	vars := map[string]string{"A": "1", "B": "2"}
	p := "~/.test.env"
	if err := repo.Save(p, vars); err != nil {
		t.Fatalf("save: %v", err)
	}

	abs := filepath.Join(dir, ".test.env")
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("env file not created: %v", err)
	}

	loaded, err := repo.Load(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if len(loaded) != len(vars) || loaded["A"] != "1" || loaded["B"] != "2" {
		t.Fatalf("loaded mismatch: %#v", loaded)
	}
}
