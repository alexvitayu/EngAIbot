package groq_ai

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

type Groq struct {
	client *groq.Client
}

func NewGroq(ctx context.Context) (*Groq, error) {
	client, err := groq.NewClient(
		groq.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)
	if err != nil {
		panic(err)
	}

	// Создание чат-комплишена
	resp, err := client.Chat.Create(ctx, &types.CreateChatCompletionRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Объясни квантовые вычисления в одном предложении",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	g := &Groq{
		client: client,
	}
	fmt.Println(resp.Choices[0].Message.Content)
	return g, nil
}
