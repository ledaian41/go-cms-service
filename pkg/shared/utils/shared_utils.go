package shared_utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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

func RandomID() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("❌ Failed to generate random ID: %v", err)
	}
	return hex.EncodeToString(b)
}
