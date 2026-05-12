package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APPEnv     string
	TgAPI      string
	GeminiAPI  string
	GroqAPI    string
	DBConfig   DBConfig
	PoolConfig PoolConfig
}

type DBConfig struct {
	DATABASE_URL string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	DBSSLMode    string
}

type PoolConfig struct {
	DBMaxConns        string
	DBMaxIdleConns    string
	DBConnMaxLifetime string
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

	dbConfig := DBConfig{
		DATABASE_URL: getEnv("DATABASE_URL", ""),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5434"),
		DBName:       getEnv("DB_NAME", ""),
		DBUser:       getEnv("DB_USER", ""),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBSSLMode:    getEnv("DB_SSL_MODE", "disable"),
	}

	poolConfig := PoolConfig{
		DBMaxConns:        getEnv("DB_MAX_CONNS", "10"),
		DBMaxIdleConns:    getEnv("DB_MAX_IDLE_CONNS", "5"),
		DBConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "30m"),
	}

	config := &AppConfig{
		APPEnv:     env,
		TgAPI:      getEnv("TELEGRAM_APITOKEN", ""),
		GeminiAPI:  getEnv("GEMINI_API_KEY", ""),
		GroqAPI:    getEnv("GROQ_API_KEY", ""),
		DBConfig:   dbConfig,
		PoolConfig: poolConfig,
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
