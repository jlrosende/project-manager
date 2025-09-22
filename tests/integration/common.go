//go:build integration
// +build integration

package integration_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
	"github.com/jlrosende/project-manager/internal/core/services"
	"github.com/jlrosende/project-manager/mocks"
)

func buildService(t *testing.T) ports.ProjectService {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	mockEnv := mocks.NewMockEnvVarsRepository(ctrl)
	mockGit := mocks.NewMockGitRepository(ctrl)

	p1 := &domain.Project{
		Name:       "INDITEX",
		Path:       "/tmp/inditex",
		DefaultEnv: "dev",
		Environments: []*domain.Environment{
			{Name: "dev", Color: "240", EnvVarsFile: ".env.dev"},
			{Name: "pre", Color: "240", EnvVarsFile: ".env.pre"},
			{Name: "pro", Color: "240", EnvVarsFile: ".env.pro"},
		},
	}
	p2 := &domain.Project{Name: "Mahou"}
	p3 := &domain.Project{Name: "Accenture"}
	p4 := &domain.Project{Name: "Personal"}
	p5 := &domain.Project{Name: "test"}
	projects := []*domain.Project{p1, p2, p3, p4, p5}

	mockRepo.EXPECT().List().Return(projects, nil).AnyTimes()
	mockRepo.EXPECT().AddEnvironment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockRepo.EXPECT().UpdateEnvironment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockRepo.EXPECT().UpdateProject(gomock.Any()).DoAndReturn(func(p *domain.Project) error {
		for _, pr := range projects {
			if pr.Path == p.Path {
				pr.Name = p.Name
			}
		}
		return nil
	}).AnyTimes()
	mockEnv.EXPECT().Load(gomock.Any()).Return(domain.EnvVars{}, nil).AnyTimes()
	mockGit.EXPECT().Load(gomock.Any()).Return(&domain.GitConfig{}, nil).AnyTimes()

	return services.NewProjectService(mockRepo, mockEnv, mockGit)
}

func normalize(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}

	return strings.Join(lines, "\n")
}

func render(t *testing.T, svc ports.ProjectService, width int, focusRight bool) string {
	t.Helper()

	w, err := tui.NewWindow(svc.(*services.ProjectService), tui.Options{})
	if err != nil {
		t.Fatalf("v2 window: %v", err)
	}

	cmd := w.Init()
	if cmd != nil {
		msg := cmd()
		w.Update(msg)
	}

	w.Update(tea.WindowSizeMsg{Width: width, Height: 24})

	if focusRight {
		w.Update(tea.KeyMsg{Type: tea.KeyRight})
	}

	return w.View()
}
