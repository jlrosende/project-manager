//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"
)

func TestCLIDelete_BackupDestinationRequiresBackup(t *testing.T) {
	svc := &fakeDeleteService{}

	destination := "/backups/custom.zip"
	code, stdout, stderr, err := executeDeleteCommand(t, svc, nil, "sample-app", "--backup-destination", destination)
	if err == nil {
		t.Fatalf("expected error when using --backup-destination without --backup (code=%d, stdout=%s)", code, stdout)
	}

	if code == 0 {
		t.Fatalf("expected non-zero exit code when backup destination is missing --backup, got %d", code)
	}

	if svc.calls != 0 {
		t.Fatalf("expected service not to be called, got %d call(s)", svc.calls)
	}

	expected := "--backup-destination requires --backup"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error %q, got %v", expected, err)
	}

	if !strings.Contains(stderr, expected) {
		t.Fatalf("expected stderr to mention %q, got %s", expected, stderr)
	}
}
