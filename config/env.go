package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var Env *AppConfig

type AppConfig struct {
	CachePath              string
	DbHost                 string
	DbUser                 string
	DbPwd                  string
	RedisHost              string
	MaxUploadFileSize      int64
	MaxTotalUploadFileSize int64
	AppHost                string
}

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("❌ Error loading .env file: %v", err)
	}

	maxUploadFileSize, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_FILE_SIZE"), 10, 64)
	maxTotalUploadFileSize, err := strconv.ParseInt(os.Getenv("MAX_TOTAL_UPLOAD_FILE_SIZE"), 10, 64)

	Env = &AppConfig{
		DbHost:                 os.Getenv("DATABASE_HOST"),
		DbUser:                 os.Getenv("DATABASE_USER"),
		DbPwd:                  os.Getenv("DATABASE_PWD"),
		CachePath:              os.Getenv("CACHE_PATH"),
		RedisHost:              os.Getenv("REDIS_HOST"),
		MaxUploadFileSize:      maxUploadFileSize,
		MaxTotalUploadFileSize: maxTotalUploadFileSize,
		AppHost:                os.Getenv("APP_HOST"),
	}
}
