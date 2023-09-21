package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), os.Getenv("PG_DSN"))
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
