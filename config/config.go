package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error(".env file doesn't exist")
		slog.Info("reading env variables from os")

		// If .env file doesn't exist then load env vars from os
		switch {
		case os.Getenv("APP_PORT") == "":
			return fmt.Errorf("please provide APP_PORT env variable")
		case os.Getenv("PG_DSN") == "":
			return fmt.Errorf("please provide PG_DSN env variables")
		case os.Getenv("ENV") == "":
			return fmt.Errorf("please provide ENV env variables")
		}
	}

	return nil
}
