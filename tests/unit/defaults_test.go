package unit_test

import (
	"testing"

	cli "github.com/jlrosende/project-manager/internal/adapters/handlers/cli"
)

func TestCLINew_DefaultFlagValues(t *testing.T) {
	root := cli.Root()

	newCmd, _, err := root.Find([]string{"new"})
	if err != nil {
		t.Fatalf("find new command: %v", err)
	}

	cases := []struct {
		name     string
		flagName string
		want     string
	}{
		{name: "dry-run", flagName: "dry-run", want: "false"},
		{name: "force", flagName: "force", want: "false"},
		{name: "allow-unknown", flagName: "allow-unknown", want: "false"},
		{name: "output", flagName: "output", want: "text"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			flag := newCmd.Flag(tc.flagName)
			if flag == nil {
				t.Fatalf("flag %s not defined", tc.flagName)
			}

			if flag.DefValue != tc.want {
				t.Fatalf("default for %s = %s, want %s", tc.flagName, flag.DefValue, tc.want)
			}
		})
	}
}
