package main

import (
	"log/slog"
	"os"

	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/alexvitayu/EngAIbot/internal/logger"
	"github.com/alexvitayu/EngAIbot/internal/tg_bot"
)

func main() {
	//ctx := context.Background()

	cfg, err := config.LoadCfg()
	if err != nil {
		slog.Error("LoadConfig", "error", err)
		os.Exit(1)
	}

	customLogger := logger.SetupLogger(cfg)
	slog.SetDefault(customLogger)
	slog.Debug("APP_ENV", "app_env", cfg.APPEnv)

	tg_bot.TgBot(*cfg)
}
