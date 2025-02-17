package config

import (
	"Hekzory/tg-llm-bot/go/shared/config"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"fmt"
	"time"
)

type ServiceConfig struct {
	Config    config.Config         `toml:"default_config"`
	TGSConfig TelegramServiceConfig `toml:"tg_config"`
}

type TelegramServiceConfig struct {
	ActionPollingInterval time.Duration `toml:"action_polling_interval"`
	TelegramToken         string        `toml:"tg_token"`
}

func (cfg *ServiceConfig) DefaultTGConfig() {
	duration, _ := time.ParseDuration("5s")
	*cfg = ServiceConfig{
		Config: config.DefaultConfig(),

		TGSConfig: TelegramServiceConfig{
			ActionPollingInterval: duration,
			TelegramToken:         "",
		},
	}
}

func (cfg *ServiceConfig) LoadConfig(logger *logging.Logger) error {
	logger.Debug("Initializing default telegram service config")
	cfg.DefaultTGConfig()
	
	logger.Debug("Loading telegram service config")
	err := config.LoadConfig(logger, "telegram-service.toml", cfg)
	if err != nil {
		logger.Error("Failed to load telegram service config: %v", err)
		return err
	}
	logger.Debug("Full config dump - Default config: %+v, TG config: %+v", 
		cfg.Config, 
		cfg.TGSConfig)
		
	if cfg.TGSConfig.TelegramToken == "" {
		logger.Fatal("Telegram token is empty after config load!")
		return fmt.Errorf("telegram token is empty after config load")
	}
	
	return nil
}
