package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	HTTPAddr       string
	DatabaseDriver string
	DatabaseDSN    string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	StorageRoot    string
	QueueName      string
	ScanMode       string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:         getEnv("APP_ENV", "dev"),
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		DatabaseDriver: getEnv("DATABASE_DRIVER", "postgres"),
		DatabaseDSN:    getEnv("DATABASE_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=release_platform sslmode=disable"),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		StorageRoot:    getEnv("STORAGE_ROOT", "./data"),
		QueueName:      getEnv("QUEUE_NAME", "release_scan"),
		ScanMode:       getEnv("SCAN_MODE", "async"),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}
