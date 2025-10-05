//go:build integration
// +build integration

package integration_test

import (
	"context"
	"testing"

	deletecmd "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

type fakeDeleteService struct {
	calls   int
	lastOpt domain.ProjectDeleteOptions
	handler func(domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error)
	result  *domain.ProjectDeleteResult
	err     error
}

func (f *fakeDeleteService) DeleteProject(ctx context.Context, opts domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error) {
	f.calls++
	f.lastOpt = opts

	if f.handler != nil {
		return f.handler(opts)
	}

	if f.err != nil {
		return nil, f.err
	}

	return f.result, nil
}

func executeDeleteCommand(t *testing.T, svc *fakeDeleteService, confirm deletecmd.ConfirmationFunc, args ...string) (int, string, string, error) {
	t.Helper()

	if confirm == nil {
		confirm = func(context.Context, deletecmd.ConfirmationRequest) (bool, error) { return true, nil }
	}

	code, stdout, stderr, err := deletecmd.ExecuteForTesting(deletecmd.Options{
		Service: svc,
		Confirm: confirm,
	}, args)

	return code, stdout, stderr, err
}
