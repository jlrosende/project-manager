package domain

// ArtifactType identifies a class of resource affected during deletion.
type ArtifactType string

const (
	ArtifactRegistry   ArtifactType = "registry"
	ArtifactEnvVars    ArtifactType = "env-vars"
	ArtifactConfig     ArtifactType = "config"
	ArtifactGitInclude ArtifactType = "git-include"
	ArtifactGitIgnore  ArtifactType = "git-ignore"
	ArtifactWorkspace  ArtifactType = "workspace"
	ArtifactSkeleton   ArtifactType = "skeleton"
	ArtifactCache      ArtifactType = "cache"
	ArtifactHook       ArtifactType = "hook"
	ArtifactBackup     ArtifactType = "backup"
)

// DeletionArtifact describes an individual resource planned for deletion. It is
// used both for previews and execution summaries.
type DeletionArtifact struct {
	Type        ArtifactType
	Path        string
	Description string
}
