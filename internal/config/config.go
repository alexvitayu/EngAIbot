package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APPEnv    string
	TgAPI     string
	GeminiAPI string
	GroqAPI   string
}

func LoadCfg() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")

	// Загружаем env файл для конкретной среды
	if env == "development" {
		if _, err := os.Stat(".env.dev"); err == nil {
			if loadErr := godotenv.Load(".env.dev"); loadErr != nil {
				slog.Debug("failed to load .env.dev", "error", loadErr)
				return nil, fmt.Errorf("loading %s: %w", ".env.dev", loadErr)
			}
		}
	}

	if env == "test" {
		if _, err := os.Stat("../../../.env.test"); err == nil {
			if loadErr := godotenv.Load("../../../.env.test"); loadErr != nil {
				slog.Debug("failed to load .env.test", "error", loadErr)
				return nil, fmt.Errorf("loading %s: %w", ".env.test", loadErr)
			}
		}
	}

	config := &AppConfig{
		APPEnv:    env,
		TgAPI:     getEnv("TELEGRAM_APITOKEN", ""),
		GeminiAPI: getEnv("GEMINI_API_KEY", ""),
		GroqAPI:   getEnv("GROQ_API_KEY", ""),
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
