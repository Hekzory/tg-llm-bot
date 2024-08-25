package config

import (
	"Hekzory/tg-llm-bot/go/shared/logging"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all configuration settings.
type Config struct {
	LogLevel          string        `toml:"log_level"`
	DatabaseUrl       string        `toml:"database_url"`
	ServerPort        int           `toml:"port"`
	DBPollingInterval time.Duration `toml:"db_polling_interval"`
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	duration, _ := time.ParseDuration("5s")
	return Config{
		LogLevel:          "DEBUG",
		DatabaseUrl:       "postgresql://myuser:secret@db:5432/mydatabase",
		ServerPort:        1111,
		DBPollingInterval: duration,
	}
}

// NewConfig initializes and decodes the configuration.
func LoadConfig(logger *logging.Logger, fileDir string, cfg any) (error) {
	configDir := "config"
	configFilePath := filepath.Join(configDir, fileDir)

	if _, err := os.Stat(configFilePath); err != nil {
		serviceConfigText, err := toml.Marshal(cfg)
		if err != nil {
			return err
		}
		if err = os.WriteFile(configFilePath, serviceConfigText, 0644); err != nil {
			return  err
		}
	}

	if _, err := toml.DecodeFile(configFilePath, &cfg); err != nil {
		return err
	}
	return nil
}
