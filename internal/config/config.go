package config

import "os"

type Config struct {
	ServerPort  string
	LogLevel    string
	DatabaseURL string
}

func Load() *Config {
	cfg := &Config{
		ServerPort:  getEnvOrDefault("SERVER_PORT", ":8080"),
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "INFO"),
		DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
