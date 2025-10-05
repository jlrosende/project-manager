package ports

import "io/fs"

// Filesystem defines operating system interactions required by the application
// layer. Concrete adapters should provide the actual OS access.
type Filesystem interface {
	EnsureDir(path string, mode fs.FileMode) error
	IsDirEmpty(path string) (bool, error)
	Rename(oldPath, newPath string) error
	WriteFile(path string, data []byte, mode fs.FileMode) error
	Join(elem ...string) string
	IsAbs(path string) bool
	Abs(path string) (string, error)
	ExpandHome(path string) string
	Remove(path string) error
	UserHomeDir() (string, error)
}
