package config

import (
	"Hekzory/tg-llm-bot/go/shared/logging"
	"time"
)

type Config struct {
	ServerPort 			int
	DatabaseUrl			string
	LogLevel 			string
	DBpollingInterval 	time.Duration
	TelegramToken 		string
}

func LoadConfig(logger *logging.Logger) (*Config, error) {
	duration, _ := time.ParseDuration("30s")
	cfg := &Config{
		ServerPort: 		1111,
		DatabaseUrl:		"postgresql://myuser:secret:5432/mydatabase",
		LogLevel: 			"DEBUG",
		DBpollingInterval: 	duration,
		TelegramToken: 		"",
	}
	return cfg, nil
}
