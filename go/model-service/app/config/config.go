package config

import (
	"Hekzory/tg-llm-bot/go/shared/config"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"fmt"
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
			// "gemma2:27b-instruct-q6_K", "phi3:14b-medium-128k-instruct-q8_0", "command-r:35b-08-2024-q6_K"
			ModelName: "gemma2:9b-instruct-q4_K_S",
		},
	}
}

func (cfg *ServiceConfig) LoadConfig(logger *logging.Logger) error {
	logger.Debug("Initializing default model service config")
	cfg.DefaultModelConfig()
	
	logger.Debug("Loading model service config")
	err := config.LoadConfig(logger, "model-service.toml", cfg)
	if err != nil {
		logger.Error("Failed to load model service config: %v", err)
		return err
	}
	
	logger.Debug("Loaded config values - Model API URL: %s, Model Name: %s, Timeout: %v", 
		cfg.ModelConfig.ModelApiUrl,
		cfg.ModelConfig.ModelName,
		cfg.ModelConfig.ModelAnswerTimeout)
		
	if cfg.ModelConfig.ModelApiUrl == "" {
		logger.Fatal("Model API URL is empty after config load!")
		return fmt.Errorf("model API URL is empty after config load")
	}
	
	if cfg.ModelConfig.ModelName == "" {
		logger.Fatal("Model name is empty after config load!")
		return fmt.Errorf("model name is empty after config load")
	}
	
	return nil
}
