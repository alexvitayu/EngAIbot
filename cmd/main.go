package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexvitayu/EngAIbot/internal/ai_agents/groq_ai"
	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/alexvitayu/EngAIbot/internal/db"
	"github.com/alexvitayu/EngAIbot/internal/db/repository"
	"github.com/alexvitayu/EngAIbot/internal/logger"
	"github.com/alexvitayu/EngAIbot/internal/service"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadCfg()
	if err != nil {
		slog.Error("LoadConfig", "error", err)
		os.Exit(1)
	}

	customLogger := logger.SetupLogger(cfg)
	slog.SetDefault(customLogger)
	slog.Debug("APP_ENV", "app_env", cfg.APPEnv)

	//tg_bot.TgBot(*cfg)

	conn, err := db.NewPgxPool(ctx, cfg)
	if err != nil {
		slog.Error("failed to create connections pool", "error", err)
		os.Exit(1)
	}
	slog.Debug("connections pool initialized!")

	repo := repository.Repository{
		Conn: conn,
	}

	ai, err := groq_ai.NewGroq(cfg)
	if err != nil {
		slog.Error("failed to initialize groq AI", "error", err)
		os.Exit(1)
	}

	s := service.NewPhraseService(&repo, ai)
	err = s.AddPhrases(ctx, service_dto.UserSettings{})
	if err != nil {
		slog.Error("failed to add phrases to DB", "error", err)
	}
	slog.Debug("phrases saved to DD successfully!!! my congrats!!!")
}
