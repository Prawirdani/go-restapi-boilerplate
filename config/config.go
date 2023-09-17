package config

import (
	"fmt"
	"os"
)

func LoadEnv() error {
	switch {
	case os.Getenv("APP_PORT") == "":
		return fmt.Errorf("please provide APP_PORT env variable")
	case os.Getenv("PG_DSN") == "":
		return fmt.Errorf("please provide PG_DSN env variables")
	}
	return nil
}
