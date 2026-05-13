package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/alexvitayu/EngAIbot/internal/ai_agents"
	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PhraseKeeper interface {
	CreatePhrasesBatch(ctx context.Context, dtos []*db_dto.GetPhrasesDTO) error
	ExistsUser(ctx context.Context, tgID int64) (int64, bool, error)
	CreateUser(ctx context.Context, dto db_dto.UserDTO) (int64, error)
	CreateSettings(ctx context.Context, dto db_dto.SettingsDTO) error
	ExistsSetting(ctx context.Context, userID int64, language string) (int64, bool, error)
	UpdateSettings(ctx context.Context, dto db_dto.SettingsDTO, settingID int64) error
	GetPhrase(ctx context.Context, dto db_dto.GetPhrasesDTO) (*db_dto.FetchPhraseDTO, error)
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

	phrasesDTOs := make([]*db_dto.GetPhrasesDTO, 0, len(resp.Phrases))

	for _, v := range resp.Phrases {
		phrasesDTOs = append(phrasesDTOs, &db_dto.GetPhrasesDTO{
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

func (p *PhraseService) AddUser(ctx context.Context, dto service_dto.User) (int64, error) {
	id, exists, err := p.DB.ExistsUser(ctx, dto.TelegramUserID)
	if err != nil {
		return 0, err
	}
	if exists {
		slog.Warn("user already exists", "user", dto.UserName)

		return id, nil
	}
	userDTO := db_dto.UserDTO{
		TelegramUserID: dto.TelegramUserID,
		ChatID:         dto.ChatID,
		UserName:       dto.UserName,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
	}
	userID, err := p.DB.CreateUser(ctx, userDTO)
	if err != nil {
		return 0, fmt.Errorf("AddUsser:%w", err)
	}
	return userID, nil
}

func (p *PhraseService) AddOrUpdateSettings(ctx context.Context, dto service_dto.UserSettings) error {
	interval, err := strconv.Atoi(dto.Interval)
	if err != nil {
		return fmt.Errorf("failed to convert string to int:%w", err)
	}
	settingsDTO := db_dto.SettingsDTO{
		UserID:   dto.UserID,
		Language: dto.Language,
		Level:    dto.Level,
		Topic:    dto.Topic,
		Interval: interval,
	}

	settingID, exists, err := p.DB.ExistsSetting(ctx, dto.UserID, dto.Language)
	if err != nil {
		return fmt.Errorf("exists setting error:%w", err)
	}
	if exists {
		err = p.DB.UpdateSettings(ctx, settingsDTO, settingID)
		if err != nil {
			return fmt.Errorf("update settings error:%w", err)
		}
		return nil
	}
	err = p.DB.CreateSettings(ctx, settingsDTO)
	if err != nil {
		return fmt.Errorf("AddSettings: %w", err)
	}
	return nil
}

func (p *PhraseService) SendPhraseNow(ctx context.Context, dto service_dto.UserSettings) error {
	getPhraseDTO := db_dto.GetPhrasesDTO{
		TargetLanguage: dto.Language,
		Level:          dto.Level,
		Topic:          dto.Topic,
	}
	phrases, err := p.DB.GetPhrase(ctx, getPhraseDTO)
	if err != nil {
		fmt.Errorf("get phrase from DB error:%w", err)
	}
	return nil
}
