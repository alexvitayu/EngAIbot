package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

func (r *Repository) CreatePhrasesBatch(ctx context.Context, dtos []*db_dto.PhrasesDTO) error {
	if len(dtos) == 0 {
		return nil
	}
	query := `INSERT INTO phrases (
    target_language,
	level,
	topic,
	in_language_text,
	in_russian_text,
	generated_by
	) VALUES `

	var values []any
	var placeholders []string

	for i, dto := range dtos {
		offset := i * 6
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6),
		)
		values = append(values,
			dto.TargetLanguage, dto.Level, dto.Topic,
			dto.InLanguageText, dto.InRussianText, dto.GeneratedBy,
		)
	}

	query += strings.Join(placeholders, ",")

	_, err := r.Conn.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to add phrases to DB: %w", err)
	}
	return nil
}
