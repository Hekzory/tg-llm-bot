package config

import "time"

type Config struct {
	ServerPort        int
	DatabaseURL       string
	ModelApiUrl       string
	LogLevel          string
	DBPollingInterval time.Duration
	ModelName         string
}

func LoadConfig() (*Config, error) {

	duration, _ := time.ParseDuration("5s")

	cfg := &Config{
		ServerPort:        1111,
		DatabaseURL:       "postgresql://myuser:secret@db:5432/mydatabase",
		ModelApiUrl:       "http://ollama:11434/",
		LogLevel:          "DEBUG",
		DBPollingInterval: duration,
		ModelName:         "llama3.1:8b-instruct-q8_0",
	}
	return cfg, nil
}
