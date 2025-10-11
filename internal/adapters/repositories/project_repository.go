package repositories

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type ProjectRepository struct {
	git *config.Config
	fs  ports.Filesystem
}

type renameOperation struct {
	from string
	to   string
}

var (
	_             ports.ProjectRepository = (*ProjectRepository)(nil)
	equalsSpacing                         = regexp.MustCompile(`\s*=\s*`)
)

func NewProjectRepository(fsys ports.Filesystem) (*ProjectRepository, error) {
	git, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return nil, err
	}

	if fsys == nil {
		fsys = NewFilesystem()
	}

	return &ProjectRepository{
		git: git,
		fs:  fsys,
	}, nil
}

func (p *ProjectRepository) Get(name string) (*domain.Project, error) {
	for _, section := range p.git.Raw.Sections {
		if section.IsName("includeIf") {
			for _, sub := range section.Subsections {
				// Read subsections and get path of the project
				if path, ok := strings.CutPrefix(sub.Name, "gitdir/i:"); ok {
					path = filepath.Clean(path)
					projetPath := filepath.Join(path, ".project.hcl")

					project, err := p.loadDotProject(projetPath)
					if err != nil {
						continue
					}

					if project.Name != name {
						continue
					}

					project.Path = path

					return project, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("project '%s' not found, projects are case sensitive", name)
}

func (p *ProjectRepository) List() ([]*domain.Project, error) {
	projects := []*domain.Project{}

	for _, section := range p.git.Raw.Sections {
		if section.IsName("includeIf") {
			for _, sub := range section.Subsections {
				if path, ok := strings.CutPrefix(sub.Name, "gitdir/i:"); ok {
					path = filepath.Clean(path)
					projetPath := filepath.Join(path, ".project.hcl")

					project, err := p.loadDotProject(projetPath)
					if err != nil {
						continue
					}

					path = p.fs.ExpandHome(path)
					project.Path = path

					projects = append(projects, project)
				}
			}
		}
	}

	return projects, nil
}

/*
NOTE: future work
  - Create directory if not exist
  - Create .env .project.hcl and .<project>.gitconfig files
  - If .env exist warn and continue
  - If .<project>.gitconfig exist warn and continue
  - If .project.hcl exist
*/
func (p *ProjectRepository) Create(
	name, path, _ string,
	shell, envFile string,
	_ domain.EnvVars,
	_ *domain.GitConfig,
) (*domain.Project, error) {
	// Expand '~' to user home
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	// Ensure path exists and is a directory (create if missing)
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !info.IsDir() {
		return nil, errors.New("the path must be a directory")
	}

	// Normalize to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	path = absPath

	// Create paths (idempotent)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}

	// .project.hcl
	projPath := filepath.Join(path, ".project.hcl")
	if _, err = os.Stat(projPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exists in directory %s", filepath.Base(projPath), path)
	}

	fpProj, err := os.OpenFile(projPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}

	defer fpProj.Close()

	doc := hclwrite.NewEmptyFile()
	body := doc.Body()
	proj := &domain.Project{
		Name:         name,
		Description:  "",
		Shell:        strings.TrimSpace(shell),
		EnvVarsFile:  envFile,
		Environments: []*domain.Environment{},
	}
	proj.Path = path
	gohcl.EncodeIntoBody(proj, body)

	rendered := hclwrite.Format(doc.Bytes())

	rendered = equalsSpacing.ReplaceAll(rendered, []byte(" = "))
	if _, err := fpProj.Write(rendered); err != nil {
		return nil, err
	}

	return proj, nil
}

func (p *ProjectRepository) Delete(_ string) error {
	return nil
}

func (p *ProjectRepository) UpdateProject(project *domain.Project) error {
	projPath := filepath.Join(project.Path, ".project.hcl")
	if strings.HasPrefix(projPath, "~/") {
		h, _ := os.UserHomeDir()
		projPath = filepath.Join(h, projPath[2:])
	}

	// reload global git config to ensure we operate on latest state
	if gg, err := config.LoadConfig(config.GlobalScope); err == nil {
		p.git = gg
	}

	projectDir := filepath.Dir(projPath)
	gitdir := fmt.Sprintf("gitdir/i:%s/", projectDir)
	includeIf := p.git.Raw.Section("includeIf").Subsection(gitdir)
	oldGitPath := includeIf.Option("path")

	newGitPath := filepath.Join(projectDir, fmt.Sprintf(".%s.gitconfig", project.Name))
	if oldGitPath != "" && oldGitPath != newGitPath {
		if _, err := os.Stat(oldGitPath); err == nil {
			_ = os.Rename(oldGitPath, newGitPath)
		}

		includeIf.SetOption("path", newGitPath)

		gitConf, _ := p.git.Marshal()

		if home, err := os.UserHomeDir(); err == nil {
			if fpGit, err := os.OpenFile(filepath.Join(home, ".gitconfig"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600); err == nil {
				_, _ = fpGit.Write(gitConf)
				_ = fpGit.Close()
				_ = p.git.Validate()
			}
		}
	}

	// merge with existing .project.hcl to preserve fields like environments
	cur, err := p.loadDotProject(projPath)
	if err == nil && cur != nil {
		if project.Description == "" {
			project.Description = cur.Description
		}

		if strings.TrimSpace(project.Shell) == "" {
			project.Shell = cur.Shell
		}

		if strings.TrimSpace(project.EnvVarsFile) == "" {
			project.EnvVarsFile = cur.EnvVarsFile
		}

		if len(project.Environments) == 0 && len(cur.Environments) > 0 {
			project.Environments = cur.Environments
		}

		if strings.TrimSpace(project.DefaultEnv) == "" {
			project.DefaultEnv = cur.DefaultEnv
		}
	}

	doc := hclwrite.NewEmptyFile()
	body := doc.Body()
	gohcl.EncodeIntoBody(project, body)

	if err := os.WriteFile(projPath, doc.Bytes(), 0o600); err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) UpdateEnvironment(projectName, originalEnvName string, env *domain.Environment) error {
	project, err := p.Get(projectName)
	if err != nil {
		return err
	}

	idx := -1

	for i, e := range project.Environments {
		if e.Name == originalEnvName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("environment %s not found", originalEnvName)
	}

	old := project.Environments[idx]

	// Update basic fields
	project.Environments[idx].Name = env.Name

	project.Environments[idx].Color = env.Color
	if env.EnvVarsMode == "" {
		project.Environments[idx].EnvVarsMode = domain.EnvVarsModeMerge
	} else {
		project.Environments[idx].EnvVarsMode = env.EnvVarsMode
	}

	// Handle env file rename if filename changed
	if strings.TrimSpace(env.EnvVarsFile) != "" && env.EnvVarsFile != old.EnvVarsFile {
		oldPath := old.EnvVarsFile
		newPath := env.EnvVarsFile

		projPath := project.Path
		if strings.HasPrefix(projPath, "~/") {
			h, _ := os.UserHomeDir()
			projPath = filepath.Join(h, projPath[2:])
		}

		if !filepath.IsAbs(oldPath) {
			oldPath = filepath.Join(projPath, oldPath)
		}

		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(projPath, newPath)
		}

		if _, err := os.Stat(oldPath); err == nil {
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				_ = p.fs.Rename(oldPath, newPath)
			}
		}

		project.Environments[idx].EnvVarsFile = env.EnvVarsFile
	}

	return p.UpdateProject(project)
}

func (p *ProjectRepository) AddEnvironment(projectName string, env *domain.Environment, _ domain.EnvVars) error {
	project, err := p.Get(projectName)
	if err != nil {
		return err
	}

	for _, e := range project.Environments {
		if e.Name == env.Name {
			return fmt.Errorf("environment %s already exists", env.Name)
		}
	}

	if env.EnvVarsMode == "" {
		env.EnvVarsMode = domain.EnvVarsModeMerge
	}

	if strings.HasPrefix(project.Path, "~/") {
		h, _ := os.UserHomeDir()
		project.Path = filepath.Join(h, project.Path[2:])
	}

	envFile := strings.TrimSpace(env.EnvVarsFile)
	if envFile == "" {
		slug := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(env.Name, " ", "-")))
		if slug != "" {
			envFile = "." + slug + ".env"
		} else {
			envFile = ".env"
		}
	}

	projPath := filepath.Join(project.Path, ".project.hcl")

	fp, err := os.OpenFile(projPath, os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	defer fp.Close()

	b := &strings.Builder{}
	b.WriteString("\n")
	b.WriteString("environment \"")
	b.WriteString(env.Name)
	b.WriteString("\" {\n")

	if env.Color != "" {
		b.WriteString("  color = \"")
		b.WriteString(env.Color)
		b.WriteString("\"\n")
	}

	b.WriteString("  env_vars_mode = \"")
	b.WriteString(env.EnvVarsMode)
	b.WriteString("\"\n")
	b.WriteString("  env_vars_file = \"")
	b.WriteString(envFile)
	b.WriteString("\"\n")
	b.WriteString("}\n")

	if _, err := io.WriteString(fp, b.String()); err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) ResolveIdentifier(
	_ context.Context,
	lookup domain.ProjectIdentifier,
) (domain.ProjectIdentifier, error) {
	name := strings.TrimSpace(lookup.Name)
	targetPath := strings.TrimSpace(lookup.Path)

	if name == "" && targetPath == "" {
		return domain.ProjectIdentifier{}, domain.ErrProjectNotFound
	}

	if targetPath != "" {
		targetPath = p.fs.ExpandHome(targetPath)
		if !filepath.IsAbs(targetPath) {
			if abs, err := filepath.Abs(targetPath); err == nil {
				targetPath = abs
			}
		}

		targetPath = filepath.Clean(targetPath)
	}

	for _, section := range p.git.Raw.Sections {
		if !section.IsName("includeIf") {
			continue
		}

		for _, sub := range section.Subsections {
			if !strings.HasPrefix(sub.Name, "gitdir/i:") {
				continue
			}

			rawPath := strings.TrimPrefix(sub.Name, "gitdir/i:")
			rawPath = strings.TrimSuffix(rawPath, "/")
			expanded := p.fs.ExpandHome(rawPath)
			expanded = filepath.Clean(expanded)

			if targetPath != "" && expanded != targetPath {
				continue
			}

			project, err := p.loadDotProject(filepath.Join(expanded, ".project.hcl"))
			if err != nil {
				continue
			}

			if name != "" && project.Name != name {
				continue
			}

			return domain.ProjectIdentifier{
				Name:       project.Name,
				Path:       expanded,
				RegistryID: sub.Name,
				Status:     domain.ProjectStatus{},
			}, nil
		}
	}

	return domain.ProjectIdentifier{}, domain.ErrProjectNotFound
}

func (p *ProjectRepository) FinalizeDeletion(
	_ context.Context,
	identifier domain.ProjectIdentifier,
	_ domain.DeleteScope,
) error {
	includeIf := p.git.Raw.Section("includeIf")
	if includeIf != nil {
		targetPath := filepath.Clean(p.fs.ExpandHome(identifier.Path))

		subsections := includeIf.Subsections
		writeIdx := 0

		for _, sub := range subsections {
			if strings.HasPrefix(sub.Name, "gitdir/i:") {
				rawPath := strings.TrimPrefix(sub.Name, "gitdir/i:")
				rawPath = strings.TrimSuffix(rawPath, "/")
				expanded := filepath.Clean(p.fs.ExpandHome(rawPath))

				if expanded == targetPath {
					continue
				}
			}

			subsections[writeIdx] = sub
			writeIdx++
		}

		includeIf.Subsections = subsections[:writeIdx]
	}

	if err := p.git.Validate(); err != nil {
		return err
	}

	data, err := p.git.Marshal()
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(filepath.Join(home, ".gitconfig"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(data); err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) loadDotProject(path string) (*domain.Project, error) {
	project := &domain.Project{}

	path = p.fs.ExpandHome(path)

	err := hclsimple.DecodeFile(path, nil, project)
	if err != nil {
		return nil, err
	}

	if project.Shell == "" {
		if shell, ok := os.LookupEnv("SHELL"); ok {
			project.Shell = shell
		} else {
			return nil, fmt.Errorf("unable to load default shell. Add SHELL environment variable or set the `shell` variable inside your %s file", filepath.Join(path, ".project.hcl"))
		}
	}

	return project, nil
}

func (p *ProjectRepository) LoadProjectDefinition(
	ctx context.Context,
	identifier domain.ProjectIdentifier,
) (*domain.Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	configPath, root, err := p.projectConfigPath(identifier)
	if err != nil {
		return nil, err
	}

	project, err := p.loadDotProject(configPath)
	if err != nil {
		return nil, err
	}

	project.Path = root

	return project, nil
}

func (p *ProjectRepository) ApplyEditChangeSet(
	ctx context.Context,
	identifier domain.ProjectIdentifier,
	changeSet *domain.EditChangeSet,
) (*domain.Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if changeSet == nil || (!changeSet.HasProjectChanges() && !changeSet.HasEnvironmentChanges()) {
		return nil, errors.New("edit change set has no mutations")
	}

	configPath, root, err := p.projectConfigPath(identifier)
	if err != nil {
		return nil, err
	}

	project, err := p.loadDotProject(configPath)
	if err != nil {
		return nil, err
	}

	project.Path = root
	preview := cloneProject(project)
	preview.Path = root

	var renameOps []renameOperation

	if changeSet.Project != nil {
		if err := applyProjectMutations(preview, changeSet.Project.Fields); err != nil {
			return nil, err
		}
	}

	if changeSet.Environment != nil {
		op, err := p.planEnvironmentMutations(preview, root, changeSet.Environment)
		if err != nil {
			return nil, err
		}

		renameOps = append(renameOps, op...)
	}

	if errs := validateEditSnapshot(preview, root, changeSet); len(errs) > 0 {
		return nil, errs
	}

	performed := make([]renameOperation, 0, len(renameOps))
	for _, op := range renameOps {
		if err := p.performRename(op); err != nil {
			for i := len(performed) - 1; i >= 0; i-- {
				_ = p.performRename(renameOperation{from: performed[i].to, to: performed[i].from})
			}

			return nil, err
		}

		performed = append(performed, op)
	}

	if err := p.persistProjectDefinition(configPath, preview); err != nil {
		for i := len(performed) - 1; i >= 0; i-- {
			_ = p.performRename(renameOperation{from: performed[i].to, to: performed[i].from})
		}

		return nil, err
	}

	return preview, nil
}

func (p *ProjectRepository) AcquireEditLock(
	ctx context.Context,
	identifier domain.ProjectIdentifier,
) (*domain.ProjectLock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	lockPath, _, err := p.projectLockPath(identifier)
	if err != nil {
		return nil, err
	}

	fp, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, domain.ErrProjectLocked
		}

		return nil, err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	metadata := fmt.Sprintf("pid=%d\ntimestamp=%s\n", os.Getpid(), timestamp)

	if _, err := fp.WriteString(metadata); err != nil {
		fp.Close()
		_ = os.Remove(lockPath)
		return nil, err
	}

	if err := fp.Close(); err != nil {
		_ = os.Remove(lockPath)
		return nil, err
	}

	release := func() error {
		if err := os.Remove(lockPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		return nil
	}

	return domain.NewProjectLock(lockPath, release), nil
}

func (p *ProjectRepository) projectRoot(identifier domain.ProjectIdentifier) (string, error) {
	root := strings.TrimSpace(identifier.Path)
	if root == "" {
		return "", errors.New("project path is required")
	}

	root = p.fs.ExpandHome(root)
	if !p.fs.IsAbs(root) {
		abs, err := p.fs.Abs(root)
		if err != nil {
			return "", err
		}

		root = abs
	}

	return filepath.Clean(root), nil
}

func (p *ProjectRepository) projectConfigPath(identifier domain.ProjectIdentifier) (string, string, error) {
	root, err := p.projectRoot(identifier)
	if err != nil {
		return "", "", err
	}

	return filepath.Join(root, ".project.hcl"), root, nil
}

func (p *ProjectRepository) projectLockPath(identifier domain.ProjectIdentifier) (string, string, error) {
	root, err := p.projectRoot(identifier)
	if err != nil {
		return "", "", err
	}

	return filepath.Join(root, "._pm.edit.lock"), root, nil
}

func (p *ProjectRepository) persistProjectDefinition(configPath string, project *domain.Project) error {
	dir := filepath.Dir(configPath)

	tmp, err := os.CreateTemp(dir, ".project-edit-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	doc := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(project, doc.Body())
	rendered := equalsSpacing.ReplaceAll(hclwrite.Format(doc.Bytes()), []byte(" = "))

	if _, err := tmp.Write(rendered); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), configPath)
}

func applyProjectMutations(project *domain.Project, fields map[string]*domain.FieldMutation) error {
	if len(fields) == 0 {
		return nil
	}

	for key, mutation := range fields {
		if mutation == nil {
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return fmt.Errorf("project mutation %s: %w", key, err)
		}

		switch key {
		case domain.ProjectFieldDescription:
			project.Description = value
		case domain.ProjectFieldShell:
			project.Shell = strings.TrimSpace(value)
		case domain.ProjectFieldEnvVarsFile:
			project.EnvVarsFile = strings.TrimSpace(value)
		case domain.ProjectFieldDefaultEnv:
			project.DefaultEnv = strings.TrimSpace(value)
		default:
			return fmt.Errorf("unsupported project field mutation %q", key)
		}
	}

	return nil
}

func (p *ProjectRepository) planEnvironmentMutations(
	project *domain.Project,
	root string,
	changeSet *domain.EnvironmentChangeSet,
) ([]renameOperation, error) {
	if changeSet == nil || changeSet.Name == "" {
		return nil, errors.New("environment change set requires a target name")
	}

	env := findEnvironment(project, changeSet.Name)
	if env == nil {
		return nil, fmt.Errorf("environment %s not found", changeSet.Name)
	}

	renames := []renameOperation{}

	for key, mutation := range changeSet.Fields {
		if mutation == nil {
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return nil, fmt.Errorf("environment mutation %s: %w", key, err)
		}

		switch key {
		case domain.EnvironmentFieldColor:
			env.Color = value
		case domain.EnvironmentFieldEnvVarsMode:
			mode := strings.TrimSpace(strings.ToLower(value))
			if mode == "" {
				env.EnvVarsMode = domain.EnvVarsModeMerge
			} else {
				env.EnvVarsMode = mode
			}
		case domain.EnvironmentFieldEnvVarsFile:
			trimmed := strings.TrimSpace(value)
			oldFile := strings.TrimSpace(env.EnvVarsFile)
			env.EnvVarsFile = trimmed

			if trimmed == "" || oldFile == trimmed {
				continue
			}

			from := resolveEnvironmentPath(root, oldFile)
			to := resolveEnvironmentPath(root, trimmed)
			if from != "" && to != "" && from != to {
				renames = append(renames, renameOperation{from: from, to: to})
			}
		default:
			return nil, fmt.Errorf("unsupported environment field mutation %q", key)
		}
	}

	return renames, nil
}

func cloneProject(in *domain.Project) *domain.Project {
	if in == nil {
		return nil
	}

	out := *in
	if len(in.Environments) > 0 {
		out.Environments = make([]*domain.Environment, len(in.Environments))
		for i, env := range in.Environments {
			if env == nil {
				continue
			}

			copyEnv := *env
			out.Environments[i] = &copyEnv
		}
	}

	return &out
}

func findEnvironment(project *domain.Project, name string) *domain.Environment {
	if project == nil {
		return nil
	}

	for _, env := range project.Environments {
		if env != nil && env.Name == name {
			return env
		}
	}

	return nil
}

func resolveEnvironmentPath(root, target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return ""
	}

	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}

	root = strings.TrimSpace(root)
	if root == "" {
		return ""
	}

	return filepath.Clean(filepath.Join(root, target))
}

func validateEditSnapshot(project *domain.Project, root string, changeSet *domain.EditChangeSet) domain.ProjectValidationErrors {
	var errs domain.ProjectValidationErrors

	if project == nil || changeSet == nil {
		return errs
	}

	if changeSet.Project != nil {
		if _, ok := changeSet.Project.Fields[domain.ProjectFieldShell]; ok {
			if strings.TrimSpace(project.Shell) == "" {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.ProjectFieldShell,
					Message: "shell cannot be empty",
				})
			}
		}

		if _, ok := changeSet.Project.Fields[domain.ProjectFieldEnvVarsFile]; ok {
			trimmed := strings.TrimSpace(project.EnvVarsFile)
			if trimmed == "" {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.ProjectFieldEnvVarsFile,
					Message: "env vars file cannot be empty",
				})
			} else if ok, _, err := pathWithinRoot(root, trimmed); err != nil {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.ProjectFieldEnvVarsFile,
					Message: fmt.Sprintf("validate env vars file: %v", err),
				})
			} else if !ok {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.ProjectFieldEnvVarsFile,
					Message: fmt.Sprintf("env vars file must be within project root (%s)", root),
				})
			}
		}

		if _, ok := changeSet.Project.Fields[domain.ProjectFieldDefaultEnv]; ok {
			value := strings.TrimSpace(project.DefaultEnv)
			if value != "" && findEnvironment(project, value) == nil {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.ProjectFieldDefaultEnv,
					Message: fmt.Sprintf("environment %q does not exist", value),
				})
			}
		}
	}

	if changeSet.Environment != nil {
		env := findEnvironment(project, changeSet.Environment.Name)
		if env == nil {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsFile,
				Message: fmt.Sprintf("environment %q not found", changeSet.Environment.Name),
			})
			return errs
		}

		if _, ok := changeSet.Environment.Fields[domain.EnvironmentFieldEnvVarsMode]; ok {
			mode := strings.TrimSpace(env.EnvVarsMode)
			if mode == "" {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.EnvironmentFieldEnvVarsMode,
					Message: "env vars mode cannot be empty",
				})
			} else if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.EnvironmentFieldEnvVarsMode,
					Message: fmt.Sprintf("env vars mode must be %q or %q", domain.EnvVarsModeMerge, domain.EnvVarsModeReplace),
				})
			}
		}

		if _, ok := changeSet.Environment.Fields[domain.EnvironmentFieldEnvVarsFile]; ok {
			trimmed := strings.TrimSpace(env.EnvVarsFile)
			if trimmed == "" {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.EnvironmentFieldEnvVarsFile,
					Message: "env vars file cannot be empty",
				})
			} else if ok, _, err := pathWithinRoot(root, trimmed); err != nil {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.EnvironmentFieldEnvVarsFile,
					Message: fmt.Sprintf("validate env vars file: %v", err),
				})
			} else if !ok {
				errs = append(errs, domain.ProjectValidationError{
					Field:   domain.EnvironmentFieldEnvVarsFile,
					Message: fmt.Sprintf("env vars file must be within project root (%s)", root),
				})
			}
		}
	}

	return errs
}

func pathWithinRoot(root, candidate string) (bool, string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return false, "", errors.New("project root is empty")
	}

	resolvedRoot := root
	if !filepath.IsAbs(resolvedRoot) {
		abs, err := filepath.Abs(resolvedRoot)
		if err != nil {
			return false, "", err
		}
		resolvedRoot = abs
	}
	resolvedRoot = filepath.Clean(resolvedRoot)

	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		return false, "", nil
	}

	resolved := candidate
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(resolvedRoot, resolved)
	}
	resolved = filepath.Clean(resolved)

	rel, err := filepath.Rel(resolvedRoot, resolved)
	if err != nil {
		return false, resolved, err
	}

	if rel == "." {
		return true, resolved, nil
	}

	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return false, resolved, nil
	}

	return true, resolved, nil
}

func (p *ProjectRepository) performRename(op renameOperation) error {
	if op.from == "" || op.to == "" || op.from == op.to {
		return nil
	}

	if _, err := os.Stat(op.from); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if _, err := os.Stat(op.to); err == nil {
		return fmt.Errorf("destination %s already exists", op.to)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return p.fs.Rename(op.from, op.to)
}
