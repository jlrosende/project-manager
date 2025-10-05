//go:build integration
// +build integration

package integration_test

import (
	"errors"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestCLIDelete_BackupFailureAborts(t *testing.T) {
	expectedErr := errors.New("backup failed")
	svc := &fakeDeleteService{
		handler: func(opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
			if opts.Backup == nil {
				t.Fatal("expected backup request")
			}
			return nil, expectedErr
		},
	}

	code, _, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--backup", "--all")
	if err == nil {
		t.Fatalf("expected error, got nil (code=%d)", code)
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected backup error, got %v", err)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit code on backup failure, got %d", code)
	}

	if len(stderr) == 0 {
		t.Fatal("expected stderr output for backup failure")
	}
}
