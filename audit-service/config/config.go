package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	Debug bool
}

func NewConfig() *Config {
	// Load .env file if present
	_ = godotenv.Load()

	return &Config{
		Port:  getEnv("PORT", "8081"),
		Debug: getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
