package repositories

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

var _ ports.Filesystem = (*Filesystem)(nil)

// Filesystem implements the ports.Filesystem interface using the local OS.
type Filesystem struct{}

// NewFilesystem constructs a filesystem adapter backed by the standard library.
func NewFilesystem() *Filesystem { return &Filesystem{} }

func (Filesystem) EnsureDir(path string, mode fs.FileMode) error {
	if st, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, mode)
		}

		return err
	} else if !st.IsDir() {
		return &fs.PathError{Op: "mkdir", Path: path, Err: fs.ErrInvalid}
	}

	return nil
}

func (Filesystem) IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func (fsys Filesystem) Rename(oldPath, newPath string) error {
	if err := fsys.EnsureDir(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}

	return os.Rename(oldPath, newPath)
}

func (fsys Filesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	if err := fsys.EnsureDir(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, mode)
}

func (Filesystem) Join(elem ...string) string { return filepath.Join(elem...) }

func (Filesystem) IsAbs(path string) bool { return filepath.IsAbs(path) }

func (Filesystem) Abs(path string) (string, error) { return filepath.Abs(path) }

func (Filesystem) ExpandHome(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}

	return path
}

func (Filesystem) Remove(path string) error { return os.Remove(path) }

func (Filesystem) UserHomeDir() (string, error) { return os.UserHomeDir() }

func (fsys Filesystem) PlanDeletion(ctx context.Context, target domain.ProjectIdentifier, scope domain.DeleteScope, backup *domain.BackupRequest) (*domain.ProjectDeletePlan, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	root := fsys.ExpandHome(target.Path)
	if !filepath.IsAbs(root) {
		abs, err := filepath.Abs(root)
		if err != nil {
			return nil, err
		}
		root = abs
	}

	project, err := loadProjectDefinition(root)
	if err != nil {
		return nil, err
	}

	plan := &domain.ProjectDeletePlan{Scope: scope}

	if backup != nil {
		plan.Backup = &domain.BackupArtifact{FinalPath: backup.Destination}
		plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
			Type:        domain.ArtifactBackup,
			Path:        root,
			Description: "project backup",
		})
	}

	if scope.IncludesMetadata() {
		plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
			Type:        domain.ArtifactRegistry,
			Path:        target.RegistryID,
			Description: "registry entry",
		})

		configPath := filepath.Join(root, ".project.hcl")
		plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
			Type:        domain.ArtifactConfig,
			Path:        configPath,
			Description: ".project.hcl",
		})

		if include := gitIncludePath(root, project, target); include != "" {
			plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
				Type:        domain.ArtifactGitInclude,
				Path:        include,
				Description: "git include",
			})
		}

		if needs, err := fsys.gitignoreNeedsCleanup(root); err != nil {
			return nil, err
		} else if needs {
			plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
				Type:        domain.ArtifactGitIgnore,
				Path:        filepath.Join(root, ".gitignore"),
				Description: ".gitignore",
			})
		}
	}

	if scope.IncludesEnvironment() {
		envFiles := fsys.collectEnvFiles(root, project)
		for _, envPath := range envFiles {
			plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
				Type:        domain.ArtifactEnvVars,
				Path:        envPath,
				Description: filepath.Base(envPath),
			})
		}
	}

	if scope.IncludesWorkspace() {
		plan.Artifacts = append(plan.Artifacts, domain.DeletionArtifact{
			Type:        domain.ArtifactWorkspace,
			Path:        root,
			Description: "workspace",
		})
	}

	return plan, nil
}

func (Filesystem) PlanBackup(ctx context.Context, target domain.ProjectIdentifier, req *domain.BackupRequest, dryRun bool) (*domain.BackupArtifact, error) {
	if req == nil {
		return nil, nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	destination := strings.TrimSpace(req.Destination)
	if destination == "" {
		return nil, fmt.Errorf("backup destination is empty")
	}

	artifact := &domain.BackupArtifact{FinalPath: destination}
	if dryRun {
		return artifact, nil
	}

	artifact.TempPath = filepath.Join(filepath.Dir(destination), fmt.Sprintf(".%s.partial", filepath.Base(destination)))

	return artifact, nil
}

func (fsys Filesystem) CreateBackup(ctx context.Context, target domain.ProjectIdentifier, req *domain.BackupRequest) (*domain.BackupArtifact, error) {
	if req == nil {
		return nil, nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	destination := strings.TrimSpace(req.Destination)
	if destination == "" {
		return nil, fmt.Errorf("backup destination is empty")
	}

	destination = fsys.ExpandHome(destination)
	if err := fsys.EnsureDir(filepath.Dir(destination), 0o755); err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(destination), "pm-backup-*.zip")
	if err != nil {
		return nil, err
	}

	cleanup := true
	defer func() {
		if cleanup {
			_ = tmpFile.Close()
			_ = os.Remove(tmpFile.Name())
		}
	}()

	zipWriter := zip.NewWriter(tmpFile)

	root := fsys.ExpandHome(target.Path)
	if !filepath.IsAbs(root) {
		abs, err := filepath.Abs(root)
		if err != nil {
			zipWriter.Close()
			return nil, err
		}
		root = abs
	}

	project, err := loadProjectDefinition(root)
	if err != nil {
		zipWriter.Close()
		return nil, err
	}

	if req.IncludeWorkspace {
		err = fsys.addWorkspaceToZip(ctx, zipWriter, root)
	} else {
		err = fsys.addMetadataToZip(ctx, zipWriter, root, project)
	}

	if cerr := zipWriter.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return nil, err
	}

	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	if err := os.Rename(tmpFile.Name(), destination); err != nil {
		return nil, err
	}

	cleanup = false

	info, err := os.Stat(destination)
	if err != nil {
		return nil, err
	}

	return &domain.BackupArtifact{
		FinalPath: destination,
		Created:   true,
		SizeBytes: info.Size(),
	}, nil
}

func (fsys Filesystem) ExecuteDeletion(ctx context.Context, plan *domain.ProjectDeletePlan) ([]domain.DeletionArtifact, error) {
	if plan == nil {
		return nil, nil
	}

	var removed []domain.DeletionArtifact

	for _, artifact := range plan.Artifacts {
		if err := ctx.Err(); err != nil {
			return removed, err
		}

		path := fsys.ExpandHome(artifact.Path)

		switch artifact.Type {
		case domain.ArtifactConfig:
			if err := removeFile(path); err != nil {
				return removed, err
			}
			removed = append(removed, artifact)
		case domain.ArtifactGitIgnore:
			changed, err := fsys.cleanupGitignore(path)
			if err != nil {
				return removed, err
			}
			if changed {
				removed = append(removed, artifact)
			}
		case domain.ArtifactSkeleton, domain.ArtifactCache, domain.ArtifactWorkspace:
			if err := os.RemoveAll(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return removed, err
			}
			removed = append(removed, artifact)
		default:
			// handled elsewhere
		}
	}

	return removed, nil
}

func loadProjectDefinition(root string) (*domain.Project, error) {
	projectPath := filepath.Join(root, ".project.hcl")
	if _, err := os.Stat(projectPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var project domain.Project
	if err := hclsimple.DecodeFile(projectPath, nil, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func gitIncludePath(root string, project *domain.Project, target domain.ProjectIdentifier) string {
	name := strings.TrimSpace(target.Name)
	if project != nil && strings.TrimSpace(project.Name) != "" {
		name = project.Name
	}

	if name == "" {
		name = filepath.Base(root)
	}

	if name == "" {
		return ""
	}

	return filepath.Join(root, fmt.Sprintf(".%s.gitconfig", name))
}

func (fsys Filesystem) collectEnvFiles(root string, project *domain.Project) []string {
	paths := map[string]struct{}{}

	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" {
			return
		}

		candidate := fsys.ExpandHome(p)
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(root, candidate)
		}

		candidate = filepath.Clean(candidate)
		paths[candidate] = struct{}{}
	}

	if project == nil {
		add(".env")
	} else {
		add(project.EnvVarsFile)
		for _, env := range project.Environments {
			add(env.EnvVarsFile)
		}
	}

	results := make([]string, 0, len(paths))
	for path := range paths {
		results = append(results, path)
	}

	sort.Strings(results)

	return results
}

func (fsys Filesystem) gitignoreNeedsCleanup(root string) (bool, error) {
	path := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == ".env" {
			return true, nil
		}
	}

	return false, nil
}

func (fsys Filesystem) cleanupGitignore(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	lines := strings.Split(string(data), "\n")
	filtered := make([]string, 0, len(lines))
	removed := false

	for _, line := range lines {
		if strings.TrimSpace(line) == ".env" {
			removed = true
			continue
		}

		filtered = append(filtered, line)
	}

	if !removed {
		return false, nil
	}

	for len(filtered) > 0 && strings.TrimSpace(filtered[len(filtered)-1]) == "" {
		filtered = filtered[:len(filtered)-1]
	}

	if len(filtered) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return false, err
		}

		return true, nil
	}

	content := strings.Join(filtered, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, err
	}

	return true, nil
}

func (fsys Filesystem) addMetadataToZip(ctx context.Context, writer *zip.Writer, root string, project *domain.Project) error {
	files := fsys.collectEnvFiles(root, project)
	projectFile := filepath.Join(root, ".project.hcl")
	if _, err := os.Stat(projectFile); err == nil {
		files = append(files, projectFile)
	}

	dedup := map[string]struct{}{}

	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return err
		}

		if _, seen := dedup[file]; seen {
			continue
		}
		dedup[file] = struct{}{}

		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		rel, err := filepath.Rel(root, file)
		if err != nil || strings.HasPrefix(rel, "..") {
			rel = filepath.Base(file)
		}

		if err := writePathToZip(writer, file, rel, info); err != nil {
			return err
		}
	}

	return nil
}

func (fsys Filesystem) addWorkspaceToZip(ctx context.Context, writer *zip.Writer, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		return writePathToZip(writer, path, rel, info)
	})
}

func writePathToZip(writer *zip.Writer, absPath, rel string, info fs.FileInfo) error {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.ToSlash(rel)
	if info.IsDir() {
		header.Name += "/"
		_, err = writer.CreateHeader(header)
		return err
	}

	header.Method = zip.Deflate

	entry, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(entry, file)
	return err
}

func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}
