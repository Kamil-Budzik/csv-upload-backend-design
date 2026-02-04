package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBName     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env files not provided", err)
	}

	return Config{
		Port:       os.Getenv("PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
