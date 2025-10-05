package ports

//go:generate go tool mockgen -source=filesystem_port.go -destination=../../../mocks/mock_filesystem_port.go -package=mocks

import (
	"context"
	"io/fs"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

// Filesystem defines operating system interactions required by the application
// layer. Concrete adapters should provide the actual OS access.
type Filesystem interface {
	EnsureDir(path string, mode fs.FileMode) error
	IsDirEmpty(path string) (bool, error)
	Rename(oldPath, newPath string) error
	WriteFile(path string, data []byte, mode fs.FileMode) error
	Join(elem ...string) string
	IsAbs(path string) bool
	Abs(path string) (string, error)
	ExpandHome(path string) string
	Remove(path string) error
	UserHomeDir() (string, error)

	PlanDeletion(
		ctx context.Context,
		target domain.ProjectIdentifier,
		scope domain.DeleteScope,
		backup *domain.BackupRequest,
	) (*domain.ProjectDeletePlan, error)
	PlanBackup(
		ctx context.Context,
		target domain.ProjectIdentifier,
		req *domain.BackupRequest,
		dryRun bool,
	) (*domain.BackupArtifact, error)
	CreateBackup(
		ctx context.Context,
		target domain.ProjectIdentifier,
		req *domain.BackupRequest,
	) (*domain.BackupArtifact, error)
	ExecuteDeletion(ctx context.Context, plan *domain.ProjectDeletePlan) ([]domain.DeletionArtifact, error)
}
