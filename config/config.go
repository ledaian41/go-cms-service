package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	CachePath              string
	MaxUploadFileSize      int64
	MaxTotalUploadFileSize int64
}

func LoadConfig() *AppConfig {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("❌ Error loading .env file: %v", err)
	}

	maxUploadFileSize, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_FILE_SIZE"), 10, 64)
	maxTotalUploadFileSize, err := strconv.ParseInt(os.Getenv("MAX_TOTAL_UPLOAD_FILE_SIZE"), 10, 64)

	// Create and return AppConfig instance
	return &AppConfig{
		CachePath:              os.Getenv("CACHE_PATH"),
		MaxUploadFileSize:      maxUploadFileSize,
		MaxTotalUploadFileSize: maxTotalUploadFileSize,
	}
}
