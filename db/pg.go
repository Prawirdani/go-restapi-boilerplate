package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/prawirdani/go-restapi-boilerplate/config"
)

func NewPostgreSQL(c config.Config) *pgx.Conn {
	db, err := pgx.Connect(context.Background(), c.Postgres.DSN)
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

func InitSchema(pgConn *pgx.Conn) {
	slog.Info("Intializing PG Schema...")
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(50) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT current_timestamp
	)
    `
	_, err := pgConn.Exec(context.Background(), createUserTableSQL)
	if err != nil {
		slog.Error("PG Init Schema Failed", "cause", err)
		os.Exit(1)
	}
	slog.Info("Intializing PG Schema Success")
}
