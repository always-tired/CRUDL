package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	URL string
}

type Config struct {
	Env  string
	HTTP HTTPConfig
	DB   DBConfig
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Env: "dev",
		HTTP: HTTPConfig{
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	if v := os.Getenv("ENV"); v != "" {
		cfg.Env = v
	}
	if v := os.Getenv("HTTP_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return cfg, errors.New("invalid HTTP_PORT")
		}
		cfg.HTTP.Port = p
	}
	if v := os.Getenv("HTTP_READ_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.HTTP.ReadTimeout = d
		}
	}
	if v := os.Getenv("HTTP_WRITE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.HTTP.WriteTimeout = d
		}
	}
	if v := os.Getenv("HTTP_IDLE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.HTTP.IdleTimeout = d
		}
	}
	if v := os.Getenv("DB_URL"); v != "" {
		cfg.DB.URL = v
	}

	if cfg.DB.URL == "" {
		return cfg, errors.New("DB_URL is required")
	}
	if cfg.HTTP.Port == 0 {
		return cfg, errors.New("HTTP_PORT is required")
	}

	return cfg, nil
}
