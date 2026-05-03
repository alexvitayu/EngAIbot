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
	config := &AppConfig{
		APPEnv:    env,
		TgAPI:     "TELEGRAM_APITOKEN",
		GeminiAPI: "GEMINI_API_KEY",
		GroqAPI:   "GROQ_API_KEY",
	}
	return config, nil
}
