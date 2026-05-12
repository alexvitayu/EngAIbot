package tg_bot

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/alexvitayu/EngAIbot/internal/service"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TgBot(ctx context.Context, cfg config.AppConfig, count int, service service.PhraseService) {
	// 1. Инициализация бота
	bot, err := tgbotapi.NewBotAPI(cfg.TgAPI)
	if err != nil {
		slog.Error("TgBot: %w", err)
		log.Panic(err) //TODO
	}

	slog.Info("EngAIbot authorized successfully", "username", bot.Self.UserName)

	// 2. Настраиваем получение обновлений (long polling)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	var wg sync.WaitGroup

	pool := NewWorkerPool(bot, updates, count, &wg, service)

	pool.Start(ctx) // запускаем воркеров для обработки сообщений из канала updates

	<-ctx.Done() // получаем сигнал на остановку
	pool.Stop()  // говорим воркерам остановиться
	wg.Wait()    // ждём, пока все воркеры завершатся
}

func HandleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	callbackConfig := tgbotapi.NewCallback(query.ID, "")
	bot.Request(callbackConfig)

	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	switch {
	case query.Data == "menu_start":
		sendMainMenu(bot, chatID, messageID)
	case query.Data == "settings_language":
		showLanguageSelector(bot, chatID, messageID)
	case query.Data == "settings_level":
		showLevelSelector(bot, chatID, messageID)
	case query.Data == "settings_topic":
		showTopicSelector(bot, chatID, messageID)
	case query.Data == "settings_interval":
		showIntervalSelector(bot, chatID, messageID)
	case query.Data == "settings_confirm":
		confirmSettings(bot, chatID, messageID)
	case query.Data == "menu_back":
		SendStartMenu(bot, chatID, messageID) // Возврат в главное меню
	case strings.HasPrefix(query.Data, "lang_"):
		saveLanguage(bot, chatID, messageID, query.Data[5:])
	case strings.HasPrefix(query.Data, "level_"):
		saveLevel(bot, chatID, messageID, query.Data[6:])
	case strings.HasPrefix(query.Data, "topic_"):
		saveTopic(bot, chatID, messageID, query.Data[6:])
	case strings.HasPrefix(query.Data, "interval_"):
		saveInterval(bot, chatID, messageID, query.Data[9:])
	case query.Data == "confirm_yes":
		startLearning(bot, chatID, messageID)
	case query.Data == "confirm_no":
		clearSettings(bot, chatID, messageID)
	}
}

// Функция для показа главного меню (с возможностью редактировать сообщение)
func SendStartMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "👋 Добро пожаловать!\n\nВыберите действие:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚀 Начать учиться", "menu_start"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ О боте", "menu_about"),
		),
	)

	// Если передан messageID — редактируем существующее сообщение
	if messageID != 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
		editMsg.ReplyMarkup = &keyboard
		bot.Send(editMsg)
	} else {
		// Иначе отправляем новое
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}

// Хранилище временных настроек (в памяти)
var userTempSettings = make(map[int64]*service_dto.UserSettings)
var mu sync.RWMutex

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "🎯 *Настройка обучения*\n\nВыберите параметр для настройки:"

	// Кнопки для каждого параметра
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 Язык: не выбран", "settings_language"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Уровень: не выбран", "settings_level"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Тема: не выбрана", "settings_topic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏰ Периодичность: не выбрана", "settings_interval"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Начать обучение", "settings_confirm"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "menu_back"),
		),
	)

	// Получаем текущие настройки пользователя
	settings := GetTempSettings(chatID)
	if settings != nil {
		// Обновляем текст кнопок с выбранными значениями
		updateButtonsWithSettings(&keyboard, settings)
	}

	if messageID != 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
		editMsg.ParseMode = "Markdown"
		editMsg.ReplyMarkup = &keyboard
		bot.Send(editMsg)
	} else {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}

// updateButtonsWithSettings обновляет текст кнопок в существующей клавиатуре
func updateButtonsWithSettings(keyboard *tgbotapi.InlineKeyboardMarkup, settings *service_dto.UserSettings) {
	if settings == nil {
		return
	}

	// Обновляем каждую кнопку в зависимости от её CallbackData
	for i, row := range keyboard.InlineKeyboard {
		for j, button := range row {
			switch *button.CallbackData {
			case "settings_language":
				if settings.Language != "" {
					keyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("🌍 Язык: %s", settings.Language)
				}
			case "settings_level":
				if settings.Level != "" {
					levelDisplay := formatLevel(settings.Level)
					keyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("📊 Уровень: %s", levelDisplay)
				}
			case "settings_topic":
				if settings.Topic != "" {
					topicDisplay := getTopicDisplay(settings.Topic)
					keyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("📚 Тема: %s", topicDisplay)
				}
			case "settings_interval":
				if settings.Interval != "" {
					intervalDisplay := getIntervalDisplay(settings.Interval)
					keyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("⏰ Периодичность: %s", intervalDisplay)
				}
			}
		}
	}
}

// formatLevel преобразует код уровня в читаемый вид
func formatLevel(level string) string {
	// a1 → A1, b1 → B1
	return strings.ToUpper(level)
}

// getTopicDisplay преобразует код темы в читаемое название
func getTopicDisplay(topicCode string) string {
	topics := map[string]string{
		"travel":        "✈️ Путешествия",
		"work":          "💼 Работа",
		"business":      "💼 Бизнес",
		"family":        "👨‍👩‍👧 Семья",
		"food":          "🍕 Еда",
		"health":        "❤️ Здоровье",
		"hobbies":       "🎮 Хобби",
		"technology":    "💻 Технологии",
		"tech":          "💻 Технологии",
		"weather":       "🌤️ Погода",
		"education":     "📚 Образование",
		"shopping":      "🛍️ Покупки",
		"daily life":    "🏠 Повседневность",
		"daily_life":    "🏠 Повседневность",
		"home_routines": "🏠 Домашние дела",
		"any subject":   "🎯 Разное",
	}

	if display, exists := topics[topicCode]; exists {
		return display
	}
	// Если не нашли, возвращаем с заглавной буквы
	if len(topicCode) > 0 {
		return strings.ToUpper(topicCode[:1]) + topicCode[1:]
	}
	return topicCode
}

// getIntervalDisplay преобразует код интервала в читаемое название
func getIntervalDisplay(intervalCode string) string {
	intervals := map[string]string{
		"1":  "⏰ Каждый час",
		"3":  "🕒 Каждые 3 часа",
		"6":  "📅 Каждые 6 часов",
		"12": "🌙 Каждые 12 часов",
		"24": "🌟 Раз в день",
	}

	if display, exists := intervals[intervalCode]; exists {
		return display
	}
	return fmt.Sprintf("⏰ Раз в %s часов", intervalCode)
}

// saveLevel сохраняет выбранный уровень
func saveLevel(bot *tgbotapi.BotAPI, chatID int64, messageID int, level string) {
	formattedLevel := strings.ToUpper(level)
	saveTempSetting(chatID, "level", formattedLevel)
	bot.Send(tgbotapi.NewCallback(strconv.Itoa(messageID), "✅ Уровень сохранён!"))

	// Возвращаемся в главное меню настроек
	sendMainMenu(bot, chatID, messageID)
}

// saveTopic сохраняет выбранную тему
func saveTopic(bot *tgbotapi.BotAPI, chatID int64, messageID int, topic string) {
	// Сохраняем тему во временные настройки
	saveTempSetting(chatID, "topic", topic)

	// Показываем уведомление
	bot.Send(tgbotapi.NewCallback(strconv.Itoa(messageID), "✅ Тема сохранена!"))

	// Возвращаемся в главное меню настроек
	sendMainMenu(bot, chatID, messageID)
}

// saveInterval сохраняет выбранный интервал
func saveInterval(bot *tgbotapi.BotAPI, chatID int64, messageID int, interval string) {
	// Сохраняем интервал во временные настройки
	saveTempSetting(chatID, "interval", interval)

	// Показываем уведомление
	bot.Send(tgbotapi.NewCallback(strconv.Itoa(messageID), "✅ Периодичность сохранена!"))

	// Возвращаемся в главное меню настроек
	sendMainMenu(bot, chatID, messageID)
}

// Показываем выбор языка
func showLanguageSelector(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "🌍 *Выберите язык для изучения:*"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang_english"),
			tgbotapi.NewInlineKeyboardButtonData("🇵🇱 Polish", "lang_polish"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇩🇪 German", "lang_german"),
			tgbotapi.NewInlineKeyboardButtonData("🇮🇹 Italian", "lang_italian"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к настройкам", "back_to_main_menu"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)
}

// Сохраняем выбранный язык
func saveLanguage(bot *tgbotapi.BotAPI, chatID int64, messageID int, language string) {
	saveTempSetting(chatID, "language", language)

	// Показываем уведомление
	bot.Send(tgbotapi.NewCallback(strconv.Itoa(messageID), "✅ Язык сохранён!"))

	// Возвращаемся в главное меню настроек
	sendMainMenu(bot, chatID, messageID)
}

func showLevelSelector(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "📊 *Выберите уровень владения языком:*"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔰 A1 (Начальный)", "level_a1"),
			tgbotapi.NewInlineKeyboardButtonData("⭐ A2 (Элементарный)", "level_a2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚀 B1 (Средний)", "level_b1"),
			tgbotapi.NewInlineKeyboardButtonData("💪 B2 (Выше среднего)", "level_b2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к настройкам", "back_to_main_menu"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)
}

func showTopicSelector(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "📚 *Выберите тему для фраз:*"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✈️ Путешествия", "topic_travel"),
			tgbotapi.NewInlineKeyboardButtonData("💼 Работа", "topic_work"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍👩‍👧 Семья", "topic_family"),
			tgbotapi.NewInlineKeyboardButtonData("🍕 Еда", "topic_food"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к настройкам", "back_to_main_menu"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)
}

func showIntervalSelector(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	text := "⏰ *Как часто присылать новые фразы?*"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каждый час", "interval_1"),
			tgbotapi.NewInlineKeyboardButtonData("Каждые 3 часа", "interval_3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каждые 6 часов", "interval_6"),
			tgbotapi.NewInlineKeyboardButtonData("Раз в день", "interval_24"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к настройкам", "back_to_main_menu"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)
}

func confirmSettings(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	settings := GetTempSettings(chatID)
	if settings.Language == "" || settings.Level == "" || settings.Topic == "" || settings.Interval == "" {
		text := "⚠️ *Не все параметры выбраны!*\n\nПожалуйста, заполните все настройки перед началом обучения."

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Вернуться к настройкам", "back_to_main_menu"),
			),
		)

		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
		editMsg.ParseMode = "Markdown"
		editMsg.ReplyMarkup = &keyboard
		bot.Send(editMsg)
		return
	}
	text := fmt.Sprintf(
		"📋 *Проверьте ваши настройки:*\n\n"+
			"🌍 Язык: *%s*\n"+
			"📊 Уровень: *%s*\n"+
			"📚 Тема: *%s*\n"+
			"⏰ Периодичность: *%s*\n\n"+
			"❓ Всё верно?",
		settings.Language, settings.Level, settings.Topic, settings.Interval,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, всё верно", "confirm_yes"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, изменить", "confirm_no"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)
}

func startLearning(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	settings := GetTempSettings(chatID)
	text := fmt.Sprintf(
		"🎉 *Отлично! Настройки сохранены!*\n\n"+
			"🌍 Язык: *%s*\n"+
			"📊 Уровень: *%s*\n"+
			"📚 Тема: *%s*\n"+
			"⏰ Периодичность: *%s*\n\n"+
			"✨ Скоро вы получите первую фразу!\n"+
			"А пока можете начать изучение прямо сейчас:",
		settings.Language, settings.Level, settings.Topic, settings.Interval,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Получить фразу сейчас", "get_phrase_now"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "menu_back"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &keyboard
	bot.Send(editMsg)

	// Запускаем планировщик для пользователя
	//startSchedulerForUser(chatID, settings)
}

func clearSettings(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	// Очищаем временные настройки
	clearTempSettings(chatID)

	// Возвращаем в главное меню настроек
	sendMainMenu(bot, chatID, messageID)

	// Отправляем уведомление
	bot.Send(tgbotapi.NewMessage(chatID, "🔄 Настройки сброшены. Выберите параметры заново."))
}

func GetTempSettings(userID int64) *service_dto.UserSettings {
	mu.RLock()
	defer mu.RUnlock()
	return userTempSettings[userID]
}

func saveTempSetting(userID int64, key, value string) {
	mu.Lock()
	defer mu.Unlock()

	if userTempSettings[userID] == nil {
		userTempSettings[userID] = &service_dto.UserSettings{}
	}

	switch key {
	case "language":
		userTempSettings[userID].Language = value
	case "level":
		userTempSettings[userID].Level = value
	case "topic":
		userTempSettings[userID].Topic = value
	case "interval":
		userTempSettings[userID].Interval = value
	}
}

func clearTempSettings(userID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(userTempSettings, userID)
}
