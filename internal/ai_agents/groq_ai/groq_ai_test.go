package groq_ai

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/jeffphp/lingua-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var detector lingua.LanguageDetector

func init() {
	detector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.English, lingua.Polish, lingua.German).
		Build()
}

func TestGroq_GenerateInfo(t *testing.T) {
	cfg, err := config.LoadCfg()
	require.NoError(t, err)
	slog.Info("APP_ENV", "app_env", cfg.APPEnv)
	slog.Info("GroqAPIKey", "API_Key", cfg.GroqAPI)

	var testCases = []struct {
		Name     string
		Quota    int
		Lang     string
		Level    string
		Subject  string
		WantLang string
		WantLen  int
	}{
		{
			Name:     "4_phrases_in_english",
			Quota:    4,
			Lang:     "English",
			Level:    "B1",
			Subject:  "Home routines",
			WantLang: "English",
			WantLen:  4,
		},
		{
			Name:     "5_phrases_in_german",
			Quota:    5,
			Lang:     "German",
			Level:    "A1",
			Subject:  "Any subject",
			WantLang: "German",
			WantLen:  5,
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			groq, err := NewGroq(cfg)
			require.NoError(t, err)
			resp, err := groq.GeneratePhrases(t.Context(), tc.Quota, tc.Lang, tc.Level, tc.Subject)
			require.NoError(t, err)

			f, err := os.Create(fmt.Sprintf("./phrases%d.txt", i+1))
			require.NoError(t, err)
			defer f.Close()
			for _, str := range resp.Phrases {
				f.WriteString(str.PhraseInLanguage + "\n")
				f.WriteString(str.PhraseInRussian + "\n")
			}

			if len(resp.Phrases) > 0 {
				text := strings.Builder{}
				for _, str := range resp.Phrases {
					text.WriteString(str.PhraseInLanguage + " ")
				}
				fmt.Println(text.String())
				lang, exists := detector.DetectLanguageOf(text.String())
				require.True(t, exists)
				assert.Equal(t, tc.WantLang, lang.String())
			}
			assert.Equal(t, tc.WantLen, len(resp.Phrases))
		})
	}
}
