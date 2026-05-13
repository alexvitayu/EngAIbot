package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/alexvitayu/EngAIbot/internal/ai_agents/groq_ai"
	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/alexvitayu/EngAIbot/internal/db"
	"github.com/alexvitayu/EngAIbot/internal/db/repository"
	"github.com/alexvitayu/EngAIbot/internal/logger"
	"github.com/alexvitayu/EngAIbot/internal/service"
	"github.com/alexvitayu/EngAIbot/internal/tg_bot"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadCfg()
	if err != nil {
		slog.Error("LoadConfig", "error", err)
		os.Exit(1)
	}

	customLogger, err := logger.SetupLogger(cfg)
	if err != nil {
		fmt.Println("logger error")
	}
	slog.SetDefault(customLogger)
	slog.Debug("APP_ENV", "app_env", cfg.APPEnv)

	conn, err := db.NewPgxPool(ctx, cfg)
	if err != nil {
		slog.Error("failed to create connections pool", "error", err)
		os.Exit(1)
	}
	slog.Debug("connections pool initialized!")

	repo := repository.NewRepository(conn)

	ai, err := groq_ai.NewGroq(cfg)
	if err != nil {
		slog.Error("failed to initialize groq AI", "error", err)
		os.Exit(1)
	}
	slog.Debug("AI agent initialized successfully!")

	service := service.NewPhraseService(repo, ai)

	err = tg_bot.TgBot(ctx, *cfg, *service)
	if err != nil {
		slog.Error("TgBot: %w", err)
	}
}
