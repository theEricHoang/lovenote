package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecretKey string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("	no .env file found, using system env vars")
	}

	config := Config{
		JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),
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
