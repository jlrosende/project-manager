//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_MissingTargetError(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			return nil, domain.ErrProjectNotFound
		},
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, nil, "ghost-app")
	if err == nil {
		t.Fatalf("expected missing-target error, got nil (code=%d)", code)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit code for missing target, got %d", code)
	}

	if stderr == "" {
		t.Fatal("expected descriptive missing target error message")
	}
}
