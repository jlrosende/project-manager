package domain

import "fmt"

// DeleteScope enumerates the different deletion behaviours supported by the
// project manager. Scopes determine which artifacts are removed when executing
// a delete request.
type DeleteScope int

const (
	// DeleteScopeMetadata removes registry metadata and configuration files while
	// keeping workspace contents intact. This is the default when no explicit
	// flag is provided.
	DeleteScopeMetadata DeleteScope = iota
	// DeleteScopeKeepFiles mirrors DeleteScopeMetadata but is used when the user
	// explicitly requests "--keep-files". Keeping a dedicated constant simplifies
	// downstream logging and summaries.
	DeleteScopeKeepFiles
	// DeleteScopeEnvOnly removes only environment data linked to the project.
	DeleteScopeEnvOnly
	// DeleteScopeAll removes every artifact associated with the project including
	// workspace directories and caches.
	DeleteScopeAll
)

// String implements fmt.Stringer for improved logging output.
func (s DeleteScope) String() string {
	switch s {
	case DeleteScopeMetadata:
		return "metadata"
	case DeleteScopeKeepFiles:
		return "keep-files"
	case DeleteScopeEnvOnly:
		return "env-only"
	case DeleteScopeAll:
		return "all"
	default:
		return fmt.Sprintf("unknown-scope-%d", int(s))
	}
}

// IsValid reports whether the scope is recognised.
func (s DeleteScope) IsValid() bool {
	switch s {
	case DeleteScopeMetadata, DeleteScopeKeepFiles, DeleteScopeEnvOnly, DeleteScopeAll:
		return true
	default:
		return false
	}
}

// IncludesWorkspace indicates whether the scope includes removing workspace
// directories and user-authored files.
func (s DeleteScope) IncludesWorkspace() bool {
	return s == DeleteScopeAll
}

// IncludesEnvironment indicates whether environment secrets should be removed
// for the provided scope.
func (s DeleteScope) IncludesEnvironment() bool {
	switch s {
	case DeleteScopeMetadata, DeleteScopeKeepFiles, DeleteScopeEnvOnly, DeleteScopeAll:
		return true
	default:
		return false
	}
}

// IncludesMetadata signals if registry metadata must be cleared.
func (s DeleteScope) IncludesMetadata() bool {
	switch s {
	case DeleteScopeMetadata, DeleteScopeKeepFiles, DeleteScopeAll:
		return true
	default:
		return false
	}
}
