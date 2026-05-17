package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
	"github.com/alexvitayu/EngAIbot/internal/db/repository"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	bot  *tgbotapi.BotAPI
	repo *repository.Repository
	mu   sync.RWMutex
	jobs map[int64]cron.EntryID
}

func NewScheduler(repo *repository.Repository) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
		repo: repo,
		jobs: make(map[int64]cron.EntryID),
	}
}

func (s *Scheduler) SetBot(bot *tgbotapi.BotAPI) {
	s.bot = bot
}

func (s *Scheduler) Start() {
	s.cron.Start()
	slog.Info("scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	slog.Info("scheduler stopped")
}

func (s *Scheduler) LoadAndStart(ctx context.Context) error {
	users, err := s.repo.GetUsersForScheduler(ctx)
	if err != nil {
		return fmt.Errorf("failed to load scheduled users: %w", err)
	}

	slog.Info("loading scheduled jobs", "users_found", len(users))

	for _, user := range users {
		interval, err := strconv.Atoi(user.Settings.Interval)
		if err != nil {
			slog.Error("invalid interval",
				"user_id", user.TgUserID,
				"interval", user.Settings.Interval,
				"error", err)
			continue
		}

		if err = s.AddJob(user.TgUserID, user.ChatID, interval, user.Settings); err != nil {
			slog.Error("failed to restore job",
				"user_id", user.TgUserID,
				"error", err)
		}
	}

	s.Start()
	slog.Info("scheduler started", "active_jobs", len(s.jobs))

	return nil
}

// AddJob добавляет задание в планировщик для пользователя
func (s *Scheduler) AddJob(userID, chatID int64, intervalHours int, settings service_dto.UserSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если уже есть задание для этого пользователя - удаляем старое
	if oldJobID, exists := s.jobs[userID]; exists {
		s.cron.Remove(oldJobID)
		delete(s.jobs, userID)
		slog.Info("removed existing job for user", "user_id", userID)
	}

	cronSpec := fmt.Sprintf("@every 20s")

	// Добавляем функцию в cron
	jobID, err := s.cron.AddFunc(cronSpec, func() {
		// Контекст для фоновой задачи
		bgCtx := context.Background()

		// Отправляем фразу
		if err := s.sendPhraseToUser(bgCtx, chatID, userID, settings); err != nil {
			slog.Error("failed to send scheduled phrase",
				"user_id", userID,
				"error", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	// Сохраняем ID задания
	s.jobs[userID] = jobID

	slog.Info("scheduled job added",
		"user_id", userID,
		"chat_id", chatID,
		"interval_hours", intervalHours,
		"cron_spec", cronSpec)

	return nil
}

func (s *Scheduler) sendPhraseToUser(ctx context.Context, chatID, userID int64, settings service_dto.UserSettings) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in scheduled job", "user_id", userID, "recover", r)
		}
	}()

	if s.bot == nil {
		return fmt.Errorf("bot not set in scheduler")
	}

	// Получаем случайную фразу
	phrase, err := s.repo.GetRandomPhrase(ctx, db_dto.GetPhrasesDTO{
		TargetLanguage: settings.Language,
		Level:          settings.Level,
		Topic:          settings.Topic,
	})
	if err != nil {
		return fmt.Errorf("failed to get random phrase: %w", err)
	}

	// Отправляем сообщение
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("📖 *Ваша фраза для изучения:*\n\n%s\n\n_%s_",
			phrase.InLanguageText,
			phrase.InRussianText))
	msg.ParseMode = "Markdown"

	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	slog.Info("scheduled phrase sent",
		"user_id", userID,
		"chat_id", chatID)

	return nil
}
