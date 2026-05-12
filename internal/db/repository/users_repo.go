package repository

import (
	"context"
	"fmt"

	"github.com/alexvitayu/EngAIbot/internal/db/db_dto"
)

func (r *Repository) Exists(ctx context.Context, tgID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE tg_user_id = $1 LIMIT 1)"
	var exists bool
	err := r.Conn.QueryRow(ctx, query, tgID).Scan(&exists)
	return exists, err
}

func (r *Repository) CreateUser(ctx context.Context, dto db_dto.UserDTO) error {
	query := `INSERT INTO users (
                   tg_user_id,
                   user_name,
                   first_name,
                   last_name,
                   chat_id)
                   VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Conn.Exec(ctx, query, dto.TelegramUserID, dto.UserName, dto.FirstName, dto.LastName, dto.ChatID)
	if err != nil {
		return fmt.Errorf("failed to add user to DB: %w", err)
	}
	return nil
}
