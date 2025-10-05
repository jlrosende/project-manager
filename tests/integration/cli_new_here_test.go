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

func TestCLINew_InitHere(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(dir)
	cmd := cli.Root()
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetArgs([]string{"new", "myproj", "--here"})
	if e := cmd.Execute(); e != nil {
		t.Fatalf("execute: %v; stderr=%s", e, err.String())
	}
	if _, e := os.Stat(filepath.Join(dir, ".project.hcl")); e != nil {
		t.Fatalf(".project.hcl not created: %v", e)
	}
}
