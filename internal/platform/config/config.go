package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	Port      string
	SecretKey string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("can't load .env file")
	}

	cfg := &Config{
		DBUrl:     os.Getenv("DB_URL"),
		Port:      os.Getenv("PORT"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DB_URL is required")
	}

	if cfg.SecretKey == "" {
		log.Fatal("SECRET_KEY is required")
	}

	return cfg
}
