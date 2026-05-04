package ai_agents

import (
	"context"
)

type AIAgents interface {
	GenerateInfo(ctx context.Context, quota, level int, lang, subject string) (*AIResponse, error)
}

type AIResponse struct {
	Phrases []Phrase `json:"phrases"`
}

type Phrase struct {
	PhraseInLanguage string `json:"phrase_in_language"`
	PhraseInRussian  string `json:"phrase_in_russian"`
	Level            string `json:"level"`
	Topic            string `json:"topic"`
}
