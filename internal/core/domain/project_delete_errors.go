package domain

import "errors"

var (
	// ErrProjectNotFound indicates the requested project could not be resolved.
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectLocked signals the project is currently locked by another
	// operation.
	ErrProjectLocked = errors.New("project is locked")
	// ErrWorkspaceReadOnly indicates the underlying filesystem is read-only.
	ErrWorkspaceReadOnly = errors.New("workspace is read-only")
	// ErrPartialDeletion is returned when some artifacts fail to delete.
	ErrPartialDeletion = errors.New("project deletion completed with errors")
)
