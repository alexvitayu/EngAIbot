package service_dto

import "time"

type User struct {
	TelegramUserID int64
	UserName       string
	FirstName      string
	LastName       string
	Settings       []UserSettings //пользователь может изучать не один язык
	ChatID         int64
	CreatedAt      time.Time
}

type UserSettings struct {
	UserID   int64
	TgUserID int64
	Language string
	Level    string
	Topic    string
	Interval string
}
