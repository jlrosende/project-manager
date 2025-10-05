//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLINew_DryRun(t *testing.T) {
	dir := t.TempDir()
	cmd := cli.Root()
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetArgs([]string{"new", "dry", dir, "--dry-run"})
	if e := cmd.Execute(); e != nil {
		t.Fatalf("execute: %v; stderr=%s", e, err.String())
	}
	if _, e := os.Stat(filepath.Join(dir, ".project.hcl")); !os.IsNotExist(e) {
		t.Fatalf("dry-run should not create files")
	}
}
