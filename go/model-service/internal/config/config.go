package config

import (
	"Hekzory/tg-llm-bot/go/shared/config"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"time"
)

type ServiceConfig struct {
	Config      config.Config      `toml:"default_config"`
	ModelConfig ModelServiceConfig `toml:"model_config"`
}

type ModelServiceConfig struct {
	ModelApiUrl        string        `toml:"model_api_url"`
	ModelAnswerTimeout time.Duration `toml:"model_answer_timeout"`
	ModelName          string        `toml:"model_name"`
}

func (cfg *ServiceConfig) DefaultModelConfig() {
	modelAnswerTimeout, _ := time.ParseDuration("10m")
	*cfg = ServiceConfig{
		Config: config.DefaultConfig(),

		ModelConfig: ModelServiceConfig{
			ModelApiUrl:        "http://ollama:11434/",
			ModelAnswerTimeout: modelAnswerTimeout,
			// other considered choices are "mistral-nemo:12b-instruct-2407-q6_K", "llama3.1:70b-instruct-q4_K_S"
			// "gemma2:27b-instruct-q6_K", "phi3:14b-medium-128k-instruct-q8_0"
			ModelName: "gemma2:27b-instruct-q6_K",
		},
	}
}

func (cfg *ServiceConfig) LoadConfig(logger *logging.Logger) error {
	cfg.DefaultModelConfig()
	err := config.LoadConfig(logger, "model-service.toml", cfg)
	if err != nil {
		return err
	}
	return nil
}
