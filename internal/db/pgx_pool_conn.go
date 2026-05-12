package db

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPgxPool(ctx context.Context, cfg *config.AppConfig) (*pgxpool.Pool, error) {
	// Парсим конфиг из DSN
	conf, err := pgxpool.ParseConfig(cfg.DBConfig.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("parse conf: %w", err)
	}
	maxConns, minConns, maxLifetime := parseCfg(cfg)
	conf.MaxConns = maxConns
	conf.MinConns = minConns
	conf.MaxConnIdleTime = maxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

func parseCfg(cfg *config.AppConfig) (int32, int32, time.Duration) {
	maxConns, err := strconv.ParseInt(cfg.PoolConfig.DBMaxConns, 10, 32)
	if err != nil {
		slog.Error("ParseInt", "error", err)
	}
	minConns, err := strconv.ParseInt(cfg.PoolConfig.DBMaxIdleConns, 10, 32)
	if err != nil {
		slog.Error("ParseInt", "error", err)
	}
	maxLifetime, err := time.ParseDuration(cfg.PoolConfig.DBConnMaxLifetime)
	if err != nil {
		slog.Error("ParseDuration", "error", err)
	}
	return int32(maxConns), int32(minConns), maxLifetime
}
