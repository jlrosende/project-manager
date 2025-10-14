//go:build unit
// +build unit

package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestEnsureProjectUniqueness_NameConflict(t *testing.T) {
	dir := t.TempDir()
	writeProjectFile(t, dir)

	existing := &domain.Project{Name: "demo", Path: dir}

	candidate := domain.ProjectDefinition{
		Name: "demo",
		Path: filepath.Join(t.TempDir(), "other"),
	}

	err := services.EnsureProjectUniqueness([]*domain.Project{existing}, candidate)
	if err == nil {
		t.Fatal("expected name conflict error")
	}

	conflict, ok := err.(*services.ProjectConflictError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}

	if conflict.Code() != "NEW-CONFLICT-NAME" {
		t.Fatalf("unexpected conflict code: %s", conflict.Code())
	}

	if conflict.Field() != "name" {
		t.Fatalf("unexpected conflict field: %s", conflict.Field())
	}

	if conflict.ExitCode() != 3 {
		t.Fatalf("unexpected exit code: %d", conflict.ExitCode())
	}
}

func TestEnsureProjectUniqueness_PathConflict(t *testing.T) {
	dir := t.TempDir()
	writeProjectFile(t, dir)

	existing := &domain.Project{Name: "demo", Path: dir}

	candidate := domain.ProjectDefinition{
		Name: "other",
		Path: dir,
	}

	err := services.EnsureProjectUniqueness([]*domain.Project{existing}, candidate)
	if err == nil {
		t.Fatal("expected path conflict error")
	}

	conflict, ok := err.(*services.ProjectConflictError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}

	if conflict.Code() != "NEW-CONFLICT-PATH" {
		t.Fatalf("unexpected conflict code: %s", conflict.Code())
	}

	if conflict.Field() != "path" {
		t.Fatalf("unexpected conflict field: %s", conflict.Field())
	}
}

func TestEnsureProjectUniqueness_IgnoresMissingMetadata(t *testing.T) {
	dir := t.TempDir()
	existing := &domain.Project{Name: "demo", Path: dir}

	candidate := domain.ProjectDefinition{
		Name: "demo",
		Path: dir,
	}

	if err := services.EnsureProjectUniqueness([]*domain.Project{existing}, candidate); err != nil {
		t.Fatalf("expected stale registry entry to be ignored: %v", err)
	}
}

func writeProjectFile(t *testing.T, dir string) {
	t.Helper()

	path := filepath.Join(dir, ".project.hcl")
	if err := os.WriteFile(path, []byte("project { name = \"demo\" }\n"), 0o600); err != nil {
		t.Fatalf("write project file: %v", err)
	}
}
