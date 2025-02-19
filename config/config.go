package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppConfig struct {
	DatabaseUrl string
}

func LoadConfig() *AppConfig {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create and return AppConfig instance
	return &AppConfig{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}
}
