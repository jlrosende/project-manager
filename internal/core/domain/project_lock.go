package domain

import "errors"

// ProjectLock represents an acquired edit lock for a project. The lock is
// backed by a filesystem sentinel and must be released once the edit completes
// to avoid blocking subsequent commands.
type ProjectLock struct {
	Path    string
	release func() error
}

var errNilRelease = errors.New("project lock has no release handler")

// NewProjectLock constructs a ProjectLock with the supplied release handler.
// A nil release handler results in a lock that cannot be released and will
// surface an error when Release is invoked.
func NewProjectLock(path string, release func() error) *ProjectLock {
	return &ProjectLock{
		Path:    path,
		release: release,
	}
}

// Release clears the underlying filesystem sentinel. Calls on a nil lock are
// treated as no-ops for convenience.
func (l *ProjectLock) Release() error {
	if l == nil {
		return nil
	}

	if l.release == nil {
		return errNilRelease
	}

	return l.release()
}
