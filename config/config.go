package config

import (
	"log"
	"os"
)

type Config struct {
	AppName  string
	Port     string
	LogLevel string
	LogPath  string
	Version  string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppName:  getEnv("APP_NAME", "coreBanking"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "DEBUG"),
		LogPath:  getEnv("LOG_PATH", "log/transactions.log"),
		Version:  getEnv("VERSION", "v1"),
	}

	log.Printf("Config loaded: %+v\n", cfg)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
