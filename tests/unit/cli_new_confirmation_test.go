//go:build unit
// +build unit

package unit_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLINew_ProjectPromptDecline(t *testing.T) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	target := filepath.Join(home, "demo-project")

	cmd := cli.Root()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetIn(strings.NewReader("n\n"))
	cmd.SetArgs([]string{"new", "demo", target})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute command: %v", err)
	}

	if !strings.Contains(stdout.String(), "Operation cancelled.") {
		t.Fatalf("expected cancellation message, got %q", stdout.String())
	}

	if _, err := os.Stat(filepath.Join(target, ".project.hcl")); !os.IsNotExist(err) {
		t.Fatalf("expected project file to be absent, got err=%v", err)
	}
}

func TestCLINew_EnvironmentPromptDecline(t *testing.T) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectPath := filepath.Join(home, "demo-project")

	create := cli.Root()
	create.SetOut(&bytes.Buffer{})
	create.SetErr(&bytes.Buffer{})
	create.SetIn(strings.NewReader("y\n"))
	create.SetArgs([]string{"new", "demo", projectPath})

	if err := create.Execute(); err != nil {
		t.Fatalf("create project: %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectPath, ".project.hcl")); err != nil {
		t.Fatalf("expected project file to exist, got %v", err)
	}

	addEnv := cli.Root()
	var stdout, stderr bytes.Buffer
	addEnv.SetOut(&stdout)
	addEnv.SetErr(&stderr)
	addEnv.SetIn(strings.NewReader("n\n"))
	addEnv.SetArgs([]string{"new", "demo", "staging"})

	if err := addEnv.Execute(); err != nil {
		t.Fatalf("add environment: %v", err)
	}

	if !strings.Contains(stdout.String(), "Operation cancelled.") {
		t.Fatalf("expected cancellation message, got %q", stdout.String())
	}

	if _, err := os.Stat(filepath.Join(projectPath, ".staging.env")); !os.IsNotExist(err) {
		t.Fatalf("expected environment file to be absent, got err=%v", err)
	}
}
