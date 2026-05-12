package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/alexvitayu/EngAIbot/internal/ai_agents"
	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
)

type PhraseKeeper interface {
	CreatePhrasesBatch(ctx context.Context, dtos []*db_dto.PhrasesDTO) error
	Exists(ctx context.Context, tgID int64) (bool, error)
	CreateUser(ctx context.Context, dto db_dto.UserDTO) error
	CreateSettings(ctx context.Context, dto db_dto.SettingsDTO) error
}

type PhraseService struct {
	DB PhraseKeeper
	AI ai_agents.PhraseGenerator
}

func NewPhraseService(db PhraseKeeper, ai ai_agents.PhraseGenerator) *PhraseService {
	return &PhraseService{
		DB: db,
		AI: ai,
	}
}

func (p *PhraseService) AddPhrases(ctx context.Context, dto service_dto.UserSettings) error {
	resp, err := p.AI.GeneratePhrases(ctx, 20, dto.Level, dto.Language, dto.Topic)
	if err != nil {
		return fmt.Errorf("failed to generate phrases: %w", err)
	}

	phrasesDTOs := make([]*db_dto.PhrasesDTO, 0, len(resp.Phrases))

	for _, v := range resp.Phrases {
		phrasesDTOs = append(phrasesDTOs, &db_dto.PhrasesDTO{
			TargetLanguage: "English",
			Level:          "A1",
			Topic:          "Travel",
			InLanguageText: v.PhraseInLanguage,
			InRussianText:  v.PhraseInRussian,
			GeneratedBy:    "groq",
		})
	}
	err = p.DB.CreatePhrasesBatch(ctx, phrasesDTOs)
	if err != nil {
		return fmt.Errorf("AddPhrases: %w", err)
	}
	return nil
}

func (p *PhraseService) AddUser(ctx context.Context, dto service_dto.User) error {
	exists, err := p.DB.Exists(ctx, dto.TelegramUserID)
	if err != nil {
		return err
	}
	if exists {
		slog.Warn("user already exists", "user", dto.UserName)
		return nil
	}
	userDTO := db_dto.UserDTO{
		TelegramUserID: dto.TelegramUserID,
		ChatID:         dto.ChatID,
		UserName:       dto.UserName,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
	}
	err = p.DB.CreateUser(ctx, userDTO)
	if err != nil {
		return fmt.Errorf("AddUsser:%w", err)
	}
	return nil
}

func (p *PhraseService) AddSettings(ctx context.Context, dto service_dto.UserSettings) error {
	interval, err := strconv.Atoi(dto.Interval)
	if err != nil {
		return fmt.Errorf("failed to convert string to int:%w", err)
	}
	settingsDTO := db_dto.SettingsDTO{
		Language: dto.Language,
		Level:    dto.Level,
		Topic:    dto.Topic,
		Interval: interval,
	}
	err = p.DB.CreateSettings(ctx, settingsDTO)
	if err != nil {
		return fmt.Errorf("AddSettings: %w", err)
	}
	return nil
}

func (p *PhraseService) ProcessPhrase() {

}
