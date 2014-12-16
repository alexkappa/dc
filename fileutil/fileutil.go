package fileutil

import (
	"os"
	"path/filepath"
)

// IsDir will return true if path is a directory on the servers file system. It
// returns false otherwise.
func IsDir(path string) bool {
	if f, err := os.Open(path); err == nil {
		if i, err := f.Stat(); err == nil {
			return i.IsDir()
		}
	}
	return false
}

// Abs will try to make a relative path into an absolute path using the current
// working directory. If an error occurs during the process path will be
// returned as the absolute path.
func Abs(path string) string {
	if !filepath.IsAbs(path) {
		currentDir, err := os.Getwd()
		if err != nil {
			return path
		}
		return filepath.Join(currentDir, path)
	}
	return path
}
