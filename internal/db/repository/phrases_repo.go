package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

var ErrPhraseNotExists = errors.New("phrase not found")

func (r *Repository) CreatePhrasesBatch(ctx context.Context, dtos []*db_dto.GetPhrasesDTO) error {
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

func (r *Repository) GetPhrase(ctx context.Context, dto db_dto.GetPhrasesDTO) (*db_dto.FetchPhraseDTO, error) {
	query := `SELECT 
    in_language_text,
    in_russian_text 
	FROM phrases WHERE target_language=$1, level=$2, topic=$3;`
	var phrases db_dto.FetchPhraseDTO
	err := r.Conn.QueryRow(ctx, query, dto.TargetLanguage, dto.Level, dto.Topic).Scan(&phrases)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &db_dto.FetchPhraseDTO{}, ErrPhraseNotExists
		}
		return &db_dto.FetchPhraseDTO{}, fmt.Errorf("failed to get phrases from phrases table: %w", err)
	}
	return &phrases, nil
}
