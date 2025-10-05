//go:build integration
// +build integration

package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCLINew_FailurePaths(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		dir := t.TempDir()
		oldWD, getErr := os.Getwd()
		if getErr != nil {
			t.Fatalf("getwd: %v", getErr)
		}

		if chErr := os.Chdir(dir); chErr != nil {
			t.Fatalf("chdir: %v", chErr)
		}
		t.Cleanup(func() {
			_ = os.Chdir(oldWD)
		})

		_, stderr, err := runNewCommand(t, "new")
		if err == nil {
			t.Fatalf("expected error for missing name; stderr=%s", stderr)
		}

		if _, statErr := os.Stat(filepath.Join(dir, ".project.hcl")); !os.IsNotExist(statErr) {
			t.Fatalf("unexpected .project.hcl created on failure")
		}

		if _, statErr := os.Stat(filepath.Join(dir, ".env")); !os.IsNotExist(statErr) {
			t.Fatalf("unexpected .env created on failure")
		}
	})

	t.Run("conflicting location", func(t *testing.T) {
		dir := t.TempDir()
		oldWD, getErr := os.Getwd()
		if getErr != nil {
			t.Fatalf("getwd: %v", getErr)
		}

		if chErr := os.Chdir(dir); chErr != nil {
			t.Fatalf("chdir: %v", chErr)
		}
		t.Cleanup(func() {
			_ = os.Chdir(oldWD)
		})

		_, stderr, err := runNewCommand(t, "new", "conflict", dir, "--here")
		if err == nil {
			t.Fatalf("expected error for conflicting path and --here; stderr=%s", stderr)
		}

		if _, statErr := os.Stat(filepath.Join(dir, ".project.hcl")); !os.IsNotExist(statErr) {
			t.Fatalf("unexpected .project.hcl created when flags conflict")
		}

		if _, statErr := os.Stat(filepath.Join(dir, ".env")); !os.IsNotExist(statErr) {
			t.Fatalf("unexpected .env created when flags conflict")
		}
	})
}
