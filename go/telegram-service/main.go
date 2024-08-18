package main

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/config"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/handler"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/repository"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/service"

	"fmt"
)

func main() {
	fmt.Println("Hello, telegram-service and David and Oleg!")

	logger, _ := logging.NewLogger("DEBUG")
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config successfully loaded")

	db, err := database.NewDatabase("postgresql://myuser:secret@db:5432/mydatabase", logger)
	if err != nil {
		logger.Fatal("Error while loading database: %s", err)
	}

	userRepo := repository.NewUserRepository(db, logger)
	messageRepo := repository.NewMessageRepository(db, logger)

	svc := service.NewTelegramService(userRepo, messageRepo, logger)

	handler := handler.NewTelegramHandler(svc, logger, cfg)

	handler.StartServer(cfg.ServerPort)

}
