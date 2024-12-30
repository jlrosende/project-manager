package repositories

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-git/go-git/v5/config"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/ports"
)

type GitRepository struct {
	git *config.Config
}

var _ ports.GitRepository = (*GitRepository)(nil)

func NewGitRepository() (*GitRepository, error) {
	return &GitRepository{
		git: config.NewConfig(),
	}, nil
}

func (g *GitRepository) Load(path string) (*domain.GitConfig, error) {

	gitFileContent, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err := g.git.Unmarshal(gitFileContent); err != nil {
		return nil, err
	}

	if err := g.git.Validate(); err != nil {
		return nil, err
	}

	tagGPGSing, err := strconv.ParseBool(g.git.Raw.Section("tag").Option("gpgsign"))

	if err != nil {
		return nil, err
	}

	commitGPGSing, err := strconv.ParseBool(g.git.Raw.Section("tag").Option("gpgsign"))
	if err != nil {
		return nil, err
	}

	return &domain.GitConfig{
		User: domain.User{
			Name:       g.git.User.Name,
			Email:      g.git.User.Email,
			SigningKey: g.git.Raw.Section("user").Option("signingkey"),
		},
		Tag: domain.Tag{
			GPGSign: tagGPGSing,
		},
		Commit: domain.Commit{
			GPGSign: commitGPGSing,
		},
	}, nil
}

func (g *GitRepository) Save(path string, gitConfig *domain.GitConfig) error {

	g.git.User.Email = gitConfig.User.Email
	g.git.User.Name = gitConfig.User.Name
	g.git.Raw.AddOption("user", "", "signingkey", gitConfig.User.SigningKey)

	g.git.Raw.AddOption("tag", "", "gpgsign", "true")
	g.git.Raw.AddOption("commit", "", "gpgsign", "true")

	newGitConf, err := g.git.Marshal()

	if err != nil {
		return err
	}

	_, err = os.Stat(path)

	if !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists in directory %s", path, filepath.Dir(path))
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	n_bytes, err := fp.Write(newGitConf)
	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("%d Bytes written in %s", n_bytes, path))

	return nil
}
