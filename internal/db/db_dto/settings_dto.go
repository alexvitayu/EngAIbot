package db_dto

import "time"

type SettingsDTO struct {
	Language  string
	Level     string
	Topic     string
	Interval  int
	IsActive  bool
	CreateAt  time.Time
	UpdatedAt time.Time
}
