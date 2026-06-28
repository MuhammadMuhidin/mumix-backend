package config

import "os"

// Config holds all configuration for the application
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/mumix?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-change-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
