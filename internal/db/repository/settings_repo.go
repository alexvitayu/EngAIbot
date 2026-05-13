package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

func (r *Repository) CreateSettings(ctx context.Context, dto db_dto.SettingsDTO) error {
	query := `INSERT INTO user_settings (
                           user_id,
                           language,
                           level,
                           topic,
                           interval_hours)
                           VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Conn.Exec(ctx, query, dto.UserID, dto.Language, dto.Level, dto.Topic, dto.Interval)
	if err != nil {
		return fmt.Errorf("failed to add user settings to DB: %w", err)
	}
	return nil
}

func (r *Repository) ExistsSetting(ctx context.Context, userID int64, language string) (int64, bool, error) {
	var settingID int64
	query := `SELECT id FROM user_settings WHERE user_id=$1 AND language=$2;`
	err := r.Conn.QueryRow(ctx, query, userID, language).Scan(&settingID)

	if err == nil {
		return settingID, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return 0, false, err
}

func (r *Repository) UpdateSettings(ctx context.Context, dto db_dto.SettingsDTO, settingID int64) error {
	query := `UPDATE user_settings SET 
                         level=$1,
                         topic=$2,
                         interval_hours=$3
                     WHERE id=$4;`
	_, err := r.Conn.Exec(ctx, query, dto.Level, dto.Topic, dto.Interval, settingID)
	if err != nil {
		return fmt.Errorf("failed to update user settings into DB: %w", err)
	}
	return nil
}
