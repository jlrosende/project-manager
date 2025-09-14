package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Setup(levelText string, filePath string) (*slog.Logger, error) {
	var lv slog.LevelVar

	if levelText == "" {
		levelText = "info"
	}

	if err := lv.UnmarshalText([]byte(levelText)); err != nil {
		return nil, err
	}

	if filePath == "" {
		cache, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(cache, "pm.log")
	}

	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	lg := slog.New(slog.NewTextHandler(fp, &slog.HandlerOptions{Level: lv.Level()}))

	lg = lg.With(
		slog.Group("ps",
			slog.Int("pid", os.Getpid()),
			slog.Int("ppid", os.Getppid()),
			slog.String("project", os.Getenv("PM_ACTIVE_PROJECT")),
		),
	)
	slog.SetDefault(lg)

	return lg, nil
}
