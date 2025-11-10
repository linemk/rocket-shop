package path

import (
	"os"
	"path/filepath"
)

// GetProjectRoot ищет корневую директорию проекта по наличию .git директории
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("не удалось получить рабочую директорию: " + err.Error())
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil && info.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("не удалось найти корень проекта (.git)")
		}

		dir = parent
	}
}
