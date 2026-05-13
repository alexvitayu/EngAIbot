package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		Conn: conn,
	}
}
