package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	URL string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port    string
	GinMode string
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	AllowedOrigins string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_EXPIRY: %w", err)
	}

	refreshExpiry, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("parse REFRESH_TOKEN_EXPIRY: %w", err)
	}

	cfg := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://naseej:naseej@localhost:5432/naseej?sslmode=disable"),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Auth: AuthConfig{
			JWTSecret:          getEnv("JWT_SECRET", ""),
			JWTExpiry:          jwtExpiry,
			RefreshTokenExpiry: refreshExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		},
	}

	if cfg.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
