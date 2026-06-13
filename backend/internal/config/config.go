package config

import "os"

type Config struct {
	AppName     string
	AppPort     string
	AppBaseURL  string
	FrontendURL string
	JWTSecret   string
	DatabaseURL string
	RedisAddr   string
	AppEnv      string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "LinkHub API"),
		AppPort:     getEnv("APP_PORT", "8081"),
		AppBaseURL:  getEnv("APP_BASE_URL", "http://localhost:8081"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3002"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisAddr:   getEnv("REDIS_ADDR", ""),
		AppEnv:      getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
