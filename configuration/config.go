package configuration

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT  string
	HOSTS string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err

	}

	return &Config{
		PORT:  getEnv("PORT", "8080"),
		HOSTS: getEnv("HOSTS", "localhost"),
	}, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
