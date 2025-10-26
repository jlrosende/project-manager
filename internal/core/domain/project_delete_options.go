package domain

import "errors"

// ProjectDeleteOptions aggregates all inputs required to perform a delete
// operation.
type ProjectDeleteOptions struct {
	Target ProjectIdentifier
	Scope  DeleteScope
	DryRun bool
	Force  bool
	Backup *BackupRequest
}

// BackupRequest captures user intent for creating an archive before deleting.
type BackupRequest struct {
	Destination      string
	IncludeWorkspace bool
}

// ErrInvalidDeleteScope indicates mutually exclusive flags were combined or an
// unsupported scope was requested.
var ErrInvalidDeleteScope = errors.New("invalid delete scope")

// Validate ensures the delete request is coherent before execution.
func (o ProjectDeleteOptions) Validate() error {
	if o.Target.IsZero() {
		return errors.New("delete target is required")
	}

	if !o.Scope.IsValid() {
		return ErrInvalidDeleteScope
	}

	if o.Backup != nil && o.Backup.Destination == "" {
		return errors.New("backup destination must not be empty")
	}

	return nil
}

// RequiresBackup reports whether the request includes a backup step.
func (o ProjectDeleteOptions) RequiresBackup() bool {
	return o.Backup != nil
}
