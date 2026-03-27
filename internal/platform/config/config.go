package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("can't load .env file")
	}

	cfg := &Config{
		DBUrl: os.Getenv("DB_URL"),
		Port:  os.Getenv("PORT"),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DB_URL is required")
	}

	return cfg
}
