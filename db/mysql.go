package db

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/prawirdani/go-restapi-boilerplate/config"
)

func NewMySQL(c config.Config) *sql.DB {
	db, err := sql.Open("mysql", c.Mysql.DSN)
	if err != nil {
		slog.Error("Mysql init error", "cause", err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		slog.Error("Mysql Ping error", "cause", err)
		os.Exit(1)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	slog.Info("MySQL DB Connection Established")
	return db
}
