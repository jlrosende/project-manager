package state

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/router"
)

func TestMessagesCompile(_ *testing.T) {
	_ = NavigateToMsg{Route: router.RouteProjects}
	_ = BackMsg{}
	_ = ErrorMsg{}
	_ = NewProjectMsg{}
	_ = SelectProjectMsg{ID: "x"}
	_ = SaveProjectMsg{Name: "a", Path: "b"}
	_ = CancelMsg{}
	_ = AddEnvVarMsg{}
	_ = EditEnvVarMsg{Name: "dev"}
	_ = RemoveEnvVarMsg{Name: "dev"}
	_ = SaveEnvVarMsg{Name: "k", Value: "v"}
	_ = SaveEnvFormMsg{OriginalName: "dev", Name: "dev", Color: "#fff", Mode: "file", EnvVarsRaw: "X=1"}
	_ = PostCreateChoiceMsg{Choice: 0, ProjectName: "p"}

	var _ tea.Msg = NavigateToMsg{}
}
