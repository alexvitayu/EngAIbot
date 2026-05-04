package tg_bot

//
//import (
//	"log"
//	"log/slog"
//	"time"
//
//	"github.com/alexvitayu/EngAIbot/internal/config"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//type User struct {
//	UserID    int64 //это уникальный числовой идентификатор, который не меняется даже при смене username
//	UserName  string
//	FirstName string
//	LastName  string
//	Settings  []Settings //пользователь может изучать не один язык
//	ChatID    int64
//}
//
//type Settings struct {
//	Language string
//	Level    string
//	Topic    string
//	Interval time.Duration
//}
//
//func TgBot(cfg config.AppConfig) {
//	// 1. Инициализация бота
//	bot, err := tgbotapi.NewBotAPI(cfg.TgAPI)
//	if err != nil {
//		slog.Error("TgBot: %w", err)
//		log.Panic(err) //TODO
//	}
//
//	slog.Info("EngAIbot authorized successfully", "username", bot.Self.UserName)
//
//	// 2. Настраиваем получение обновлений (long polling)
//	updateConfig := tgbotapi.NewUpdate(0)
//	updateConfig.Timeout = 60
//	updates := bot.GetUpdatesChan(updateConfig)
//
//	for update := range updates {
//		ExtractUser(update) //собираем информацию о новом пользователе
//		// Обработка команды /start
//		if update.Message != nil && update.Message.IsCommand() {
//			switch update.Message.Command() {
//			case "start":
//				sendStartMenu(bot, update.Message.Chat.ID)
//			}
//			continue
//		}
//
//		// Обработка нажатий на кнопки
//		if update.CallbackQuery != nil {
//			handleCallback(bot, update.CallbackQuery)
//		}
//	}
//}
//
//// Функция для показа главного меню (с возможностью редактировать сообщение)
//func sendStartMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
//	text := "👋 Добро пожаловать!\n\nВыберите действие:"
//
//	keyboard := tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🚀 Начать учиться", "menu_start"),
//		),
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("ℹ️ О боте", "menu_about"),
//		),
//	)
//
//	// Если передан messageID — редактируем существующее сообщение
//	if messageID != 0 {
//		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
//		editMsg.ReplyMarkup = &keyboard
//		bot.Send(editMsg)
//	} else {
//		// Иначе отправляем новое
//		msg := tgbotapi.NewMessage(chatID, text)
//		msg.ReplyMarkup = keyboard
//		bot.Send(msg)
//	}
//}
//
//// handleCallback обрабатывает нажатия кнопок
//func handleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
//	// Отвечаем на callback
//	callbackConfig := tgbotapi.NewCallback(query.ID, "")
//	bot.Request(callbackConfig)
//
//	// Обрабатываем разные кнопки
//	switch query.Data {
//	case "menu_start":
//		// Показываем меню выбора языка
//		showLanguageMenu(bot, query.Message.Chat.ID, query.Message.MessageID)
//
//	case "menu_about":
//		response := "ℹ️ О боте\n\nЭто учебный бот для изучения иностранных языков. Выберите, какой" +
//			"язык хотите изучать, выберите уровень языка, тему, на которую хотите получать фразы. Не забудьте" +
//			"настроить время, через которое будут приходить новые фразы.\n\nВерсия: 1.0.0"
//
//		// Отправляем новое сообщение (или редактируем текущее)
//		msg := tgbotapi.NewMessage(query.Message.Chat.ID, response)
//
//		// Добавляем кнопку "Назад" в конец сообщения о боте
//		keyboard := tgbotapi.NewInlineKeyboardMarkup(
//			tgbotapi.NewInlineKeyboardRow(
//				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в главное меню", "menu_back_to_main"),
//			),
//		)
//		msg.ReplyMarkup = keyboard
//		bot.Send(msg)
//
//	case "menu_back_to_main":
//		// Возвращаемся в главное меню (редактируем текущее сообщение)
//		sendStartMenu(bot, query.Message.Chat.ID, query.Message.MessageID)
//
//	case "menu_back":
//		// Возврат на предыдущий шаг (например, из меню языков в главное меню)
//		sendStartMenu(bot, query.Message.Chat.ID, query.Message.MessageID)
//	}
//}
//
//func showLanguageMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
//	text := "🌍 Выберите язык для изучения:"
//
//	keyboard := tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang_english"),
//			tgbotapi.NewInlineKeyboardButtonData("🇵🇱 Polish", "lang_polish"),
//		),
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🇩🇪 German", "lang_german"),
//			tgbotapi.NewInlineKeyboardButtonData("🇪🇸 Spanish", "lang_spanish"),
//		),
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "menu_back"),
//		),
//	)
//
//	// Редактируем текущее сообщение (заменяем главное меню на меню языков)
//	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
//	editMsg.ReplyMarkup = &keyboard
//	bot.Send(editMsg)
//}
//
//func handleSettings(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
//	// Отвечаем на callback
//	callbackConfig := tgbotapi.NewCallback(query.ID, "")
//	bot.Request(callbackConfig)
//
//	switch query.Data {
//	case "select_language":
//
//	}
//}
//
//func ExtractUser(update tgbotapi.Update) *User {
//	var user *tgbotapi.User
//	var chatID int64
//
//	switch {
//	case update.Message != nil:
//		user = update.Message.From
//		chatID = update.Message.Chat.ID
//
//	case update.CallbackQuery != nil:
//		user = update.CallbackQuery.From
//		chatID = update.CallbackQuery.Message.Chat.ID
//
//	default:
//		return nil
//	}
//
//	return &User{
//		UserID:    user.ID,
//		UserName:  user.UserName,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		ChatID:    chatID,
//	}
//}
//
//func SetUpSettings() {
//
//}
