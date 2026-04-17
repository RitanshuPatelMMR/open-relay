package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl                 string
	RedisURL              string
	WorkerCount           int
	MaxRetryAttempts      int
	RequestTimeoutSeconds int
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no .env file, using system env")
	}

	return &Config{
		DBUrl:                 getEnv("DB_URL", ""),
		RedisURL:              getEnv("REDIS_URL", "redis://localhost:6379"),
		WorkerCount:           getEnvInt("WORKER_COUNT", 5),
		MaxRetryAttempts:      getEnvInt("MAX_RETRY_ATTEMPTS", 5),
		RequestTimeoutSeconds: getEnvInt("REQUEST_TIMEOUT_SECONDS", 10),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}