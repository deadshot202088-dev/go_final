package config

import (
	"os"
	"strconv"
)

const (
	defaultPort   = 7540
	defaultDBFile = "scheduler.db"
)

// Config хранит настройки приложения, считанные один раз при старте
type Config struct {
	Port     int
	DBFile   string
	Password string
}

// Load считывает конфигурацию из переменных окружения
func Load() Config {
	cfg := Config{
		Port:     defaultPort,
		DBFile:   defaultDBFile,
		Password: os.Getenv("TODO_PASSWORD"),
	}

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = p
		}
	}

	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		cfg.DBFile = envFile
	}

	return cfg
}
