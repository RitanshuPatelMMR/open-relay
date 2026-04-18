package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl    string
	RedisUrl string
	Port     string
}

func Load() *Config {
	_ = godotenv.Load("../.env")
	cfg := &Config{
		DBUrl:    os.Getenv("DB_URL"),
		RedisUrl: os.Getenv("REDIS_URL"),
		Port:     os.Getenv("API_PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = "8081"
	}
	if cfg.DBUrl == "" || cfg.RedisUrl == "" {
		log.Fatal("DB_URL and REDIS_URL required")
	}
	return cfg
}