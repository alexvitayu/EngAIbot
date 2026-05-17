package tg_bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (wp *WorkerPool) worker(ctx context.Context, idx int) {
	defer func() {
		wp.wg.Done()
		slog.Debug(fmt.Sprintf("worker %d finished", idx))
	}()

	for {
		select {
		case <-wp.stopChan:
			return
		case update, ok := <-wp.updates:
			if !ok {
				return
			}
			var userSettings *service_dto.UserSettings
			var TGUserID int64
			var chatID int64
			var messageID int
			var UserName string
			var FirstName string
			var LastName string
			if update.Message != nil {
				TGUserID = update.Message.From.ID
				chatID = update.Message.Chat.ID
				messageID = update.Message.MessageID
				UserName = update.Message.From.UserName
				FirstName = update.Message.From.FirstName
				LastName = update.Message.From.LastName
			}
			if update.CallbackQuery != nil {
				TGUserID = update.CallbackQuery.From.ID
				chatID = update.CallbackQuery.Message.Chat.ID
				messageID = update.CallbackQuery.Message.MessageID
				UserName = update.CallbackQuery.From.UserName
				FirstName = update.CallbackQuery.From.FirstName
				LastName = update.CallbackQuery.From.LastName
			}

			userDTO := service_dto.User{
				TelegramUserID: TGUserID,
				UserName:       UserName,
				FirstName:      FirstName,
				LastName:       LastName,
				Settings:       nil,
				ChatID:         chatID,
				CreatedAt:      time.Time{},
			}

			if update.Message != nil && update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					SendStartMenu(wp.bot, update.Message.Chat.ID, 0)
				}
			}
			if update.CallbackQuery != nil && update.CallbackQuery.Data == "confirm_yes" {
				settings := GetTempSettings(TGUserID)
				userSettings = &service_dto.UserSettings{
					TgUserID: TGUserID,
					Language: settings.Language,
					Level:    settings.Level,
					Topic:    settings.Topic,
					Interval: settings.Interval,
				}

				userID, err := wp.service.AddUser(ctx, userDTO)
				if err != nil {
					wp.errChan <- err
				}
				userSettings.UserID = userID
				err = wp.service.AddOrUpdateSettings(ctx, *userSettings)
				if err != nil {
					wp.errChan <- err
				}
				startLearning(wp.bot, chatID, messageID, wp.Sched, userSettings)
			}
			if update.CallbackQuery != nil && update.CallbackQuery.Data == "get_phrase_now" {
				settings := GetTempSettings(TGUserID)
				userSettings = &service_dto.UserSettings{
					TgUserID: TGUserID,
					Language: settings.Language,
					Level:    settings.Level,
					Topic:    settings.Topic,
					Interval: settings.Interval,
				}
				phrase, err := wp.service.SendPhraseNow(ctx, *userSettings)
				if err != nil {
					wp.errChan <- err
				}
				editMsg := tgbotapi.NewEditMessageText(chatID, messageID, phrase.InLanguageText)
				editMsg.ParseMode = "Markdown"
				wp.bot.Send(editMsg)
			}
			if update.CallbackQuery != nil && update.CallbackQuery.Data != "confirm_yes" {
				settings := GetTempSettings(TGUserID)
				HandleCallback(wp.bot, update.CallbackQuery, wp.Sched, settings)
			}
		}
	}
}
