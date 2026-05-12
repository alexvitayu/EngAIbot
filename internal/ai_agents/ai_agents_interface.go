package ai_agents

import (
	"context"
)

type PhraseGenerator interface {
	GeneratePhrases(ctx context.Context, quota int, level, lang, topic string) (*AIResponse, error)
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
