package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig struct {
	CachePath string
}

func LoadConfig() *AppConfig {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("❌ Error loading .env file: %v", err)
	}

	// Create and return AppConfig instance
	return &AppConfig{
		CachePath: os.Getenv("CACHE_PATH"),
	}
}
