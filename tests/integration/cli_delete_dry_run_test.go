//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_DryRunAllScope(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if !opts.DryRun {
				t.Fatalf("expected DryRun true, got false")
			}
			if opts.Scope != domain.DeleteScopeAll {
				t.Fatalf("expected DeleteScopeAll, got %v", opts.Scope)
			}
			plan := &domain.ProjectDeletePlan{
				Scope: domain.DeleteScopeAll,
				Artifacts: []domain.DeletionArtifact{
					{Type: domain.ArtifactWorkspace, Description: "workspace"},
					{Type: domain.ArtifactEnvVars, Description: "env"},
				},
			}
			return &domain.ProjectDeleteResult{Plan: plan, Summary: "dry run"}, nil
		},
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--dry-run", "--all")
	if err != nil {
		t.Fatalf("execute dry-run: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected dry-run exit code 0, got %d", code)
	}

	if !strings.Contains(stdout, "workspace") || !strings.Contains(stdout, "No changes were applied") {
		t.Fatalf("expected dry-run preview output, got: %s", stdout)
	}
}
