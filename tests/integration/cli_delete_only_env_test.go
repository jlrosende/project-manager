//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_OnlyEnvScope(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if opts.Scope != domain.DeleteScopeEnvOnly {
				t.Fatalf("expected env-only scope, got %v", opts.Scope)
			}
			return &domain.ProjectDeleteResult{
				Scope:            opts.Scope,
				ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactEnvVars, Description: "env"}},
				ArtifactsSkipped: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry, Description: "registry"}},
				Summary:          "Environment secrets removed",
			}, nil
		},
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--only-env")
	if err != nil {
		t.Fatalf("only-env execution: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(stdout, "Environment secrets removed") {
		t.Fatalf("stdout missing env summary: %s", stdout)
	}
}
