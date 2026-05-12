package db_dto

import "time"

type PhrasesDTO struct {
	TargetLanguage string
	Level          string
	Topic          string
	InLanguageText string
	InRussianText  string
	GeneratedBy    string
	UsageCount     int
	CreatedAt      time.Time
}
