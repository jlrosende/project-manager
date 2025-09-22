package repositories

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
	"github.com/jlrosende/project-manager/internal/tools"
)

type ProjectRepository struct {
	git *config.Config
}

var _ ports.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository() (*ProjectRepository, error) {
	git, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return nil, err
	}

	return &ProjectRepository{
		git: git,
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

					path = tools.ExpandHome(path)
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

	// Require empty directory to avoid clobbering existing projects
	if isEmpty, err := tools.IsDirEmpty(path); err != nil {
		return nil, err
	} else if !isEmpty {
		return nil, fmt.Errorf("directory %s, is not empty", path)
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

	if _, err := fpProj.Write(doc.Bytes()); err != nil {
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
				_ = tools.Rename(oldPath, newPath)
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

func (p *ProjectRepository) loadDotProject(path string) (*domain.Project, error) {
	project := &domain.Project{}

	path = tools.ExpandHome(path)

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
