package configs

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	Debug       bool
}

func NewConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:root@localhost:5432/seo_audit?sslmode=disable"),
		Debug:       getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
