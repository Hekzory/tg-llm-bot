package config

import (
	"Hekzory/tg-llm-bot/go/shared/config"
	"Hekzory/tg-llm-bot/go/shared/logging"
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
	cfg.DefaultTGConfig()
	err := config.LoadConfig(logger, "telegram-service.toml", cfg)
	if err != nil {
		return err
	}
	return nil
}
