package db_dto

import "time"

type SettingsDTO struct {
	UserID    int64
	Language  string
	Level     string
	Topic     string
	Interval  int
	IsActive  bool
	CreateAt  time.Time
	UpdatedAt time.Time
}
