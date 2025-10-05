//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_KeepFilesScope(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if opts.Scope != domain.DeleteScopeKeepFiles {
				t.Fatalf("expected keep-files scope, got %v", opts.Scope)
			}
			return &domain.ProjectDeleteResult{
				Scope:            opts.Scope,
				ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry}},
				ArtifactsSkipped: []domain.DeletionArtifact{{Type: domain.ArtifactWorkspace, Description: "preserved"}},
				Summary:          "Metadata removed, files kept",
			}, nil
		},
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--keep-files")
	if err != nil {
		t.Fatalf("execute keep-files: %v (stderr=%s)", err, stderr)
	}

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(stdout, "files kept") {
		t.Fatalf("expected output noting preserved files: %s", stdout)
	}
}
