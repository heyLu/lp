package fileutil

import (
	"os"
	"path"
	"path/filepath"
)

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func MatchExists(glob string) bool {
	matches, _ := filepath.Glob(glob)
	return len(matches) > 0
}

func Join(fileOrDir string, elem ...string) string {
	dir := fileOrDir
	if IsFile(fileOrDir) {
		dir = path.Dir(fileOrDir)
	}
	return path.Join(dir, path.Join(elem...))
}

func IsExecutable(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}

	isExecutable := info.Mode() & 0111
	return isExecutable != 0 && !info.IsDir()
}
