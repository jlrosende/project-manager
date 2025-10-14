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

	if _, exists := payload["environments"]; exists {
		t.Fatalf("legacy environments map should not be present: %v", payload)
	}

	envPayload, ok := payload["environment"].(map[string]any)
	if !ok {
		t.Fatalf("expected environment object in skeleton: %v", payload)
	}

	if _, ok := envPayload["env_vars"].(map[string]any); !ok {
		t.Fatalf("expected env_vars map in skeleton environment: %v", envPayload)
	}

	if _, ok := envPayload["env_vars_mode"].(string); !ok {
		t.Fatalf("expected env_vars_mode string in skeleton environment: %v", envPayload)
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
	if !strings.Contains(content, "environment:") {
		t.Fatalf("environment block missing from YAML skeleton: %s", content)
	}

	if strings.Contains(content, "environments:") {
		t.Fatalf("legacy environments map should not appear in YAML skeleton: %s", content)
	}

	if !strings.Contains(content, "env_vars:") {
		t.Fatalf("env_vars section missing from YAML skeleton: %s", content)
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

	if !strings.Contains(stdout, "\"environment\"") {
		t.Fatalf("stdout missing environment section: %s", stdout)
	}

	if strings.Contains(stdout, "\"environments\"") {
		t.Fatalf("legacy environments map detected in stdout skeleton: %s", stdout)
	}
}
