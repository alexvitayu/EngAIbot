package repository

import (
	"context"
	"fmt"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

func (r *Repository) CreateSettings(ctx context.Context, dto db_dto.SettingsDTO) error {
	query := `INSERT INTO user_settings (
                           language,
                           level,
                           topic,
                           interval_hours)
                           VALUES ($1, $2, $3, $4);`
	_, err := r.Conn.Exec(ctx, query, dto.Language, dto.Level, dto.Topic, dto.Interval)
	if err != nil {
		return fmt.Errorf("failed to add user settings to DB: %w", err)
	}
	return nil
}
