package config

import "os"

type Config struct {
	Port      string
	DBURL     string
	JWTSecret string
}

func Load() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/messenger"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
