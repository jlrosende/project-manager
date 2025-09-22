package tools

import (
	"io"
	"os"
)

func EnsureDir(path string, mode os.FileMode) error {
	if st, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, mode)
		}

		return err
	} else if !st.IsDir() {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrInvalid}
	}

	return nil
}

func IsDirEmpty(path string) (bool, error) {
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

func Rename(oldPath, newPath string) error {
	if err := EnsureDir(filepathDir(newPath), 0o755); err != nil {
		return err
	}

	return os.Rename(oldPath, newPath)
}

func WriteFile(path string, data []byte, mode os.FileMode) error {
	if err := EnsureDir(filepathDir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, mode)
}

func filepathDir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[:i]
		}
	}

	return "."
}
