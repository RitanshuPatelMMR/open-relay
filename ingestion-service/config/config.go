package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	RedisURL  string
	Port      string
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no .env file, using system env")
	}

	return &Config{
		DBUrl:    getEnv("DB_URL", ""),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}