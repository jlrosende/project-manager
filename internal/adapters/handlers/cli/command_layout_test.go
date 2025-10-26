package cli_test

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cmdDelete "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/delete"
	cmdEdit "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/edit"
	cmdInit "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/init"
	cmdList "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/list"
	cmdNew "github.com/jlrosende/project-manager/internal/adapters/handlers/cli/new"
)

func TestCommandsExposeStandardLayout(t *testing.T) {
	tests := []struct {
		name      string
		factory   func() *cobra.Command
		expectUse string
	}{
		{name: "delete", factory: cmdDelete.Command, expectUse: "delete"},
		{name: "edit", factory: cmdEdit.Command, expectUse: "edit"},
		{name: "init", factory: cmdInit.Command, expectUse: "init"},
		{name: "list", factory: cmdList.Command, expectUse: "list"},
		{name: "new", factory: cmdNew.Command, expectUse: "new"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.factory()

			if cmd == nil {
				t.Fatalf("Command() returned nil")
			}

			if cmd.Run != nil {
				t.Fatalf("expected Run to be nil; command should use RunE")
			}

			if cmd.RunE == nil {
				t.Fatalf("RunE must be defined")
			}

			if !strings.HasPrefix(cmd.Use, tt.expectUse) {
				t.Fatalf("expected Use to start with %q, got %q", tt.expectUse, cmd.Use)
			}

			if !cmd.SilenceUsage {
				t.Fatalf("SilenceUsage must be true")
			}

			if cmd.TraverseChildren {
				t.Fatalf("TraverseChildren should be false for leaf commands")
			}
		})
	}
}
