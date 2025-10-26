//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLINew_GenerateSkeletonJSON(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SHELL", "/bin/sh")

	dir := t.TempDir()
	skeleton := filepath.Join(dir, "project.json")

	stdout, stderr, err := runNewCommand(t, "new", "--generate-cli-skeleton-json", skeleton)
	if err != nil {
		t.Fatalf("generate CLI skeleton JSON: %v; stderr=%s", err, stderr)
	}

	if !strings.Contains(stdout, "CLI input skeleton written") {
		t.Fatalf("expected confirmation message, got %s", stdout)
	}

	data, readErr := os.ReadFile(skeleton)
	if readErr != nil {
		t.Fatalf("read skeleton file: %v", readErr)
	}

	var payload map[string]any
	if decodeErr := json.Unmarshal(data, &payload); decodeErr != nil {
		t.Fatalf("decode skeleton JSON: %v", decodeErr)
	}

	if _, exists := payload["environment"]; exists {
		t.Fatalf("environment block should not be present: %v", payload)
	}

	if _, exists := payload["metadata"]; exists {
		t.Fatalf("metadata key should not be present: %v", payload)
	}

	description, ok := payload["description"].(string)
	if !ok || description != "Describe your project" {
		t.Fatalf("expected description placeholder, got %v", payload["description"])
	}

	shell, ok := payload["shell"].(string)
	if !ok || shell == "" {
		t.Fatalf("expected shell value, got %v", payload["shell"])
	}

	envFile, ok := payload["env-file"].(string)
	if !ok || envFile != ".env" {
		t.Fatalf("expected env-file placeholder .env, got %v", payload["env-file"])
	}

	envVars, ok := payload["env_vars"].(map[string]any)
	if !ok {
		t.Fatalf("expected env_vars map in skeleton: %v", payload)
	}

	if value, ok := envVars["EXAMPLE_KEY"].(string); !ok || value != "VALUE" {
		t.Fatalf("expected env_vars placeholder EXAMPLE_KEY=VALUE, got %v", envVars)
	}
}

func TestCLINew_GenerateSkeletonYAML(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SHELL", "/bin/sh")

	dir := t.TempDir()
	skeleton := filepath.Join(dir, "project.yaml")

	stdout, stderr, err := runNewCommand(t, "new", "--generate-cli-skeleton-yaml", skeleton)
	if err != nil {
		t.Fatalf("generate CLI skeleton YAML: %v; stderr=%s", err, stderr)
	}

	if !strings.Contains(stdout, "CLI input skeleton written") {
		t.Fatalf("expected confirmation message, got %s", stdout)
	}

	data, readErr := os.ReadFile(skeleton)
	if readErr != nil {
		t.Fatalf("read skeleton YAML: %v", readErr)
	}

	content := string(data)
	if strings.Contains(content, "environment:") {
		t.Fatalf("environment block should not appear in YAML skeleton: %s", content)
	}

	if strings.Contains(content, "metadata:") {
		t.Fatalf("metadata block should not appear in YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "description:") {
		t.Fatalf("description missing from YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "shell:") {
		t.Fatalf("shell missing from YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "env-file:") {
		t.Fatalf("env-file missing from YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "env_vars:") {
		t.Fatalf("env_vars section missing from YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "EXAMPLE_KEY: VALUE") {
		t.Fatalf("placeholder env var missing from YAML skeleton: %s", content)
	}
}

func TestCLINew_GenerateSkeletonJSONStdout(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SHELL", "/bin/sh")

	stdout, stderr, err := runNewCommand(t, "new", "--generate-cli-skeleton-json")
	if err != nil {
		t.Fatalf("generate CLI skeleton JSON stdout: %v; stderr=%s", err, stderr)
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr output: %s", stderr)
	}

	if strings.Contains(stdout, "\"environment\"") {
		t.Fatalf("stdout skeleton should not include environment block: %s", stdout)
	}

	if strings.Contains(stdout, "\"metadata\"") {
		t.Fatalf("stdout skeleton should not include metadata block: %s", stdout)
	}

	if !strings.Contains(stdout, "\"description\"") {
		t.Fatalf("stdout skeleton missing description field: %s", stdout)
	}

	if !strings.Contains(stdout, "\"shell\"") {
		t.Fatalf("stdout skeleton missing shell field: %s", stdout)
	}

	if !strings.Contains(stdout, "\"env-file\"") {
		t.Fatalf("stdout skeleton missing env-file field: %s", stdout)
	}

	if !strings.Contains(stdout, "\"env_vars\"") {
		t.Fatalf("stdout skeleton missing env_vars section: %s", stdout)
	}

	if !strings.Contains(stdout, "EXAMPLE_KEY") {
		t.Fatalf("stdout skeleton missing placeholder env var: %s", stdout)
	}
}

func TestCLINew_GenerateSkeletonExistingProjectPrefillsPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("create project: %v; stderr=%s", err, stderr)
	}

	stdout, stderr, err := runNewCommand(t, "new", "demo", "Staging Env", "--generate-cli-skeleton-json")
	if err != nil {
		t.Fatalf("generate skeleton for existing project: %v; stderr=%s", err, stderr)
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr output: %s", stderr)
	}

	var payload map[string]any
	if decodeErr := json.Unmarshal([]byte(stdout), &payload); decodeErr != nil {
		t.Fatalf("decode skeleton JSON: %v", decodeErr)
	}

	pathValue, ok := payload["path"].(string)
	if !ok {
		t.Fatalf("expected path string in skeleton payload: %v", payload)
	}

	if pathValue != projectDir {
		t.Fatalf("expected skeleton path %s, got %s", projectDir, pathValue)
	}

	envVars, ok := payload["env_vars"].(map[string]any)
	if !ok || len(envVars) == 0 {
		t.Fatalf("expected env_vars map in skeleton payload: %v", payload)
	}

	if _, exists := payload["environment"]; exists {
		t.Fatalf("environment block should not be present in skeleton payload: %v", payload)
	}
}

func TestCLINew_GenerateSkeletonExistingProjectInvalidEnvironmentName(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := t.TempDir()

	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("create project: %v; stderr=%s", err, stderr)
	}

	_, stderr, err := runNewCommand(t, "new", "demo", "invalid/name", "--generate-cli-skeleton-json")
	if err == nil {
		t.Fatalf("expected error when environment name contains path separators; stderr=%s", stderr)
	}

	if !strings.Contains(err.Error(), "environment name") {
		t.Fatalf("unexpected error: %v", err)
	}
}
