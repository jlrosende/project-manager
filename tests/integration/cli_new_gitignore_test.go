//go:build integration
// +build integration

package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLINew_GitignoreNoSecrets(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := filepath.Join(home, "project")

	if _, stderr, err := runNewCommand(t, "new", "secrets", projectDir); err != nil {
		t.Fatalf("initial project creation failed: %v; stderr=%s", err, stderr)
	}

	cfgPath := filepath.Join(t.TempDir(), "env.yaml")
	cfgContents := fmt.Sprintf(`name: secrets
path: %q
environment:
  name: staging
  env_vars_file: ".env.staging"
  env_vars_mode: merge
  env_vars:
    API_KEY: super-secret
`, projectDir)

	if err := os.WriteFile(cfgPath, []byte(cfgContents), 0o600); err != nil {
		t.Fatalf("write environment config: %v", err)
	}

	if _, stderr, err := runNewCommand(t, "new", "secrets", "staging", "--cli-input", cfgPath); err != nil {
		t.Fatalf("environment addition failed: %v; stderr=%s", err, stderr)
	}

	gitignore := filepath.Join(projectDir, ".gitignore")
	data, readErr := os.ReadFile(gitignore)
	if readErr != nil {
		t.Fatalf("read .gitignore: %v", readErr)
	}

	if !strings.Contains(string(data), ".env") {
		t.Fatalf(".gitignore missing .env entry: %s", string(data))
	}

	envFile := filepath.Join(projectDir, ".env.staging")
	envData, readErr := os.ReadFile(envFile)
	if readErr != nil {
		t.Fatalf("read .env.staging: %v", readErr)
	}

	if !strings.Contains(string(envData), "API_KEY=super-secret") {
		t.Fatalf(".env.staging missing expected variable: %s", string(envData))
	}

	projectData, readErr := os.ReadFile(filepath.Join(projectDir, ".project.hcl"))
	if readErr != nil {
		t.Fatalf("read .project.hcl: %v", readErr)
	}

	if strings.Contains(string(projectData), "super-secret") {
		t.Fatalf("secret leaked into .project.hcl: %s", string(projectData))
	}
}
