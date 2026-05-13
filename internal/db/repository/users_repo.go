package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

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
