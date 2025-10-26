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

func TestCLINew_AddEnvironmentRequiresExistingProject(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	stdout, stderr, err := runNewCommand(t, "new", "demo", "staging", "--environment-env-var", "API_URL=https://api.example.com")
	if err == nil {
		t.Fatalf("expected failure when adding environment without project; stdout=%s stderr=%s", stdout, stderr)
	}

	if !strings.Contains(err.Error(), "requires an existing project") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLINew_AddEnvironmentWithFlags(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("create project: %v; stderr=%s", err, stderr)
	}

	envArgs := []string{
		"new", "demo", "staging",
		"--environment-env-file", ".env.staging",
		"--environment-mode", "replace",
		"--environment-color", "160",
		"--environment-env-var", "API_URL=https://api.example.com",
		"--environment-env-var", "FEATURE_X=true",
	}

	if stdout, stderr, err := runNewCommand(t, envArgs...); err != nil {
		t.Fatalf("add environment: %v; stdout=%s stderr=%s", err, stdout, stderr)
	}

	hclPath := filepath.Join(projectDir, ".project.hcl")
	hclData, err := os.ReadFile(hclPath)
	if err != nil {
		t.Fatalf("read project file: %v", err)
	}

	hcl := string(hclData)
	if !strings.Contains(hcl, "environment \"staging\"") {
		t.Fatalf("environment block missing from project file:\n%s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_mode = \"replace\"") {
		t.Fatalf("expected env_vars_mode override in project file:\n%s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_file = \".env.staging\"") {
		t.Fatalf("expected env_vars_file override in project file:\n%s", hcl)
	}

	if !strings.Contains(hcl, "color = \"160\"") {
		t.Fatalf("expected color override in project file:\n%s", hcl)
	}

	envPath := filepath.Join(projectDir, ".env.staging")
	envData, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read environment file: %v", err)
	}

	envLines := strings.Split(strings.TrimSpace(string(envData)), "\n")
	expectedLines := map[string]bool{
		"API_URL=https://api.example.com": true,
		"FEATURE_X=true":                  true,
	}

	if len(envLines) != len(expectedLines) {
		t.Fatalf("unexpected env file contents: %q", envLines)
	}

	for _, line := range envLines {
		if !expectedLines[line] {
			t.Fatalf("unexpected env entry %q", line)
		}
	}
}

func TestCLINew_AddEnvironmentForceReplaces(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("create project: %v; stderr=%s", err, stderr)
	}

	if _, stderr, err := runNewCommand(t, "new", "demo", "staging", "--environment-env-var", "OLD=1"); err != nil {
		t.Fatalf("initial environment add failed: %v; stderr=%s", err, stderr)
	}

	forceArgs := []string{
		"new", "demo", "staging",
		"--environment-env-file", ".env.force",
		"--environment-mode", "replace",
		"--environment-color", "200",
		"--environment-env-var", "NEW=1",
		"--force",
	}

	if stdout, stderr, err := runNewCommand(t, forceArgs...); err != nil {
		t.Fatalf("force environment update failed: %v; stdout=%s stderr=%s", err, stdout, stderr)
	}

	hclPath := filepath.Join(projectDir, ".project.hcl")
	hclData, err := os.ReadFile(hclPath)
	if err != nil {
		t.Fatalf("read project file: %v", err)
	}

	hcl := string(hclData)
	if !strings.Contains(hcl, "env_vars_file = \".env.force\"") {
		t.Fatalf("expected env_vars_file to update after force:\n%s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_mode = \"replace\"") {
		t.Fatalf("expected env_vars_mode to update after force:\n%s", hcl)
	}

	if !strings.Contains(hcl, "color         = \"200\"") {
		t.Fatalf("expected color to update after force:\n%s", hcl)
	}

	oldEnvPath := filepath.Join(projectDir, ".staging.env")
	if _, err := os.Stat(oldEnvPath); !os.IsNotExist(err) {
		t.Fatalf("expected old env file %s to be removed or renamed", oldEnvPath)
	}

	newEnvPath := filepath.Join(projectDir, ".env.force")
	envData, err := os.ReadFile(newEnvPath)
	if err != nil {
		t.Fatalf("read new env file: %v", err)
	}

	envContent := strings.TrimSpace(string(envData))
	if envContent != "NEW=1" {
		t.Fatalf("expected env file to be replaced with single entry, got %q", envContent)
	}
}

func TestCLINew_AddEnvironment_ConfigAndFlagsMerge(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("create project: %v; stderr=%s", err, stderr)
	}

	cfgPath := filepath.Join(t.TempDir(), "env-config.json")
	config := fmt.Sprintf(`{
		"name": "demo",
		"path": %q,
		"environment": {
			"name": "staging",
			"env_vars_file": ".env.config",
			"env_vars_mode": "merge",
			"color": "42",
			"env_vars": {
				"FOO": "from-config",
				"CONFIG_ONLY": "value"
			}
		}
	}`, projectDir)

	if err := os.WriteFile(cfgPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	envArgs := []string{
		"new", "demo", "staging",
		"--cli-input", cfgPath,
		"--environment-env-file", ".env.cli",
		"--environment-mode", "replace",
		"--environment-color", "160",
		"--environment-env-var", "FOO=from-cli",
		"--environment-env-var", "CLI_ONLY=1",
	}

	if stdout, stderr, err := runNewCommand(t, envArgs...); err != nil {
		t.Fatalf("add environment using config: %v; stdout=%s stderr=%s", err, stdout, stderr)
	}

	hclPath := filepath.Join(projectDir, ".project.hcl")
	hclData, err := os.ReadFile(hclPath)
	if err != nil {
		t.Fatalf("read project file: %v", err)
	}

	hcl := string(hclData)
	if !strings.Contains(hcl, "env_vars_file = \".env.cli\"") {
		t.Fatalf("expected CLI env file to override config:\n%s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_mode = \"replace\"") {
		t.Fatalf("expected CLI env mode to override config:\n%s", hcl)
	}

	if !strings.Contains(hcl, "color = \"160\"") {
		t.Fatalf("expected CLI color to override config:\n%s", hcl)
	}

	envPath := filepath.Join(projectDir, ".env.cli")
	envData, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read merged env file: %v", err)
	}

	envEntries := strings.Split(strings.TrimSpace(string(envData)), "\n")
	expected := map[string]string{
		"FOO":         "from-cli",
		"CLI_ONLY":    "1",
		"CONFIG_ONLY": "value",
	}

	for _, entry := range envEntries {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			t.Fatalf("unexpected env entry %q", entry)
		}

		key := parts[0]
		value := parts[1]
		expectedValue, ok := expected[key]
		if !ok {
			t.Fatalf("unexpected env key %q", key)
		}

		if value != expectedValue {
			t.Fatalf("expected %s=%s, got %s", key, expectedValue, value)
		}

		delete(expected, key)
	}

	if len(expected) != 0 {
		t.Fatalf("missing env entries: %v", expected)
	}
}

func TestCLINew_ConfigLegacyEnvironmentsRequiresAllowUnknown(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := filepath.Join(t.TempDir(), "legacy-project")
	cfgPath := filepath.Join(t.TempDir(), "legacy-config.json")
	config := fmt.Sprintf(`{
		"name": "demo",
		"path": %q,
		"environments": {
			"staging": {}
		}
	}`, projectDir)

	if err := os.WriteFile(cfgPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	if stdout, stderr, err := runNewCommand(t, "new", "demo", "--cli-input", cfgPath); err == nil {
		t.Fatalf("expected legacy environments error; stdout=%s stderr=%s", stdout, stderr)
	}

	if _, stderr, err := runNewCommand(t, "new", "demo", "--cli-input", cfgPath, "--allow-unknown"); err != nil {
		t.Fatalf("expected success when ignoring legacy environments: %v; stderr=%s", err, stderr)
	}

	hclPath := filepath.Join(projectDir, ".project.hcl")
	if _, err := os.Stat(hclPath); err != nil {
		t.Fatalf("expected project file after allowing unknown fields: %v", err)
	}

	hclData, err := os.ReadFile(hclPath)
	if err != nil {
		t.Fatalf("read project file: %v", err)
	}

	if strings.Contains(string(hclData), "environment \"") {
		t.Fatalf("legacy environments map should be ignored, but environment block detected:\n%s", string(hclData))
	}
}
