//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_ConflictingScopes(t *testing.T) {
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			t.Fatalf("service should not be called when flags conflict: %+v", opts)
			return nil, nil
		},
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--all", "--keep-files")
	if err == nil {
		t.Fatalf("expected conflict error, got nil (code=%d)", code)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit on conflict, got %d", code)
	}

	if stderr == "" {
		t.Fatal("expected informative conflict error message")
	}
}
