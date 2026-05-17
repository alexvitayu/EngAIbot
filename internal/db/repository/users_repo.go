package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
	"github.com/alexvitayu/EngAIbot/internal/service/service_dto"
)

type ScheduledUser struct {
	TgUserID int64
	ChatID   int64
	Settings service_dto.UserSettings
}

func (r *Repository) ExistsUser(ctx context.Context, tgID int64) (int64, bool, error) {
	var userID int64
	query := "SELECT id FROM users WHERE tg_user_id = $1"
	err := r.Conn.QueryRow(ctx, query, tgID).Scan(&userID)

	if err == nil {
		return userID, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return 0, false, err
}

func (r *Repository) CreateUser(ctx context.Context, dto db_dto.UserDTO) (int64, error) {
	var userID int64
	query := `INSERT INTO users (
                   tg_user_id,
                   user_name,
                   first_name,
                   last_name,
                   chat_id)
                   VALUES ($1, $2, $3, $4, $5) RETURNING(id);`
	err := r.Conn.QueryRow(ctx, query, dto.TelegramUserID, dto.UserName, dto.FirstName, dto.LastName, dto.ChatID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to add user to DB: %w", err)
	}
	return userID, nil
}

func (r *Repository) GetUsersForScheduler(ctx context.Context) ([]ScheduledUser, error) {
	query := `
        SELECT 
            u.tg_user_id,
            u.chat_id,
            us.language,
            us.level,
            us.topic,
            us.interval_hours
        FROM users u
        INNER JOIN user_settings us ON u.id = us.user_id
        WHERE us.interval_hours IS NOT NULL`

	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query scheduled users: %w", err)
	}
	defer rows.Close()

	var users []ScheduledUser
	for rows.Next() {
		var su ScheduledUser
		err = rows.Scan(
			&su.TgUserID,
			&su.ChatID,
			&su.Settings.Language,
			&su.Settings.Level,
			&su.Settings.Topic,
			&su.Settings.Interval,
		)
		if err != nil {
			slog.Error("failed to scan scheduled user", "error", err)
			continue
		}
		users = append(users, su)
	}

	return users, rows.Err()
}
