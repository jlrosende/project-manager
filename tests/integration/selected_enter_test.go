//go:build integration
// +build integration

package integration_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func TestSelectProjectWithDefaultOnEnter(t *testing.T) {
	svc := buildService(t)

	w, err := tui.NewWindow(svc.(*services.ProjectService), tui.Options{})
	if err != nil {
		t.Fatalf("v2 window: %v", err)
	}

	if cmd := w.Init(); cmd != nil {
		if msg := cmd(); msg != nil {
			_, _ = w.Update(msg)
		}
	}

	_, _ = w.Update(tea.WindowSizeMsg{Width: 49, Height: 24})

	msgs := []tea.Msg{tea.KeyMsg{Type: tea.KeyEnter}}
	for _, m := range msgs {
		_, cmd := w.Update(m)
		if cmd != nil {
			if msg := cmd(); msg != nil {
				_, _ = w.Update(msg)
			}
		}
	}

	sp := w.SelectedProject()
	if sp == nil || sp.Name != "INDITEX" {
		t.Fatalf("expected selected project 'INDITEX'; got %#v", sp)
	}

	if env := w.SelectedEnvironment(); env != "dev" {
		t.Fatalf("expected selected environment 'dev' on default; got %q", env)
	}
}

func TestSelectProjectWithoutDefaultOnEnter(t *testing.T) {
	svc := buildService(t)

	w, err := tui.NewWindow(svc.(*services.ProjectService), tui.Options{})
	if err != nil {
		t.Fatalf("v2 window: %v", err)
	}

	if cmd := w.Init(); cmd != nil {
		if msg := cmd(); msg != nil {
			_, _ = w.Update(msg)
		}
	}

	_, _ = w.Update(tea.WindowSizeMsg{Width: 49, Height: 24})

	msgs := []tea.Msg{tea.KeyMsg{Type: tea.KeyDown}, tea.KeyMsg{Type: tea.KeyEnter}}
	for _, m := range msgs {
		_, cmd := w.Update(m)
		if cmd != nil {
			if msg := cmd(); msg != nil {
				_, _ = w.Update(msg)
			}
		}
	}

	sp := w.SelectedProject()
	if sp == nil || sp.Name != "Mahou" {
		t.Fatalf("expected selected project 'Mahou'; got %#v", sp)
	}

	if env := w.SelectedEnvironment(); env != "" {
		t.Fatalf("expected empty environment when no default; got %q", env)
	}
}

func TestSelectEnvironmentOnEnter(t *testing.T) {
	svc := buildService(t)

	w, err := tui.NewWindow(svc.(*services.ProjectService), tui.Options{})
	if err != nil {
		t.Fatalf("v2 window: %v", err)
	}

	if cmd := w.Init(); cmd != nil {
		if msg := cmd(); msg != nil {
			_, _ = w.Update(msg)
		}
	}

	_, _ = w.Update(tea.WindowSizeMsg{Width: 49, Height: 24})

	msgs := []tea.Msg{tea.KeyMsg{Type: tea.KeyRight}, tea.KeyMsg{Type: tea.KeyDown}, tea.KeyMsg{Type: tea.KeyEnter}}
	for _, m := range msgs {
		_, cmd := w.Update(m)
		if cmd != nil {
			if msg := cmd(); msg != nil {
				_, _ = w.Update(msg)
			}
		}
	}

	sp := w.SelectedProject()
	if sp == nil || sp.Name != "INDITEX" {
		t.Fatalf("expected selected project 'INDITEX'; got %#v", sp)
	}

	if env := w.SelectedEnvironment(); env != "pre" {
		t.Fatalf("expected selected environment 'pre'; got %q", env)
	}
}
