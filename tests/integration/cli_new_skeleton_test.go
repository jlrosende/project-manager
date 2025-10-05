//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLINew_GenerateSkeletonJSON(t *testing.T) {
	dir := t.TempDir()
	skeleton := filepath.Join(dir, "project.json")

	cmd := cli.Root()
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.SetArgs([]string{"new", "--generate-cli-skeleton-json", skeleton})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("generate CLI skeleton JSON: %v", err)
	}

	if _, err := os.Stat(skeleton); err != nil {
		t.Fatalf("expected skeleton file: %v", err)
	}
}

func TestCLINew_GenerateSkeletonYAML(t *testing.T) {
	dir := t.TempDir()
	skeleton := filepath.Join(dir, "project.yaml")

	cmd := cli.Root()
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.SetArgs([]string{"new", "--generate-cli-skeleton-yaml", skeleton})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("generate CLI skeleton YAML: %v", err)
	}

	if _, err := os.Stat(skeleton); err != nil {
		t.Fatalf("expected skeleton file: %v", err)
	}
}

func TestCLINew_GenerateSkeletonJSONStdout(t *testing.T) {
	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{"new", "--generate-cli-skeleton-json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("generate CLI skeleton JSON stdout: %v", err)
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	if !strings.Contains(out.String(), "\"name\"") {
		t.Fatalf("stdout missing skeleton content: %s", out.String())
	}
}
