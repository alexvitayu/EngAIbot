package db_dto

import "time"

type GetPhrasesDTO struct {
	TargetLanguage string
	Level          string
	Topic          string
	InLanguageText string
	InRussianText  string
	GeneratedBy    string
	UsageCount     int
	CreatedAt      time.Time
}

type FetchPhraseDTO struct {
	InLanguageText string
	InRussianText  string
}
