package config

import "os"

type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
}

func Load() Config {
	return Config{
		MongoURI:  os.Getenv("mongodb://localhost:27017"),
		DBName:    os.Getenv("docker-herevemarket"),
		JWTSecret: os.Getenv("MgGYQBvqZodV4sPFJaC6XRbspRiklcs6OmHG714ynxk="),
	}
}
