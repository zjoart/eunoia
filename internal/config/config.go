package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type AIConfig struct {
	GeminiAPIKey string
}

type Config struct {
	AppEnv string
	Port   string
	DB     DBConfig
	AI     AIConfig
}

func LoadConfig() *Config {
	config := &Config{
		Port: getEnv("PORT"),
		DB: DBConfig{
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASS"),
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Name:     getEnv("DB_NAME"),
		},

		AppEnv: getEnv("APP_ENV"),
		AI: AIConfig{
			GeminiAPIKey: getEnv("GEMINI_API_KEY"),
		},
	}

	return config
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	panic(fmt.Sprintf("%s is required", key))
}
