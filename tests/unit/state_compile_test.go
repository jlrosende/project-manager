//go:build unit
// +build unit

package unit_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/router"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/state"
)

func TestMessagesCompile(_ *testing.T) {
	_ = state.NavigateToMsg{Route: router.RouteProjects}
	_ = state.BackMsg{}
	_ = state.ErrorMsg{}
	_ = state.NewProjectMsg{}
	_ = state.SelectProjectMsg{ID: "x"}
	_ = state.SaveProjectMsg{Name: "a", Path: "b"}
	_ = state.CancelMsg{}
	_ = state.AddEnvVarMsg{}
	_ = state.EditEnvVarMsg{Name: "dev"}
	_ = state.RemoveEnvVarMsg{Name: "dev"}
	_ = state.SaveEnvVarMsg{Name: "k", Value: "v"}
	_ = state.SaveEnvFormMsg{OriginalName: "dev", Name: "dev", Color: "#fff", Mode: "file", EnvVarsRaw: "X=1"}
	_ = state.PostCreateChoiceMsg{Choice: 0, ProjectName: "p"}

	var _ tea.Msg = state.NavigateToMsg{}
}
