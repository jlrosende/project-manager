//go:build integration
// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLINew_DryRunProjectCreationText(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	target := filepath.Join(home, "dry-run-project")
	stdout, stderr, err := runNewCommand(t, "new", "dry", target, "--dry-run")
	if err != nil {
		t.Fatalf("dry-run execution failed: %v; stderr=%s", err, stderr)
	}

	if _, statErr := os.Stat(filepath.Join(target, ".project.hcl")); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run should not create project file, stat err=%v", statErr)
	}

	if _, statErr := os.Stat(filepath.Join(target, ".env")); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run should not create env file, stat err=%v", statErr)
	}

	if !strings.Contains(stdout, "pm new (dry-run)") {
		t.Fatalf("stdout missing dry-run header: %s", stdout)
	}

	if !strings.Contains(stdout, "name: dry") {
		t.Fatalf("stdout missing project name: %s", stdout)
	}

	expectedPath := fmt.Sprintf("path: %s", target)
	if !strings.Contains(stdout, expectedPath) {
		t.Fatalf("stdout missing project path (%s): %s", expectedPath, stdout)
	}

	if strings.Contains(stdout, "environment:") {
		t.Fatalf("unexpected environment section in project creation dry-run: %s", stdout)
	}
}

func TestCLINew_DryRunEnvironmentJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SHELL", "/bin/sh")

	projectDir := filepath.Join(home, "demo-project")
	if _, stderr, err := runNewCommand(t, "new", "demo", projectDir); err != nil {
		t.Fatalf("project creation failed: %v; stderr=%s", err, stderr)
	}

	stdout, stderr, err := runNewCommand(
		t,
		"new", "demo", "staging",
		"--environment-env-file", ".env.staging",
		"--environment-mode", "replace",
		"--environment-color", "160",
		"--environment-env-var", "API_URL=https://api.example.com",
		"--dry-run",
		"--output", "json",
	)
	if err != nil {
		t.Fatalf("environment dry-run failed: %v; stderr=%s", err, stderr)
	}

	var payload map[string]any
	if decodeErr := json.Unmarshal([]byte(stdout), &payload); decodeErr != nil {
		t.Fatalf("decode dry-run JSON: %v", decodeErr)
	}

	if _, exists := payload["env_vars"]; exists {
		t.Fatalf("unexpected top-level env_vars in environment dry-run: %v", payload)
	}

	envPayload, ok := payload["environment"].(map[string]any)
	if !ok {
		t.Fatalf("environment payload missing or invalid: %v", payload)
	}

	if name, ok := envPayload["name"].(string); !ok || name != "staging" {
		t.Fatalf("expected environment name staging, got %v", envPayload["name"])
	}

	if file, ok := envPayload["env_vars_file"].(string); !ok || file != ".env.staging" {
		t.Fatalf("expected env_vars_file .env.staging, got %v", envPayload["env_vars_file"])
	}

	if mode, ok := envPayload["env_vars_mode"].(string); !ok || mode != "replace" {
		t.Fatalf("expected env_vars_mode replace, got %v", envPayload["env_vars_mode"])
	}

	if color, ok := envPayload["color"].(string); !ok || color != "160" {
		t.Fatalf("expected color 160, got %v", envPayload["color"])
	}

	varsPayload, ok := envPayload["env_vars"].(map[string]any)
	if !ok {
		t.Fatalf("expected env_vars map, got %T", envPayload["env_vars"])
	}

	if value, ok := varsPayload["API_URL"].(string); !ok || value != "https://api.example.com" {
		t.Fatalf("expected API_URL entry, got %v", varsPayload["API_URL"])
	}

	countValue, ok := envPayload["env_vars_count"].(float64)
	if !ok {
		t.Fatalf("expected env_vars_count number, got %T", envPayload["env_vars_count"])
	}

	if expected := float64(len(varsPayload)); countValue != expected {
		t.Fatalf("expected env_vars_count %v, got %v", expected, countValue)
	}

	if _, statErr := os.Stat(filepath.Join(projectDir, ".env.staging")); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run should not create environment file, stat err=%v", statErr)
	}
}
