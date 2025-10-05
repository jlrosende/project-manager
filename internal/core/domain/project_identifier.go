package domain

// ProjectIdentifier captures the metadata resolved prior to executing a delete
// operation. It is a distilled view of a registered project that exposes only
// the information required for destructive workflows.
type ProjectIdentifier struct {
	Name       string
	Path       string
	RegistryID string
	Status     ProjectStatus
}

// IsZero reports whether the identifier contains meaningful data.
func (id ProjectIdentifier) IsZero() bool {
	return id.Name == "" && id.Path == "" && id.RegistryID == ""
}

// ProjectStatus represents current lifecycle flags for a project.
type ProjectStatus struct {
	Locked   bool
	ReadOnly bool
}

// IsActionable reports whether the project is safe to modify.
func (s ProjectStatus) IsActionable() bool {
	return !s.Locked && !s.ReadOnly
}
