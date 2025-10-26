//go:build contract
// +build contract

package contract_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLINew_FlagsPrecedence(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "project.yaml")
	_ = os.WriteFile(cfg, []byte("name: cfgname\npath: "+dir+"\n"), 0o600)
	cmd := cli.Root()
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetArgs([]string{"new", "override", dir, "--cli-input", cfg})
	e := cmd.Execute()
	if e != nil {
		t.Fatalf("execute: %v; stderr=%s", e, err.String())
	}
	hcl := filepath.Join(dir, ".project.hcl")
	b, r := os.ReadFile(hcl)
	if r != nil {
		t.Fatalf("read .project.hcl: %v", r)
	}
	if !strings.Contains(string(b), "name = \"override\"") {
		t.Fatalf("flags did not override config; got: %s", string(b))
	}

	cmd = cli.Root()
	out.Reset()
	err.Reset()
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetArgs([]string{"new", "n", dir, "--here"})
	e = cmd.Execute()
	if e == nil {
		t.Fatalf("expected error on conflicting flags")
	}
}
