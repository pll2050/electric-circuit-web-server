package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
	// Firebase configuration
	FirebaseProjectID             string
	FirebaseServiceAccountKeyPath string
	FirebaseDatabaseURL           string
}

func Load() *Config {
	return &Config{
		Port:                          getEnv("PORT", "8080"),
		DatabaseURL:                   getEnv("DATABASE_URL", "postgres://postgres:q1w2e3r4@localhost:5432/electric_circuit?sslmode=disable"),
		Environment:                   getEnv("ENVIRONMENT", "development"),
		FirebaseProjectID:             getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseServiceAccountKeyPath: getEnv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH", ""),
		FirebaseDatabaseURL:           getEnv("FIREBASE_DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
