package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppName string
	Port    int
	Env     string

	DB struct {
		URL string
	}

	JWT struct {
		Secret string
	}

	Log struct {
		Level string
	}
}

func Load() *Config {
	port := getInt("APP_PORT", 0)
	if port == 0 {
		port = getInt("PORT", 8080)
	}

	cfg := &Config{
		AppName: getString("APP_NAME", "Research-api"),
		Port:    port,
		Env:     getString("APP_ENV", "development"),
	}
	cfg.DB.URL = getString("DB_URL", "")

	cfg.JWT.Secret = getString("JWT_SECRET", "")

	cfg.Log.Level = getString("LOG_LEVEL", "info")

	validate(cfg)

	return cfg
}

func getString(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid value for %s: %v", key, err)
		}
		return i
	}
	return defaultVal
}

func validate(cfg *Config) {

	if cfg.DB.URL == "" {
		log.Fatal("DB_URL is required")
	}
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
}
