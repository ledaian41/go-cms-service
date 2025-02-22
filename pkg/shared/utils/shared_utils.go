package shared_utils

import (
	"os"
	"strings"
)

func IsJsonPath(path string) bool {
	return strings.HasSuffix(path, ".json")
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
