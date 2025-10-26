//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLINew_PerProjectGitConfigIncludeIf(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))

	initialGit := []byte("[user]\n\tname = Existing User\n\temail = user@example.com\n")
	gitconfig := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(gitconfig, initialGit, 0o600); err != nil {
		t.Fatalf("write seed gitconfig: %v", err)
	}

	projectDir := filepath.Join(home, "workspace", "gitcfg")
	args := []string{"new", "gitcfg", projectDir}

	_, stderr, err := runNewCommand(t, args...)
	if err != nil {
		t.Fatalf("pm new failed: %v; stderr=%s", err, stderr)
	}

	data, readErr := os.ReadFile(gitconfig)
	if readErr != nil {
		t.Fatalf("read gitconfig: %v", readErr)
	}

	includeHeader := `[includeIf "gitdir/i:` + projectDir + `/"]`
	if !bytes.Contains(data, []byte(includeHeader)) {
		t.Fatalf("missing includeIf header for project: got %s", string(data))
	}

	perProject := filepath.Join(projectDir, ".gitcfg.gitconfig")
	expectPathLine := "path = " + perProject
	if !strings.Contains(string(data), expectPathLine) {
		t.Fatalf("missing per-project path entry: expected line %q in %s", expectPathLine, string(data))
	}

	if _, statErr := os.Stat(perProject); statErr != nil {
		t.Fatalf("per-project gitconfig not created: %v", statErr)
	}
}
