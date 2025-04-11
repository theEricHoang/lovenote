package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecretKey       string
	DatabaseURL        string
	IsProduction       bool
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("	no .env file found, using system env vars")
	}

	config := Config{
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", ""),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/lovenote?sslmode=disable"),
		IsProduction:       getEnvAsBool("PRODUCTION", false),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "us-east-2"),
	}

	if config.JWTSecretKey == "" {
		log.Fatal("missing required environment variable: JWT_SECRET_KEY")
	}

	return config
}

// check if env vars are present
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Convert string env variable to bool
func getEnvAsBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	// Check if the value is "1", "true", or "yes" (case-insensitive)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
