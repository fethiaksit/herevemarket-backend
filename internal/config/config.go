package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	MongoURI        string
	DBName          string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() Config {
	return Config{

		MongoURI:        getEnvOrDefault("MONGO_URI", AppEnv.MongoURI),
		DBName:          getEnvOrDefault("DB_NAME", AppEnv.DBName),
		JWTSecret:       getEnvOrDefault("JWT_SECRET", AppEnv.JWTSecret),
		AccessTokenTTL:  getDurationEnv("ACCESS_TOKEN_TTL_MINUTES", AppEnv.AccessTokenTTLMinutes, time.Minute),
		RefreshTokenTTL: getDurationEnv("REFRESH_TOKEN_TTL_DAYS", AppEnv.RefreshTokenTTLDays, 24*time.Hour),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue int, unit time.Duration) time.Duration {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return time.Duration(parsed) * unit
		}
	}
	return time.Duration(defaultValue) * unit
}
