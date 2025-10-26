package repositories

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jlrosende/project-manager/internal/core/ports"
)

type loggerAdapter struct {
	logger *slog.Logger
}

var _ ports.Logger = (*loggerAdapter)(nil)

// NewLogger wraps a slog logger in the ports.Logger interface.
func NewLogger(logger *slog.Logger) ports.Logger {
	if logger == nil {
		logger = slog.Default()
	}

	return &loggerAdapter{logger: logger}
}

func (l *loggerAdapter) Debug(msg string, fields ...ports.LogField) {
	l.log(func(m string, args ...any) { l.logger.Debug(m, args...) }, msg, fields...)
}

func (l *loggerAdapter) Info(msg string, fields ...ports.LogField) {
	l.log(func(m string, args ...any) { l.logger.Info(m, args...) }, msg, fields...)
}

func (l *loggerAdapter) Warn(msg string, fields ...ports.LogField) {
	l.log(func(m string, args ...any) { l.logger.Warn(m, args...) }, msg, fields...)
}

func (l *loggerAdapter) Error(msg string, fields ...ports.LogField) {
	l.log(func(m string, args ...any) { l.logger.Error(m, args...) }, msg, fields...)
}

func (l *loggerAdapter) log(fn func(string, ...any), msg string, fields ...ports.LogField) {
	if l == nil || l.logger == nil {
		return
	}

	args := make([]any, 0, len(fields))
	for _, f := range fields {
		if f.Key == "" {
			continue
		}

		args = append(args, makeAttr(f))
	}

	fn(msg, args...)
}

func makeAttr(f ports.LogField) slog.Attr {
	switch v := f.Value.(type) {
	case []ports.LogField:
		return slog.Group(f.Key, fieldsToArgs(v)...)
	case map[string]any:
		return slog.Group(f.Key, mapToArgs(v)...)
	default:
		return slog.Any(f.Key, v)
	}
}

func fieldsToArgs(fields []ports.LogField) []any {
	if len(fields) == 0 {
		return nil
	}

	args := make([]any, 0, len(fields))
	for _, field := range fields {
		if field.Key == "" {
			continue
		}

		args = append(args, makeAttr(field))
	}

	return args
}

func mapToArgs(values map[string]any) []any {
	if len(values) == 0 {
		return nil
	}

	args := make([]any, 0, len(values))
	for key, value := range values {
		if key == "" {
			continue
		}

		args = append(args, slog.Any(key, value))
	}

	return args
}

// SetupLogger configures a slog logger writing to the provided path (or a default
// location) and sets it as the process default.
func SetupLogger(levelText, filePath string) (*slog.Logger, error) {
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

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return nil, err
	}

	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
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
