//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_PathResolution(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return &domain.ProjectDeleteResult{
				Scope:            domain.DeleteScopeMetadata,
				ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry}},
				Summary:          "path delete",
			}, nil
		},
	}

	path := "/workspaces/sample-app"
	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, path)
	if err != nil {
		t.Fatalf("execute delete: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	if svc.lastOpt.Target.Path != path {
		t.Fatalf("expected path to propagate into service, got %+v", svc.lastOpt.Target)
	}

	if svc.lastOpt.Target.Name != "" {
		t.Fatalf("path invocation should not infer name yet, got %q", svc.lastOpt.Target.Name)
	}

	if !strings.Contains(stdout, "path delete") {
		t.Fatalf("missing summary in stdout: %s", stdout)
	}
}
