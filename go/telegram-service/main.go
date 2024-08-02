package main

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
	//"Hekzory/tg-llm-bot/go/telegram-service/internal/bot"
	"Hekzory/tg-llm-bot/go/shared/config"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/handler"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/repository"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/service"

	"fmt"


)

func main() {
	fmt.Println("Hello, telegram-service and David and Oleg!")
	
	logger, _ := logging.NewLogger("DEBUG")
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config successfully loaded")

	db := database.NewDatabase(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
	cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database))
	logger.Info("Database successfuly loaded")

	repo := repository.NewModelRepository(db, logger)

	svc := service.NewModelService(repo, logger)

	handler := handler.NewModelHandler(svc, logger)

	handler.StartServer(cfg.Server.Port)

}