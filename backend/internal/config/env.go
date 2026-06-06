package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads backend/.env when invoked from the repo root or .env from backend/.
func LoadEnv() {
	for _, path := range []string{".env", "backend/.env"} {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		_ = godotenv.Load(path)
		return
	}
}
