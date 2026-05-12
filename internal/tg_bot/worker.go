package tg_bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
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
			if update.Message != nil && update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					SendStartMenu(wp.bot, update.Message.Chat.ID, 0)
				}
				continue
			}

			if update.CallbackQuery != nil && update.CallbackQuery.Data == "confirm_yes" {
				userID := update.CallbackQuery.From.ID
				chatID := update.CallbackQuery.Message.Chat.ID
				messageID:= update.CallbackQuery.Message.MessageID
				userDTO := service_dto.User{
					TelegramUserID: userID,
					UserName:       update.CallbackQuery.From.UserName,
					FirstName:      update.CallbackQuery.From.FirstName,
					LastName:       update.CallbackQuery.From.LastName,
					Settings:       nil,
					ChatID:         chatID,
					CreatedAt:      time.Time{},
				}
				err := wp.service.AddUser(ctx, userDTO)
				if err != nil {
					wp.errChan <- err
				}
				settings := GetTempSettings(userID)
				err = wp.service.AddSettings(ctx, *settings)
				if err != nil {
					wp.errChan <- err
				}
				startLearning(wp.bot, chatID, messageID)
			}
			// Обработка нажатий на кнопки
			if update.CallbackQuery != nil && update.CallbackQuery.Data != "confirm_yes" {
				HandleCallback(wp.bot, update.CallbackQuery)
			}
		}
	}
}
