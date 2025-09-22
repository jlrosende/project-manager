//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/state"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestTUI_EditProject_RenameReloadsList(t *testing.T) {
	svc := buildService(t)

	w, err := tui.NewWindow(svc.(*services.ProjectService), tui.Options{})
	if err != nil {
		t.Fatalf("window: %v", err)
	}

	if cmd := w.Init(); cmd != nil {
		w.Update(cmd())
	}

	w.Update(tea.WindowSizeMsg{Width: 120, Height: 24})

	w.Update(state.EditProjectMsg{ID: "INDITEX"})
	w.Update(state.SaveProjectMsg{
		Name:        "INDITEX-NEW",
		Path:        "/tmp/inditex",
		Shell:       "/bin/sh",
		EnvVarsFile: ".env",
	})

	view := normalize(w.View())
	if !strings.Contains(view, "INDITEX-NEW") {
		t.Fatalf("projects list not refreshed with new name; view=\n%s", view)
	}
}
