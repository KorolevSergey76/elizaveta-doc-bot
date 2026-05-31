package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken         string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	RedisAddr        string
}

func NewConfig() *Config { // Переименовали из Load
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("cannot load .env file")
	}

	return &Config{
		BotToken:         os.Getenv("BOT_TOKEN"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
	}
}
