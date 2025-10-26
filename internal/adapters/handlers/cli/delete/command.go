package delete

import (
	"bytes"
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/configs"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

const (
	flagAll               = "all"
	flagKeepFiles         = "keep-files"
	flagOnlyEnv           = "only-env"
	flagDryRun            = "dry-run"
	flagForce             = "force"
	flagBackup            = "backup"
	flagBackupDestination = "backup-destination"
)

// Service defines the subset of the project service required by the CLI
// delete command.
type Service interface {
	DeleteProject(context.Context, domain.ProjectDeleteOptions) (*domain.ProjectDeleteResult, error)
}

// Options configures command execution for testing and dependency injection.
type Options struct {
	Service Service
	Confirm ConfirmationFunc
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
	Palette map[string]string
	Config  *configs.Config
}

// Command constructs the `pm delete` Cobra command using the shared CLI layout.
func Command() *cobra.Command {
	cmd := newBaseCommand()
	cmd.RunE = run

	return cmd
}

func newBaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete <target>",
		Short:        "Delete a registered project and its artifacts",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().Bool(flagAll, false, "Remove project metadata, env vars, and workspace files")
	cmd.Flags().Bool(flagKeepFiles, false, "Remove metadata but keep workspace files")
	cmd.Flags().Bool(flagOnlyEnv, false, "Remove stored environment variables only")
	cmd.Flags().Bool(flagDryRun, false, "Preview deletion steps without making changes")
	cmd.Flags().Bool(flagForce, false, "Skip confirmation prompt")
	cmd.Flags().Bool(flagBackup, false, "Create a backup archive before deleting")
	cmd.Flags().String(flagBackupDestination, "", "Custom destination for the backup archive (requires --backup)")

	return cmd
}

// ExecuteForTesting executes the command using injected dependencies and
// returns the exit code, stdout, stderr, and error for assertions in tests.
func ExecuteForTesting(opts Options, args []string) (int, string, string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	opts.Out = &stdoutBuf

	opts.Err = &stderrBuf
	if opts.In == nil {
		opts.In = bytes.NewBuffer(nil)
	}

	cmd := newBaseCommand()
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runWithOptions(cmd, args, opts)
	}

	cmd.SetOut(&stdoutBuf)
	cmd.SetErr(&stderrBuf)
	cmd.SetIn(opts.In)
	cmd.SetArgs(args)

	err := cmd.Execute()
	if err != nil {
		return 1, stdoutBuf.String(), stderrBuf.String(), err
	}

	return 0, stdoutBuf.String(), stderrBuf.String(), nil
}
