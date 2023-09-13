package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Cors   CorsConfig
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
		log.Panicf("failed decode config into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
