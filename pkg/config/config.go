package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	HTTPPort    string

	PostgresURL  string
	RedisURL     string
	KafkaBrokers string
	OTLPURL      string
}

func Load() *Config {
	// load .env if exists (dev only)
	_ = godotenv.Load()

	cfg := &Config{
		ServiceName:  getEnv("SERVICE_NAME", "unknown-service"),
		HTTPPort:     getEnv("HTTP_PORT", "8080"),
		PostgresURL:  getEnv("POSTGRES_URL", ""),
		RedisURL:     getEnv("REDIS_URL", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		OTLPURL:      getEnv("OTLP_URL", "localhost:4317"),
	}

	validate(cfg)
	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid int env %s", key)
	}
	return val
}

func validate(c *Config) {
	if c.ServiceName == "" {
		log.Fatal("SERVICE_NAME must be set")
	}
}
