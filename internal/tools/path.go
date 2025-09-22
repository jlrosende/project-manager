package tools

import (
	"os"
	"path/filepath"
)

func ExpandHome(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}

	return path
}

func AbsPath(path string) (string, error) {
	p := ExpandHome(path)
	return filepath.Abs(p)
}

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}
