//go:build integration
// +build integration

package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLINew_GitignoreNoSecrets(t *testing.T) {
	temp := t.TempDir()
	projectDir := filepath.Join(temp, "project")

	configDir := filepath.Join(temp, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	cfg := filepath.Join(configDir, "project.yaml")
	cfgContents := "name: secrets\npath: " + projectDir + "\nenvironments:\n  API_KEY: super-secret\n"
	if err := os.WriteFile(cfg, []byte(cfgContents), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, stderr, err := runNewCommand(t, "new", "secrets", "--cli-input", cfg)
	if err != nil {
		t.Fatalf("pm new failed: %v; stderr=%s", err, stderr)
	}

	gitignore := filepath.Join(projectDir, ".gitignore")
	data, readErr := os.ReadFile(gitignore)
	if readErr != nil {
		t.Fatalf("read .gitignore: %v", readErr)
	}

	if !strings.Contains(string(data), ".env") {
		t.Fatalf(".gitignore missing .env entry: %s", string(data))
	}

	envData, readErr := os.ReadFile(filepath.Join(projectDir, ".env"))
	if readErr != nil {
		t.Fatalf("read .env: %v", readErr)
	}

	if !strings.Contains(string(envData), "API_KEY=super-secret") {
		t.Fatalf(".env missing expected variable: %s", string(envData))
	}

	projectData, readErr := os.ReadFile(filepath.Join(projectDir, ".project.hcl"))
	if readErr != nil {
		t.Fatalf("read .project.hcl: %v", readErr)
	}

	if strings.Contains(string(projectData), "super-secret") {
		t.Fatalf("secret leaked into .project.hcl: %s", string(projectData))
	}
}
