package main

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/app/config"
	"Hekzory/tg-llm-bot/go/telegram-service/app/handler"
	"Hekzory/tg-llm-bot/go/telegram-service/app/repository"
	"Hekzory/tg-llm-bot/go/telegram-service/app/service"

	"fmt"
)

func main() {
	fmt.Println("Hello, telegram-service and David and Oleg!")

	logger, _ := logging.NewLogger("DEBUG")
	
	var cfg config.ServiceConfig
	err := cfg.LoadConfig(logger)
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config successfully loaded")

	db, err := database.NewDatabase(cfg.Config.DatabaseUrl, logger)
	if err != nil {
		logger.Fatal("Error while loading database: %s", err)
	}

	userRepo := repository.NewUserRepository(db, logger)
	messageRepo := repository.NewMessageRepository(db, logger)
	conversationRepo := repository.NewConversationRepository(db, logger)

	svc := service.NewTelegramService(userRepo, messageRepo, conversationRepo, logger)

	handler := handler.NewTelegramHandler(svc, logger, &cfg)

	handler.StartServer(cfg.Config.ServerPort)

}
