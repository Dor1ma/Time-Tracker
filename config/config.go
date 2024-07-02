package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DbHost string
	DbUser string
	DbPort string
	DbName string
	DbPass string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	config := &Config{
		DbHost: os.Getenv("DB_HOST"),
		DbUser: os.Getenv("DB_USER"),
		DbPort: os.Getenv("DB_PORT"),
		DbName: os.Getenv("DB_NAME"),
		DbPass: os.Getenv("DB_PASS"),
	}

	return config, nil
}
