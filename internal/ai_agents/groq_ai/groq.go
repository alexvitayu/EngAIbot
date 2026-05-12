package groq_ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/types"
	"github.com/alexvitayu/EngAIbot/internal/ai_agents"
	"github.com/alexvitayu/EngAIbot/internal/config"
)

type Groq struct {
	client *groq.Client
}

func NewGroq(cfg *config.AppConfig) (*Groq, error) {
	client, err := groq.NewClient(
		groq.WithAPIKey(cfg.GroqAPI),
	)
	if err != nil {
		slog.Error("failed to create a new groq client", "error", err)
		return &Groq{}, fmt.Errorf("NewGroq: %w", err)
	}
	g := &Groq{
		client: client,
	}
	return g, nil
}

func (g *Groq) GeneratePhrases(ctx context.Context, quota int, level, lang, topic string) (*ai_agents.AIResponse, error) {
	prompt := fmt.Sprintf(`
Ты - генератор учебных материалов для изучения %s языка.
Сгенерируй %d фраз на языке=%s для студентов уровня=%s на тему=%s.

ТРЕБОВАНИЯ:
1. Ответ ДОЛЖЕН быть строго в формате JSON
2. Используй структуру:
{
"phrases": [{
			"phrase_in_language": "здесь находится сгенерированная тобой фраза на языке %s",
			"phrase_in_russian": "здесь точный перевод сгенерированной фразы на русский язык",
			"level": "здесь языковой уровень сгенерированной фразы по международной системе CEFR (А1, А2, В1, B2, C1, C2) %s",
			"topic": "%s"}]
}
3. Перевод должен быть точным и естественным на русском языке
4. Не добавляй никакой текст до или после JSON
5. Только валидный JSON без лишних символов`,
		lang,
		quota, lang, level, topic,
		lang, level, topic,
	)

	resp, err := g.client.Chat.Create(ctx, &types.CreateChatCompletionRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleSystem,
				Content: "Ты - полезный ассистент, который всегда отвечает в формате JSON без markdown форматирования.",
			},
			{
				Role:    types.RoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &types.ResponseFormat{
			Type: "json_object",
		},
	})

	if err != nil {
		slog.Error("failed to get response from groq", "error", err)
		return &ai_agents.AIResponse{}, fmt.Errorf("GeneratePhrases: %w", err)
	}

	if len(resp.Choices) == 0 {
		slog.Info("no answers in response from groq")
	}

	jsonResponse := resp.Choices[0].Message.Content

	var phrases ai_agents.AIResponse
	if err := json.Unmarshal([]byte(jsonResponse), &phrases); err != nil {
		slog.Error("failed to unmarshal json", "error", err)
		return &ai_agents.AIResponse{}, fmt.Errorf("GeneratePhrases: %w", err)
	}

	slog.Debug("Quota of generated phrases", "quota", len(phrases.Phrases))
	return &phrases, nil
}
