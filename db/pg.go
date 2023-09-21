package db

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL() *pgxpool.Pool {
	pgConf, err := pgxpool.ParseConfig(os.Getenv("PG_DSN"))
	if err != nil {
		slog.Error("Error parsing postgres dns address", err)
	}
	pgConf.MaxConns = 5
	pgConf.MinConns = 0
	pgConf.MaxConnLifetime = time.Hour * 1
	pgConf.MaxConnIdleTime = time.Minute * 15

	db, err := pgxpool.NewWithConfig(context.Background(), pgConf)
	if err != nil {
		slog.Error("PGSQL Init Failed", "cause", err)
		os.Exit(1)
	}

	if err := db.Ping(context.Background()); err != nil {
		slog.Error("PostgreSQL Ping error", "cause", err)
	}

	slog.Info("PostgreSQL DB Connection Established")
	return db
}
