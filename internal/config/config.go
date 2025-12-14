package config

import (
	"os"
	"strings"
)

type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
}

func Load() Config {
	return Config{

		MongoURI:  getEnvOrDefault("MONGO_URI", AppEnv.MongoURI),
		DBName:    getEnvOrDefault("DB_NAME", AppEnv.DBName),
		JWTSecret: getEnvOrDefault("JWT_SECRET", AppEnv.JWTSecret),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}
