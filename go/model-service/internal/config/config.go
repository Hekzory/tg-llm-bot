package config

type Config struct {
	ServerPort  int
	DatabaseURL string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort:  1111,
		DatabaseURL: "None",
		LogLevel:    "DEBUG",
	}
	return cfg, nil
}
