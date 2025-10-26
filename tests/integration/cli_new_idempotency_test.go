//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/bootstrap"
)

func TestCLINew_IdempotentReRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	dir := t.TempDir()
	creationArgs := []string{"new", "demo", dir}

	stdout, stderr, err := runNewCommand(t, creationArgs...)
	if err != nil {
		t.Fatalf("first run failed: %v; stderr=%s", err, stderr)
	}

	container, err := bootstrap.New(bootstrap.Options{Logger: nil})
	if err != nil {
		t.Fatalf("bootstrap after first run: %v", err)
	}

	probe, err := container.ProjectService.Probe("demo")
	if err != nil {
		t.Fatalf("probe existing project: %v", err)
	}

	if !probe.RegistryHit || !probe.ProjectFileExists {
		t.Fatalf("expected probe to report existing project; probe=%+v", probe)
	}

	projects, err := container.ProjectService.List()
	if err != nil {
		t.Fatalf("list projects: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected exactly one project after first run, got %d", len(projects))
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

	rerunArgs := []string{"new", "demo"}

	stdout, stderr, err = runNewCommand(t, rerunArgs...)
	if err != nil {
		t.Fatalf("rerun should succeed as a no-op: %v", err)
	}

	expectedMessage := "project demo already exists; nothing to do"
	if !strings.Contains(stdout, expectedMessage) {
		t.Fatalf("rerun stdout missing no-op message: %q", stdout)
	}

	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("unexpected stderr on rerun: %q", stderr)
	}

	rerunHCL, readErr := os.ReadFile(hclPath)
	if readErr != nil {
		t.Fatalf("read .project.hcl after rerun: %v", readErr)
	}

	if !bytes.Equal(rerunHCL, initialHCL) {
		t.Fatalf(".project.hcl changed after no-op rerun")
	}

	rerunEnv, readErr := os.ReadFile(envPath)
	if readErr != nil {
		t.Fatalf("read .env after rerun: %v", readErr)
	}

	if !bytes.Equal(rerunEnv, initialEnv) {
		t.Fatalf(".env changed after no-op rerun")
	}
}

func TestCLINew_RecreateAfterStaleMetadata(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	dir := t.TempDir()
	args := []string{"new", "demo", dir}

	if _, stderr, err := runNewCommand(t, args...); err != nil {
		t.Fatalf("initial run failed: %v; stderr=%s", err, stderr)
	}

	gitconfigPath := filepath.Join(home, ".gitconfig")
	initialConfig, err := os.ReadFile(gitconfigPath)
	if err != nil {
		t.Fatalf("read initial gitconfig: %v", err)
	}

	hclPath := filepath.Join(dir, ".project.hcl")
	envPath := filepath.Join(dir, ".env")

	if removeErr := os.Remove(hclPath); removeErr != nil {
		t.Fatalf("remove .project.hcl: %v", removeErr)
	}

	if removeErr := os.Remove(envPath); removeErr != nil {
		t.Fatalf("remove .env: %v", removeErr)
	}

	stdout, stderr, err := runNewCommand(t, args...)
	if err != nil {
		t.Fatalf("recreation run failed: %v; stdout=%s stderr=%s", err, stdout, stderr)
	}

	newHCL, readErr := os.ReadFile(hclPath)
	if readErr != nil {
		t.Fatalf("read recreated .project.hcl: %v", readErr)
	}

	if !strings.Contains(string(newHCL), "name = \"demo\"") {
		t.Fatalf("recreated .project.hcl missing project name; contents=%s", string(newHCL))
	}

	if _, readErr := os.Stat(envPath); readErr != nil {
		t.Fatalf("stat recreated .env: %v", readErr)
	}

	gitConfig, readErr := os.ReadFile(gitconfigPath)
	if readErr != nil {
		t.Fatalf("read gitconfig after recreation: %v", readErr)
	}

	gitdirKey := fmt.Sprintf("gitdir/i:%s/", dir)
	if strings.Count(string(gitConfig), gitdirKey) != 1 {
		t.Fatalf("expected single includeIf entry for %s; gitconfig=%s", dir, string(gitConfig))
	}

	if !bytes.Equal(gitConfig, initialConfig) {
		t.Fatalf("gitconfig changed unexpectedly during recreation\ninitial=%s\nupdated=%s", string(initialConfig), string(gitConfig))
	}
}
