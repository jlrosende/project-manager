//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLIEditProject_DefaultEnvUnknownValidation(t *testing.T) {
	fx := newEditProjectFixture(t, "invalid-default-env")

	original, err := os.ReadFile(filepath.Join(fx.projectDir, ".project.hcl"))
	if err != nil {
		t.Fatalf("read .project.hcl: %v", err)
	}

	cmd := cli.Root()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"edit", fx.name, "--project-default-env", "qa"})

	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	exitCode := extractExitCode(t, err)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "Default Environment") {
		t.Fatalf("expected default environment validation message, stderr=%s", stderr.String())
	}

	updated, readErr := os.ReadFile(filepath.Join(fx.projectDir, ".project.hcl"))
	if readErr != nil {
		t.Fatalf("read updated .project.hcl: %v", readErr)
	}

	if !bytes.Equal(original, updated) {
		t.Fatal("project configuration should remain unchanged after validation failure")
	}
}

func TestCLIEditEnvironment_InvalidModeValidation(t *testing.T) {
	fx := newEditProjectFixture(t, "invalid-env-mode")

	original, err := os.ReadFile(filepath.Join(fx.projectDir, ".project.hcl"))
	if err != nil {
		t.Fatalf("read .project.hcl: %v", err)
	}

	cmd := cli.Root()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"edit", fx.name, "staging", "--env-env-vars-mode", "invalid"})

	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	exitCode := extractExitCode(t, err)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "Env Vars Mode") {
		t.Fatalf("expected env vars mode validation message, stderr=%s", stderr.String())
	}

	updated, readErr := os.ReadFile(filepath.Join(fx.projectDir, ".project.hcl"))
	if readErr != nil {
		t.Fatalf("read updated .project.hcl: %v", readErr)
	}

	if !bytes.Equal(original, updated) {
		t.Fatal("project configuration should remain unchanged after validation failure")
	}
}

func extractExitCode(t *testing.T, err error) int {
	t.Helper()

	type exitCoder interface {
		ExitCode() int
	}

	var ec exitCoder
	if !errors.As(err, &ec) {
		t.Fatalf("error does not expose exit code: %v", err)
	}

	return ec.ExitCode()
}
