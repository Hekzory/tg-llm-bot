package config

import (
	"Hekzory/tg-llm-bot/go/shared/logging"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all configuration settings.
type Config struct {
	Logger   LoggerConfig   `toml:"logger"`
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
}

// LoggerConfig defines logging level
type LoggerConfig struct {
	Level string `toml:"log_level"`
}

// DatabaseConfig defines database connection parameters.
type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"db_name"`
}

// ServerConfig specifies server port.
type ServerConfig struct {
	Port int `toml:"port"`
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		LoggerConfig{
			Level: "DEBUG",
		},
		DatabaseConfig{
			Host:     "default",
			Port:     5432,
			Username: "user",
			Password: "default",
			Database: "default",
		},
		ServerConfig{Port: 1111},
	}
}

// NewConfig initializes and decodes the configuration.
func NewConfig(logger *logging.Logger) (*Config, error) {
	var config Config
	configDir := "config"
	configFilePath := filepath.Join(configDir, "config.toml")

	// Ensure config directory exists.
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0644); err != nil {
			logger.Fatal("Failed to create config directory")
		}
	}

	// Create default config if it doesn't exist.
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		defaultConfigText, err := toml.Marshal(DefaultConfig())
		if err != nil {
			logger.Fatal("Error encoding default config")
		}
		if err = os.WriteFile(configFilePath, defaultConfigText, 0644); err != nil {
			logger.Fatal("Failed to write default config file")
		}
	}

	// Decode the configuration file.
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		logger.Fatal("Configuration file could not be decoded")
	}

	return &config, nil
}
