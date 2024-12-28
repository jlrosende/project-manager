package domain

type GitConfig struct {
	User struct {
		Name       string
		Email      string
		SigningKey string
	}
	Commit struct {
		GPGSign bool
	}
	Tag struct {
		GPGSign bool
	}
}
