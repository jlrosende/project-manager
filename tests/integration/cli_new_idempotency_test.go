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

func runNewCommand(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)

	executeErr := cmd.Execute()

	return out.String(), errBuf.String(), executeErr
}

func TestCLINew_IdempotentReRun(t *testing.T) {
	dir := t.TempDir()
	args := []string{"new", "demo", dir}

	_, stderr, err := runNewCommand(t, args...)
	if err != nil {
		t.Fatalf("first run failed: %v; stderr=%s", err, stderr)
	}

	hclPath := filepath.Join(dir, ".project.hcl")
	envPath := filepath.Join(dir, ".env")

	initialHCL, readErr := os.ReadFile(hclPath)
	if readErr != nil {
		t.Fatalf("read .project.hcl after first run: %v", readErr)
	}

	initialEnv, readErr := os.ReadFile(envPath)
	if readErr != nil {
		t.Fatalf("read .env after first run: %v", readErr)
	}

	_, stderr, err = runNewCommand(t, args...)
	if err == nil {
		t.Fatalf("expected failure on rerun without --force; stderr=%s", stderr)
	}

	rerunHCL, readErr := os.ReadFile(hclPath)
	if readErr != nil {
		t.Fatalf("read .project.hcl after failed rerun: %v", readErr)
	}

	if !bytes.Equal(rerunHCL, initialHCL) {
		t.Fatalf(".project.hcl changed after failed rerun")
	}

	rerunEnv, readErr := os.ReadFile(envPath)
	if readErr != nil {
		t.Fatalf("read .env after failed rerun: %v", readErr)
	}

	if !bytes.Equal(rerunEnv, initialEnv) {
		t.Fatalf(".env changed after failed rerun")
	}

	if writeErr := os.WriteFile(hclPath, []byte("stale"), 0o600); writeErr != nil {
		t.Fatalf("prepare stale .project.hcl: %v", writeErr)
	}

	if writeErr := os.WriteFile(envPath, []byte("stale"), 0o600); writeErr != nil {
		t.Fatalf("prepare stale .env: %v", writeErr)
	}

	forcedArgs := append(append([]string{}, args...), "--force")
	_, stderr, err = runNewCommand(t, forcedArgs...)
	if err != nil {
		t.Fatalf("force rerun failed: %v; stderr=%s", err, stderr)
	}

	forcedHCL, readErr := os.ReadFile(hclPath)
	if readErr != nil {
		t.Fatalf("read .project.hcl after force rerun: %v", readErr)
	}

	if bytes.Equal(forcedHCL, []byte("stale")) {
		t.Fatalf("force rerun did not rewrite .project.hcl")
	}

	forcedEnv, readErr := os.ReadFile(envPath)
	if readErr != nil {
		t.Fatalf("read .env after force rerun: %v", readErr)
	}

	if bytes.Equal(forcedEnv, []byte("stale")) {
		t.Fatalf("force rerun did not rewrite .env")
	}
}
