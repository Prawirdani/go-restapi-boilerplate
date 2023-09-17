package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Cors     CorsConfig
	Postgres PostgreSQLConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Env          string
}

type CorsConfig struct {
	Credentials bool
}

type PostgreSQLConfig struct {
	DSN string
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath("./config/")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	var c Config

	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
