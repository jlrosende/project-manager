package domain

// BackupArtifact captures metadata about the archive produced during deletion.
type BackupArtifact struct {
	TempPath  string
	FinalPath string
	Created   bool
	SizeBytes int64
}

// Completed reports whether the archive was successfully created.
func (b *BackupArtifact) Completed() bool {
	return b != nil && b.Created && b.FinalPath != ""
}
