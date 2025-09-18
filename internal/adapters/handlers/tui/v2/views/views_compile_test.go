package views

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/v2/components"
)

func TestProjectsViewCompile(_ *testing.T) {
	_ = NewProjectsView(nil, lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.Color("240"))
}

func TestProjectFormViewCompile(_ *testing.T) {
	_ = NewProjectFormView(
		"",
		"",
		lipgloss.NewStyle(),
		lipgloss.NewStyle(),
		lipgloss.NewStyle(),
		lipgloss.NewStyle(),
		lipgloss.NewStyle(),
		lipgloss.NewStyle(),
	)
}

func TestEnvVarsViewCompile(_ *testing.T) {
	items := []components.ListItem{{ID: "__add__", Label: "+ Add environment"}}
	_ = NewEnvVarsView(items, lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.Color("240"))
}

func TestEnvFormViewCompile(_ *testing.T) {
	_ = NewEnvFormView("", "", "", "file", "", lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle())
}

func TestPostCreateViewCompile(_ *testing.T) {
	_ = NewPostCreateView("proj", lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle())
}
