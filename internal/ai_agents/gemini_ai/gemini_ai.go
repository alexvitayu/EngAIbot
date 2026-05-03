package gemini_ai

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//
//	"github.com/google/generative-ai-go/genai"
//	"google.golang.org/api/option"
//)
//
//type Gemini struct {
//	client *genai.Client
//}
//
//func NewGemini(ctx context.Context, key string) (*Gemini, error) {
//	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	model := client.GenerativeModel("gemini-2.0-flash-lite")
//	model.SetTemperature(0.9)
//
//	prompt := `Сгенерируй простую фразу на английском для изучения и её перевод на русский.
//               Формат: "EN: [фраза] | RU: [перевод]"`
//
//	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(resp.Candidates[0].Content.Parts[0])
//}
