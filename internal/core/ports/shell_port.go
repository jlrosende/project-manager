package ports

import "github.com/jlrosende/project-manager/internal/core/domain"

type ShellService interface {
	Start(project domain.Project) error
}

type ShellRepository interface {
	Start(project domain.Project) error
}
