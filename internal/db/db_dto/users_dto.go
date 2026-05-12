package db_dto

import "time"

type UserDTO struct {
	TelegramUserID int64
	ChatID         int64
	UserName       string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
}
