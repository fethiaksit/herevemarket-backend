package config

type Env struct {
	MongoURI              string
	DBName                string
	JWTSecret             string
	AccessTokenTTLMinutes int
	RefreshTokenTTLDays   int
}

var AppEnv = Env{
	MongoURI:              "mongodb://localhost:27017",
	DBName:                "docker-herevemarket",
	JWTSecret:             "MgGYQBvqZodV4sPFJaC6XRbspRiklcs6OmHG714ynxk",
	AccessTokenTTLMinutes: 20,
	RefreshTokenTTLDays:   7,
}
