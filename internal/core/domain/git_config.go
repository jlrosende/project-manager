package domain

type GitConfig struct {
	User   User
	Commit Commit
	Tag    Tag
}

type User struct {
	Name       string
	Email      string
	SigningKey string
}

type Commit struct {
	GPGSign bool
}

type Tag struct {
	GPGSign bool
}

type Option func(*GitConfig)

func New(options ...Option) *GitConfig {
	gitConfig := &GitConfig{}
	for _, o := range options {
		o(gitConfig)
	}
	return gitConfig
}

func WithName(name string) func(*GitConfig) {
	return func(g *GitConfig) {
		g.User.Name = name
	}
}

func WithEmail(email string) func(*GitConfig) {
	return func(g *GitConfig) {
		g.User.Email = email
	}
}

func WithSigningKey(signingKey string) func(*GitConfig) {
	return func(g *GitConfig) {
		g.User.SigningKey = signingKey
	}
}

func WithSign(enable bool) func(*GitConfig) {
	return func(g *GitConfig) {
		g.Commit.GPGSign = enable
		g.Tag.GPGSign = enable
	}
}

func WithCommitSign(enable bool) func(*GitConfig) {
	return func(g *GitConfig) {
		g.Commit.GPGSign = enable
	}
}

func WithTagSign(enable bool) func(*GitConfig) {
	return func(g *GitConfig) {
		g.Tag.GPGSign = enable
	}
}
