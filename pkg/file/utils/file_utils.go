package file_utils

import (
	"fmt"
	"go-cms-service/config"
	"go-cms-service/pkg/shared/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	basename := strings.TrimSuffix(originalName, ext)
	return fmt.Sprintf("%s_%s%s", shared_utils.RandomID(4), sanitizeFileName(basename), ext)
}

func sanitizeFileName(fileName string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9-_.]`)
	fileName = reg.ReplaceAllString(fileName, "_")
	fileName = strings.ToLower(fileName)
	return fileName
}

func GenerateUploadPath(fileName string) (string, error) {
	now := time.Now()

	uploadDir := fmt.Sprintf("%s/files/%d/%02d/%02d",
		config.LoadConfig().CachePath,
		now.Year(),
		now.Month(),
		now.Day(),
	)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir error: %v", err)
	}

	safeFileName := GenerateUniqueFilename(fileName)
	fullPath := filepath.Join(uploadDir, safeFileName)
	return fullPath, nil
}
