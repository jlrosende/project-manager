//go:build unit
// +build unit

package unit_test

import (
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestMergeProjectInputs_EnvironmentFlagsOverrideConfig(t *testing.T) {
	t.Helper()

	name := "demo"
	path := "/tmp/demo"
	cfgEnvFile := ".env.config"
	cfgEnvMode := "merge"
	cfgColor := "42"

	cfg := &domain.ConfigInput{
		Name: ptrString(name),
		Path: ptrString(path),
		Environment: &domain.EnvironmentInput{
			EnvVarsFile: ptrString(cfgEnvFile),
			EnvVarsMode: ptrString(cfgEnvMode),
			Color:       ptrString(cfgColor),
			EnvVars: map[string]string{
				"FOO":         "from-config",
				"CONFIG_ONLY": "value",
			},
		},
	}

	envVars := map[string]string{
		"FOO":      "from-cli",
		"CLI_ONLY": "1",
	}

	flags := services.ProjectConfigFlags{
		Name: name,
		Path: path,
		Environment: services.EnvironmentConfigFlags{
			Name:        "staging",
			NameSet:     true,
			EnvVarsFile: ".env.cli",
			EnvFileSet:  true,
			EnvVarsMode: "replace",
			ModeSet:     true,
			Color:       "160",
			ColorSet:    true,
			EnvVars:     envVars,
		},
	}

	merged, mergedEnvVars, err := services.MergeProjectInputs(cfg, flags)
	if err != nil {
		t.Fatalf("merge project inputs: %v", err)
	}

	if merged.Environment == nil {
		t.Fatalf("expected environment to be present in merged definition")
	}

	if merged.Environment.EnvVarsFile == nil || *merged.Environment.EnvVarsFile != ".env.cli" {
		t.Fatalf("expected env vars file to be overridden, got %v", merged.Environment.EnvVarsFile)
	}

	if merged.Environment.EnvVarsMode == nil || *merged.Environment.EnvVarsMode != "replace" {
		t.Fatalf("expected env vars mode to be overridden, got %v", merged.Environment.EnvVarsMode)
	}

	if merged.Environment.Color == nil || *merged.Environment.Color != "160" {
		t.Fatalf("expected color to be overridden, got %v", merged.Environment.Color)
	}

	if merged.Environment.Name == nil || *merged.Environment.Name != "staging" {
		t.Fatalf("expected environment name to be propagated, got %v", merged.Environment.Name)
	}

	expectedVars := map[string]string{
		"FOO":         "from-cli",
		"CLI_ONLY":    "1",
		"CONFIG_ONLY": "value",
	}

	for key, value := range expectedVars {
		if mergedVal, ok := merged.Environment.EnvVars[key]; !ok || mergedVal != value {
			t.Fatalf("expected merged env var %s=%s, got %s", key, value, merged.Environment.EnvVars[key])
		}

		if mergedEnvVars[key] != value {
			t.Fatalf("expected returned env vars map to contain %s=%s, got %s", key, value, mergedEnvVars[key])
		}
	}

	if len(merged.Environment.EnvVars) != len(expectedVars) {
		t.Fatalf("unexpected merged env vars length, got %v", merged.Environment.EnvVars)
	}
}
