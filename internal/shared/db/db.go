package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(dsn string) (*sqlx.DB, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Avoid server-side prepared statements which can break behind transaction poolers (pgbouncer).
	if cfg.DefaultQueryExecMode == pgx.QueryExecModeCacheStatement {
		cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}
	cfg.StatementCacheCapacity = 0

	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
