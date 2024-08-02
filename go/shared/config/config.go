package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Bot      BotConfig      `toml:"telegram"`
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
	Ollama   OllamaConfig   `toml:"ollama"`
}

type BotConfig struct {
	Token string `toml:"TELEGRAM_BOT_TOKEN"`
}

type DatabaseConfig struct {
	Host     string `toml:"DATABASE_HOST"`
	Port     int    `toml:"DATABASE_PORT"`
	Username string `toml:"DATABASE_USERNAME"`
	Password string `toml:"DATABASE_PASSWORD"`
	Database string `toml:"DATABASE_NAME"`
}

type ServerConfig struct {
	Port int `toml:"SERVER_PORT"`
}

type OllamaConfig struct {
	Count int `toml:"LAMA_COUNT"`
}

func NewConfig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("../shared/config/config.toml", &config)
	if err != nil {
		panic(err)
	}
	return &config, nil
}
