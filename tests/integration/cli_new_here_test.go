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

func TestCLINew_InitHere(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(dir)

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

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

func TestCLINew_HereRejectsPathArgument(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(dir)

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{"new", "demo", dir, "--here"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error when combining --here with path; stdout=%s stderr=%s", out.String(), errBuf.String())
	} else if !strings.Contains(err.Error(), "--here") {
		t.Fatalf("unexpected error message: %v", err)
	}

	if _, statErr := os.Stat(filepath.Join(dir, ".project.hcl")); !os.IsNotExist(statErr) {
		t.Fatalf(".project.hcl should not be created; statErr=%v", statErr)
	}
}

func TestCLINew_HereRejectsEnvironmentAddition(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	bootstrapCmd := cli.Root()
	bootstrapCmd.SetOut(&bytes.Buffer{})
	bootstrapCmd.SetErr(&bytes.Buffer{})
	bootstrapCmd.SetArgs([]string{"new", "demo", projectDir})
	if err := bootstrapCmd.Execute(); err != nil {
		t.Fatalf("bootstrap project failed: %v", err)
	}

	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(projectDir)

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{"new", "demo", "staging", "--here"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error when combining --here with environment; stdout=%s stderr=%s", out.String(), errBuf.String())
	} else if !strings.Contains(err.Error(), "--here") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
