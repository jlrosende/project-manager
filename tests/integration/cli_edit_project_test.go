//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

type editProjectFixture struct {
	home       string
	projectDir string
	name       string
	gitConfig  string
}

func newEditProjectFixture(t *testing.T, name string) editProjectFixture {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, ".cache"))

	projectDir := filepath.Join(home, "workspace", name)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	projectHCL := `name = "` + name + `"
` +
		`description = "Initial project description"
` +
		`shell = "/bin/bash"
` +
		`env_vars_file = ".env"
` +
		`default_env = "staging"

` +
		`environment "staging" {
` +
		`  color = "240"
` +
		`  env_vars_mode = "merge"
` +
		`  env_vars_file = ".env.staging"
` +
		`}
` +
		`environment "production" {
` +
		`  color = "120"
` +
		`  env_vars_mode = "merge"
` +
		`  env_vars_file = ".env.production"
` +
		`}
`

	if err := os.WriteFile(filepath.Join(projectDir, ".project.hcl"), []byte(projectHCL), 0o600); err != nil {
		t.Fatalf("write .project.hcl: %v", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, ".env"), []byte("KEY=VALUE\n"), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, ".env.staging"), []byte("STAGING=1\n"), 0o600); err != nil {
		t.Fatalf("write staging env file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, ".env.production"), []byte("PRODUCTION=1\n"), 0o600); err != nil {
		t.Fatalf("write production env file: %v", err)
	}

	perProjectGit := filepath.Join(projectDir, "."+name+".gitconfig")
	if err := os.WriteFile(perProjectGit, []byte("[user]\n\tname = tester\n"), 0o600); err != nil {
		t.Fatalf("write per-project git config: %v", err)
	}

	gitConfig := filepath.Join(home, ".gitconfig")
	gitContent := `[includeIf "gitdir/i:` + projectDir + `/"]
	path = ` + perProjectGit + `
`
	if err := os.WriteFile(gitConfig, []byte(gitContent), 0o600); err != nil {
		t.Fatalf("write global git config: %v", err)
	}

	return editProjectFixture{
		home:       home,
		projectDir: projectDir,
		name:       name,
		gitConfig:  gitConfig,
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}

	return string(data)
}

func TestCLIEditProject_FlagsApply(t *testing.T) {
	fx := newEditProjectFixture(t, "demo")

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"--project-description", "Updated description",
		"--project-shell", "/bin/zsh",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute edit: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if !strings.Contains(hcl, "description = \"Updated description\"") {
		t.Fatalf("project description not updated: %s", hcl)
	}

	if !strings.Contains(hcl, "shell = \"/bin/zsh\"") {
		t.Fatalf("project shell not updated: %s", hcl)
	}

	stdout := out.String()
	if !strings.Contains(stdout, "Updated project \"demo\" (2 changes):") {
		t.Fatalf("stdout missing header: %s", stdout)
	}

	if !strings.Contains(stdout, "Description: \"Initial project description\" -> \"Updated description\"") {
		t.Fatalf("stdout missing description change: %s", stdout)
	}
	if !strings.Contains(stdout, "Shell: \"/bin/bash\" -> \"/bin/zsh\"") {
		t.Fatalf("stdout missing shell change: %s", stdout)
	}
}

func TestCLIEditProject_CliInputDryRun(t *testing.T) {
	fx := newEditProjectFixture(t, "sample")

	input := filepath.Join(fx.projectDir, "edit.yaml")
	payload := "project:\n  description: New description from file\n"
	if err := os.WriteFile(input, []byte(payload), 0o600); err != nil {
		t.Fatalf("write cli input: %v", err)
	}

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"--cli-input", input,
		"--dry-run",
		"--output", "json",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute edit dry-run: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if strings.Contains(hcl, "New description from file") {
		t.Fatalf("dry-run should not persist changes: %s", hcl)
	}

	trimmed := strings.TrimSpace(out.String())
	if trimmed == "" {
		t.Fatal("expected json output")
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(trimmed), &data); err != nil {
		t.Fatalf("parse json output: %v\noutput=%s", err, trimmed)
	}

	dryRun, _ := data["dry_run"].(bool)
	if !dryRun {
		t.Fatalf("expected dry_run true in output: %s", trimmed)
	}

	scope, _ := data["scope"].(string)
	if scope != "project" {
		t.Fatalf("expected scope 'project', got %q", scope)
	}

	project, _ := data["project"].(string)
	if project != fx.name {
		t.Fatalf("expected project %q, got %q", fx.name, project)
	}

	changes, ok := data["changes"].([]any)
	if !ok || len(changes) == 0 {
		t.Fatalf("expected at least one change in output: %s", trimmed)
	}
}

func TestCLIEditProject_DryRunTextOutput(t *testing.T) {
	fx := newEditProjectFixture(t, "sample-text")

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"--project-description", "Preview description",
		"--dry-run",
		"--output", "text",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute edit dry-run text: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if strings.Contains(hcl, "Preview description") {
		t.Fatalf("dry-run should not persist changes: %s", hcl)
	}

	stdout := out.String()
	header := "Previewing project edit for \"" + fx.name + "\" (dry-run)"
	if !strings.Contains(stdout, header) {
		t.Fatalf("stdout missing dry-run header: %s", stdout)
	}

	if !strings.Contains(stdout, "No changes have been written; showing preview only.") {
		t.Fatalf("stdout missing dry-run disclaimer: %s", stdout)
	}

	if !strings.Contains(stdout, "Planned changes (1 change):") {
		t.Fatalf("stdout missing change summary: %s", stdout)
	}

	if !strings.Contains(stdout, "Description: \"Initial project description\" -> \"Preview description\"") {
		t.Fatalf("stdout missing description change: %s", stdout)
	}
}

func TestCLIEditProject_NoChanges(t *testing.T) {
	fx := newEditProjectFixture(t, "plain")

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{"edit", fx.name})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute edit without changes: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	stdout := out.String()
	if !strings.Contains(stdout, "No changes applied; nothing to update.") {
		t.Fatalf("expected no-change message, got: %s", stdout)
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if !strings.Contains(hcl, "Initial project description") {
		t.Fatalf("project file unexpectedly modified: %s", hcl)
	}
}

func TestCLIEditEnvironment_FlagsApply(t *testing.T) {
	fx := newEditProjectFixture(t, "envflags")

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"staging",
		"--env-color", "245",
		"--env-env-vars-mode", "replace",
		"--env-env-vars-file", ".env.stage.updated",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute env edit: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if !strings.Contains(hcl, "environment \"staging\"") {
		t.Fatalf("staging environment block missing: %s", hcl)
	}

	if !strings.Contains(hcl, "color = \"245\"") {
		t.Fatalf("environment color not updated: %s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_mode = \"replace\"") {
		t.Fatalf("environment mode not updated: %s", hcl)
	}

	if !strings.Contains(hcl, "env_vars_file = \".env.stage.updated\"") {
		t.Fatalf("environment env vars file not updated: %s", hcl)
	}

	if !strings.Contains(hcl, "environment \"production\"") || !strings.Contains(hcl, "env_vars_file = \".env.production\"") {
		t.Fatalf("other environments unexpectedly modified: %s", hcl)
	}

	if _, err := os.Stat(filepath.Join(fx.projectDir, ".env.stage.updated")); err != nil {
		t.Fatalf("expected renamed env file: %v", err)
	}

	if _, err := os.Stat(filepath.Join(fx.projectDir, ".env.staging")); !os.IsNotExist(err) {
		t.Fatalf("expected original env file renamed, err=%v", err)
	}

	stdout := out.String()
	if !strings.Contains(stdout, "Updated environment \"staging\" in project \"envflags\" (3 changes):") {
		t.Fatalf("stdout missing environment header: %s", stdout)
	}

	if !strings.Contains(stdout, "Color: \"240\" -> \"245\"") {
		t.Fatalf("stdout missing color change: %s", stdout)
	}

	if !strings.Contains(stdout, "Env Vars Mode: \"merge\" -> \"replace\"") {
		t.Fatalf("stdout missing mode change: %s", stdout)
	}

	if !strings.Contains(stdout, "Env Vars File: \".env.staging\" -> \".env.stage.updated\"") {
		t.Fatalf("stdout missing env vars file change: %s", stdout)
	}

	if !strings.Contains(stdout, "Other environments unchanged; project metadata untouched.") {
		t.Fatalf("stdout missing unchanged scope message: %s", stdout)
	}
}

func TestCLIEditEnvironment_CliInputDryRun(t *testing.T) {
	fx := newEditProjectFixture(t, "envdryrun")

	input := filepath.Join(fx.projectDir, "env-edit.yaml")
	payload := "environment:\n  env_vars_file: .env.stage.file\n  color: 250\n"
	if err := os.WriteFile(input, []byte(payload), 0o600); err != nil {
		t.Fatalf("write cli input: %v", err)
	}

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"staging",
		"--cli-input", input,
		"--dry-run",
		"--output", "json",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute env dry-run: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if strings.Contains(hcl, ".env.stage.file") || strings.Contains(hcl, "color = \"250\"") {
		t.Fatalf("dry-run should not persist environment changes: %s", hcl)
	}

	if _, err := os.Stat(filepath.Join(fx.projectDir, ".env.stage.file")); !os.IsNotExist(err) {
		t.Fatalf("dry-run should not create new env file: %v", err)
	}

	trimmed := strings.TrimSpace(out.String())
	if trimmed == "" {
		t.Fatal("expected json output")
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(trimmed), &data); err != nil {
		t.Fatalf("parse json output: %v\noutput=%s", err, trimmed)
	}

	dryRun, _ := data["dry_run"].(bool)
	if !dryRun {
		t.Fatalf("expected dry_run true in output: %s", trimmed)
	}

	scope, _ := data["scope"].(string)
	if scope != "environment" {
		t.Fatalf("expected scope 'environment', got %q", scope)
	}

	envName, _ := data["environment"].(string)
	if envName != "staging" {
		t.Fatalf("expected environment 'staging', got %q", envName)
	}

	changes, ok := data["changes"].([]any)
	if !ok || len(changes) == 0 {
		t.Fatalf("expected at least one change in output: %s", trimmed)
	}
}

func TestCLIEditEnvironment_DryRunTextOutput(t *testing.T) {
	fx := newEditProjectFixture(t, "envdryruntext")

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs([]string{
		"edit",
		fx.name,
		"staging",
		"--env-color", "241",
		"--dry-run",
		"--output", "text",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute env dry-run text: %v (stderr=%s)", err, errBuf.String())
	}

	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %s", errBuf.String())
	}

	hcl := readFile(t, filepath.Join(fx.projectDir, ".project.hcl"))
	if strings.Contains(hcl, "color = \"241\"") {
		t.Fatalf("dry-run should not persist environment changes: %s", hcl)
	}

	stdout := out.String()
	header := "Previewing environment edit for \"staging\" in project \"" + fx.name + "\" (dry-run)"
	if !strings.Contains(stdout, header) {
		t.Fatalf("stdout missing dry-run header: %s", stdout)
	}

	if !strings.Contains(stdout, "No changes have been written; showing preview only.") {
		t.Fatalf("stdout missing dry-run disclaimer: %s", stdout)
	}

	if !strings.Contains(stdout, "Planned changes (1 change):") {
		t.Fatalf("stdout missing change summary: %s", stdout)
	}

	if !strings.Contains(stdout, "Color: \"240\" -> \"241\"") {
		t.Fatalf("stdout missing color change: %s", stdout)
	}

	if !strings.Contains(stdout, "Other environments unchanged; project metadata untouched.") {
		t.Fatalf("stdout missing unchanged scope message: %s", stdout)
	}
}
