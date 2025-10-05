package domain

// ProjectDeletePlan represents the set of actions computed for a delete
// operation. It is used for both dry-run previews and as an execution recipe.
type ProjectDeletePlan struct {
	Scope     DeleteScope
	Artifacts []DeletionArtifact
	Backup    *BackupArtifact
}

// WithArtifact appends a new artifact to the plan and returns the updated plan.
func (p ProjectDeletePlan) WithArtifact(artifact DeletionArtifact) ProjectDeletePlan {
	p.Artifacts = append(p.Artifacts, artifact)
	return p
}

// HasBackup reports whether the plan includes backup work.
func (p ProjectDeletePlan) HasBackup() bool {
	return p.Backup != nil
}
