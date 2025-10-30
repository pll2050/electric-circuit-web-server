package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

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
