package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	SessionTTL  time.Duration
	LogLevel    slog.Level
}

func Load() Config {
	ttl := 12 * time.Hour
	if raw := os.Getenv("SESSION_TTL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			ttl = parsed
		}
	}
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}
	return Config{DatabaseURL: env("DATABASE_URL", "postgres://airbridge:airbridge@localhost:55433/airbridge?sslmode=disable"), HTTPAddr: env("HTTP_ADDR", ":8080"), SessionTTL: ttl, LogLevel: level}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func parseBool(value string) bool { parsed, _ := strconv.ParseBool(value); return parsed }
