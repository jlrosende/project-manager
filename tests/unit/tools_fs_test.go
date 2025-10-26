package unit

import (
	"os"
	"path/filepath"
	"testing"

	repositories "github.com/jlrosende/project-manager/internal/adapters/repositories"
)

func TestIsDirEmpty(t *testing.T) {
	dir := t.TempDir()

	fsys := repositories.NewFilesystem()

	empty, err := fsys.IsDirEmpty(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !empty {
		t.Fatalf("expected empty dir")
	}

	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("x"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	empty, err = fsys.IsDirEmpty(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if empty {
		t.Fatalf("expected non-empty dir")
	}
}

func TestEnsureDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b")

	fsys := repositories.NewFilesystem()

	if err := fsys.EnsureDir(dir, 0o755); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}

	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		t.Fatalf("dir not created: %v", err)
	}
}

func TestRenameAndWriteFile(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "old", "f.txt")
	newp := filepath.Join(dir, "new", "f.txt")

	fsys := repositories.NewFilesystem()

	if err := fsys.WriteFile(old, []byte("hi"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := fsys.Rename(old, newp); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if _, err := os.Stat(newp); err != nil {
		t.Fatalf("renamed file missing: %v", err)
	}
}
