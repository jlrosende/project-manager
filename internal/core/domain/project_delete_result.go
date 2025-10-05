package domain

// ProjectDeleteResult summarises the outcome of a delete operation.
type ProjectDeleteResult struct {
	Scope            DeleteScope
	Plan             *ProjectDeletePlan
	Summary          string
	BackupPath       string
	ArtifactsRemoved []DeletionArtifact
	ArtifactsSkipped []DeletionArtifact
	Errors           []error
}

// HasErrors reports whether the delete encountered issues.
func (r *ProjectDeleteResult) HasErrors() bool {
	if r == nil {
		return false
	}

	return len(r.Errors) > 0
}
