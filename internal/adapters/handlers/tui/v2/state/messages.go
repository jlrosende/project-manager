package state

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/router"
)

type NavigateToMsg struct {
	Route  router.Route
	Params map[string]string
}

type BackMsg struct{}

type ErrorMsg struct{ Err error }

type NewProjectMsg struct{}

type SelectProjectMsg struct{ ID string }

type SaveProjectMsg struct {
	Name          string
	Path          string
	Subproject    string
	Shell         string
	GitUserName   string
	GitUserEmail  string
	GitSigningKey string
	CommitGPGSign string
	TagGPGSign    string
	EnvVarsFile   string
	EnvVarsRaw    string
}

type CancelMsg struct{}

type ExitMsg struct{}

type AddEnvVarMsg struct{}

type EditEnvVarMsg struct{ Name string }

type EditProjectMsg struct{ ID string }

type RemoveEnvVarMsg struct{ Name string }

type SaveEnvVarMsg struct{ Name, Value string }

type SaveEnvFormMsg struct {
	OriginalName string
	Name         string
	Color        string
	Mode         string
	EnvVarsFile  string
	EnvVarsRaw   string
}

type PostCreateChoiceMsg struct {
	Choice      int
	ProjectName string
}

var _ tea.Msg = SaveEnvFormMsg{}

var (
	_ tea.Msg = NavigateToMsg{}
	_ tea.Msg = BackMsg{}
	_ tea.Msg = ErrorMsg{}
	_ tea.Msg = NewProjectMsg{}
	_ tea.Msg = SelectProjectMsg{}
	_ tea.Msg = SaveProjectMsg{}
	_ tea.Msg = CancelMsg{}
	_ tea.Msg = ExitMsg{}
	_ tea.Msg = AddEnvVarMsg{}
	_ tea.Msg = EditEnvVarMsg{}
	_ tea.Msg = RemoveEnvVarMsg{}
	_ tea.Msg = SaveEnvVarMsg{}
	_ tea.Msg = EditProjectMsg{}
	_ tea.Msg = PostCreateChoiceMsg{}
)
