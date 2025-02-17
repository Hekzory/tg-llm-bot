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

	logger.Debug("Attempting to load config from: %s", configFilePath)

	if _, err := os.Stat(configFilePath); err != nil {
		logger.Debug("Config file not found, creating default config")
		serviceConfigText, err := toml.Marshal(cfg)
		if err != nil {
			logger.Error("Failed to marshal default config: %v", err)
			return err
		}
		if err = os.WriteFile(configFilePath, serviceConfigText, 0644); err != nil {
			logger.Error("Failed to write default config: %v", err)
			return  err
		}
		logger.Debug("Created default config file")
	}

	logger.Debug("Reading config file")
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		logger.Error("Failed to read config file: %v", err)
		return err
	}
	logger.Debug("Raw config content: %s", string(content))
	
	meta, err := toml.Decode(string(content), cfg)
	if err != nil {
		logger.Error("Failed to decode config file: %v", err)
		return err
	}
	logger.Debug("Successfully loaded config. Undecoded keys: %v", meta.Undecoded())
	return nil
}
