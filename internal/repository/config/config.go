// Package config provides application configuration.
package config

import (
	"os"
	"strings"
)

// Config holds all configuration for the application.
type Config struct {
	Port        string
	DBURL       string
	CORSOrigins []string
	LogLevel    string
}

// Load loads the configuration from environment variables.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBURL:       getEnv("DB_URL", "postgres://studytrack:studytrack@localhost:5432/studytrack?sslmode=disable"),
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
