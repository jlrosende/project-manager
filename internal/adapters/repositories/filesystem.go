package repositories

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

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
