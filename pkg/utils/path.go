package utils

import (
	"os"
	"path"
	"strings"
)

// GetRootPath Set the project root path name, run it in any directory under the project, and return the project root path address
func GetRootPath(rootDirName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	arr := strings.Split(strings.ReplaceAll(dir, "\\", "/"), rootDirName)
	return path.Join(
		arr[0],
		rootDirName,
	), nil
}
