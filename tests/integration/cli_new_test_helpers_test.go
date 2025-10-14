//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func runNewCommand(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	cmd := cli.Root()
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)

	executeErr := cmd.Execute()

	return out.String(), errBuf.String(), executeErr
}
