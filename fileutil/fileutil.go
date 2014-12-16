package fileutil

import (
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	if f, err := os.Open(path); err == nil {
		if i, err := f.Stat(); err == nil {
			return i.IsDir()
		}
	}
	return false
}

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
