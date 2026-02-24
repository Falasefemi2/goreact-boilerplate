package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Env           string
	DatabaseURL   string
	JWTSecret     string
	AllowedOrigin string
	ResendAPIKey  string
	FromEmail     string
}

func Load() *Config {
	// In production, real env vars are set by the platform.
	// In local dev, we load from .env file.
	// If .env doesn't exist, we just continue — not an error.
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("APP_ENV", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),
		FromEmail:     getEnv("FROM_EMAIL", "onboarding@resend.dev"),
	}

	return cfg
}

// getEnv reads an env var, returns a fallback if not set.
// If fallback is empty string and var is missing, it exits.
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		if fallback == "" {
			log.Fatalf("required environment variable %s is not set", key)
		}
		return fallback
	}
	return val
}
