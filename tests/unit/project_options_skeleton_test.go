//go:build unit
// +build unit

package unit_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestGenerateProjectSkeletonJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project.json")

	if err := services.GenerateProjectSkeleton(path, services.SkeletonFormatJSON); err != nil {
		t.Fatalf("GenerateProjectSkeleton JSON: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read skeleton: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal skeleton: %v", err)
	}

	if payload["name"] == nil || payload["path"] == nil {
		t.Fatalf("expected name and path fields, got %#v", payload)
	}
}

func TestGenerateProjectSkeletonYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project.yaml")

	if err := services.GenerateProjectSkeleton(path, services.SkeletonFormatYAML); err != nil {
		t.Fatalf("GenerateProjectSkeleton YAML: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read skeleton: %v", err)
	}

	var payload map[string]any
	if err := yaml.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal skeleton: %v", err)
	}

	if payload["name"] == nil || payload["path"] == nil {
		t.Fatalf("expected name and path fields, got %#v", payload)
	}
}

func TestGenerateProjectSkeletonExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.json")

	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if err := services.GenerateProjectSkeleton(path, services.SkeletonFormatJSON); err == nil {
		t.Fatalf("expected error when file exists")
	}
}
