package domain

// ProjectExistence describes whether a project entry is present in the registry
// and whether its on-disk metadata file still exists.
type ProjectExistence struct {
	Name              string
	RegistryHit       bool
	ProjectPath       string
	ProjectFileExists bool
}
