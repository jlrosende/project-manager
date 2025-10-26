//go:build integration
// +build integration

package integration_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_PartialCleanupFailure(t *testing.T) {
	reason := errors.New("failed to remove workspace")
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return &domain.ProjectDeleteResult{
				Scope:            opts.Scope,
				ArtifactsRemoved: []domain.DeletionArtifact{{Type: domain.ArtifactRegistry}},
				ArtifactsSkipped: []domain.DeletionArtifact{{Type: domain.ArtifactWorkspace, Description: "workspace"}},
				Errors:           []error{reason},
				Summary:          "partial failure",
			}, domain.ErrPartialDeletion
		},
	}

	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--all")
	if err == nil {
		t.Fatalf("expected error for partial failure, got nil (code=%d)", code)
	}

	if !errors.Is(err, domain.ErrPartialDeletion) {
		t.Fatalf("expected ErrPartialDeletion, got %v", err)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit code, got %d", code)
	}

	joined := stdout + stderr
	if !strings.Contains(joined, reason.Error()) {
		t.Fatalf("expected failure reason in output, got: %s | %s", stdout, stderr)
	}
}
