package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBName     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string

	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetimeMin int

	MinioUser        string
	MinioBucketName  string
	MinioPassword    string
	MinioApiPort     string
	MinioConsolePort string
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", key, val)
	}
	return num
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

		DBMaxOpenConns:       getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:       getEnvInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxLifetimeMin: getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 60),

		MinioBucketName:  os.Getenv("MINIO_BUCKET_NAME"),
		MinioUser:        os.Getenv("MINIO_USER"),
		MinioPassword:    os.Getenv("MINIO_PASSWORD"),
		MinioApiPort:     os.Getenv("MINIO_API_PORT"),
		MinioConsolePort: os.Getenv("MINIO_CONSOLE_PORT"),
	}
}
