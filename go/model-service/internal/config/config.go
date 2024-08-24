package config

import "time"

type Config struct {
	ServerPort         int
	DatabaseURL        string
	ModelApiUrl        string
	LogLevel           string
	DBPollingInterval  time.Duration
	ModelAnswerTimeout time.Duration
	ModelName          string
}

func LoadConfig() (*Config, error) {

	pollingInterval, _ := time.ParseDuration("5s")
	modelAnswerTimeout, _ := time.ParseDuration("10m")

	cfg := &Config{
		ServerPort:         1111,
		DatabaseURL:        "postgresql://myuser:secret@db:5432/mydatabase",
		ModelApiUrl:        "http://ollama:11434/",
		LogLevel:           "DEBUG",
		DBPollingInterval:  pollingInterval,
		ModelAnswerTimeout: modelAnswerTimeout,
		// other considered choices are "mistral-nemo:12b-instruct-2407-q6_K", "llama3.1:70b-instruct-q4_K_S"
		ModelName: "gemma2:27b-instruct-q6_K",
	}
	return cfg, nil
}
