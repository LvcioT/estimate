package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	GinMode          string
	DBPath           string
	SessionSecret    string
	HashCost         int
	DefaultAdminUser string
	DefaultAdminPass string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	hashCost, err := strconv.Atoi(getEnv("HASH_COST", "10"))
	if err != nil {
		log.Printf("Warning: Invalid HASH_COST value, using default: 10")
		hashCost = 10
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		GinMode:          getEnv("GIN_MODE", "debug"),
		DBPath:           getEnv("DB_PATH", "estimate.db"),
		SessionSecret:    getEnv("SESSION_SECRET", "change_this_in_production"),
		HashCost:         hashCost,
		DefaultAdminUser: getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultAdminPass: getEnv("DEFAULT_ADMIN_PASSWORD", "admin"),
	}
}

// getEnv retrieves an environment variable value or returns the default if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
